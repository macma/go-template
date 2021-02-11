# Makefile
ENV ?= dev
SHELL:=/bin/bash
-include $(ENV).env
export MIGRATION_PATH:=./db/migrations/
	
# create database on dev machine
.PHONY: db-create
db-create:
	@env "PGPASSWORD=${POSTGRES_PASSWORD}" psql -h "${POSTGRES_ADDRESS}" -p "${POSTGRES_PORT}" \
			-U "${POSTGRES_USERNAME}" -w -c "SELECT 1 FROM pg_database WHERE datname = '${POSTGRES_DATABASE}';" \
			| grep 1 \
			|| env "PGPASSWORD=${POSTGRES_PASSWORD}" psql -h "${POSTGRES_ADDRESS}" -p "${POSTGRES_PORT}" \
				-U "${POSTGRES_USERNAME}" -w -c "create database \"${POSTGRES_DATABASE}\";"

.PHONY: db-drop
db-drop:
	@env "PGPASSWORD=${POSTGRES_PASSWORD}" psql -h "${POSTGRES_ADDRESS}" -p "${POSTGRES_PORT}" \
				-U "${POSTGRES_USERNAME}" -w -c "drop database \"${POSTGRES_DATABASE}\";"


.PHONY: build
build:
	@CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v -o cd-assets/fiat-integration cmd/fiat-integration-app-server/main.go # incremental build, should be faster

### SQL
# make migrate-create name=add-this-column
migrate-create:
	@mkdir -p ./state/$(MIGRATION_PATH)
	@migrate create -dir ./state/$(MIGRATION_PATH) -ext sql $(name)
migrate:
	@migrate -path ./state/$(MIGRATION_PATH) \
		-database "postgres://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)@$(POSTGRES_ADDRESS)/$(POSTGRES_DATABASE)?sslmode=disable" up
# make migrate version=version
migrate-revert:
	@migrate -path ./state/$(MIGRATION_PATH) \
		-database "postgres://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)@$(POSTGRES_ADDRESS)/$(POSTGRES_DATABASE)?sslmode=disable" force $(version)
	@migrate -path ./state/$(MIGRATION_PATH) \
		-database "postgres://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)@$(POSTGRES_ADDRESS)/$(POSTGRES_DATABASE)?sslmode=disable" down 1
migrate-rollback:
	@migrate -path ./state/$(MIGRATION_PATH) \
		-database "postgres://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)@$(POSTGRES_ADDRESS)/$(POSTGRES_DATABASE)?sslmode=disable" down 1
migrate-drop:
	@migrate -path ./state/$(MIGRATION_PATH) \
		-database "postgres://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)@$(POSTGRES_ADDRESS)/$(POSTGRES_DATABASE)?sslmode=disable" down

.PHONY: test-coverage
test-coverage:
	cc-test-reporter before-build; \
	./cc.sh; \
	retval=$$?; \
	if [ $$retval -ne 0 ]; then exit 1; fi; \
	cc-test-reporter upload-coverage;