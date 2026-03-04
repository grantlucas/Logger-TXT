BINARY_NAME := logger-txt
BUILD_DIR := build
VERSION ?= dev
DATE := $(shell date +%d/%m/%Y)
LDFLAGS := -ldflags "-X github.com/grantlucas/Logger-TXT/internal/cmd.version=$(VERSION) -X github.com/grantlucas/Logger-TXT/internal/cmd.date=$(DATE)"
COVERAGE_FILE := coverage.out
COVERAGE_THRESHOLD := 100

.PHONY: build test lint clean install coverage vet

build:
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/logger-txt/

test:
	go test ./...

vet:
	go vet ./...

lint: vet
	golangci-lint run

coverage:
	go test -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./internal/...
	@COVERAGE=$$(go tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print $$3}' | tr -d '%'); \
	echo "Coverage: $${COVERAGE}%"; \
	if [ "$$(echo "$${COVERAGE} < $(COVERAGE_THRESHOLD)" | bc)" -eq 1 ]; then \
		echo "FAIL: Coverage $${COVERAGE}% is below $(COVERAGE_THRESHOLD)%"; \
		exit 1; \
	fi

clean:
	rm -rf $(BUILD_DIR) $(COVERAGE_FILE)

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME) 2>/dev/null || \
	cp $(BUILD_DIR)/$(BINARY_NAME) $(HOME)/go/bin/$(BINARY_NAME)
