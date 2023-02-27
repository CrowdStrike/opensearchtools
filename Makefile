.PHONY: lint
lint: 
	@echo "Linting..."
	@golangci-lint run

.PHONY: format
format: 
	@echo "Formatting..."
	@gofmt -w -l -e .

.PHONY: test
test: 
	@echo "Testing..."
	@go test ./...

.PHONY: help
help:
	@echo "make targets:"
	@echo "lint"
	@echo "format"
	@echo "test"
