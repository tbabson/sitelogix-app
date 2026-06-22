.PHONY: build run dev test migrate-up migrate-down docker-up docker-down lint tidy

APP_NAME=sitelogix
BINARY=bin/$(APP_NAME)
MIGRATE=migrate -path ./migrations -database "$(DATABASE_URL)"

build:
	go build -o $(BINARY) ./cmd/server

run: build
	./$(BINARY)

dev:
	go run ./cmd/server

test:
	go test ./... -v -cover

tidy:
	go mod tidy

lint:
	golangci-lint run ./...

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down

migrate-create:
	migrate create -ext sql -dir ./migrations -seq $(name)

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-build:
	docker compose build

docker-logs:
	docker compose logs -f app

docker-reset:
	docker compose down -v && docker compose up -d

db-shell:
	docker compose exec postgres psql -U sitelogix -d sitelogix
