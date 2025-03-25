.PHONY: build install clean

# Binary name
BINARY=slark
# Go module path
MODULE=$(shell grep module go.mod | awk '{print $$2}')
# Main entry point
MAIN_FILE=cmd/slark/main.go
# Version from git tag or default
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
# Build directory
BUILD_DIR=build

build:
	@echo "Building $(BINARY)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-X $(MODULE)/internal/version.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY) $(MAIN_FILE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY)"

install: build
	@echo "Installing $(BINARY)..."
	@cp $(BUILD_DIR)/$(BINARY) ${GOPATH}/bin/$(BINARY)
	@echo "$(BINARY) installed successfully to ${GOPATH}/bin/$(BINARY)"

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# For users using go install directly
go-install:
	@echo "Running go install..."
	go install -ldflags "-X $(MODULE)/internal/version.Version=$(VERSION)" ./cmd/slark

# Help output
help:
	@echo "Available commands:"
	@echo "  make build     - Build the binary to ./build directory"
	@echo "  make install   - Install the binary to GOPATH/bin"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make go-install - Use go install to install directly"
	@echo "  make help      - Show this help" 