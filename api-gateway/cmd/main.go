package main

import (
	"log"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

	// Структура api-gateway:
	// cmd - минимальная точка входа, чтобы запуск не смешивался со сборкой приложения.
	// config - чтение env и адресов внешних микросервисов.
	// infra/clients - gRPC-обертки над сгенерированными contract-клиентами.
	// internal/app - DI-сборка Fiber, middleware, handlers, services и gRPC connections.
	// internal/handlers - HTTP-слой: парсит DTO, достает auth context и вызывает сервисы.
	// internal/services - слой сценариев gateway, работающий только с domain-моделями.
	// internal/middlewares - сквозные HTTP-проверки: JWT, роли, rate limit, blacklist, logger.
	// internal/router - единая регистрация публичных и защищенных маршрутов.
	// internal/models - DTO для HTTP, domain-модели и константы контекста/ролей.
	// internal/errors - общий формат ошибок и HTTP-ответов.
}
