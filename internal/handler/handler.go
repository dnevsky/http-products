package handler

import (
	"github.com/dnevsky/http-products/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	logger   *zap.SugaredLogger
	services *service.Service
}

func NewHandler(logger *zap.SugaredLogger, services *service.Service) *Handler {
	return &Handler{logger: logger, services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/", h.getAll)

	return router
}
