include .envrc
MIGRATIONS_PATH=./cmd/migrate/migrations/

.PHONY: migrate-create
migration:
	@/home/kenzo/.goose/bin/goose -dir ${MIGRATIONS_PATH} -s create ${filter-out $@,${MAKECMDGOALS}} sql

.PHONY: migrate-up
migrate-up:
		@/home/kenzo/.goose/bin/goose  -dir ${MIGRATIONS_PATH}  postgres ${MIGR_ADDR} up

.PHONY: migrate-down
migrate-down:
		@/home/kenzo/.goose/bin/goose  -dir ${MIGRATIONS_PATH} postgres ${MIGR_ADDR} down

.PHONY: seed
seed:
	@go run cmd/migrate/migrations/main.go

.PHONY: gen-docs
gen-docs:
	@swag init -g main.go -d ./cmd/api,db,store && swag fmt


.PHONY: test
test:
		@go test -v ./...
