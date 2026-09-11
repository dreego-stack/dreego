# Roadmap

Dreego is an SSR-first Go application framework growing toward multi-language
`.dreego` components, a first-party Wails host, and optional browser behavior
through DreeJS. Production SSR remains the primary target.

This roadmap uses named capability phases rather than semantic versions.
Release numbers report shipped changes; they do not promise that an entire
phase fits into one release. A phase advances only when tests and real
applications prove its contracts.

The current release plan uses v0.8 for coordinated module boundaries, v0.9 to
begin Wails v3 Phase 1, and v0.10 to begin DreeJS. These are planned slices rather
than compatibility promises; evidence may move unfinished work to a later
release.

Static site generation is not planned. Dreego prioritizes dynamic SSR with
explicit caching and invalidation. See the
[SSR over static site generation decision](decisions/ssr-over-ssg.md).

## Product boundaries

- The root package owns target-neutral application and rendering contracts.
- SSR and Wails are first-party hosts because they share compiler, render,
  asset, diagnostic, and compatibility contracts.
- Provider integrations, caches, SSE, and WebSockets remain external plugins.
- The built-in processor matrix normalizes multiple input languages into HTML,
  JavaScript, or Go output without adding a universal processor interface.
- Applications do not operate npm or Node. Managed TypeScript tooling uses a
  pinned native compiler and requires explicit installation.
- Business state remains authoritative in Go or an explicit backend. DreeJS may
  own local presentation state and explicitly load private data islands.
- Unsupported host or processor combinations fail during generation or build.

Detailed plans live in [`_plan/`](../_plan/README.md).

## Completed foundation: SSR

The released foundation provides:

- explicit `App` ownership and generated registration;
- typed components, props, routes, forms, and actions;
- sessions, CSRF, output safety, recovery, request limits, timeouts, and
  graceful shutdown;
- deterministic generation, diagnostics, accessibility gates, race tests, and
  release automation;
- a target-neutral render boundary with an explicit SSR host.

SSR is the production baseline and remains the primary deployment model.

## Completed phase: multi-language Dreego

Root sections express purpose while `lang` selects the input language:

```html
<server lang="go"></server>
<head lang="html"></head>
<body lang="html"></body>
<style lang="css"></style>
<client lang="js"></client>
```

The processor matrix names normalized output first and input second:

```text
html/html   HTML input       -> HTML IR
html/md     Markdown input   -> HTML IR
js/js       JavaScript input -> JavaScript output
js/ts       TypeScript input -> JavaScript output
js/lua      Lua input        -> JavaScript output
html/output HTML IR          -> generated Go renderer
```

Stable client languages and their version targets are tracked in the
[client-language compatibility matrix](../cmd/dreego/_docs/client-languages.md).

Markdown-to-HTML, TypeScript-to-JavaScript, and Browser Lua are shipped. Browser
Lua includes isolated scripts, safe emission, source-accurate diagnostics,
local functions, loops, deterministic tables, precise string escapes,
feature-linked helpers, and source-aware runtime stacks. Its supported subset
is covered through black-box generation tests and may grow through compatible,
tested patch releases.

Lua-to-Go is not planned. Dreego does not embed a Lua plugin VM or load native
Go plugins. Raw JavaScript, Go, HTML, and CSS remain dependency-free defaults.
Browser Lua also has no external compiler dependency.

JavaScript, TypeScript, and Browser Lua are the complete client-language set
for the foreseeable future. Additional client languages are outside the
roadmap so the existing processors can be hardened instead of widened.

## Current phase: coordinated module boundaries — planned for v0.8

Split Core, SSR, testing, and the CLI into independently consumable Go modules
with one coordinated version and module-specific tags. Rename the HTTP host to
`adapter/ssr` and the CLI source to `cmd/dreego`. Wails is intentionally not
part of this release.

## Phase: Wails v3 Phase 1 — planned for v0.9

Add an experimental, opt-in first-party host against one explicitly pinned
Wails v3 beta release:

