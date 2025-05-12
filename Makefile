include .env
MIGRATIONS_DIR := ./db/migrations
PROTO_SRC := 3-api-grpc-Homework/grpc/ikakbolit.proto
PROTO_OUT := 3-api-grpc-Homework/grpc/ikakbolit

POSTGRES_DSN := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSL)

run:
	docker-compose up --build -d   

start-infra:
	docker-compose up --build -d

stop-infra:
	docker-compose down
	
unit-test:
	go test -short -v -timeout=2m -cover -coverpkg=./... -coverprofile=cover.out ./...

test:
	$(MAKE) goose-up
	go test -v -timeout=2m -coverpkg=./... -coverprofile=cover.out ./...
	grep -v "\.gen\.go\>" tests/cover.out | grep -v '_test\>' | grep -v '\<tests\>' > tests/cover.skipgen.out
	go tool cover -func=tests/cover.skipgen.out #go tool cover -html=tests/cover.skipgen.out
	$(MAKE) goose-down

lint: nilaway
	golangci-lint run --config .golangci.yml

nilaway:
	nilaway -include-pkgs="ikakbolit" ./...

goose-up:
	goose -allow-missing -dir $(MIGRATIONS_DIR) postgres "$(POSTGRES_DSN)" up

goose-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(POSTGRES_DSN)" down

install-deps:
	go install -tags='no_clickhouse no_libsql no_mssql no_mysql no_sqlite3 no_vertica no_ydb' github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install go.uber.org/nilaway/cmd/nilaway@latest

# кодогенерация на основе openapi 
swagger-gen:
	swagger generate server -f internal/3-api-grpc-Homework/docs/swagger.yaml -A ikakbolit --target internal/3-api-grpc-Homework/server

proto-gen:
	protoc \
		--proto_path=3-api-grpc-Homework/grpc \
		--go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
		$(PROTO_SRC)