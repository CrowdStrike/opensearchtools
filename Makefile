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
	@echo "format"
	@echo "test"
