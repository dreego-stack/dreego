
---
type: Decision
title: Routing, Plugin Routes & Component System
description: Routing, plugin route registration, and component system with namespace hierarchy
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Routing, Plugin Routes & Component System

**Date:** 2026-07-28 (updated after review)
**Status:** Accepted (plugin discovery in Decision 2 superseded by [monorepo-plugin-layout](monorepo-plugin-layout.md) v0.0.21)
**Replaces:** [file-based-routing](file-based-routing.md) (updated)

## Context / Open Questions

1. Go package system: Each directory = one package. `dreego/routes/about/` is package `about`. Currently `main.go` must import each route package individually (`demo/main.go:4-7`). That doesn't scale.

2. Plugin routes: `dreego generate` scans `"."` — never finds external packages in the module cache or vendor. How do plugin routes get into the binary?

3. Component namespace: If the user has `components/Button.dreego` AND `dreego-ui` also offers `Button` — how is it resolved?

## Decision 1: Generated Route Import File Instead of Manual Imports

`dreego generate` creates ONE central import file that imports all route packages. The user only imports this one file.

```
dreego/
├── routes/                    ← User writes here
│   ├── get.dreego           → dreego → init() registers GET /
│   ├── about/get.dreego       → dreego → init() registers GET /about
│   ├── users/[id]/get.dreego  → dreego → init() registers GET /users/{id}
│   └── ...
├── gen/                       ← GENERATED (committed)
│   └── routes.go              → imports ALL route packages
│
main.go imports ONLY `_ "myapp/dreego/gen"`
```

```go
// dreego/gen/routes.go  (GENERATED)
package gen

import (
    _ "myapp/dreego/routes"          // index, about, ...
    _ "myapp/dreego/routes/about"    // if about/ is a subdirectory
    _ "myapp/dreego/routes/users/_id_"
    _ "github.com/dreego/dreego-auth" // Plugin with init() registration
)
```

Each route package (including plugin packages) registers itself via `init()` → `runtime.Register()`. `gen/routes.go` imports all — `main.go` imports only `gen`.

### Why Keep `init()`?

- Go-idiomatic: `database/sql` drivers do exactly the same
- Plugin packages need no special treatment
- No runtime scanning, no reflection
- `go build` only links imported packages → tree shaking

## Decision 2: Plugin Routes via init() — No dreego generate Needed

> **Historical note:** the v0.0.21 monorepo experiment was later superseded. Every optional plugin now lives in a separate repository and module. Registration and discovery will be revised by `app-runtime.1` and `plugin-contract.1`; the pattern below is not a current compatibility promise.

The plugin author commits generated `dree.go` files IN the plugin module.

```
plugins/auth/
├── routes/
│   ├── login.go          ← pre-generated (contains init() + runtime.Register)
│   └── ...
├── go.mod                ← module github.com/dreego-stack/dreego/plugins/auth
```

Plugin developer workflow (official, monorepo):
```bash
cd plugins/auth
dreego generate               # generates routes/*.go
git add routes/*.go && git commit
```

User workflow (official, monorepo):
```bash
# go.work already links plugins/auth for local dev
# main.go imports _ "github.com/dreego-stack/dreego/plugins/auth"
```

Community plugin (separate repo) workflow:
```bash
go get github.com/dreego-stack/dreego-community-auth@v0.1.0
# dreego generate adds the import to gen/routes.go
```

`dreego generate` detects official plugin packages via filesystem scan under `plugins/<name>/` and community plugins via `go.mod` + `go list -m -json all`. `go build` automatically links the correct version.

## Decision 3: Routing Conventions

| Syntax                            | Method | Path                     | Go Param              |
|-----------------------------------|--------|--------------------------|-----------------------|
| `get.dreego`                      | GET    | `/`                      | —                     |
| `about/get.dreego`                | GET    | `/about`                 | —                     |
| `users/[id]/get.dreego`           | GET    | `/users/{id}`            | `c.Param("id")`       |
| `blog/[...catchall]/get.dreego`   | GET    | `/blog/{catchall...}`    | `c.Param("catchall")` |
| `(auth)/login/get.dreego`         | GET    | `/login`                 | —                     |

Optional segments are deliberately unsupported. Each method file owns one
explicit route pattern.

Priority: Static > Dynamic > Catch-All

Conflict detection: `dreego generate` throws an error if two routes claim the same pattern:
```
error: route conflict: /auth/login
  dreego/routes/auth/login.dreego
  plugin: dreego-auth (github.com/dreego/dreego-auth)
```

### API Routes & HTTP Methods

Directories define the URL path. The method-only filename keeps each HTTP
operation in a separate, focused file:

