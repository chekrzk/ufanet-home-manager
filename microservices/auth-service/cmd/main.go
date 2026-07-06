package main

import (
	"log"

	"github.com/chekrzk/ufanet-home-manager/auth-service/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

	// Структура auth-service:
	// cmd - минимальная точка запуска сервиса.
	// api/interceptors - gRPC middleware для логирования и zero-trust role checks.
	// api/server - gRPC procedures и конвертация transport-моделей в service commands.
	// internal/app - сборка ресурсов, репозитория, service layer и gRPC server.
	// internal/service - бизнес-сценарии авторизации: register, login, refresh.
	// internal/repository - GORM-доступ к users и миграции таблиц.
	// internal/models - доменные модели и команды сервиса.
	// internal/hasher и internal/jwt - инфраструктурные политики паролей и токенов.
	// resources - env config, logger, DB и общий lifecycle внешних ресурсов.
}
