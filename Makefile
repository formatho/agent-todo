# Agent Todo Platform - Docker Commands
# Usage: make <command>

.PHONY: help build up down logs clean rebuild ps

help:
	@echo "Agent Todo Platform - Docker Commands"
	@echo ""
	@echo "  make build     - Build all Docker images"
	@echo "  make up        - Start all services"
	@echo "  make down      - Stop all services"
	@echo "  make logs      - Show logs from all services"
	@echo "  make logs-b    - Show backend logs"
	@echo "  make logs-f    - Show frontend logs"
	@echo "  make logs-db   - Show database logs"
	@echo "  make rebuild   - Rebuild all images (no cache)"
	@echo "  make clean     - Remove all containers, images, and volumes"
	@echo "  make ps        - Show running containers"
	@echo "  make shell-b   - Open shell in backend container"
	@echo "  make shell-db  - Open psql in database container"
	@echo ""

build:
	docker compose build

up:
	docker compose up -d
	@echo ""
	@echo "Services starting..."
	@echo "Frontend: http://localhost:3000"
	@echo "Backend:  http://localhost:8080"
	@echo "Swagger:  http://localhost:8080/docs"
	@echo ""

up-logs:
	docker compose up

down:
	docker compose down

logs:
	docker compose logs -f

logs-b:
	docker compose logs -f backend

logs-f:
	docker compose logs -f frontend

logs-db:
	docker compose logs -f postgres

rebuild:
	docker compose down
	docker compose build --no-cache
	docker compose up -d

clean:
	docker compose down -v --rmi local
	@echo "Cleaned up all containers, volumes, and local images"

ps:
	docker compose ps

shell-b:
	docker compose exec backend sh

shell-db:
	docker compose exec postgres psql -U agent_todo -d agent_todo

# Production commands
prod-up:
	docker compose -f docker-compose.prod.yml up -d

prod-down:
	docker compose -f docker-compose.prod.yml down

prod-logs:
	docker compose -f docker-compose.prod.yml logs -f

prod-rebuild:
	docker compose -f docker-compose.prod.yml down
	docker compose -f docker-compose.prod.yml build --no-cache
	docker compose -f docker-compose.prod.yml up -d
