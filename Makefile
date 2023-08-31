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
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o magecomm-linux-amd64 \
    && GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o magecomm-linux-arm64 \
	&& GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o magecomm-windows-amd64.exe \
	&& GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o magecomm-darwin-amd64 \
	&& GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o magecomm-darwin-arm64

install: ## Install go binary for your platform
	go install
	@echo 'If $$GOPATH/bin is in your PATH, you can run magecomm from anywhere'

test: ## Run tests
	go test -v ./...

help: ## Additional Details of what this project is
	@$(cat) README.md
