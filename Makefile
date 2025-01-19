.PHONY: build install uninstall clean test

BINARY_NAME := serverstatus
BINARY_PATH := /opt/$(BINARY_NAME)
BUILD_DIR := .
INSTALL_DIR := /opt/$(BINARY_NAME)
LOG_DIR := /var/log/$(BINARY_NAME)
CONFIG_DIR := /etc/$(BINARY_NAME)

build:
	@./scripts/build.sh

install: build
	@sudo ./scripts/install.sh

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