ifneq (,$(wildcard ./.env))
	include .env
	export
endif

run:
	go run ./cmd/api/

build:
	go build -o bin/api ./cmd/api/

lint:
	golangci-lint run

fmt:
	gofmt -w .
	goimports -w .

tidy:
	go mod tidy

test:
	go test ./... -v

test-integration:
	go test ./... -tags=integration -v

test-all:
	go test ./... -tags=integration -v -cover

migrate-up:
	migrate -path migrations -database "mysql://${DATABASE_USER}:${DATABASE_PASSWORD}@tcp(${DATABASE_HOST}:${DATABASE_PORT})/${DATABASE_NAME}" up

migrate-down:
	migrate -path migrations -database "mysql://${DATABASE_USER}:${DATABASE_PASSWORD}@tcp(${DATABASE_HOST}:${DATABASE_PORT})/${DATABASE_NAME}" down 1
