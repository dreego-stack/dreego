FROM golang:1.22-alpine AS build

WORKDIR /app
COPY go.mod ./
RUN go mod download 2>/dev/null || true

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /dreego ./cmd/dreego

FROM golang:1.22-alpine AS app

WORKDIR /app
COPY --from=build /dreego /usr/local/bin/dreego

RUN dreego init demo
WORKDIR /app/demo
RUN dreego generate
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/demo/bin/app .

FROM scratch
COPY --from=app /app/demo/bin/app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
