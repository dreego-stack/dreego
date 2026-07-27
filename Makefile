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
	@mkdir -p $(HOME)/.vscode/extensions
	@rm -rf $(HOME)/.vscode/extensions/dreego
	@ln -sf $(PWD)/.vscode/dreego-extension $(HOME)/.vscode/extensions/dreego
	@echo "dreego extension installed → ~/.vscode/extensions/dreego"
	@echo "restart VS Code or run 'Developer: Reload Window'"

dx-clean:
	rm -rf $(HOME)/.vscode/extensions/dreego
	@echo "dreego extension removed"

clean:
	rm -f *_dreego.go
	rm -rf bin/
