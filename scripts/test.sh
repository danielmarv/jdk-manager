#!/bin/bash

# Test script for JDK Manager
set -e

echo "Running JDK Manager tests..."

# Run unit tests
echo "Running unit tests..."
go test -v ./...

# Run integration tests
echo "Running integration tests..."
go test -v ./test/integration/...

# Run tests with coverage
echo "Running tests with coverage..."
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

echo "Tests completed successfully!"
echo "Coverage report generated: coverage.html"
