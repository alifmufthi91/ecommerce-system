new_migration:
	@read -p "Enter migration name: " migration_name; \
	migrate create -ext sql -dir migrations -seq $$migration_name

migrate-up:
	migrate -source file://migrations -database postgresql://postgres:postgres@localhost:5432/ecommerce?sslmode=disable up

migrate-down:
	migrate -source file://migrations -database postgresql://postgres:postgres@localhost:5432/ecommerce?sslmode=disable down