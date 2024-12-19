all: help

## Local development:

dev-cover: cover ## Run go test on all modules with coverage and open the report in the browser
	go tool cover -html=coverage.out

bench: ## Run go test on all modules with benchmarks
	cd benchmark && go test -bench=. ./...

## Continuous integration:

TPARSE := go run github.com/mfridman/tparse@latest

test: ## Run the binary with dev config
	go test ./... -json | ${TPARSE} -all -progress

cover: ## Run go test on all modules with coverage
	go test -coverpkg=./... -coverprofile coverage.out ./... -json | ${TPARSE} -all -progress
	go tool cover -html=coverage.out -o coverage.html

release: test ## Create a new release
	@git diff-index --quiet HEAD -- || (echo "\n\033[31mWorking directory is not clean, please commit your changes before creating a new release.\033[0m" && exit 1)
	@echo "Creating a new release..."
	@echo "Please enter the new version number: "
	@read version; \
	git tag -a $$version -m "Release $$version"; \
	git push origin

## Help:

define print_help
	# Self generating help
	# Inspired from https://gist.github.com/thomaspoignant/5b72d579bd5f311904d973652180c705#file-makefile-L89
	echo 'Usage:'
	echo '  make [target]...'
	echo ''
	echo 'Targets:'
	awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "        %-20s%s\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "\n    %s\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
	echo ''
endef

help: ## Show this help.
	@$(call print_help)

%: # Fallback rule to print help when typo
	@echo "make: *** No rule to make target '$@'.  Stop."
	@$(call print_help)