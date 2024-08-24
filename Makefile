#!make
include .make.env

DOCKER_COMPOSE ?= "docker-compose"
DOCKER_COMPOSE_FILE ?= docker-compose.yaml

# colors
CYAN=\033[0;36m
NC=\033[0m

.PHONY: help

## Show help screen
help:
	@echo "[~~~~~ $(CYAN)HELP$(NC) ~~~~~]"
	@printf "Available targets:\n\n"
	@awk '/^[a-zA-Z\-\_0-9%:\\]+/ { \
	  helpMessage = match(lastLine, /^## (.*)/); \
	  if (helpMessage) { \
		helpCommand = $$1; \
		helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
  gsub("\\\\", "", helpCommand); \
  gsub(":+$$", "", helpCommand); \
		printf "  \x1b[32;01m%-35s\x1b[0m %s\n", helpCommand, helpMessage; \
	  } \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST) | sort -u
	@printf "\n"

## Install Dependencies
install:
	make install/db-migrator
	make install/sqlc

## Install db migrator
install/db-migrator:
	curl -sSf https://atlasgo.sh | sh

## Install sqlc
install/sqlc:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

## Running docker compose for dependencies only
dev/dep:
	@echo "[~~~~~ $(CYAN)DEV DEPENDENCIES$(NC) ~~~~~]"
	@echo "running service dependencies using $(DOCKER_COMPOSE) 💨"
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d

## Stop docker container
dev/down:
	@echo "[~~~~~ $(CYAN)DEV CLEAR$(NC) ~~~~~]"
	@echo "Removing docker container service 🧹"
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down

## Migrate database
run/migrate:
	atlas migrate apply --dir "file://db/migrations" --url "postgres://local:supersecret@localhost:5477/wallet?sslmode=disable"

## Run rest server
run/bin:
	@go run cmd/rest/main.go

## Run all
run:
	@go mod tidy
	@make dev/dep
	@make run/migrate
