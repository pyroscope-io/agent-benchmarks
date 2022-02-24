.PHONY: all
all: build-docker run ## Build with docker and run the benchmark suite

.PHONY: build
build: runner/*.go runner/go.* ## Build the runner
	cd runner ; go build -v

.PHONY: build-docker
build-docker: ## Build the runner, using a dockerized golang compiler
	docker run --rm -v ${PWD}:/opt/agent-benchmarks/ \
		-w /opt/agent-benchmarks/runner \
		golang:latest go build -v

.PHONY: run
run: ## Run the benchmark suite
	./runner/runner $(shell ls -d fibonacci-*)

.PHONY: test
test: build ## Test the runner with a dummy benchmark
	./runner/runner dummy

help: ## Show this help
	@egrep '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | sed 's/Makefile://' | awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-z0-9A-Z_-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 }'
