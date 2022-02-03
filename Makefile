PACKAGES_PATH = $(shell go list -f '{{ .Dir }}' ./...)

.PHONY: all
all: check_tools ensure-deps fmt imports vet test

.PHONY: check_tools
check_tools:
	@type "goimports" > /dev/null 2>&1 || (echo 'Please install goimports: go install golang.org/x/tools/cmd/goimports@latest'; exit 1)

.PHONY: ensure-deps
ensure-deps:
	@echo "=> Syncing dependencies with go mod tidy"
	@go mod tidy

.PHONY: fmt
fmt:
	@echo "=> Executing go fmt"
	@go fmt ./...

.PHONY: imports
imports:
	@echo "=> Executing goimports"
	@goimports -w $(PACKAGES_PATH)

.PHONY: vet
vet:
	@echo "=> Executing go vet"
	@go vet ./...

.PHONY: test
test:
	@echo "=> Running tests"
	@go test ./... -covermode=atomic -coverpkg=./... -count=1 -race

.PHONY: test-cover
test-cover:
	@echo "=> Running tests and generating report"
	@go test ./... -covermode=atomic -coverprofile=/tmp/coverage.out -coverpkg=./... -count=1 -race
	@go tool cover -html=/tmp/coverage.out
