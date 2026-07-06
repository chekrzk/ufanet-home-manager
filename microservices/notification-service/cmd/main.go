package main

import (
	"log"

	"github.com/chekrzk/ufanet-home-manager/notification-service/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

	// Структура notification-service:
	// cmd - минимальная точка запуска сервиса.
	// api/interceptors - gRPC middleware для логирования и zero-trust role checks.
	// api/server - gRPC procedures и конвертация transport-запросов в service commands.
	// infra/streams - Redis Stream publisher для доставки событий уведомлений.
	// internal/app - сборка resources, repository, service layer и gRPC server.
	// internal/service - бизнес-сценарии устройств, уведомлений и отметки прочтения.
	// internal/repository - GORM-доступ к devices/notifications и миграции таблиц.
	// internal/models - доменные модели уведомлений, устройств и pagination.
	// resources - env config, logger, DB, Redis client и общий lifecycle ресурсов.
}
