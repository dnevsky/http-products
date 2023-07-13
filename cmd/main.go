package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"net/http"
	_ "net/http/pprof"

	server "github.com/dnevsky/http-products"
	"github.com/dnevsky/http-products/cache"
	red "github.com/dnevsky/http-products/cache/redis"
	"github.com/dnevsky/http-products/internal/handler"
	"github.com/dnevsky/http-products/internal/repository"
	"github.com/dnevsky/http-products/internal/repository/postgres"
	"github.com/dnevsky/http-products/internal/service"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {

	// profiler, metrics
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe("localhost:8081", nil)
	}()

	godotenv.Load(".env")
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("error while init logger: %s", err.Error())
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		sugar.Fatalf("error while init cache. check .env: %s", err.Error())
	}

	updateCacheTime, err := strconv.Atoi(os.Getenv("APP_CACHE_UPDATE"))
	if err != nil {
		sugar.Fatalf("error while init update cache gourutine. check .env: %s", err.Error())
	}

	cfgRedis := redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWD"),
		DB:       redisDB,
	}

	redis, err := red.NewRedisCache(&cfgRedis)
	if err != nil {
		sugar.Errorw("error while init cache: %s", err.Error())
	}

	db, err := postgres.NewPostgresDB(os.Getenv("DB_URI"))
	if err != nil {
		sugar.Errorw("error while init database: %s", err.Error())
	}

	cache := cache.NewCache(sugar, redis)
	repository := repository.NewRepository(sugar, db)
	services := service.NewService(sugar, cache)
	handlers := handler.NewHandler(sugar, services)

	// инициализируем кэш (собираем информацию с базы данных)
	// первый раз мы будем ждать инициализацию базы кэша, что-бы у нас она уже была, а дальше мы в фоне будем обновлять кэш
	waitFirstCache := make(chan struct{}, 1)
	updateCacheCtx, cancelCacheCtx := context.WithCancel(context.Background())
	go updateCache(updateCacheCtx, sugar, waitFirstCache, updateCacheTime, cache, repository)

	<-waitFirstCache

	srv := new(server.Server)

	go func() {
		if err := srv.Run(os.Getenv("APP_PORT"), handlers.InitRoutes()); err != nil {
			sugar.Errorf("%s", err.Error())
		}
	}()

	sugar.Infoln("http-products is work...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	sugar.Infoln("http-products shutting down")

	cancelCacheCtx()

	if err := srv.Shutdown(context.Background()); err != nil {
		sugar.Errorw("error while stop server: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		sugar.Errorw("error while stop db: %s", err.Error())
	}

	if err := redis.Close(); err != nil {
		sugar.Errorw("error while stop redis cache: %s", err.Error())
	}
}

// Обновляем кэш. Мы стартуем вечный цикл, в котором пытаемся считать данные из базы данных. Если не получается, то пытаемся через 5 сек
func updateCache(ctx context.Context, logger *zap.SugaredLogger, waitFirstCache chan struct{}, updateTimeCache int, cache *cache.Cache, repository *repository.Repository) {
	writeInChannel := false

	for {
		logger.Infoln("Start update cache...")
		products, err := repository.Product.GetAll(ctx)
		if err != nil {
			logger.Errorw(
				"error while get update from database & update the cache. Trying after 5 sec...",
				"error", err.Error(),
			)
			time.Sleep(time.Second * 5)
			continue
		}

		data := make([]string, 0, len(products))

		// JSON
		for _, p := range products {
			r, err := json.Marshal(&p)
			if err != nil {
				logger.Errorw(
					"error while marshalling data. Trying after 5 sec...",
					"error", err.Error(),
				)
				time.Sleep(time.Second * 5)
				continue
			}

			data = append(data, string(r))
		}

		// // SPLIT
		// for _, p := range products {
		// 	data = append(data, p.Id+":"+fmt.Sprint(p.Price))
		// }

		err = cache.Product.UpdateData(ctx, data)
		if err != nil {
			logger.Infow(
				"error while update the cache. Trying after 5 sec...",
				"error", err.Error(),
			)
			time.Sleep(time.Second * 5)
			continue
		}

		if !writeInChannel {
			writeInChannel = true
			waitFirstCache <- struct{}{}
		}

		logger.Infoln("Cache sucefully updated. Sleep...")

		time.Sleep(time.Second * time.Duration(updateTimeCache))
	}
}
