# Load env variables from .env file
ifneq (,$(wildcard .env))
    include .env
    export
endif

# Construct DB URL from .env variables
DB_URL=postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: migrate-up migrate-down migrate-create run db-seed

# Run database migrations up
migrate-up:
	migrate -path ./src/database/migrations -database "$(DB_URL)" up

# Run database migrations down
migrate-down:
	migrate -path ./src/database/migrations -database "$(DB_URL)" down

# Create a new migration file
# Usage: make migrate-create name=create_users_table
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: Please specify the migration name, e.g., 'make migrate-create name=my_migration'"; \
		exit 1; \
	fi
	migrate create -ext sql -dir ./src/database/migrations -seq $(name)

# Seed the database
db-seed:
	go run src/main.go --seed

# Run the app with hot reloading using air
run:
	air
