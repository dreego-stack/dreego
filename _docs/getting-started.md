# Getting Started

This guide is the canonical path from installation to a running Dreego
application. The same commands run in CI as a black-box test, so the steps
below are guaranteed to work on a clean machine with Go 1.27 or newer.

## Prerequisites

Dreego requires Go 1.27 or newer. Check your installation:

```bash
go version
```

If `go` is not found, install it from https://go.dev/doc/install. If the
version is older than 1.27, upgrade before continuing — generated code uses
Go 1.27 features and `go mod` directives.

## 1. Install the CLI

```bash
go install github.com/dreego-stack/dreego/cmd/dreego@latest
```

The install path must include the `/cmd/dreego` suffix. The module root
`github.com/dreego-stack/dreego` is the shared implementation module and does
not contain an installable `main` package, so installing it directly fails with
`module … found, but does not contain package`.

This installs the `dreego` binary into your `GOPATH/bin` directory. Make sure
that directory is on your `PATH` (the Go installer usually adds it).

## 2. Create a project

```bash
dreego new myapp
cd myapp
```

`dreego new` scaffolds a project from the `web-minimal` template:

- writes `main.go`, `Taskfile.yml`, `.gitignore`
- writes the `www/` tree: `routes/+page.dreego`, `layouts/default.dreego`,
  and `dreego.config.json`
- runs `go mod init` and `go mod tidy` against the published `dreego` module
  (resolved from the public Go proxy — no `replace` directive)

Both starters style themselves with local `<style>` blocks and load no Tailwind
or any other stylesheet from a CDN. No CDN origin is added to the Content
Security Policy.

Two templates ship today. `web-minimal` is the default for `dreego new` and
stays the smallest starting point. `web-app` is a full SSR
application starter: an app shell with a `Nav` component in the layout header, a
`Card` component, `/` with a server-rendered typed form, and a nested
`/dashboard` route. `-t` (or `--template`) selects a template explicitly, and
`-l` (or `--list`) lists the available templates:

```bash
dreego new myapp -t web-app
dreego new myapp -l
```

The project name must be a valid Go module path segment (letters, digits,
hyphens, underscores; must start with a letter). `dreego new myapp` creates a
module named `myapp`; `dreego new github.com/me/myapp` is also accepted.

## 3. Generate and run

```bash
dreego generate    # transpiles .dreego files → dree.go per directory
go run .            # builds and starts the server on :8080
```

Open http://localhost:8080 in your browser. The page rendered is the one
defined in `www/routes/+page.dreego`.

For day-to-day development:

```bash
dreego build       # generate + go build → build/bin/<name>
dreego run         # build + start server (dev only)
dreego run -d      # with debug logging (JSONL)
dreego dev         # watch .dreego files, rebuild + restart on change
```

> **Note:** `dreego build` and `dreego run` are dev tools, not for production.
> Production builds use `go build` (or `dreego build --target <os/arch>` for
> cross-compilation).

## main.go

The scaffolded `main.go` uses the explicit App API — no globals, no hidden
state, no runtime registration magic:

```go
package main

import (
	"log"
	"os"

	dreego "github.com/dreego-stack/dreego/core"
	"github.com/dreego-stack/dreego/adapter/ssr"
	"myapp/www"
)

const port = "8080"

func main() {
	app := dreego.New()
	if err := www.Register(app); err != nil {
		log.Fatal(err)
	}
	addr := ":" + port
	if p := os.Getenv("DREEGO_PORT"); p != "" {
		addr = ":" + p
	}
	if err := ssr.Listen(app, addr); err != nil {
		log.Fatal(err)
	}
}
```

`dreego.New()` returns an `*App` that owns route declarations, middleware, and
session policy. `www.Register(app)` wires generated routes and components into
the `App`. `ssr.Listen(app, addr)` creates the explicit HTTP host with secure
timeout defaults. The listening port is a `const` in `main.go`, so there is one
obvious place to change it; `DREEGO_PORT` overrides it at runtime for
containers.

## Adding a Layout

Create `www/layouts/default.dreego` — wraps all pages:

```html
<head><title>My App</title></head>

<body>
    <nav><a href="/">Home</a> | <a href="/about">About</a></nav>
    <main>{#slot}</main>
</body>

<style>
    nav { padding: 1rem; background: #1e293b; }
    nav a { color: #e2e8f0; margin-right: 1rem; }
</style>
```

## Creating a Component

Create `www/components/Card.dreego`:

```
DREEFILE component (title string)

<body>
    <article class="card">
        <h2>{{ title }}</h2>
        <div>{#slot}</div>
    </article>
</body>

<style>
.card { border: 1px solid #e2e8f0; padding: 1rem; border-radius: 8px; }
</style>
```

Use it in any route or layout:

```html
COMPONENT "www/components" IMPORT { Card }

<body>
<@Card title="Welcome">
    <p>This is the card body.</p>
</@Card>
</body>
```

The `COMPONENT … IMPORT` directive is a header directive and therefore appears
before the root sections. The component name comes from its filename, so
`Card.dreego` is called as `<@Card>`.

## Dynamic Routes

Create `www/routes/users/[id]/+page.dreego`:

```html
<head><title>User {{ c.Param("id") }}</title></head>

<server>
    userID := c.Param("id")
</server>

<body>
    <h1>User: {{ userID }}</h1>
</body>
```

Visiting `/users/42` shows "User: 42".

## Troubleshooting

| Symptom | Cause / Fix |
|---------|-------------|
| `dreego: command not found` | `go install` put the binary in `$(go env GOPATH)/bin`; add it to `PATH`. |
| `go: command not found` | Install Go 1.27+ from https://go.dev/doc/install. |
| `go: go.mod requires ... but ...` | Your Go toolchain is older than 1.27. Upgrade. |
| `dreego new: invalid project name "..."` | The name must be a valid Go module path segment (start with a letter; only letters, digits, `-`, `_`, `/`, `.`). |
| `go mod tidy: ... unresolved dependency` | No network, or the CLI was built from an untagged checkout so the published tag is unknown. Set `DREEGO_LOCAL_REPO=/path/to/dreego` to point the scaffold at a local checkout. |
| `dreego generate: no routes found` | Create at least `www/routes/+page.dreego` (the scaffold already does). |

## See Also

- [Components](https://github.com/dreego-stack/dreego/blob/main/_docs/components.md) — full component docs
- [Routing](https://github.com/dreego-stack/dreego/blob/main/_docs/routing.md) — dynamic segments, groups, methods
- [Runtime API](https://github.com/dreego-stack/dreego/blob/main/_docs/runtime.md) — SSRContext, sessions, config
- [CLI Reference](https://github.com/dreego-stack/dreego/blob/main/_docs/cli.md)
- [Docs Index](https://github.com/dreego-stack/dreego/blob/main/_docs/index.md)
