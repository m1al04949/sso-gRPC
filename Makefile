# Makefile

# Команда по умолчанию
.DEFAULT_GOAL := migrate

# Переменные
STORAGE_PATH := ./storage/sso.db
MIGRATIONS_PATH := ./migrations
MIGRATIONS_TESTS_PATH = ./tests/migrations
MIGRATIONS_TESTS_TABLE = migrations_test

# Цель для миграции
migrate:
	go run ./cmd/migrator --storage-path=$(STORAGE_PATH) \
		--migrations-path=$(MIGRATIONS_PATH)
test_migrate:
	go run ./cmd/migrator --storage-path=$(STORAGE_PATH) \
		--migrations-path=$(MIGRATIONS_TESTS_PATH) \
		--migrations-table=$(MIGRATIONS_TESTS_TABLE)

# Запуск SSO
start:
	go run ./cmd/sso/main.go

test:
	go test ./tests/auth_register_login_test.go