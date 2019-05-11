
SERVICE ?= glee


.PHONY: run
run: ## Runs the service in development and watches for changes
ifeq ($(shell which modd),)
	brew install modd
endif
	modd

.PHONY: run-daemon
run-daemon:
	systemctl restart glee.service

.PHONY: test
test: ## lint ## Runs the tests and higlights race conditions
	GO111MODULE=on go test -race $$(go list ./...)

.PHONY: coverage
COVERALLS_TOKEN ?= c45yMbJJTs5jkjsc3MSMGf46JEZNVwsoO
coverage: ## Runs the tests and reports coverage to coveralls
ifeq ($(shell which goveralls),)
	GO111MODULE=off go get -u github.com/mattn/goveralls
endif
	GO111MODULE=on goveralls -v -service=Buildkite -ignore=cmd/$(SERVICE)/main.go -repotoken $(COVERALLS_TOKEN)

DIST_PATH ?= .

.PHONY: build-linux
build-linux: ## Builds the executable for linux amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -a -installsuffix cgo \
		-o $(DIST_PATH)/$(SERVICE) ./cmd/$(SERVICE)

.PHONY: build
build: build-linux ## Builds the executable for all arch (only linux for now)

.PHONY: lint
lint: ## Runs more than 20 different linters using golangci-lint to ensure consistency in code.
ifeq ($(shell which golangci-lint),)
	brew install golangci/tap/golangci-lint
endif
	golangci-lint run

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
