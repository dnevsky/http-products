.PHONY: build run shutdown postgres create-migrate migrate redis

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