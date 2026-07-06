# Ufanet Home Manager

Система для ЖЭУ/управляющей организации: жители видят новости по дому, получают уведомления, привязывают профиль к адресу и создают заявки на работы. Управляющий публикует новости, работает с домами и работниками. Работники указывают специализацию и доступность, получают назначенные заявки и меняют их статус.

## Запуск через Make

Нужны Docker Desktop и `make`.

```powershell
make up
```

Команда соберет образы и поднимет весь стек:

- `front` на `http://localhost:5173`
- `api-gateway` на `http://localhost:8086`
- gRPC микросервисы на портах `50051-50055`
- PostgreSQL на `localhost:5432`
- Redis на `localhost:6379`

Полезные команды:

```powershell
make run
make stop
make down
make ps
make logs
make logs-service SERVICE=api-gateway
make infra
make test
```

Назначение команд:

- `make up` - собрать и запустить весь проект.
- `make run` - запустить без пересборки образов.
- `make stop` - остановить контейнеры без удаления.
- `make down` - остановить и удалить контейнеры/сеть.
- `make ps` - посмотреть состояние контейнеров.
- `make logs` - смотреть логи всех сервисов.
- `make logs-service SERVICE=...` - смотреть логи одного сервиса.
- `make infra` - поднять только PostgreSQL и Redis.
- `make test` - запустить Go-тесты по workspace-модулям.

## Запуск через Docker Compose

Сборка и запуск:

```powershell
docker-compose up -d --build
```

Запуск без пересборки:

```powershell
docker-compose up -d
```

Остановка без удаления контейнеров:

```powershell
docker-compose stop
```

Остановка с удалением контейнеров и сети:

```powershell
docker-compose down
```

Логи:

```powershell
docker-compose logs -f
docker-compose logs -f api-gateway
```

Проверка gateway:

```powershell
curl http://localhost:8086/health
```

После запуска frontend доступен по адресу:

```text
http://localhost:5173
```

API gateway:

```text
http://localhost:8086
```

## Суть проекта

Проект автоматизирует взаимодействие жителей дома, управляющей организации и исполнителей работ.

Основные сценарии:

- регистрация и авторизация пользователей;
- разделение ролей: житель, управляющий/администратор, работник;
- привязка профиля жителя или работника к дому;
- публикация новостей по дому;
- получение уведомлений о новостях и изменениях заявок;
- создание заявок на электрика, монтажера и других работников;
- назначение исполнителя и ведение статуса заявки;
- управление списком работников и их доступностью.

Система построена так, чтобы авторизация, профиль, новости, уведомления и заявки были независимыми частями, но работали через единый API gateway.

## Особенности реализации

### Архитектура

Проект разделен на API gateway, frontend, gRPC contracts и микросервисы:

```text
api-gateway/
contracts/
front/
microservices/
  auth-service/
  news-service/
  notification-service/
  profile-service/
  requests-service/
```

`api-gateway` принимает HTTP-запросы от frontend/Postman и вызывает микросервисы по gRPC.

Микросервисы:

- `auth-service` - регистрация, логин, refresh/access JWT.
- `news-service` - список и создание новостей.
- `notification-service` - устройства, уведомления, Redis Stream events.
- `profile-service` - профили, привязка к дому, работники, доступность работников.
- `requests-service` - заявки, статусы, комментарии, события для уведомлений.

### Слои приложения

В gateway:

- `cmd` - точка входа.
- `config` - env-конфигурация.
- `infra/clients` - gRPC-клиенты к микросервисам.
- `internal/app` - DI-сборка приложения.
- `internal/handlers` - HTTP handlers.
- `internal/services` - сценарии gateway.
- `internal/middlewares` - JWT, роли, blacklist, rate limit, logger.
- `internal/router` - регистрация маршрутов.
- `internal/models` - DTO, domain-модели и константы.
- `internal/errors` - единый формат ошибок.

В микросервисах:

- `cmd` - точка входа.
- `api/interceptors` - gRPC middleware и zero-trust проверки ролей.
- `api/server` - gRPC procedures и конвертация proto/domain.
- `internal/app` - сборка ресурсов, repository, service и gRPC server.
- `internal/service` - бизнес-логика.
- `internal/repository` - GORM-доступ к БД и миграции.
- `internal/models` - доменные модели и команды.
- `resources` - env, logger, DB, Redis/gRPC connections.
- `infra` - внешние клиенты или Redis Stream publisher.

### Данные и инфраструктура

PostgreSQL используется как основное хранилище. В docker-compose создаются отдельные базы:

- `ufanet_auth`
- `ufanet_news`
- `ufanet_notifications`
- `ufanet_profile`
- `ufanet_requests`

Redis используется в `notification-service` для stream-событий уведомлений.

GORM выполняет миграции таблиц при старте микросервисов.

### Безопасность

- Авторизация через JWT.
- Access/refresh tokens выдаются `auth-service`.
- Пароли хешируются.
- API gateway проверяет JWT и роли на HTTP-уровне.
- Микросервисы дополнительно проверяют роли через gRPC interceptors, чтобы каждый сервис не доверял только gateway.
- Blacklist middleware позволяет отклонять отозванные токены.
- Rate limit ограничивает частоту запросов.

### Контракты

gRPC-контракты лежат в `contracts/`.

Сгенерированный Go-код используется gateway и микросервисами как общий контракт взаимодействия.

### Тесты

Добавлены unit-тесты для service layer микросервисов и для service/handler layer API gateway.

Зависимости в тестах мокируются через mockery.

Запуск:

```powershell
make test
```

Или отдельно:

```powershell
cd api-gateway
go test ./...

cd ../microservices/auth-service
go test ./...
```

## Основные адреса

```text
Frontend:    http://localhost:5173
API Gateway: http://localhost:8086
PostgreSQL:  localhost:5432
Redis:       localhost:6379
```

gRPC-порты:

```text
auth-service:         50051
news-service:         50052
requests-service:     50053
notification-service: 50054
profile-service:      50055
```
