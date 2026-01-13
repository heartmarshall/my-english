.PHONY: generate
generate:
	@echo "Generating GraphQL code..."
	@if command -v gqlgen > /dev/null; then \
		gqlgen generate; \
	else \
		echo "gqlgen not found, installing..."; \
		go install github.com/99designs/gqlgen@latest; \
		gqlgen generate; \
	fi

.PHONY: generate-verbose
generate-verbose:
	@echo "Generating GraphQL code (verbose)..."
	@if command -v gqlgen > /dev/null; then \
		gqlgen generate -v; \
	else \
		echo "gqlgen not found, installing..."; \
		go install github.com/99designs/gqlgen@latest; \
		gqlgen generate -v; \
	fi

.PHONY: migrate-up
migrate-up:
	@if [ -z "$(DB_URL)" ]; then \
		echo "Error: DB_URL environment variable is required"; \
		echo "Example: make migrate-up DB_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable"; \
		exit 1; \
	fi
	@echo "Running migrations..."
	@go run github.com/pressly/goose/v3/cmd/goose -dir migrations postgres "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:
	@if [ -z "$(DB_URL)" ]; then \
		echo "Error: DB_URL environment variable is required"; \
		echo "Example: make migrate-down DB_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable"; \
		exit 1; \
	fi
	@echo "Rolling back migrations..."
	@go run github.com/pressly/goose/v3/cmd/goose -dir migrations postgres "$(DB_URL)" down

.PHONY: migrate-status
migrate-status:
	@if [ -z "$(DB_URL)" ]; then \
		echo "Error: DB_URL environment variable is required"; \
		echo "Example: make migrate-status DB_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable"; \
		exit 1; \
	fi
	@echo "Migration status:"
	@go run github.com/pressly/goose/v3/cmd/goose -dir migrations postgres "$(DB_URL)" status

.PHONY: migrate-create
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME environment variable is required"; \
		echo "Example: make migrate-create NAME=add_new_column"; \
		exit 1; \
	fi
	@echo "Creating migration: $(NAME)"
	@go run github.com/pressly/goose/v3/cmd/goose -dir migrations create $(NAME) sql
	@echo "Migration created! Don't forget to add -- +goose Up and -- +goose Down directives."

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make generate         - Generate GraphQL code"
	@echo "  make generate-verbose - Generate GraphQL code with verbose output"
	@echo ""
	@echo "Migration commands:"
	@echo "  make migrate-up       - Apply all migrations (requires DB_URL)"
	@echo "  make migrate-down     - Rollback last migration (requires DB_URL)"
	@echo "  make migrate-status   - Show migration status (requires DB_URL)"
	@echo "  make migrate-create   - Create new migration (requires NAME)"
	@echo ""
	@echo "Examples:"
	@echo "  make migrate-up DB_URL=postgres://user:pass@localhost:5432/mydb?sslmode=disable"
	@echo "  make migrate-create NAME=add_user_table"

