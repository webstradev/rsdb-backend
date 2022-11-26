unit-test:
	@echo "Installing tparse (if needed)..."
	@go install github.com/mfridman/tparse@v0.11.1
	@echo "Running unit tests..."
	@go test ./... -cover -json -count=3 | tparse