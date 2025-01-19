.PHONY: build install uninstall clean test help format lint

BINARY_NAME := serverstatus
BINARY_PATH := /opt/$(BINARY_NAME)
BUILD_DIR := .
INSTALL_DIR := /opt/$(BINARY_NAME)
LOG_DIR := /var/log/$(BINARY_NAME)
CONFIG_DIR := /etc/$(BINARY_NAME)

help:
	@echo "Usage: make [target]"
	@echo "Available targets:"
	@echo "  build       Build the project"
	@echo "  install     Install the project"
	@echo "  uninstall   Uninstall the project"
	@echo "  clean       Clean up build files"
	@echo "  test        Run tests"
	@echo "  help        Show this help message"
	@echo "  format      Format the codebase"
	@echo "  lint        Lint the codebase"

build:
	@./scripts/build.sh

install: build
	@sudo INSTALL_DIR=$(INSTALL_DIR) LOG_DIR=$(LOG_DIR) CONFIG_DIR=$(CONFIG_DIR) ./scripts/install.sh

uninstall:
	@sudo ./scripts/uninstall.sh

clean:
	@sudo rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@sudo rm -rf $(LOG_DIR)
	@sudo rm -rf $(CONFIG_DIR)

test:
	@echo "Running tests..."
	# Placeholder
	@echo "Tests completed."

format:
	@echo "Formatting code..."
	@gofmt -w .

lint:
	@echo "Linting code..."
	@golangci-lint run