# Compatibility Policy

This document defines when and how Dreego may change its public API. It is the
contract that governs the transition from the pre-v0.1 phase into the v0.1 SSR
baseline and eventually the v1 stability promise.

## Scope

"Public API" means every exported identifier in `github.com/dreego-stack/dreego/core`
that an application or plugin can import and use: types, functions, methods,
interfaces, constants, and variables. It also covers the `.dreego` template
language, the generated `www.Register(app)` contract, and the `dreego.config.json`
schema.

## Before v0.1

Until v0.1 is released, the public API is **not** stable. Dreego may change,
rename, or remove exported identifiers without a deprecation cycle when the
change is required to harden the SSR core. This is the phase in which the API
is reviewed against real applications and plugins and prematurely frozen
interfaces are revised.

The planned semantic-section rename from `<go>`, `<div>`, and root `<script>`
to `<server>`, `<body>`, and `<client>` must complete before v0.1. It is an
intentional breaking change and does not receive compatibility aliases.

Two constraints apply even before v0.1:

1. **No speculative exports.** An exported identifier must be exercised by a
   real application or plugin, or it is removed rather than kept "for later".
   Dead exports are a form of premature freezing.
2. **Documented changes.** Every breaking change lands through a pull request
   with a `.changes/*.md` changelog line that states what changed and why. The
   changelog is the migration record.

## The v0.1 SSR promise

At v0.1, the following contracts form the supported SSR baseline. Dreego
commits to changing them only through the documented pre-v1 process below. The
v0.1 label does not freeze the current `/core` import path until v1: the planned
v0.2 target-neutral root-package migration is an explicit pre-v1 breaking
change with its own migration guide.

- `App` and its configuration, registration, and lifecycle methods
  (`New`, `Register`, `RegisterRedirect`, `RegisterRewrite`, `RegisterStatic`,
  `Use`, `SetLogging`, `SetCSRF`, `SetErrorHandler`, `SetSessionStore`,
  `SetCSP`, `SetServerConfig`, `SetReady`, `Build`, `Handler`, `ServeHTTP`,
  `Listen`, `Shutdown`).
- `SSRContext` and the `Context` interface used by generated form handlers.
- The session `Store` interface and `CookieStore`/`CookiePolicy`/`Options`.
- `ServerConfig` and `DefaultServerConfig` (the type set by `SetServerConfig`).
- The exported error sentinels used by applications: `ErrAppBuilt`,
  `ErrRouteConflict`, `ErrSessionTooLarge`, `ErrRedirect`.
- The output-safety helpers (`SafeText`, `SafeAttr`, `SafeURL`, `SafeScript`,
  `SafeStyle`, `SafeRefresh`, `SafeRaw`).
- The middleware constructors (`RequestLogging`, `Compress`, `Recovery`,
  `RequestID`, `CSRF`, `MaxBodyReader`).
- The post-migration `.dreego` template language and generated
  `www.Register(app)` contract.
- The `dreego.config.json` schema.

## After v0.1, before v1

Between v0.1 and v1, the stable contracts above may still change, but only
through a documented process:

1. The change is announced in the changelog with a clear migration note.
2. Where practical, the old identifier is deprecated before removal. A
   deliberate architecture migration may remove an obsolete package atomically
   when wrappers would preserve conflicting ownership or semantics.
3. Breaking changes are batched into minor releases, never silently shipped in
   a patch.

## The plugin contract stays provisional until v1

The plugin contract is **explicitly excluded** from the v0.1 stability promise.
There is no central `Plugin` interface before v1. Plugins are ordinary Go
packages that expose an explicit `Register(app, options) error` function and
use `app.Register` and `app.Use` directly.

The plugin contract remains provisional until v1 because it can only be
validated by real external plugins. Between v0.1 and v1, Dreego builds plugins
against the public interfaces and adjusts the pre-v1 contracts when evidence
requires it. A shared capability interface is added to core only after at
least two real implementations prove the same small contract. The final
compatibility promise for the plugin contract begins at v1, not before.

## What counts as a breaking change

A change is breaking if it would cause an existing application or plugin to
stop compiling, or to change behavior in a way the author did not opt into.
Examples:

- Removing or renaming an exported identifier.
- Changing a function or method signature.
- Adding a method to a public interface (breaks external implementations).
- Changing the meaning of a documented default or configuration field.
- Changing generated code in a way that requires regeneration with a different
  result.

## Review cadence

The exported API is reviewed before v0.1 and again before v1. Each review
checks every exported identifier against the applications and plugins that
exercise it, removes dead exports, and revises prematurely frozen interfaces.
The result of each review is recorded in the changelog and, where it changes
the policy, in this document.

## See Also

- [Roadmap](https://github.com/dreego-stack/dreego/blob/main/_docs/roadmap.md)
- [Plugins](https://github.com/dreego-stack/dreego/blob/main/_docs/plugins.md)
- [Plugin Interfaces](https://github.com/dreego-stack/dreego/blob/main/_docs/plugin-interfaces.md)
