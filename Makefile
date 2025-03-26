# these setup environment variable
GOOSE_DRIVER ?= postgres
GOOSE_DBSTRING="host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
GOOSE_MIGRATION_DIR ?= sql/schema
migrate-up:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=$(GOOSE_MIGRATION_DIR) up

migrate-down:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=$(GOOSE_MIGRATION_DIR) down

migrate-reset:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=$(GOOSE_MIGRATION_DIR) reset


.PHONY: migrate-up, migrate-down, migrate-reset