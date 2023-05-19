# Go parameters
CONFIG = config.yaml
SRV_CONFIG = config.services.yaml
SRV_CONFIG_EX = config.services.example.yaml

# Build directory
BUILD_DIR = ./build

# Binary name
BINARY_NAME = remoteLog

# Default target
.DEFAULT_GOAL := build

# Build the program
build: clean
	@echo "Creating folder $(BUILD_DIR)..."
	@mkdir -p $(BUILD_DIR)
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/remoteLog
	@echo "Copying config files..."
	cp $(CONFIG) $(BUILD_DIR)/$(CONFIG)
	cp $(SRV_CONFIG_EX) $(BUILD_DIR)/$(SRV_CONFIG)

# Clean the build
clean:
	@echo "Cleaning dependencies..."
	@go clean
	@echo "Removing $(BUILD_DIR) ..."
	@rm -rf $(BUILD_DIR)

.PHONY: build clean
