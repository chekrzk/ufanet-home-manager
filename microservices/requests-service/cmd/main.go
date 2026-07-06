package main

import (
	"log"

	"github.com/chekrzk/ufanet-home-manager/requests-service/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

	// Структура requests-service:
	// cmd - минимальная точка запуска сервиса.
	// api/interceptors - gRPC middleware для логирования и zero-trust role checks.
	// api/server - gRPC procedures и конвертация proto-запросов в service commands.
	// infra/clients - исходящие клиенты к другим микросервисам, например notifications.
	// internal/app - сборка resources, repository, service layer и gRPC server.
	// internal/service - бизнес-сценарии заявок: создание, статусы, комментарии.
	// internal/repository - GORM-доступ к maintenance_requests/comments и миграции.
	// internal/models - доменные модели заявок, команд, статусов и pagination.
	// resources - env config, logger, DB, gRPC connections и их lifecycle.
}
