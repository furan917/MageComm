.PHONY: helper sh
.DEFAULT_GOAL := helper
cat := $(if $(filter $(OS),Windows_NT),type,cat)

helper: ## Describe all commands available in Makefile
	@grep -h -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST)| awk '{FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'| sort

build: ## Build go binary for your platform
	@echo "Building binary for your platform"
	go build

build-all: ## Build go binary for all supported platforms
	@echo "Building for all supported platforms"
	GOOS=linux GOARCH=amd64 go build -o magecomm-linux-amd64
	GOOS=windows GOARCH=amd64 go build -o magecomm-windows-amd64.exe
	GOOS=darwin GOARCH=amd64 go build -o magecomm-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o magecomm-darwin-arm64

help: ## Additional Details of what this project is
	@$(cat) README.md