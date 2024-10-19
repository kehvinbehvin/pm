# Variables
APP_NAME = pm 
GO_FILES = $(find . -type f -name '*.go')
BIN_PATH = /usr/local/bin/$(APP_NAME)

# Default target, run if no target is specified
all: build

build:
	@echo "Building the binary..."
	go build -o $(APP_NAME)

register-path:
	sudo mv ./$(APP_NAME) $(BIN_PATH)

# Format code using gofmt
format:
	@echo "Formatting Go files..."
	gofmt -w -l ./

# Clean the compiled binary
clean:
	@echo "Cleaning up..."
	rm -f $(APP_NAME)
	rm -f $(BIN_PATH)

# Phony targets are not associated with actual files and are always executed
.PHONY: all build format clean
