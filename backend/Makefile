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

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make generate         - Generate GraphQL code"
	@echo "  make generate-verbose - Generate GraphQL code with verbose output"