- render initial documents without a listening TCP socket;
- embed HTML, styles, scripts, and static assets in the application;
- navigate between Dreego routes through explicit desktop semantics;
- use Wails-generated typed Go-to-client bindings;
- integrate window startup, shutdown, development reload, and accessibility;
- register filesystem, shell, clipboard, and other privileged APIs explicitly
  and with least privilege.

Wails depends on target-neutral rendering and the JavaScript processor output.
It does not depend on static site generation. The first Wails slice may use raw
JavaScript or TypeScript without waiting for DreeJS.

Phase 1 is the smallest credible desktop path: rendering, embedded assets,
literal navigation, one typed binding, lifecycle cleanup, development reload,
and an accessible reference application. It proves that the integration works;
it does not promise broad coverage of Wails APIs while Wails v3 remains beta.

## Phase: DreeJS foundation — planned for v0.10

Establish the optional component lifecycle, deterministic module emission,
local presentation state, accessible cleanup, and typed serialized props. A
component without client behavior continues to emit zero DreeJS bytes. DreeJS
is the umbrella name for modular browser and WebView capabilities, not a target
and not one release-sized feature.

## Phase: DreeJS extensions — planned after v0.10

Add a small modular browser and WebView layer:

- zero runtime for components without browser behavior;
- component-scoped lifecycle and cleanup;
- typed, safely serialized client props;
- local presentation state and events;
- one-time fetch with loading, empty, error, and cancellation behavior;
- private data islands inside otherwise publicly cacheable SSR pages;
- generated assets containing only used modules;
- compatible behavior under SSR and Wails.

A public SSR page may contain a neutral account placeholder while a DreeJS
island loads the authenticated user. Private responses remain `private` and
`no-store`; personal data never enters the shared page cache.

## Deferred phase: Wails v3 Phase 2

Broader desktop integration begins only after Wails v3 has a stable upstream
release and the Phase 1 reference application has produced practical
compatibility evidence. Phase 2 may expand bridge integration, packaging,
performance, platform behavior, and explicitly granted host capabilities. It
has no assigned Dreego release number.

Caching is developed alongside those boundaries but remains provider-owned
until multiple implementations prove a shared contract. Initial proofs should
cover an in-memory or Redis-backed cache and a CDN integration such as
Cloudflare. Candidate shared semantics include cache policy, tags, explicit
invalidation, conditional requests, and safe session-aware bypasses.

## Phase: live updates

Extend the proven DreeJS lifecycle in increasing order of complexity:

1. bounded, visibility-aware polling;
2. server-sent events through an external SSE plugin;
3. bidirectional updates through an external WebSocket plugin;
4. Wails bridge updates after the Wails v3 Phase 2 gate is met.

DreeJS owns the transport-neutral client update model: ordering, sequence
numbers, resynchronization, focus preservation, accessible announcements, and
component cleanup. SSE and WebSocket connection ownership, authentication,
backpressure, limits, and horizontal pub/sub stay in their plugins. A shared
core contract is extracted only after real transports prove the same need.

## Phase: stabilization

Before a v1 compatibility promise:

- review every exported and generated contract;
- validate processor, cache, and transport protocols with real implementations;
- operate production reference applications for SSR, Wails, and DreeJS;
- publish security, accessibility, observability, scaling, performance,
  deployment, upgrade, and rollback evidence;
- remove speculative APIs rather than freezing them;
- provide migration guidance for supported upgrade paths.

## Stable promise

v1 means that the proven contracts are supportable and documented. It is not a
feature event. Experimental or insufficiently validated capabilities remain
provisional or are excluded from the v1 promise.

## Future exploration: SPA and Wasm

A full SPA runtime and WebAssembly remain independent future investigations.
DreeJS is not a SPA target, and Cloudflare Workers would require its own Wasm
host decision. Neither is required for multi-language components, Wails,
cache-aware SSR, or live updates.

## External ecosystem

Official and community plugins may provide auth, billing, maps, SSE,
WebSockets, Tailwind, observability, i18n, search, mail, PDF, storage, cache,
jobs, analytics, and other optional capabilities. Each plugin owns its
dependencies and release lifecycle. The Dreego monorepo contains the tightly
coupled application, compiler, renderer, SSR host, Wails host, and DreeJS
foundations.
