.PHONY: up down build generate dev clean dx dx-clean

up:
	docker compose up -d

down:
	docker compose down

build:
	go build -o bin/dreego ./cmd/dreego

generate:
	go run ./cmd/dreego

dev:
	go run ./cmd/dreego && go run .

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
