.PHONY: up down test-up test

up:
	@echo "starting docker compose"
	docker compose up -d --build

down:
	@echo "shuting down docker compose"
	docker compose down

test-up:
	@echo "starting up REDIS"
	docker compose up -d redis

test:
	@echo "running tests"
	@echo "waiting redis to start up"
	sleep 3
	@go test -v ./...
	make down