```
routes/api/users/get.dreego    → GET    /api/users
routes/api/users/post.dreego   → POST   /api/users
routes/api/users/put.dreego    → PUT    /api/users
routes/api/users/delete.dreego → DELETE /api/users
```

API routes render NO layout — only the `<div>` fragment. Detection: path contains `api/` → `layout = nil`.

### Per-Route Middleware (V1)

One `_middleware.go` per directory (NOT generated — the user writes it):

```go
// dreego/routes/admin/_middleware.go
package admin

import "github.com/dreego/dreego-auth"

func init() {
    runtime.RegisterMiddleware("/admin/", auth.RequireRole("admin"))
}
```

### Redirects & Rewrites

```json
// dreego.config.json
{
  "redirects": [
    { "from": "/old-blog", "to": "/blog", "status": 301 }
  ],
  "rewrites": [
    { "from": "/api/v1/*", "to": "/api/v2/*" }
  ]
}
```

Generated in `gen/routes.go` as middleware logic before the file-based routes.

## Decision 4: Component System

### Three Sources — One Namespace Hierarchy

```
Priority (highest first):
1. dreego/components/Button.dreego     ← User component (shadows plugin)
2. dreego/layouts/default.dreego       ← Layouts (special case)
3. Plugin assets via fs.FS             ← @dreego-ui/Button
```

Explicit disambiguation:
```
{#use Button from "components/Button.dreego"}     ← explicit user
{#use Button from "@dreego-ui/Button"}            ← explicit plugin
```

Without `from` specification, search user directory first, then plugins:
```
{#use Button}   ← searches Button.dreego in components/, then in plugins
```

### How Does dreego generate Find Plugin Components?

Plugins place `.dreego` components in a known path:

```
dreego-ui/
├── components/               ← CONVENTION
│   ├── Button.dreego
│   ├── Card.dreego
│   └── Alert.dreego
├── dreego.go                 ← Plugin interface impl
├── go.mod
```

`dreego generate`:
1. Reads `go.mod` → finds `github.com/dreego/dreego-ui`
2. `go list -m -json github.com/dreego/dreego-ui` → `Dir: /home/.../pkg/mod/...`
3. Searches `<Dir>/components/*.dreego`
4. In vendor mode: `vendor/github.com/dreego/dreego-ui/components/*.dreego`

No plugin loading, no reflection, no import in the CLI binary. Only filesystem scan.

### Component Usage in Templates

```html
<!-- dreego/routes/get.dreego -->
{#use Button}              <!-- finds components/Button.dreego -->
{#use Alert}               <!-- finds @dreego-ui/Alert (no user Alert) -->

<div>
    <Alert type="success">Done!</Alert>
    <Button class="primary" disabled={isLoading}>
        Submit
    </Button>
</div>
```

### Component CodeGen

The transpiler generates:
```go
// <Button class="primary" disabled={isLoading}>Submit</Button>
// becomes:
renderButton(c, ButtonProps{Class: "primary", Disabled: isLoading}, func(c *Context) string {
    return "Submit"
})
```

- Props: All HTML attributes are bundled as a `ComponentProps` struct
- Children: Content between tags as a closure (corresponds to `{#slot}`)
- Self-closing: `<Alert type="info" />` → no children parameter
- Scoping: Each component has its own scope hash (no CSS leaking)

### Plugin Components Without .dreego Source (V2)

If a plugin provides only Go functions for performance reasons (no `.dreego` transpilation needed):

```go
// dreego-ui registers a named components map
func (p *UIPlugin) Components() map[string]ComponentFunc {
    return map[string]ComponentFunc{
        "Button": renderButton,
        "Card":   renderCard,
    }
}
```

`{#use Button from "@dreego-ui/Button"}` → direct function call, no transpiler pass needed.

## Decision 5: dreego/routes/ vs dreego/pages/

`dreego/routes/` stays. The name is established (SvelteKit, SolidStart, Next.js Pages Router). `/pages/` would be equally valid, but we stick with the existing convention.

Configurable via `dreego.config.json`:
```json
{ "routeDir": "dreego/pages" }
```

## Consequences

- `dreego/gen/routes.go` is generated and committed
- Contains ALL route imports (file-based + plugin)
- `main.go` imports only `_ "myapp/dreego/gen"`
- Plugin developers commit `dree.go` files (pre-generated)
- Plugin components are located under `<plugin>/components/` (convention)
- `dreego generate` finds them via `go list` + filesystem
- User components shadow plugin components (explicit namespace fallback)
- `[...catchall]` and `(group)/` are added in lexer/parser
- Duplicate route detection throws build errors
