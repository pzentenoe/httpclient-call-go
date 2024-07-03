.PHONY: install-tools check test

# Go tool paths
GOLINT = $(shell go env GOPATH)/bin/golint
INEFFASSIGN = $(shell go env GOPATH)/bin/ineffassign
MISSPELL = $(shell go env GOPATH)/bin/misspell
GOCYCLO = $(shell go env GOPATH)/bin/gocyclo

install-tools:
	@echo "Installing tools..."
	go install golang.org/x/lint/golint@latest
	go install github.com/gordonklaus/ineffassign@latest
	go install github.com/client9/misspell/cmd/misspell@latest
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

check: install-tools
	@echo "Running checks..."
	go fmt ./...
	go vet ./...
	$(GOLINT) ./...
	$(MISSPELL) -w .
	$(GOCYCLO) -over 10 .
	$(INEFFASSIGN) .

test:
	go test ./...
