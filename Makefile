.DEFAULT_GOAL = test
.PHONY: FORCE

PROJECT_NAME := "typedyaml"
PKG := "github.com/etecs-ru/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: clean dep vet format lint test test-coverage build

clean:
	rm -rf tools/bin

dep: ## Get the dependencies
	@go mod download

vet: ## Run go vet
	@go vet ${PKG_LIST}

format: tools/gofumpt tools/gofumports ## format the package
	@./tools/bin/gofumpt -s -w .
	@./tools/bin/gofumports -w .
	@echo formatted

lint: .golangci.yml tools/golangci-lint ## lint the package
	@./tools/bin/golangci-lint run ./...
	@echo lint passed

test: ## Run unittests
	@go test -short ${PKG_LIST}

test-coverage: ## Run tests with coverage
	@go test -short -coverprofile cover.out -covermode=atomic ${PKG_LIST}
	@cat cover.out >> coverage.txt

build: ## build tool
	@go build -ldflags "-s -w" -o ./build/typedyaml .


# Source: https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


# Non-PHONY targets (real files)

tools/golangci-lint:
tools/golangci-lint: tools/go.mod tools/go.sum
	@cd tools && go build -o ./bin/golangci-lint \
		-ldflags "-s -w" github.com/golangci/golangci-lint/cmd/golangci-lint

tools/gofumpt:
tools/gofumpt: tools/go.mod tools/go.sum
	cd tools && go build -o ./bin/gofumpt -ldflags "-s -w" mvdan.cc/gofumpt

tools/gofumports:
tools/gofumports: tools/go.mod tools/go.sum
	cd tools && go build -o ./bin/gofumports -ldflags "-s -w" mvdan.cc/gofumpt/gofumports

go.mod: FORCE
	go mod tidy
	go mod verify
go.sum: go.mod