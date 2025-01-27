.PHONY: build install uninstall update clean test help format lint

BINARY_NAME := serverstatus
BUILD_DIR := .
INSTALL_DIR := /opt/$(BINARY_NAME)
LOG_DIR := /var/log/$(BINARY_NAME)
CONFIG_DIR := /etc/$(BINARY_NAME)
SERVICE_DIR := /etc/systemd/system

help:
	@echo "Usage: make [target]"
	@echo "Available targets:"
	@echo "  build       Build the project"
	@echo "  install     Install the project"
	@echo "  uninstall   Uninstall the project"
	@echo "  update      Update dependencies"
	@echo "  clean       Clean up build files"
	@echo "  test        Run tests"
	@echo "  help        Show this help message"
	@echo "  format      Format the codebase"
	@echo "  lint        Lint the codebase"

build:
	@echo "Tidying up go modules..."
	@go mod tidy
	@echo "Building..."
	@./scripts/build.sh

install:
	@INSTALL_DIR=$(INSTALL_DIR) LOG_DIR=$(LOG_DIR) CONFIG_DIR=$(CONFIG_DIR) SERVICE_DIR=$(SERVICE_DIR) ./scripts/install.sh

uninstall:
	@./scripts/uninstall.sh

update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "Dependencies updated."

clean:
	@echo "Cleaning up..."
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@rm -rf $(LOG_DIR)
	@rm -rf $(CONFIG_DIR)

test:
	@echo "Running tests..."
	@go test ./...

format:
	@echo "Formatting code..."
	@gofmt -w .

lint:
	@echo "Linting code..."
	@golangci-lint run
