.PHONY: help
help: ## Prints help (only for targets with comments)
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

SRC_PACKAGES=$(shell go list ./... | grep -v "vendor")
GOLINT:=$(shell command -v golint 2> /dev/null)
RICHGO=$(shell command -v richgo 2> /dev/null)

ifeq ($(RICHGO),)
	GOBIN=go
else
	GOBIN=richgo
endif

all: setup test

ensure-build-dir:
	mkdir -p out

update-deps: ## Update dependencies
	GO111MODULE=on go get -u

build-deps: ## Install dependencies
	GO111MODULE=on go mod tidy -v

fmt:
	$(GOBIN) fmt $(SRC_PACKAGES)

vet:
	$(GOBIN) vet $(SRC_PACKAGES)

setup: build-deps
ifeq ($(GOLINT),)
	$(GOBIN) get -u golang.org/x/lint/golint
endif
ifeq ($(RICHGO),)
	$(GOBIN) get -u github.com/kyoh86/richgo
endif

test: ensure-build-dir ## Run tests
	GO111MODULE=on ENVIRONMENT=test $(GOBIN) test $(SRC_PACKAGES) -p=1 -coverprofile ./out/coverage -short -v

test-cover-html: ## Run tests with coverage
	mkdir -p ./out
	@echo "mode: count" > coverage-all.out
	$(foreach pkg, $(SRC_PACKAGES),\
	ENVIRONMENT=test $(GOBIN) test -coverprofile=coverage.out -covermode=count $(pkg);\
	tail -n +2 coverage.out >> coverage-all.out;)
	$(GOBIN) tool cover -html=coverage-all.out -o out/coverage.html
