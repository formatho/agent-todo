.PHONY: help build run test clean docker-up docker-down migrate-up migrate-down

help:
	@echo "Available commands:"
	@echo "  make build        - Build the backend"
	@echo "  make run          - Run the backend server"
	@echo "  make test         - Run tests"
	@echo "  make migrate-up   - Run database migrations"
	@echo "  make migrate-down - Rollback database migrations"
	@echo "  make docker-up    - Start all services with docker"
	@echo "  make docker-down  - Stop all services"
	@echo "  make swagger      - Generate Swagger documentation"

build:
	go build -o main ./cmd/api

run:
	go run cmd/api/main.go

test:
	go test -v ./tests/...

migrate-up:
	migrate -path db/migrations -database "postgres://agent_todo:agent_todo_pass@localhost:5432/agent_todo?sslmode=disable" up

migrate-down:
	migrate -path db/migrations -database "postgres://agent_todo:agent_todo_pass@localhost:5432/agent_todo?sslmode=disable" down

swagger:
	swag init -g cmd/api/main.go

docker-up:
	cd .. && docker compose up -d

docker-down:
	cd .. && docker compose down

clean:
	rm -f main
	go clean
