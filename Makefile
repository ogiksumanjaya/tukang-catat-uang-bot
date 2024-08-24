#!/bin/bash

export REPO_NAME=twee-reseller-app
export NOW=$(shell date +"%Y/%m/%d %T")

# Create DB Connection String from env
DB_CONNECTION_STRING := $(TGBOT_DB_POSTGRESQL_URL)

include-env:
	@export $(shell sed 's/=.*//' .env) && echo "Environment variables loaded."

migration-create:
	migrate create -ext sql -dir migrations/postgres -seq $(name)

migrate-up:
	migrate -database '$(DB_CONNECTION_STRING)' -path migrations/postgres up

migrate-down:
	migrate -database '$(DB_CONNECTION_STRING)' -path migrations/postgres goto $(version)

migrate-down-all:
	migrate -database '$(DB_CONNECTION_STRING)' -path migrations/postgres down

migration-fix:
	migrate -path migrations -database '$(DB_CONNECTION_STRING)' force $(version)