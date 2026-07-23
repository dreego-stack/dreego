.PHONY: up down build generate dev clean

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

clean:
	rm -f *_dreego.go
	rm -rf bin/
