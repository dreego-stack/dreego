---
type: Decision
title: App-Bound Routing and Components
description: Explicit generated registration, one route file per URL, and typed component calls
tags: [pre-v0.1, routing, components, app]
timestamp: 2026-08-14T00:00:00Z
---
# App-Bound Routing and Components

**Status:** Accepted target for v0.1; implementation is pending

## Context

The released pre-v0.1 implementation uses method filenames, generated package
initializers, global runtime registration, implicit component discovery, and
several historical plugin-discovery experiments. Those mechanisms conflict with
explicit App ownership and make route and component behavior harder to reason
about.

This decision defines the migration target. Public documentation may continue
to describe released behavior only when it clearly labels that behavior as
current and provisional.

## Generated registration

Generated code registers explicitly with its owning App:

```go
app := dreego.New()
gen.Register(app)
app.Listen(":8080")
```

Generated packages do not register routes through `init`, blank imports, or
package-global state. Two App instances can therefore own independent routes,
middleware, sessions, and plugins in one process.

## Route sources

One `.dreego` route file owns all declared HTTP methods for one URL.

| Source | URL |
|---|---|
| `routes/+page.dreego` | `/` |
| `routes/about.dreego` | `/about` |
| `routes/about/+page.dreego` | `/about` |
| `routes/users/[id]/+page.dreego` | `/users/{id}` |
| `routes/blog/[...path]/+page.dreego` | `/blog/{path...}` |
| `routes/(auth)/login.dreego` | `/login` |

A flat file and `+page.dreego` resolving to the same URL conflict.
`index.dreego` and optional segments are unsupported. Static segments take
priority over dynamic segments, which take priority over catch-all segments.
Route groups organize source without adding a URL segment.

## HTTP method sections

`<go>` and `<div>` default to GET. An explicit `method` attribute binds the
section to another HTTP method:

```dreego
import UserResult "components/UserResult.dreego"

<go method="post">
result, err := createUser(c)
if err != nil {
    return "", err
}
</go>

<div method="post">
    <@UserResult result={result} />
</div>
```

Successfully reaching the end of method logic renders the matching `<div>`.
An action returning `nil` continues to rendering. Dreego never inserts an
automatic redirect: `c.Redirect` is explicit and suppresses rendering. A normal
Go error enters the App error path.

Duplicate generated, plugin, user, and reserved framework routes fail during
registration or generation with both source locations. Plugins register routes
explicitly against their owning App and cannot silently override another route.

## Component imports and calls

Component declarations and imports are the only directives allowed outside the
five root sections: `<go>`, `<head>`, `<div>`, `<style>`, and `<script>`.

```dreego
import Button "components/Button.dreego"

<div>
    <@Button class="primary" disabled={isLoading}>
        Submit
    </@Button>
</div>
```

Imports are explicit. There is no implicit user-versus-plugin namespace
fallback. A plugin that provides components exposes a documented import path;
the application chooses the imported name.

Component props are named and order-independent. Unknown, duplicate, and
missing required props fail generation with a source location. Literal and
expression props retain their declared Go types. Nested components and
lexically scoped default or named slots are required. Component styles remain
scoped to the component instance contract.

## Consequences

- `app-runtime.1` replaces global registration and `init` side effects.
- `routing-correctness.1` implements flat files, `+page.dreego`, conflicts, and
  method-specific sections.
- `component-correctness.1` implements explicit imports, typed named props,
  nesting, and lexical slots.
- The current fat Plugin interface is not part of this target.
- SSG, Wails expansion, and runtime client hydration do not change this v0.1
  SSR contract.
