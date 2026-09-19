# Build the Dreego CLI, generate the app, and compile a static binary.
# Pin the CLI to the same release as the module in go.mod when needed:
#   docker build --build-arg DREEGO_CLI_VERSION=v0.1.0 .
FROM golang:1.27-alpine AS build
ARG DREEGO_CLI_VERSION=latest
WORKDIR /src
RUN go install github.com/dreego-stack/dreego/cmd/dreego@${DREEGO_CLI_VERSION}
COPY . .
RUN dreego generate && CGO_ENABLED=0 go build -o /app -ldflags="-s -w" .

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
