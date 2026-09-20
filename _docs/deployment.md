# Deployment

Dreego apps are plain Go programs. Build and deploy like any `net/http` Go server.

## Building

```sh
dreego generate
dreego build
```

`dreego build` produces a static binary into `build/bin/`.

## Cross-Compilation

`dreego build --target <os>/<arch>` sets `GOOS`/`GOARCH` for the target platform:

```sh
# Build a linux/amd64 binary (e.g. for a CI runner or production host)
dreego build --target linux/amd64

# Binary is written to build/bin/*-linux-amd64
```

For statically linked, scratch-container-friendly binaries build with `CGO_ENABLED=0`. You can set `GOOS`/`GOARCH` yourself before building:

```sh
GOOS=linux GOARCH=arm64 dreego build
```

## Container

Scaffolded projects ship a multi-stage `Dockerfile` and a `docker-compose.yml`
that build the Dreego CLI, run `dreego generate`, and compile a static binary
into a `scratch` runtime image. Run it with:

```sh
docker compose up --build
```

The generated `Dockerfile` installs the CLI from the module proxy; pin it to
your release with `--build-arg DREEGO_CLI_VERSION=v0.1.0`. The listening port is
the `port` constant in `main.go`; `DREEGO_PORT` overrides it at runtime, and the
scaffolded `Dockerfile`/`docker-compose.yml` read the same variable (build arg
`DREEGO_PORT`, default `8080`). A minimal hand-written equivalent looks like
this:

```dockerfile
FROM golang:1.27-alpine AS build
WORKDIR /src
RUN go install github.com/dreego-stack/dreego/cmd/dreego@latest
COPY . .
RUN dreego generate && CGO_ENABLED=0 go build -o /app -ldflags="-s -w" .

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
```

Generated static assets are embedded in the binary. No separate
`dreego/static` directory is copied into the runtime image.

## Runtime

- `ssr.Listen(app, ":8080")` binds and serves, with graceful shutdown on SIGINT/SIGTERM (10s drain).
- Health/ready endpoints: `GET /health` (liveness) and `GET /ready` (readiness via `app.SetReady`).
- Static assets are embedded at build time, so a single binary serves everything.
