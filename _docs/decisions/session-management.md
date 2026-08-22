
---
type: Decision
title: Session Management — Interface in Core, Store as Plugin
description: Session interface in the core with built-in cookie store and external store plugins
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Session Management — Interface in Core, Store as Plugin

**Date:** 2026-07-28
**Status:** Accepted

## Context

SSR apps need sessions: request mapping, CSRF protection, flash messages, shopping cart, form wizards. Auth (login, `c.User()`) is a DIFFERENT topic and belongs in dreego-auth.

## Decision

**Session interface in the core. Concrete stores as plugins.**

```go
// Core: dreego/session (Interface)
type Store interface {
    Get(r *http.Request, key string) (string, error)
    Set(w http.ResponseWriter, r *http.Request, key, value string, opts *Options) error
    Delete(w http.ResponseWriter, r *http.Request, key string) error
    Destroy(w http.ResponseWriter, r *http.Request) error
}
```

```go
// Core: Built-in Cookie Store (Default, no external dependency)
type CookieStore struct { ... }
func NewCookieStore(secret []byte) *CookieStore

// Plugin: dreego-session-redis
// Plugin: dreego-session-db (PostgreSQL/SQLite)
```

## dreego-auth Builds on This

```go
func (p *AuthPlugin) login(w, r) {
    session.Set(w, r, "user_id", user.ID, nil)
}

func (p *AuthPlugin) Middlewares() []func(http.Handler) http.Handler {
    return []func(http.Handler) http.Handler{
        func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w, r) {
                uid, _ := session.Get(r, "user_id")
                ctx := context.WithValue(r.Context(), "user_id", uid)
                next.ServeHTTP(w, r.WithContext(ctx))
            })
        },
    }
}
```

## Why Not Only Plugin?

- CSRF protection (Core) needs sessions
- Flash messages (Core) need sessions
- Without session interface in the core, every plugin would need its own session logic
- Cookie store is 50 lines of Go — no reason to outsource it

## Why Not Completely in Core?

- Redis/DB stores are optional — not every app needs them
- Session interface is stable, stores can grow arbitrarily
- Core stays slim (Cookie store is the only built-in store)

## Consequences

- `dreego.New()` automatically creates a cookie store
- Apps can upgrade via `app.Use(session.NewRedisStore(...))`
- dreego-auth uses `session.Store` interface — no vendor lock-in
- All plugins (not just Auth) can use sessions
