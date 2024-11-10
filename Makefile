# Makefile

# Команда по умолчанию
.DEFAULT_GOAL := migrate

# Переменные
STORAGE_PATH := ./storage/sso.db
MIGRATIONS_PATH := ./migrations

# Цель для миграции
migrate:
	go run ./cmd/migrator --storage-path=$(STORAGE_PATH) \
		--migrations-path=$(MIGRATIONS_PATH)