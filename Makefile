.PHONY: build run test test-unit test-integration test-verbose test-coverage migrate migrate-status migrate-down sqlc-generate swagger-generate generate-mocks docker-build docker-run docker-stop

build:
	go build -o bin/server cmd/main.go

run:
	go run cmd/main.go

test:
	go test ./...

test-unit:
	go test ./... -short

test-integration:
	go test ./... -run Integration

test-verbose:
	go test ./... -v

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

migrate:
	cd internal/app/db/sqlite && goose sqlite3 ../../../tasks.db up

migrate-status:
	cd internal/app/db/sqlite && goose sqlite3 ../../../tasks.db status

sqlc-generate:
	cd internal/app/db/sqlite && sqlc generate

swagger-generate:
	swag init -g cmd/main.go -o docs

generate-mocks:
	mockgen -source=internal/app/db/sqlite/sqlc/querier.go -destination=internal/app/service/mocks/mock_querier.go -package=mocks

setup: sqlc-generate swagger-generate migrate

# Docker commands
docker-build:
	docker build -t ishare-task-api .

docker-run:
	docker run -d --name ishare-task-api -p 8080:8080 ishare-task-api

docker-stop:
	docker stop ishare-task-api && docker rm ishare-task-api

docker-dev: docker-build docker-run