.PHONY: build run shutdown postgres create-migrate migrate redis test cover

test:
	go test -short -count=1 -coverprofile=coverage.out ./...

cover:
	go test -short -count=1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# build:
# 	docker build --tag dnevsky/http-products .

# run:
# 	docker-compose up -d dnevsky/http-products

shutdown:
	docker-compose down

postgres:
	docker-compose up -d postgres
	
redis:
	docker-compose up -d redis

create-migrate:
	migrate create -ext sql -dir ./migrations -seq init

migrate:
	docker-compose run migrate-postgres