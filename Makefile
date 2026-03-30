# ── Portman Makefile ──────────────────────────────────────────────

APP_NAME    := portman
MODULE      := github.com/NoaTamburrini/portman
VERSION     := $(shell git describe --tags --always --dirty 2>/dev/null | sed 's/^v//' || echo "dev")
COMMIT      := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE  := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS     := -s -w -X $(MODULE)/internal/version.Version=$(VERSION)

# ── Build ────────────────────────────────────────────────────────

.PHONY: build
build: ## Build the binary
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME) .

.PHONY: build-all
build-all: ## Build for all platforms
	GOOS=darwin  GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME)-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME)-darwin-arm64 .
	GOOS=linux   GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME)-linux-amd64 .
	GOOS=linux   GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME)-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME)-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME)-windows-arm64.exe .

.PHONY: install
install: build ## Build and install to $GOPATH/bin
	cp bin/$(APP_NAME) $(shell go env GOPATH)/bin/$(APP_NAME)

# ── Dev ──────────────────────────────────────────────────────────

.PHONY: run
run: ## Run without building a binary
	go run -ldflags "$(LDFLAGS)" .

.PHONY: fmt
fmt: ## Format code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: test
test: ## Run tests
	go test ./...

.PHONY: lint
lint: fmt vet ## Format + vet

# ── Deps ─────────────────────────────────────────────────────────

.PHONY: deps
deps: ## Tidy and verify dependencies
	go mod tidy
	go mod verify

.PHONY: deps-update
deps-update: ## Update all dependencies to latest
	go get -u ./...
	go mod tidy

# ── Release ──────────────────────────────────────────────────────

.PHONY: release-dry
release-dry: ## Test GoReleaser locally (no publish)
	goreleaser release --snapshot --clean

.PHONY: tag
tag: ## Create a new version tag (usage: make tag v=1.2.3)
	@if [ -z "$(v)" ]; then echo "Usage: make tag v=1.2.3"; exit 1; fi
	git tag -a v$(v) -m "Release v$(v)"
	@echo "Tagged v$(v) — push with: make tag-push v=$(v)"

.PHONY: tag-push
tag-push: ## Push a version tag to trigger GitHub release (usage: make tag-push v=1.2.3)
	@if [ -z "$(v)" ]; then echo "Usage: make tag-push v=1.2.3"; exit 1; fi
	git push origin v$(v)
	@echo "Pushed v$(v) — GoReleaser will build the release on GitHub"

# ── Clean ────────────────────────────────────────────────────────

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf bin/ dist/

# ── Info ─────────────────────────────────────────────────────────

.PHONY: version
version: ## Show current version info
	@echo "Version:  $(VERSION)"
	@echo "Commit:   $(COMMIT)"
	@echo "Date:     $(BUILD_DATE)"

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
