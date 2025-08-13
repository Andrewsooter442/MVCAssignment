GO := go
GOPATH := $(shell go env GOPATH)
GOPATH_BIN := $(GOPATH)/bin
GOLANGCI_LINT := $(GOPATH_BIN)/golangci-lint
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
GOIMPORTS := $(GOPATH_BIN)/goimports
GO_PACKAGES = $(shell go list ./... | grep -v vendor)
PACKAGE_BASE := github.com/Andrewsooter442/MVCAssignment

DB_HOST = $(shell grep -A6 "^db:" config.yaml | grep "host:" | head -1 | cut -d'"' -f2)
DB_PORT = $(shell grep -A6 "^db:" config.yaml | grep "port:" | head -1 | awk '{print $$2}')
DB_USER = $(shell grep -A6 "^db:" config.yaml | grep "user:" | head -1 | cut -d'"' -f2)
DB_PASS = $(shell grep -A6 "^db:" config.yaml | grep "password:" | head -1 | cut -d'"' -f2)
DB_NAME = $(shell grep -A6 "^db:" config.yaml | grep "db_name:" | head -1 | cut -d'"' -f2)

DB_INIT_FILE= migrations/000001_create_database.up.sql
CREATE_TABLES= migrations/000002_create_users_table.up.sql
CREATE_MENU=migrations/000004_create_users_table.up.sql
MAKE_ADMIN=migrations/000003_create_users_table.up.sql


DOWN_MIGRATION_FILE = migrations/000001_init_schema.down.sql

.PHONY: help vendor build run dev lint format clean

help:
	@echo "MVCAssignment make help"
	@echo ""
	@echo "vendor: Downloads the dependencies in the vendor folder"
	@echo "build: Builds the binary of the server"
	@echo "run: Runs the binary of the server"
	@echo "dev: Combines build and run commands"
	@echo "lint: Lints the code using vet and golangci-lint"
	@echo "format: Formats the code using fmt and golangci-lint"
	@echo "clean: Removes the vendor directory and binary"

vendor:
	@${GO} mod tidy
	@${GO} mod vendor
	@echo "Vendor downloaded successfully"

build:
	@${GO} build -o mvcassignment ./cmd/
	@echo "Binary built successfully"

run:
	@./mvcassignment

dev:
	@$(GOPATH_BIN)/air -c .air.toml

install-golangci-lint:
	@echo "=====> Installing golangci-lint..."
	@curl -sSfL \
	 	https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
	 	sh -s -- -b $(GOPATH_BIN)  v1.64.8

lint: install-golangci-lint
	@$(GO) vet $(GO_PACKAGES)
	@$(GOLANGCI_LINT) run -c golangci.yaml
	@echo "Lint successful"

install-goimports:
	@echo "=====> Installing formatter..."
	@$(GO) install golang.org/x/tools/cmd/goimports@latest

format: install-goimports
	@echo "=====> Formatting code..."
	@$(GOIMPORTS) -l -w -local ${PACKAGE_BASE} $(SRC)
	@echo "Format successful"

## verify: Run format and lint checks
verify: verify-format lint

## verify-format: Verify the format
verify-format: install-goimports
	@echo "=====> Verifying format..."
	$(if $(shell $(GOIMPORTS) -l -local ${PACKAGE_BASE} ${SRC}), @echo ERROR: Format verification failed! && $(GOIMPORTS) -l -local ${PACKAGE_BASE} ${SRC} && exit 1)

clean:
	@rm -f mvcassignment
	@rm -rf vendor/
	@echo "Clean successful"

install-air:
	@echo "Make sure your GOPATH and GOPATH_BIN is set"
	@curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(GOPATH_BIN)
	@echo "Air installed successfully"	

apply-migration:
	@echo "Applying migration..."
	@echo "DB_HOST: $(DB_HOST)"
	@echo "DB_PORT: $(DB_PORT)"
	@echo "DB_USER: $(DB_USER)"
	@echo "DB_NAME: $(DB_NAME)"

	mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USER) -p$(DB_PASS)  < $(DB_INIT_FILE)
	mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USER) -p$(DB_PASS) $(DB_NAME) < $(CREATE_TABLES)
	mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USER) -p$(DB_PASS) $(DB_NAME) < $(CREATE_MENU)

create-admin:
	@echo "Creating Admin and Chef users..."
	@echo "DB_HOST: $(DB_HOST)"
	@echo "DB_PORT: $(DB_PORT)"
	@echo "DB_USER: $(DB_USER)"
	@echo "DB_NAME: $(DB_NAME)"

	mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USER) -p$(DB_PASS) $(DB_NAME) < $(MAKE_ADMIN)
