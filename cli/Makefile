help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

init: build ## Initialise the database in (default) example/.argus.state
	@./bin/argus -c example/.argus-config.yaml load

build: ## Build the binary in bin/argus
	go build -o bin/argus .
