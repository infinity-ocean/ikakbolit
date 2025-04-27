include .env
export

# Директория, в которой хранятся исполняемые
# файлы проекта и зависимости, необходимые для сборки.
LOCAL_BIN := $(CURDIR)/bin
MIGRATIONS_DIR := ./migrations
PROTO_SRC := 3-api-grpc-Homework/grpc/ikakbolit.proto
PROTO_OUT := 3-api-grpc-Homework/grpc/ikakbolit

start-infra:
	docker-compose up -d

stop-infra:
	docker-compose down
	
# ДЗ №3 кодогенерация на основе openapi 
swagger-gen:
	swagger generate server -f internal/3-api-grpc-Homework/docs/swagger.yaml -A ikakbolit --target internal/3-api-grpc-Homework/server

proto-gen:
	protoc \
		--proto_path=3-api-grpc-Homework/grpc \
		--go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
		$(PROTO_SRC)
	
migration-up:
	$(LOCAL_BIN)/goose $(opts) -allow-missing -dir $(MIGRATIONS_DIR) postgres "$(POSTGRES_DSN)" up

migration-down:
	$(LOCAL_BIN)/goose $(opts) -dir $(MIGRATIONS_DIR) postgres "$(POSTGRES_DSN)" down

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.18.0
	GOBIN=$(LOCAL_BIN) go install github.com/go-swagger/go-swagger/cmd/swagger@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest