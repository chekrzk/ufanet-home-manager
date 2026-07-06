package main

import (
	"log"

	"github.com/chekrzk/ufanet-home-manager/profile-service/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

	// Структура profile-service:
	// cmd - минимальная точка запуска сервиса.
	// api/interceptors - gRPC middleware для логирования и zero-trust role checks.
	// api/server - gRPC procedures и конвертация proto-запросов в service commands.
	// internal/app - сборка resources, repository, service layer и gRPC server.
	// internal/service - бизнес-сценарии профиля, дома, работников и доступности.
	// internal/repository - GORM-доступ к profiles/workers/availability и миграции.
	// internal/models - доменные модели профиля, работника и фильтров.
	// resources - env config, logger, DB и общий lifecycle ресурсов.
}
