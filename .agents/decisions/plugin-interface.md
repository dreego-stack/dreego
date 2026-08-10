
---
type: Decision
title: Plugin Interface (Capability-based)
description: Capability-based plugin system with Go-typical implicit interfaces
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Plugin Interface (Capability-based)

**Date:** 2026-07-28
**Status:** Accepted
**Review:** GLM-5.2 Expert Review (.tmp/output2.md)

## Context

Dreego needs a plugin system for its plugin ecosystem. The interface must be stable from the start — later changes are breaking changes for all plugins.

## Decision

**Capability-based interfaces.** Instead of a fat single interface, plugins only implement the capabilities they need — via Go's implicit interface satisfaction.

## Interface Definition

```go
package dreego

import (
    "context"
    "io/fs"
    "net/http"
)

type PluginID string

// Base — EVERY plugin MUST implement this
type Plugin interface {
    ID() PluginID
    Init(app *App, cfg map[string]any) error
}

// Inject SSR middleware (optional)
type MiddlewareProvider interface {
    Middlewares() []func(http.Handler) http.Handler
}

// Register routes — SSR + SSG (optional)
type RouteRegistrar interface {
    RegisterRoutes(r Router) error
}

// Embedded assets (CSS/JS/Images) — target-agnostic via fs.FS (optional)
type AssetProvider interface {
    Assets() fs.FS
}

// Custom tags like <dreego:map /> in the transpiler (optional)
type TranspilerHook interface {
    Namespace() string
    ParseTag(node TagNode) (TagRenderer, error)
}

// Extend request context (optional)
type ContextExtender interface {
    BindContext(b *ContextBuilder)
}

// Startup/Shutdown (optional)
type Lifecycle interface {
    OnStart(ctx context.Context) error
    OnShutdown(ctx context.Context) error
}
```

## Example: dreego-auth

```go
package auth

import (
    "embed"
    "github.com/dreego/dreego"
)

//go:embed assets/*
var assets embed.FS

type AuthPlugin struct {
    secretKey string
}

func New(secretKey string) *AuthPlugin {
    return &AuthPlugin{secretKey: secretKey}
}

func (p *AuthPlugin) ID() dreego.PluginID { return "dreego.io/auth" }

func (p *AuthPlugin) Init(app *dreego.App, cfg map[string]any) error {
    p.secretKey = cfg["secret"].(string)
    return nil
}

// AuthPlugin implements MiddlewareProvider
func (p *AuthPlugin) Middlewares() []func(http.Handler) http.Handler {
    return []func(http.Handler) http.Handler{p.authMiddleware}
}

// AuthPlugin implements RouteRegistrar
func (p *AuthPlugin) RegisterRoutes(r dreego.Router) error {
    r.Get("/auth/login", p.loginPage)
    r.Post("/api/auth/login", p.handleLogin)
    return nil
}

// AuthPlugin implements AssetProvider
func (p *AuthPlugin) Assets() fs.FS { return assets }

// AuthPlugin implements ContextExtender
func (p *AuthPlugin) BindContext(b *dreego.ContextBuilder) {
    b.Set(UserKey, func(c dreego.Context) *User {
        return c.Get(UserKey).(*User)
    })
}
```

## Plugin Order

Explicit via registration order — no magic:

```go
func main() {
    app := dreego.New()
    app.Use(session.New())   // 1. Session (MUST be before Auth)
    app.Use(auth.New("key")) // 2. Auth
    app.Use(admin.New())     // 3. Admin
    app.Listen(":8080")
}
```

Middleware is wrapped in execution order (LIFO like Chi).

## Why Capability-based

| Advantage                             | Explanation                                      |
|---------------------------------------|--------------------------------------------------|
| Plugins only implement what they need | `dreego-map` doesn't need middleware             |
| New capabilities can be added additively | `Lifecycle` came later — no breaking change    |
| Go-idiomatic                          | `io.Reader`, `fs.FS` are the same pattern        |
| No empty methods                      | No `return nil` for unused features              |

## Why `fs.FS` instead of `embed.FS`

`fs.FS` is Go's standard interface for filesystems. It works target-agnostically:
- SSR: `embed.FS` implements `fs.FS`
- SSG: `os.DirFS` implements `fs.FS` (files on disk)
- Wails: Custom `fs.FS` implementation possible

## Consequences

- `dreego.Plugin` interface in the core — never change again
- New capabilities can be added as separate interfaces
- No central plugin registry needed — user registers explicitly
- Plugins start with `go get` + `app.Use()`
