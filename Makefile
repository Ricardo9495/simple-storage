# .env variables will be loaded
include .env
export

DB_URL=postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB_NAME)?sslmode=disable
MIGRATE=migrate -path internal/migrations -database "$(DB_URL)"

# Run the app
run:
	go run server.go

# Build the binary
build:
	go build -o bin/app server.go

# Run tests
test:
	go test ./...

# Run SQL migrations (up)
migrate-up:
	$(MIGRATE) up

# Rollback last migration
migrate-down:
	$(MIGRATE) down 1

# Lint (if using golangci-lint)
lint:
	golangci-lint run
