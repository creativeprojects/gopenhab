GOCMD=env go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
GOGET=$(GOCMD) get

GOBIN=$(shell env go env GOBIN)
MOCKERY=$(GOBIN)/mockery
MOCKS="openhab/mock_subscriber.go"

TESTS=./...
COVERAGE_FILE=coverage.out

.PHONY: all test build coverage clean mocks

all: test build

build:
		@echo "[*] $@"
		$(GOBUILD) -v

test: mocks
		@echo "[*] $@"
		$(GOTEST) -v $(TESTS)

coverage: mocks
		@echo "[*] $@"
		$(GOTEST) -coverprofile=$(COVERAGE_FILE) $(TESTS)
		$(GOTOOL) cover -html=$(COVERAGE_FILE)

clean:
		@echo "[*] $@"
		$(GOCLEAN)
		rm -f $(BINARY) $(COVERAGE_FILE)
		@find . -type f -name "mock_*.go" -exec rm -v {} \;

mocks: $(MOCKS)

$(MOCKERY):
		@echo "[*] $@"
		@go install -v github.com/vektra/mockery/v2@latest

$(MOCKS): $(MOCKERY)
		@echo "[*] $@"
		@go generate ./...
