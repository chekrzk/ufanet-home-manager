package main

import (
	"log"

	"github.com/chekrzk/ufanet-home-manager/news-service/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

	// Структура news-service:
	// cmd - минимальная точка запуска сервиса.
	// api/interceptors - gRPC middleware для логирования и zero-trust role checks.
	// api/server - gRPC procedures и конвертация proto-запросов в service commands.
	// infra/clients - исходящие клиенты к другим микросервисам, например notifications.
	// internal/app - сборка resources, repository, service layer и gRPC server.
	// internal/service - бизнес-сценарии ленты: список новостей и публикация новости.
	// internal/repository - GORM-доступ к news и миграции таблиц.
	// internal/models - доменные модели, фильтры и команды news-service.
	// resources - env config, logger, DB, gRPC connections и их lifecycle.
}
