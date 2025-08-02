# Build variables
BINARY_NAME=jdk
VERSION?=1.0.0
BUILD_DIR=dist
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

build: deps ## Build the binary
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) main.go

build-all: ## Build for all platforms
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/linux-amd64/$(BINARY_NAME) main.go
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/darwin-amd64/$(BINARY_NAME) main.go
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/darwin-arm64/$(BINARY_NAME) main.go
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/windows-amd64/$(BINARY_NAME).exe main.go

test: ## Run tests
	$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

deps: ## Download dependencies
	mkdir -p $(BUILD_DIR)
	$(GOMOD) tidy
	$(GOMOD) download

install: build ## Install the binary
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

uninstall: ## Uninstall the binary
	rm -f /usr/local/bin/$(BINARY_NAME)

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
