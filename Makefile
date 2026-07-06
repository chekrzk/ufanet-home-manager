COMPOSE ?= docker-compose
SERVICE ?=

.PHONY: help build run up stop down restart ps logs logs-service test clean clean-volumes infra

help:
	@echo "Targets:"
	@echo "  make build         Build all Docker images"
	@echo "  make run           Start full stack without rebuilding"
	@echo "  make up            Start full stack"
	@echo "  make stop          Stop containers without removing them"
	@echo "  make down          Stop and remove containers/network"
	@echo "  make restart       Restart full stack"
	@echo "  make ps            Show containers"
	@echo "  make logs          Follow all logs"
	@echo "  make logs-service SERVICE=api-gateway"
	@echo "  make infra         Start only Postgres and Redis"
	@echo "  make test          Run Go tests for workspace modules"
	@echo "  make clean         Stop and remove containers"
	@echo "  make clean-volumes Stop and remove containers with volumes"

build:
	$(COMPOSE) build

run:
	$(COMPOSE) up -d

up:
	$(COMPOSE) up -d --build

infra:
	$(COMPOSE) up -d postgres redis

stop:
	$(COMPOSE) stop

down:
	$(COMPOSE) down

restart: down up

ps:
	$(COMPOSE) ps

logs:
	$(COMPOSE) logs -f

logs-service:
	$(COMPOSE) logs -f $(SERVICE)

test:
	cd contracts && go test ./...
	cd api-gateway && go test ./...
	cd microservices/auth-service && go test ./...
	cd microservices/news-service && go test ./...
	cd microservices/notification-service && go test ./...
	cd microservices/profile-service && go test ./...
	cd microservices/requests-service && go test ./...

clean:
	$(COMPOSE) down --remove-orphans

clean-volumes:
	$(COMPOSE) down -v --remove-orphans
