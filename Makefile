include .env
export

GO := go
GOPATH := $(shell go env GOPATH)
GOPATH_BIN := $(GOPATH)/bin
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
GO_PACKAGES = $(shell go list ./... | grep -v vendor)
PACKAGE_BASE := $(shell head -n 1 go.mod | awk '{print $$2}')


MIGRATE := $(GOPATH_BIN)/migrate
GOLANGCI_LINT := $(GOPATH_BIN)/golangci-lint
GOIMPORTS := $(GOPATH_BIN)/goimports
AIR := $(GOPATH_BIN)/air


DB_URL = "mysql://$(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?multiStatements=true"
MIGRATIONS_PATH = "migrations"


.PHONY: help vendor build run dev lint format clean verify verify-format \
        install-golangci-lint install-goimports install-air install-migrate-cli \
        migrate-create migrate-up migrate-down migrate-down-all migrate-to migrate-status


help:
	@echo "Application Makefile"
	@echo ""
	@echo "---- Application ----"
	@echo "vendor: Downloads Go module dependencies."
	@echo "build:  Builds the application binary."
	@echo "run:    Runs the built binary."
	@echo "dev:    Runs the app with live-reloading using Air."
	@echo "clean:  Removes the binary and vendor directory."
	@echo ""
	@echo "---- Quality & Formatting ----"
	@echo "lint:   Lints the code for potential issues."
	@echo "format: Formats the Go source code."
	@echo "verify: Runs both format verification and linting."
	@echo ""
	@echo "---- Database Migrations ----"
	@echo "migrate-create: Creates new up/down migration files."
	@echo "migrate-up:     Applies all pending migrations."
	@echo "migrate-down:   Rolls back the last applied migration."
	@echo "migrate-down-all: Rolls back all migrations."
	@echo "migrate-to version=<version_num>: Migrates to a specific version."
	@echo "migrate-status: Shows the current migration status."


vendor:
	@${GO} mod tidy
	@${GO} mod vendor
	@echo "Vendor directory created successfully."

build:
	@${GO} build -o mvcassignment ./cmd/
	@echo "Binary built successfully."

run:
	@./mvcassignment

dev: install-air
	@$(AIR) -c .air.toml

clean:
	@rm -f mvcassignment
	@rm -rf vendor/
	@echo "Clean successful."


lint: install-golangci-lint
	@$(GO) vet $(GO_PACKAGES)
	@$(GOLANGCI_LINT) run -c golangci.yaml
	@echo "Lint successful."

format: install-goimports
	@echo "=====> Formatting code..."
	@$(GOIMPORTS) -l -w -local ${PACKAGE_BASE} $(SRC)
	@echo "Format successful."

verify: verify-format lint

verify-format: install-goimports
	@echo "=====> Verifying format..."
	$(if $(shell $(GOIMPORTS) -l -local ${PACKAGE_BASE} ${SRC}), @echo "ERROR: Code is not formatted. Please run 'make format'." && exit 1)
	@echo "Format verification successful."


install-golangci-lint:
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || \
		(echo "=====> Installing golangci-lint..." && \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH_BIN) v1.55.2)

install-goimports:
	@command -v $(GOIMPORTS) >/dev/null 2>&1 || \
		(echo "=====> Installing goimports..." && \
		$(GO) install golang.org/x/tools/cmd/goimports@latest)

install-air:
	@command -v $(AIR) >/dev/null 2>&1 || \
		(echo "=====> Installing Air for live-reloading..." && \
		curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(GOPATH_BIN))

install-migrate-cli:
	@command -v $(MIGRATE) >/dev/null 2>&1 || \
		(echo "=====> Installing golang-migrate/migrate CLI..." && \
		go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest)


migrate-create:
	@read -p "Enter migration name (e.g., add_price_to_items): " name; \
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_PATH) -seq $$name

migrate-up: install-migrate-cli
	@echo "Applying all up migrations..."
	@$(MIGRATE) -path $(MIGRATIONS_PATH) -database $(DB_URL) -verbose up

migrate-down: install-migrate-cli
	@echo "Rolling back last migration..."
	@$(MIGRATE) -path $(MIGRATIONS_PATH) -database $(DB_URL) -verbose down 1

migrate-down-all: install-migrate-cli
	@echo "Rolling back all migrations..."
	@$(MIGRATE) -path $(MIGRATIONS_PATH) -database $(DB_URL) -verbose down -all

migrate-to: install-migrate-cli
	@echo "Migrating to version $(version)..."
	@$(MIGRATE) -path $(MIGRATIONS_PATH) -database $(DB_URL) -verbose goto $(version)

migrate-status: install-migrate-cli
	@echo "Checking migration status..."
	@$(MIGRATE) -path $(MIGRATIONS_PATH) -database $(DB_URL) -verbose version
