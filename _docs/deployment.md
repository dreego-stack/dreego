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

For statically linked, scratch-container-friendly binaries the generated blueprints / Dockerfiles build with `CGO_ENABLED=0`. You can set `GOOS`/`GOARCH` yourself before building:

```sh
GOOS=linux GOARCH=arm64 dreego build
```

## Container

A production Dockerfile uses a 3-stage build with `FROM scratch`:

```dockerfile
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o /dreego .

FROM scratch
COPY --from=build /dreego /dreego
COPY --from=build /app/dreego/static /dreego/static
ENTRYPOINT ["/dreego"]
```

## Runtime

- `core.Listen(":8080")` binds and serves, with graceful shutdown on SIGINT/SIGTERM (10s drain).
- Health/ready endpoints: `GET /health` (liveness) and `GET /ready` (readiness via `core.SetReady`).
- Static assets are embedded at build time, so a single binary serves everything.
