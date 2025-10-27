DB_DSN=postgres://crud:crud@localhost:5432/crud?sslmode=disable
MIGRATIONS_DIR=./migrations

.PHONY: migrate-up migrate-down migrate-status migrate-version migrate-create db-up db-down

db-up:
	docker compose up -d postgres redis

db-down:
	docker compose down

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" status

migrate-version:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" version

migrate-create:
	@read -p "Migration name: " name; \
	goose -dir $(MIGRATIONS_DIR) create $$name sql
