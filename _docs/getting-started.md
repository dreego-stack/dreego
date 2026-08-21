# Getting Started

This guide is the canonical path from installation to a running Dreego
application. The same commands run in CI as a black-box test, so the steps
below are guaranteed to work on a clean machine with Go 1.22 or newer.

## Prerequisites

Dreego requires Go 1.22 or newer. Check your installation:

```bash
go version
```

If `go` is not found, install it from https://go.dev/doc/install. If the
version is older than 1.22, upgrade before continuing — generated code uses
Go 1.22 features and `go mod` directives.

## 1. Install the CLI

```bash
go install github.com/dreego-stack/dreego/cli/dreego@latest
```

This installs the `dreego` binary into your `GOPATH/bin` directory. Make sure
that directory is on your `PATH` (the Go installer usually adds it).

## 2. Create a project

```bash
dreego new myapp
cd myapp
```

`dreego new` scaffolds a project from the `landing` blueprint:

- writes `main.go`, `go.mod`, `Dockerfile`, `.gitignore`
- writes the `www/` tree: `routes/`, `layouts/`, `components/`, `dreego.config.json`
- runs `go mod init` and `go mod tidy` against the published `dreego` module
  (resolved from the public Go proxy — no `replace` directive)

The project name must be a valid Go module path segment (letters, digits,
hyphens, underscores; must start with a letter). `dreego new myapp` creates a
module named `myapp`; `dreego new github.com/me/myapp` is also accepted.

The landing starter loads Tailwind's browser script from its CDN and its
generated `main.go` explicitly allows that origin in the Content Security
Policy. This is convenient for evaluating the blueprint. For production,
replace the browser script with locally built CSS and remove the CDN origin
from the policy.

## 3. Generate and run

```bash
dreego generate    # transpiles .dreego files → dree.go per directory
go run .            # builds and starts the server on :8080
```

Open http://localhost:8080 in your browser. The landing page rendered is the
one defined in `www/routes/get.dreego`.

For day-to-day development:

```bash
dreego build       # generate + go build → build/bin/<name>
dreego run         # build + start server (dev only)
dreego run -d      # with debug logging (JSONL)
dreego dev         # watch .dreego files, rebuild + restart on change
```

> **Note:** `dreego build` and `dreego run` are dev tools, not for production.
> Production builds use `go build` (or `dreego build --target <os/arch>` for
> cross-compilation) plus the Dockerfile that `dreego new` wrote.

## main.go

The scaffolded `main.go` uses the explicit App API — no globals, no hidden
state, no runtime registration magic:

```go
package main

import (
	"log"

	dreego "github.com/dreego-stack/dreego/core"
	"myapp/www"
)

func main() {
	app := dreego.New()
	if err := www.Register(app); err != nil {
		log.Fatal(err)
	}
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
```

`dreego.New()` returns an `*App` that owns all runtime state (router,
middleware, session store, server config). `www.Register(app)` wires the
generated routes and components into the `App`. `app.Listen(":8080")` starts
the HTTP server with the configured timeouts.

## Adding a Layout

Create `www/layouts/default.dreego` — wraps all pages:

```html
<head><title>My App</title></head>

<div>
    <nav><a href="/">Home</a> | <a href="/about">About</a></nav>
    <main>{#slot}</main>
</div>

<style>
    nav { padding: 1rem; background: #1e293b; }
    nav a { color: #e2e8f0; margin-right: 1rem; }
</style>
```

## Creating a Component

Create `www/components/Card.dreego`:

```
Component Card (title string)

<div>
    <article class="card">
        <h2>{{ title }}</h2>
        <div>{#slot}</div>
    </article>
</div>

<style>
.card { border: 1px solid #e2e8f0; padding: 1rem; border-radius: 8px; }
</style>
```

Use it in any route or layout:

```html
import Card "components/Card.dreego"

<div>
<@Card title="Welcome">
    <p>This is the card body.</p>
</@Card>
</div>
```

Imports are header directives and therefore appear before the root sections.

## Dynamic Routes

Create `www/routes/users/[id]/get.dreego`:

```html
<head><title>User {{ c.Param("id") }}</title></head>

<go>
    userID := c.Param("id")
</go>

<div>
    <h1>User: {{ userID }}</h1>
</div>
```

Visiting `/users/42` shows "User: 42".

## Troubleshooting

| Symptom | Cause / Fix |
|---------|-------------|
| `dreego: command not found` | `go install` put the binary in `$(go env GOPATH)/bin`; add it to `PATH`. |
| `go: command not found` | Install Go 1.22+ from https://go.dev/doc/install. |
| `go: go.mod requires ... but ...` | Your Go toolchain is older than 1.22. Upgrade. |
| `dreego new: invalid project name "..."` | The name must be a valid Go module path segment (start with a letter; only letters, digits, `-`, `_`, `/`, `.`). |
| `go mod tidy: ... unresolved dependency` | No network, or the CLI was built from an untagged checkout so the published tag is unknown. Set `DREEGO_LOCAL_REPO=/path/to/dreego` to point the scaffold at a local checkout. |
| `dreego generate: no routes found` | Create at least `www/routes/get.dreego` (the scaffold already does). |

## See Also

- [Components](https://github.com/dreego-stack/dreego/blob/main/_docs/components.md) — full component docs
- [Routing](https://github.com/dreego-stack/dreego/blob/main/_docs/routing.md) — dynamic segments, groups, methods
- [Runtime API](https://github.com/dreego-stack/dreego/blob/main/_docs/runtime.md) — SSRContext, sessions, config
- [CLI Reference](https://github.com/dreego-stack/dreego/blob/main/_docs/cli.md)
- [Docs Index](https://github.com/dreego-stack/dreego/blob/main/_docs/index.md)
