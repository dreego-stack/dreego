# Roadmap

Dreego is growing from an SSR-first Go web framework into a Go-native
application platform. The long-term direction is one typed application and
component model with first-party SSR, SSG, and Wails targets plus optional
browser behavior through DreeJS.

This roadmap is directional, not a deadline or release promise. Dreego expects
to remain in v0.x for a long time. Version labels describe dependency order;
implementation evidence may split or reorder a phase.

Current released behavior remains SSR-first. The semantic root sections are
implemented; other planned syntax and APIs below are unavailable until their
phase is implemented and documented.

## Product boundaries

- The root package becomes target-neutral application and rendering API.
- First-party targets live in the monorepo because they share compiler, render,
  asset, diagnostic, and compatibility contracts.
- Provider integrations and optional source languages remain external plugins
  with independent Go modules and releases.
- Applications do not need to operate npm or Node for the standard workflow.
  Approved plugins may manage pinned external tools behind Dreego commands.
- Business state remains authoritative in Go or an explicit backend. DreeJS may
  own local presentation state.
- Unsupported target or plugin capability combinations fail during generation
  or build instead of silently degrading.

The detailed implementation plans live in [`_plan/`](../_plan/README.md).

## Pre-v0.1 — semantic section migration (complete)

Complete the deliberate breaking rename while the template language remains
pre-stable:

```text
<go>     -> <server>
<div>    -> <body>
<script> -> <client>
```

`<head>` and `<style>` retain their names. Purpose and source language become
separate concepts:

```html
<server lang="go"></server>
<body lang="html"></body>
<style lang="css"></style>
<client lang="js"></client>
```

The default `lang` values may be omitted. A `<script>` nested inside an HTML
body remains a normal HTML element. Legacy root names fail with migration
guidance and have no compatibility aliases.

## v0.1 — SSR foundation

Ship and validate the current production foundation:

- explicit App ownership and generated registration;
- typed components, props, routes, forms, and actions;
- sessions, CSRF, output safety, recovery, request limits, server timeouts, and
  graceful shutdown;
- deterministic generation, diagnostics, accessibility gates, reference apps,
  race tests, and release automation;
- explicit external plugin registration with typed options.
- a coherent main history in which every published v0.0.x release tag is an
  ancestor and released code, changelog, and roadmap agree.

v0.1 is SSR-first, not SSR-forever. It intentionally avoids speculative target
interfaces before current behavior is proven.

## v0.2 — target-neutral render foundation

Separate application rendering from the HTTP host:

- move toward `github.com/dreego-stack/dreego` as the target-neutral root API;
- place HTTP ownership in `github.com/dreego-stack/dreego/target/ssr`;
- render typed pages and components without `*http.Request`,
  `http.ResponseWriter`, or an HTTP server;
- model unavailable features through explicit capabilities instead of nil
  request fields;
- preserve typed generated inputs rather than introducing `map[string]any` as
  the primary render API;
- provide an atomic migration away from the public `/core` package.

No universal `Target` interface is promised. Shared contracts are extracted
only from real SSR and non-HTTP implementations.

## v0.3 — language processor boundary

Make root sections language-aware and validate compiler extension points with
separate repositories:

- raw JavaScript, Go, HTML, and CSS remain dependency-free defaults;
- an official Markdown plugin processes `<body lang="md">`;
- an official TypeScript plugin checks and compiles `<client lang="ts">`;
- later Lua support may compile `<client lang="lua">` to JavaScript;
- processors communicate through a versioned subprocess protocol with
  structured diagnostics, source maps, assets, and capability requirements;
- managed tools require approval, pinned versions, lock data, and reproducible
  CI behavior.

Dreego retains ownership of `<@Component>`, template control flow, expressions,
escaping, and source positions even inside a processed body. There is no
embedded Lua plugin VM and no native Go `plugin` loading.

## v0.4 — SSG target

Add `target/ssg` on the render foundation:

- deterministic static HTML and asset output;
- typed enumeration for dynamic route parameters;
- base-path, canonical URL, sitemap, 404, and collision behavior;
- reference deployments for GitHub Pages and Cloudflare Pages;
- build-time capability checks and secret-safe data loading;
- optional DreeJS output for dynamic regions backed by local behavior or an
  external service.

SSG means initial HTML is generated at build time. It does not prohibit client
JavaScript, but a static host cannot provide server routes by itself.

## v0.5 — Wails target

Add `target/wails` without a hidden localhost server:

- render initial documents and assets directly into the WebView host;
- share components with SSR and SSG;
- generate typed Go-to-client bridge contracts;
- support navigation, lifecycle, development reload, and accessibility;
- keep filesystem, shell, clipboard, and other privileged APIs explicitly
  registered and least-privileged.

Application developers should not need a project-owned npm pipeline for the
standard Wails path.

## v0.6 — DreeJS foundation

Add optional modular browser behavior:

- zero runtime for fully static components;
- component-scoped client lifecycle and cleanup;
- typed, safely serialized client props;
- local presentation state and events;
- generated assets containing only used modules;
- SSR, SSG, and Wails compatibility;
- a countdown reference proving that local display state can advance without a
  server round trip while authoritative time remains server-defined.

DreeJS is not a SPA target. It is the shared optional browser layer.

## v0.7 — DreeJS data and live updates

Add network strategies in increasing order of complexity:

1. one-time fetch;
2. bounded, visibility-aware polling;
3. server-sent event streams;
4. multiplexed bidirectional live updates over WebSockets.

Server-rendered updates preserve contextual escaping, ordering, focus, form
state, accessibility announcements, authentication, authorization, and
reconnect behavior. Stateful per-connection views remain optional; the first
implementation should prefer stateless reload-and-render behavior that scales
across ordinary Go instances.

Redis or another provider may supply sessions, pub/sub, caches, locks, jobs, or
presence through external plugins. Dreego does not synchronize arbitrary Go
heap state between instances.

## v0.8+ — stabilization

Before a v1 compatibility promise:

- review every exported and generated contract;
- validate processor and plugin protocols with real external implementations;
- maintain production reference applications for SSR, SSG, Wails, and DreeJS;
- publish security, accessibility, observability, scaling, performance,
  deployment, upgrade, and rollback evidence;
- remove speculative APIs rather than freezing them;
- provide complete migration guidance for v0.x breaking changes.

## v1 — stable promise

v1 means the proven contracts are supportable and documented. It is not a
feature event. Experimental or insufficiently validated capabilities remain
provisional or are excluded from the v1 promise.

## Future exploration: SPA and Wasm

A full SPA runtime and WebAssembly are outside the planned v0.x targets. They
may be explored when real applications demonstrate value not served well by
SSR or SSG plus DreeJS.

Wasm investigations distinguish Cloudflare Workers, sandboxed processor
execution, and browser-side Go. Success in one does not imply the others.
Cloudflare Pages is covered by SSG; a Go backend on Cloudflare Workers would
require a separate Wasm target decision.

## External ecosystem

Official and community plugins may provide auth, billing, Stripe, maps,
MapLibre, SSE, WebSockets, Tailwind, observability, i18n, search, mail, PDF,
storage, cache, jobs, analytics, and other optional capabilities. Each plugin
owns its dependencies and release lifecycle. The Dreego monorepo contains only
the tightly coupled application, compiler, render, target, and DreeJS
foundations.
