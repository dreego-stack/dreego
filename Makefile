.PHONY: up down build generate dev clean dx dx-clean test coverage install-hooks

ROOT_VERSION = $$(git describe --tags --match 'v[0-9]*.[0-9]*.[0-9]*' --abbrev=0 2>/dev/null || echo dev)

up:
	cd demo/demo-ssr && docker compose up -d

down:
	cd demo/demo-ssr && docker compose down

build:
	go build -ldflags "-X main.version=$(ROOT_VERSION)" -o tmp/dreego ./cmd/dreego

generate:
	go run ./cmd/dreego

dev:
	go run ./cmd/dreego && go run .

test:
	@go run ./cmd/dreego tools install typescript
	@sh _scripts/check-wails-demo.sh
	@CGO_ENABLED=1 go test -race ./internal/... ./core/... ./adapter/ssr/... ./adapter/wails/... ./dreegotest/... ./cmd/dreego/... ./_tests/go/...
	@make coverage
	@docker build \
		-q \
		--build-arg DREEGO_VERSION="$(ROOT_VERSION)" \
		-f _tests/Dockerfile \
		-t dreego-test \
		. > /dev/null 2>&1
	@docker run \
		--rm \
		-e DREEGO_FILTER="$${DREEGO_FILTER:-}" \
		-e DREEGO_RUNS="$${DREEGO_RUNS:-1}" \
		-e DREEGO_VERSION="$(ROOT_VERSION)" \
		dreego-test

coverage:
	sh _scripts/coverage-gate.sh

dx:
	@EXT_DIR="$$(pwd)/../vscode-dreego"; \
	TARGET="$$HOME/.vscode/extensions/dreego"; \
	rm -rf "$$TARGET" 2>/dev/null; \
	ln -s "$$EXT_DIR" "$$TARGET" && echo "dreego extension installed from vscode-dreego (https://github.com/dreego-stack/vscode-dreego) — restart VS Code"

dx-clean:
	@rm -rf "$$HOME/.vscode/extensions/dreego" 2>/dev/null
	@echo "dreego extension removed"

clean:
	rm -f *_dreego.go
	rm -rf tmp/

install-hooks:
	@cp _scripts/pre-commit .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@cp _scripts/pre-push .git/hooks/pre-push
	@chmod +x .git/hooks/pre-push
	@echo "git hooks installed: pre-commit, pre-push"
