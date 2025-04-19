#!/bin/bash

# Ensure Go is installed
if ! command -v go &>/dev/null; then
    echo "Go is not installed. Please install Go to build this project."
    exit 1
fi

# Update dependencies
go mod tidy

rm -rf ./bin
# Build the application
go build -o ./bin/containix ./cmd/containix

echo "Build completed successfully. Run ./containix to start the application."
