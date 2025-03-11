include .env
export

# Директория, в которой хранятся исполняемые
# файлы проекта и зависимости, необходимые для сборки.
LOCAL_BIN := $(CURDIR)/bin
MIGRATIONS_DIR := ./migrations

start-infra:
	docker-compose up -d

stop-infra:
	docker-compose down
	
print-dsn:
	echo $(POSTGRES_DSN)
	
migration-up:
	$(LOCAL_BIN)/goose $(opts) -allow-missing -dir $(MIGRATIONS_DIR) postgres "$(POSTGRES_DSN)" up

migration-down:
	$(LOCAL_BIN)/goose $(opts) -dir $(MIGRATIONS_DIR) postgres "$(POSTGRES_DSN)" down

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.18.0