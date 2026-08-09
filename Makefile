.PHONY: up down build generate dev clean dx dx-clean test install-hooks

up:
	docker compose up -d

down:
	docker compose down

build:
	go build -ldflags "-X main.version=$$(git describe --tags --abbrev=0 2>/dev/null || echo dev)" -o bin/dreego ./cli/dreego

generate:
	go run ./cli/dreego

dev:
	go run ./cli/dreego && go run .

test:
	@docker build -q --build-arg DREEGO_VERSION="$$(git describe --tags --abbrev=0 2>/dev/null || echo dev)" -f _tests/Dockerfile -t dreego-test . > /dev/null 2>&1
	@docker run --rm -e DREEGO_FILTER="$${DREEGO_FILTER:-}" -e DREEGO_VERSION="$$(git describe --tags --abbrev=0 2>/dev/null || echo dev)" dreego-test

dx:
	@EXT_DIR="$$(pwd)/.vscode/extensions/dreego"; \
	TARGET="$$HOME/.vscode/extensions/dreego"; \
	rm -rf "$$TARGET" 2>/dev/null; \
	ln -s "$$EXT_DIR" "$$TARGET" && echo "dreego extension installed — restart VS Code"

dx-clean:
	@rm -rf "$$HOME/.vscode/extensions/dreego" 2>/dev/null
	@echo "dreego extension removed"

clean:
	rm -f *_dreego.go
	rm -rf bin/

install-hooks:
	@cp _scripts/pre-commit .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@cp _scripts/pre-push .git/hooks/pre-push
	@chmod +x .git/hooks/pre-push
	@echo "git hooks installed: pre-commit, pre-push"
