# Source: https://dou.ua/forums/topic/34806/

migration_dir := ./migrations
migrations_table := schema_migrations

stage := $(or $(YAAWS_STAGE), local)

include Makefile.$(stage)

goose-install:
	go get -u github.com/pressly/goose/cmd/goose

MIGRATION_NAME=$(or $(MIGRATION), init)
migrate-create:
	mkdir -p $(migration_dir)
	goose -dir $(migration_dir) -table $(migrations_table) postgres $(POSTGRES_URI) create $(MIGRATION_NAME) sql

migrate-up:
	goose -dir $(migration_dir) -table $(migrations_table) postgres $(POSTGRES_URI) up
migrate-redo:
	goose -dir $(migration_dir) -table $(migrations_table) postgres $(POSTGRES_URI) redo
migrate-down:
	goose -dir $(migration_dir) -table $(migrations_table) postgres $(POSTGRES_URI) down
migrate-reset:
	goose -dir $(migration_dir) -table $(migrations_table) postgres $(POSTGRES_URI) reset
migrate-status:
	goose -dir $(migration_dir) -table $(migrations_table) postgres $(POSTGRES_URI) status
