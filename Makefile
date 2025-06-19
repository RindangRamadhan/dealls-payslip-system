ifneq ($(wildcard .env),)
include .env
export
else
$(warning WARNING: .env file not found! Using .env.example)
include .env.example
export
endif

LOCAL_BIN:=$(CURDIR)/bin
BASE_STACK = docker compose -f docker-compose.yml

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

compose-up: ### Run docker compose with db only
	$(BASE_STACK) up --build -d db && docker compose logs -f
.PHONY: compose-up

compose-down: ### Stop and remove containers
	$(BASE_STACK) down --remove-orphans
.PHONY: compose-down

swag-v1: ### Generate Swagger docs
	swag init -g internal/controller/http/router.go
.PHONY: swag-v1

deps: ### Run go mod tidy & verify
	go mod tidy && go mod verify
.PHONY: deps

deps-audit: ### Check dependency vulnerabilities
	govulncheck ./...
.PHONY: deps-audit

format: ### Format code
	gofumpt -l -w .
	gci write . --skip-generated -s standard -s default
.PHONY: format

run: deps swag-v1 ### Run app locally
	go mod download && \
	CGO_ENABLED=0 go run -tags migrate ./cmd/app
.PHONY: run

docker-rm-volume: ### Remove postgres volume
	docker volume rm dealls-payslip-system_pg-data
.PHONY: docker-rm-volume

linter-golangci: ### Run GolangCI-Lint
	golangci-lint run
.PHONY: linter-golangci

linter-hadolint: ### Lint Dockerfiles
	git ls-files --exclude='Dockerfile*' --ignored | xargs hadolint
.PHONY: linter-hadolint

linter-dotenv: ### Lint dotenv files
	dotenv-linter
.PHONY: linter-dotenv

test: ### Run tests
	go test -v -race -covermode atomic -coverprofile=coverage.txt ./internal/...
.PHONY: test

mock: ### Generate mocks
	mockgen -source ./internal/repo/contracts.go -package usecase_test > ./internal/usecase/mocks_repo_test.go
	mockgen -source ./internal/usecase/contracts.go -package usecase_test > ./internal/usecase/mocks_usecase_test.go
.PHONY: mock

migrate-create: ### Create new migration
	migrate create -ext sql -dir migrations '$(word 2,$(MAKECMDGOALS))'
.PHONY: migrate-create

migrate-up: ### Run database migration
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up
.PHONY: migrate-up

seed: ## Run database seeder via Go script
	go run migrations/seeders/create_users_seeder.go
.PHONY: seed

pre-commit: swag-v1 mock format linter-golangci test ### Run before commit
.PHONY: pre-commit