.PHONY: up down build generate dev clean dx dx-clean dx-pkg

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

dx: dx-pkg
	@code --install-extension .vscode/dreego-extension/dreego.vsix 2>/dev/null && \
		echo "dreego extension installed (reload VS Code)" || \
		echo "dreego extension installed → reload VS Code (Developer: Reload Window)"

dx-pkg:
	@cd .vscode/dreego-extension && \
		rm -f dreego.vsix && \
		zip -qr dreego.vsix package.json language-configuration.json syntaxes/ icons/

dx-clean:
	@code --uninstall-extension dreego 2>/dev/null
	@rm -f .vscode/dreego-extension/dreego.vsix
	@echo "dreego extension removed"

clean:
	rm -f *_dreego.go
	rm -rf bin/
