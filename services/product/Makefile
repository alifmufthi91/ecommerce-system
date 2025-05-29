.PHONY: migrate-up migrate-down new_migration

start-db:
	./scripts/start-db.sh

stop-db:
	./scripts/stop-db.sh

init-db-and-migrate-up:
	PGDATABASE=honeybee ./scripts/ensure-db-exists.sh && go run cmd/migrate/main.go up

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

start:
	`go env GOPATH`/bin/swag init --parseDependency --parseInternal &&`go env GOPATH`/bin/air

swaggo:
	@`go env GOPATH`/bin/swag init --parseDependency --parseInternal

test:
	@go test -v -tags dynamic `go list ./... | grep -i 'internal\|pkg'` -cover -coverprofile=coverage.out

new_migration:
	@read -p "Enter migration name: " migration_name; \
	migrate create -ext sql -dir database/migrations -seq $$migration_name

migrate-up1:
	PGDATABASE=honeybee ./scripts/ensure-db-exists.sh && go run cmd/migrate/main.go up 1

migrate-down1:
	go run cmd/migrate/main.go down 1