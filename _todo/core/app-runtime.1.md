---
area: architecture
phase: v0.1-blocker
priority: 1
---
# Make App own all runtime state

## Goal
Replace package-global runtime state with an explicit `App` instance before v0.1. There is no compatibility requirement for the current global API because Dreego has no established user base yet.

The intended application shape is explicit:

```go
app := dreego.New()
gen.Register(app)
app.Listen(":8080")
```

## Scope
- Move routes, redirects, rewrites, static assets, plugins, middleware, sessions, logging, CSRF, CSP, readiness, error handlers, and the built handler cache into `App`.
- Replace generated package-level `init()` registration with an explicit generated function such as `gen.Register(app)`.
- Update route, error-page, configuration, and static-asset code generation to target the supplied app.
- Update blueprints, CLI build/dev/run flows, `dreegotest`, demo, documentation, and all tests.
- Remove global registration and configuration functions instead of preserving default-app wrappers.
- Remove `Reset`; tests create isolated app instances.
- Make `App` implement `http.Handler` or expose a stable handler method so applications can compose it with standard `net/http`.

## Multi-application acceptance case
- One process can create two independent apps with different routes, middleware, sessions, CSP, and plugins.
- A small standard-library host dispatcher can route `example.com` and `app.example.com` to different app handlers.
- Dreego documents this composition without putting domain or virtual-host policy into each app.

## Acceptance criteria
- No mutable runtime configuration or route registry remains package-global.
- Building one app never freezes or changes another app.
- After an app handler is built or listening begins, every configuration mutator either fails deterministically or follows an explicitly documented rebuild model; mutations are never silently ignored.
- Generated code performs no hidden global registration during import.
- Multiple app instances work concurrently under the race detector.
- Plugin registration receives or targets a specific app instance.
- The selected provisional v0.1 plugin registration contract is implemented without preserving the current global fat interface by accident.
- CSRF remains enabled by default, while deliberate route- or method-specific exemptions cannot disable protection for unrelated routes.
- Existing SSR behavior is preserved through end-to-end tests using the new public API.
- Public documentation and new-project blueprints teach only the explicit app model.
- The API-freeze review happens only after this migration is complete.
