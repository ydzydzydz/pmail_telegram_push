.PHONY: help build build-frontend build-backend test test-verbose clean tidy install-deps install-backend install-frontend run dev lint lint-frontend lint-backend format release release-all

BINARY_NAME := pmail_telegram_push
BACKEND_DIR := backend
FRONTEND_DIR := frontend
HOOK_DIST := $(BACKEND_DIR)/hook/dist
LDFLAGS := -ldflags="-s -w"
CGO_ENABLED := 0

PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	windows/arm64

help:
	@echo "==================== pmail_telegram_push ===================="
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  build              Build full project (frontend + backend)"
	@echo "  build-frontend     Build frontend only (output to $(HOOK_DIST))"
	@echo "  build-backend      Build backend binary only (requires built frontend)"
	@echo "  run                Build and run backend"
	@echo "  dev                Start frontend dev server"
	@echo "  test               Run backend tests"
	@echo "  test-verbose       Run backend tests with verbose output"
	@echo "  lint               Run lint (frontend + backend)"
	@echo "  lint-frontend      Run frontend lint only"
	@echo "  lint-backend       Run backend lint only"
	@echo "  format             Format code (frontend + backend)"
	@echo "  tidy               Run go mod tidy and frontend deps install"
	@echo "  install-deps       Install all dependencies (go + pnpm)"
	@echo "  install-backend    Install backend (Go) dependencies only"
	@echo "  install-frontend   Install frontend dependencies only"
	@echo "  release            Build for current platform (optimized)"
	@echo "  release-all        Cross-compile for all platforms ($(PLATFORMS))"
	@echo "  clean              Remove build artifacts"
	@echo ""

install-deps: install-backend install-frontend
	@echo "All dependencies installed successfully."

install-backend:
	@echo "==> Installing backend (Go) dependencies..."
	@cd $(BACKEND_DIR) && go mod download && go mod tidy
	@echo "Backend dependencies installed."

install-frontend:
	@echo "==> Installing frontend dependencies..."
	@cd $(FRONTEND_DIR) && pnpm install --frozen-lockfile
	@echo "Frontend dependencies installed."

build-frontend:
	@echo "==> Building frontend..."
	@cd $(FRONTEND_DIR) && pnpm build
	@echo "==> Copying frontend artifacts to $(HOOK_DIST)..."
	@rm -rf $(HOOK_DIST)
	@mkdir -p $(HOOK_DIST)
	@cp -r $(FRONTEND_DIR)/dist/* $(HOOK_DIST)/
	@echo "Frontend build completed."

build-backend:
	@echo "==> Building backend binary..."
	@cd $(BACKEND_DIR) && CGO_ENABLED=$(CGO_ENABLED) go build $(LDFLAGS) -o ../$(BINARY_NAME) .
	@echo "Backend build completed: $(BINARY_NAME)"

build: build-frontend build-backend
	@echo "Full build completed successfully!"

run: build
	@echo "==> Running $(BINARY_NAME)..."
	@./$(BINARY_NAME)

dev:
	@echo "==> Starting frontend dev server..."
	@cd $(FRONTEND_DIR) && pnpm dev

test:
	@echo "==> Running backend tests..."
	@cd $(BACKEND_DIR) && CGO_ENABLED=$(CGO_ENABLED) go test ./... -count=1
	@echo "All tests passed."

test-verbose:
	@echo "==> Running backend tests (verbose)..."
	@cd $(BACKEND_DIR) && CGO_ENABLED=$(CGO_ENABLED) go test ./... -count=1 -v
	@echo "All tests passed."

lint-frontend:
	@echo "==> Running frontend lint..."
	@cd $(FRONTEND_DIR) && pnpm lint

lint-backend:
	@echo "==> Running backend lint (go vet)..."
	@cd $(BACKEND_DIR) && go vet ./...

lint: lint-backend lint-frontend
	@echo "Lint completed."

format:
	@echo "==> Formatting backend code (gofmt)..."
	@gofmt -s -w $(BACKEND_DIR)/
	@echo "==> Formatting frontend code (prettier)..."
	@cd $(FRONTEND_DIR) && pnpm format
	@echo "Format completed."

tidy:
	@echo "==> Running go mod tidy..."
	@cd $(BACKEND_DIR) && go mod tidy
	@echo "==> Installing frontend dependencies..."
	@cd $(FRONTEND_DIR) && pnpm install
	@echo "Tidy completed."

release:
	@echo "==> Building release for current platform..."
	@cd $(BACKEND_DIR) && CGO_ENABLED=$(CGO_ENABLED) go build $(LDFLAGS) -o ../$(BINARY_NAME) .
	@echo "Release build completed: $(BINARY_NAME)"

release-all: build-frontend
	@echo "==> Cross-compiling for all platforms..."
	@for platform in $(PLATFORMS); do \
		goos=$${platform%/*}; \
		goarch=$${platform#*/}; \
		ext=""; \
		if [ "$$goos" = "windows" ]; then ext=".exe"; fi; \
		output="$(BINARY_NAME)_$${goos}_$${goarch}$${ext}"; \
		echo "  -> Building $$output ..."; \
		( cd $(BACKEND_DIR) && CGO_ENABLED=$(CGO_ENABLED) GOOS=$$goos GOARCH=$$goarch go build $(LDFLAGS) -o ../$$output . ) || exit 1; \
	done
	@echo "All releases built successfully."

clean:
	@echo "==> Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_NAME)_*
	@rm -rf $(HOOK_DIST)
	@rm -rf $(FRONTEND_DIR)/dist
	@cd $(BACKEND_DIR) && go clean
	@echo "Clean completed."
