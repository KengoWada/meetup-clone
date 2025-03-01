include .env/api.env
MIGRATIONS_PATH = ./cmd/migrate/migrations
DOCKER_COMPOSE = ./docker-compose.dev.yml

.PHONY: runserver
runserver:
	@go build -o ./bin/main ./cmd/api && ./bin/main

.PHONY: migration
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path $(MIGRATIONS_PATH) -database $(DB_ADDR) up $(num)

.PHONY: migrate-down
migrate-down:
	@migrate -path $(MIGRATIONS_PATH) -database $(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-force
migrate-force:
	@migrate -path $(MIGRATIONS_PATH) -database $(DB_ADDR) force $(filter-out $@,$(MAKECMDGOALS))

.PHONY: services-up
services-up:
	@docker compose -f $(DOCKER_COMPOSE) up --build -d

.PHONY: services-down
services-down:
	@docker compose -f $(DOCKER_COMPOSE) down

.PHONY: services-kill
services-kill:
	@docker compose -f $(DOCKER_COMPOSE) down -v

.PHONY: test
test:
	@SERVER_ENVIRONMENT=test go test -v ./...

.PHONY: test-cov
test-cov:
	@SERVER_ENVIRONMENT=test go test -v -coverprofile=coverage.out/cov.out ./... && go tool cover -html=coverage.out/cov.out -o=coverage.out/index.html

.PHONY: test-migrate-up
test-migrate-up:
	@migrate -path $(MIGRATIONS_PATH) -database $(TEST_DB_ADDR) up $(num)

.PHONY: test-migrate-down
test-migrate-down:
	@migrate -path $(MIGRATIONS_PATH) -database $(TEST_DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: test-migrate-force
test-migrate-force:
	@migrate -path $(MIGRATIONS_PATH) -database $(TEST_DB_ADDR) force $(filter-out $@,$(MAKECMDGOALS))

.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt
