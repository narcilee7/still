.PHONY: help install proto infra infra-down backend mobile build lint test migrate env clean stop

help: ## Show this help
	@echo "Still — available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## Install all workspace dependencies
	yarn install

proto: ## Generate Go server code and TypeScript SDK from proto files
	yarn generate:proto

env: ## Copy per-app .env.example into dev env files (will not overwrite existing files)
	cp -n apps/backend/.env.example apps/backend/.env.development 2>/dev/null || true
	cp -n apps/mobile/.env.example apps/mobile/.env.development 2>/dev/null || true
	@echo "Dev env files created. Edit them to add OPENAI_API_KEY, CLERK_SECRET_KEY and EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY."

infra-up: ## Start local infrastructure (PostgreSQL + MinIO)
	docker compose up -d

infra-down: ## Stop local infrastructure
	docker compose down

backend: ## Start the Go backend (loads apps/backend/.env.development)
	yarn dev:backend

mobile: ## Start the Expo mobile app (loads apps/mobile/.env.development)
	cd apps/mobile && yarn start

dev: infra ## Start infrastructure and then the backend
	yarn dev:backend

build: ## Build all workspaces
	yarn build

lint: ## Run linters across workspaces
	yarn lint

test: ## Run tests across workspaces
	yarn test

migrate: ## Run database migrations up
	yarn migrate up

migrate-down: ## Rollback one database migration
	yarn migrate down

clean: ## Stop infrastructure and remove volumes + build artifacts
	docker compose down -v
	rm -rf apps/backend/bin
	@echo "Cleaned. Run 'make infra' to recreate infrastructure."

stop: ## Stop backend and infrastructure
	docker compose down
