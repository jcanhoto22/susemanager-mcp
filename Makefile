.PHONY: build run test lint clean docker-build docker-run dev

APP_NAME = mcp-susemanager
BIN_DIR = bin
CMD_DIR = cmd/server
MAIN_FILE = $(CMD_DIR)/main.go
CONFIG_FILE = config.yaml

GO ?= go
GOFLAGS ?= -ldflags="-s -w"
GOPATH ?= $(shell $(GO) env GOPATH)

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(APP_NAME) $(MAIN_FILE)

run: build
	@echo "Running $(APP_NAME)..."
	@./$(BIN_DIR)/$(APP_NAME)

test:
	@echo "Running tests..."
	$(GO) test ./... -v -count=1

test-short:
	@echo "Running short tests..."
	$(GO) test ./... -short -count=1

test-race:
	@echo "Running tests with race detector..."
	$(GO) test ./... -race -count=1

lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null 2>&1 || (echo "golangci-lint not installed"; exit 0)
	golangci-lint run ./...

vet:
	$(GO) vet ./...

clean:
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR)
	@rm -rf dist/

tidy:
	$(GO) mod tidy
	$(GO) mod verify

docker-build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):latest .

docker-run:
	@echo "Running Docker container..."
	docker run --rm -i -v $(PWD)/$(CONFIG_FILE):/etc/suse-mcp/$(CONFIG_FILE) $(APP_NAME):latest

dev:
	@echo "Running in dev mode..."
	$(GO) run $(MAIN_FILE)

help:
	@echo "Available targets:"
	@echo "  build       - Build the binary"
	@echo "  run         - Build and run"
	@echo "  test        - Run all tests"
	@echo "  test-short  - Run short tests only"
	@echo "  test-race   - Run tests with race detector"
	@echo "  lint        - Run golangci-lint"
	@echo "  vet         - Run go vet"
	@echo "  clean       - Clean build artifacts"
	@echo "  tidy        - Tidy go modules"
	@echo "  docker-build- Build Docker image"
	@echo "  docker-run  - Run Docker container"
	@echo "  dev         - Run in dev mode (go run)"
