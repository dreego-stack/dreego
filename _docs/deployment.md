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

The landing blueprint uses a two-stage build and a non-root distroless runtime:

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/bin/server .

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app/bin/server /server
EXPOSE 8080
USER nonroot
ENTRYPOINT ["/server"]
```

Generated static assets are embedded in the binary. No separate
`dreego/static` directory is copied into the runtime image.

## Runtime

- `app.Listen(":8080")` binds and serves, with graceful shutdown on SIGINT/SIGTERM (10s drain).
- Health/ready endpoints: `GET /health` (liveness) and `GET /ready` (readiness via `app.SetReady`).
- Static assets are embedded at build time, so a single binary serves everything.
