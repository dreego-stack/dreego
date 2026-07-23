FROM golang:1.22-alpine AS build

WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download 2>/dev/null || true

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /dreego ./cmd/dreego

FROM scratch
COPY --from=build /dreego /dreego
EXPOSE 8080
ENTRYPOINT ["/dreego"]
