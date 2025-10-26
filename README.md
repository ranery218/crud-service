# CRUD Users Service

Небольшой учебный сервис на Go, реализующий регистрацию, авторизацию и работу с профилем пользователя на Postgres с cookie-сессиями.

## Стек

- Go 1.25
- Chi (маршрутизация)
- pgx/pgxpool (Postgres)
- Goose (миграции БД)
- In-memory session store
- Docker Compose (локальный Postgres)

## Подготовка окружения


1. Создайте `.env` (используется VS Code/Makefile) и добавьте, например:
   ```
   POSTGRES_PASSWORD=crud
   ```

2. Поднимите Postgres:
   ```bash
   make db-up
   ```

3. Накатите миграции:
   ```bash
   make migrate-up
   ```

## Запуск

```bash
go run ./cmd
```

Сервис слушает адрес из `config.yaml` (по умолчанию `localhost:8080`). При необходимости можно добавить цель `make run`.

## Основные эндпоинты

| Метод | Путь             | Описание                                    |
|-------|------------------|---------------------------------------------|
| POST  | `/auth/register` | регистрация пользователя                    |
| POST  | `/auth/login`    | логин, выставляет cookie `session_id`       |
| POST  | `/auth/logout`   | logout, удаляет текущую сессию              |
| PATCH | `/users/me`      | обновление текущего пользователя (cookie)   |
| DELETE| `/users/me`      | удаление аккаунта (cookie)                  |

Структуры тел запросов/ответов см. в `internal/transport/http/dto.go`.

## Тесты

```bash
go test ./...
```

## Полезные команды Makefile

| Команда              | Действие                                      |
|----------------------|-----------------------------------------------|
| `make db-up`         | поднять Postgres из docker-compose            |
| `make db-down`       | остановить Postgres                           |
| `make migrate-up`    | накатывает все миграции                       |
| `make migrate-down`  | откатывает последнюю миграцию                 |
| `make migrate-status`| показывает статус миграций                    |
| `make migrate-create`| создаёт шаблон новой миграции (требуется goose) |

## Postman / curl

- Регистрация:
  ```bash
  curl -X POST http://localhost:8080/auth/register \
       -H "Content-Type: application/json" \
       -d '{"user_name":"demo","email":"demo@example.com","password":"Test1234!"}'
  ```

- Логин:
  ```bash
  curl -i -X POST http://localhost:8080/auth/login \
       -H "Content-Type: application/json" \
       -d '{"email":"demo@example.com","password":"Test1234!"}'
  ```

  Сохраните cookie `session_id` и используйте её для защищённых запросов.

## Структура проекта

```
cmd/                      – точка входа, DI и запуск сервера
internal/
  adapters/               – Postgres-репозиторий, хешер, сессии, ID-генератор
  config/                 – загрузка конфигурации
  domain/                 – сущности и контракты домена
  services/user/          – бизнес-логика юзкейсов
  transport/http/         – хендлеры, маршруты, middleware
migrations/               – SQL-миграции для Postgres
Makefile                  – команды для БД и миграций
```
