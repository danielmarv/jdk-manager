# Build variables
BINARY_NAME=jdk
VERSION?=1.0.0
BUILD_DIR=dist

# Determine binary name with extension based on OS
ifeq ($(OS),Windows_NT)
    BINARY_NAME_EXT=$(BINARY_NAME).exe
    # Command for removing directories on Windows
    RM_CMD=cmd /c rmdir /s /q
    # Command for creating directories on Windows (only if not exists)
    MKDIR_CMD=cmd /c if not exist
else
    BINARY_NAME_EXT=$(BINARY_NAME)
    # Command for removing directories on Unix-like systems
    RM_CMD=rm -rf
    # Command for creating directories on Unix-like systems
    MKDIR_CMD=mkdir -p
endif

LDFLAGS=-ldflags "-X github.com/jdk-manager/cmd.version=${VERSION}"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build targets
.PHONY: all build build-all clean test deps help install uninstall

all: clean deps test build

build: deps ## Build the binary for the current OS
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME_EXT) main.go

build-all: ## Build for all platforms
	$(MKDIR_CMD) $(BUILD_DIR)\linux-amd64 # Use Windows path separator for mkdir on Windows
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/linux-amd64/$(BINARY_NAME) main.go
	$(MKDIR_CMD) $(BUILD_DIR)\darwin-amd64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/darwin-amd64/$(BINARY_NAME) main.go
	$(MKDIR_CMD) $(BUILD_DIR)\darwin-arm64
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/darwin-arm64/$(BINARY_NAME) main.go
	$(MKDIR_CMD) $(BUILD_DIR)\windows-amd64
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/windows-amd64/$(BINARY_NAME).exe main.go # Explicitly add .exe for Windows

test: ## Run tests
	$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

clean: ## Clean build artifacts
	$(GOCLEAN)
	-$(RM_CMD) $(BUILD_DIR) # Use the OS-specific remove command. The '-' suppresses errors if dir doesn't exist.

deps: ## Download dependencies
	$(MKDIR_CMD) $(BUILD_DIR) md $(BUILD_DIR) # Ensure dist directory exists before downloading
	$(GOMOD) download
	$(GOMOD) tidy

install: build ## Install the binary
	# This install target is primarily for Unix-like systems.
	# For Windows, manual placement or a dedicated installer is usually preferred.
	cp $(BUILD_DIR)/$(BINARY_NAME_EXT) /usr/local/bin/

uninstall: ## Uninstall the binary
	# This uninstall target is primarily for Unix-like systems.
	rm -f /usr/local/bin/$(BINARY_NAME_EXT)

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
