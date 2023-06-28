
BLUE         := $(shell printf "\033[34m")
YELLOW       := $(shell printf "\033[33m")
RED          := $(shell printf "\033[31m")
GREEN        := $(shell printf "\033[32m")
CNone        := $(shell printf "\033[0m")

INFO	= echo ${TIME} ${BLUE}[ .. ]${CNone}
WARN	= echo ${TIME} ${YELLOW}[WARN]${CNone}
ERR		= echo ${TIME} ${RED}[FAIL]${CNone}
OK		= echo ${TIME} ${GREEN}[ OK ]${CNone}
FAIL	= (echo ${TIME} ${RED}[FAIL]${CNone} && false)

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

init: build ## Initialise the database in (default) example/.argus.state
	@./bin/argus -c example/.argus-config.yaml load

build: ## Build the binary in bin/argus
	go build -o bin/argus .

.PHONY: lint.check
lint.check: ## Check install of golanci-lint
	@if ! golangci-lint --version 2>&1 >> /dev/null; then \
		echo -e "\033[0;33mgolangci-lint is not installed: run \`\033[0;32mmake lint.install\033[0m\033[0;33m\` or install it from https://golangci-lint.run\033[0m"; \
		exit 1; \
	fi

.PHONY: lint.install
lint.install: ## Install golangci-lint to the go bin dir
	echo Foo
	@if ! golangci-lint --version  2>&1 >> /dev/null; then \
		$(INFO) "Installing golangci-lint"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.52.2; \
	fi

.PHONY: lint
lint: lint.check ## Run golangci-lint
	@if ! golangci-lint run; then \
		$(WARN) "golangci-lint found issues with your code. Please check and fix them before committing." \
		exit 1; \
	fi ; \
	$(OK) No linting issues found
