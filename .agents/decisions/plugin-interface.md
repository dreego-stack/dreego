---
type: Decision
title: Provisional Plugin Registration
description: Plugins use explicit App-bound registration functions before v1
tags: [pre-v0.1, plugins, app]
timestamp: 2026-08-14T00:00:00Z
---
# Provisional Plugin Registration

**Date:** 2026-08-14
**Status:** Accepted for v0.1; compatibility remains provisional until v1

## Context

The previous design standardized a base `Plugin` interface plus optional route,
middleware, asset, transpiler, context, and lifecycle capabilities before real
external plugins existed. The current all-in-one implementation also exposes
capabilities that are not fully connected to the server.

Dreego needs an explicit extension mechanism for the v0.1 `App` architecture,
but it does not need a universal plugin object or premature stability promise.

## Decision

Plugins are ordinary Go packages that expose an App-bound registration
function with typed options:

```go
package auth

type Options struct {
    LoginPath  string
    CookieName string
}

func Register(app *dreego.App, options Options) error {
    app.Use(sessionMiddleware(options))
    app.Get(options.LoginPath, loginHandler)
    app.Post(options.LoginPath, authenticateHandler)
    return nil
}
```

Applications call registration explicitly and handle errors locally:

```go
app := dreego.New()

if err := auth.Register(app, auth.Options{
    LoginPath:  "/login",
    CookieName: "session",
}); err != nil {
    log.Fatal(err)
}
```

There is no stable central `Plugin` interface before v1. Registration order is
source order. Each plugin owns its typed configuration and registers only the
routes, middleware, or other App behavior it needs.

Assets and startup or shutdown hooks are added through explicit App methods or
small capability interfaces only after real plugins prove a shared contract.
Transpiler extensions require their own validated processor boundary and are
not implied by runtime plugin registration.

## Consequences

- The current fat `Plugin` interface is removed during the App migration.
- Plugins live in separate repositories with their own dependencies and releases.
- No empty lifecycle, asset, route, or middleware methods are required.
- Registration errors are ordinary Go errors and never hidden in package `init`.
- Auth, UI, and at least one infrastructure plugin validate common needs before v1.
- Compatibility guarantees for a shared plugin contract begin at v1, not v0.1.
