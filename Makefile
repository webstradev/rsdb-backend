BINARY_NAME=api

unit-test:
	@echo "Installing tparse (if needed)..."
	@go install github.com/mfridman/tparse@v0.11.1
	@echo "Running unit tests..."
	@go test ./... -cover -json | tparse

build:
	@echo "Building binary..."
	@go build -o bin/$(BINARY_NAME)
	@echo "Done!"

run: build
	@echo "Running binary..."
	@./bin/$(BINARY_NAME)