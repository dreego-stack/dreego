---
type: Decision
title: Target-neutral application and first-party targets
description: One typed App with explicit SSR, SSG, and Wails hosts
tags: [architecture, targets, render, v0.x]
timestamp: 2026-08-24T00:00:00Z
---
# Target-neutral application and first-party targets

**Date:** 2026-08-24
**Status:** Accepted direction; public APIs remain provisional until implemented

## Context

Dreego began with SSR coupled to the public `core` package. The product goal now
includes static sites and Wails desktop applications without forcing developers
to maintain a JavaScript build system or run an HTTP server only to render a
WebView.

Treating SSR as the implicit root API while adding SSG and Wails as secondary
features would make the first implementation the permanent architecture.
Conversely, turning rendering, routing, and every host into unrelated plugins
would weaken the coherent component and compiler model.

## Decision

Dreego will evolve toward one target-neutral, typed App and render foundation
with explicit first-party target packages:

```text
github.com/dreego-stack/dreego
github.com/dreego-stack/dreego/target/ssr
github.com/dreego-stack/dreego/target/ssg
github.com/dreego-stack/dreego/target/wails
```

The root package owns application declarations and shared render contracts. The
target packages own their host-specific lifecycle and capabilities. The same
App may be used by more than one compatible target.

SSR, SSG, and Wails remain in the monorepo because they coordinate closely with
compiler output, component metadata, assets, diagnostics, and compatibility.
Provider integrations remain separate plugins.

DreeJS is an optional browser layer shared by targets, not a target. SPA and
Wasm remain future investigations.

## Capability model

Components and processors require explicit small capabilities. Builds fail
when a selected target cannot provide a requirement. Dreego will not introduce
one broad public `Target` interface before real implementations prove shared
methods.

Non-HTTP rendering is not modeled as an HTTP request with nil fields. HTTP-only
features remain owned by the SSR host. Generated page and component inputs stay
typed; a dynamic `map[string]any` is not the primary rendering contract.

## Alternatives considered

### Keep SSR in the root package

Rejected because it makes SSR privileged by history while every later target
uses a visibly different integration model.

### Move every target to external plugins

Rejected for the first-party targets. Their compiler and render coordination is
part of Dreego's core product behavior and compatibility surface.

### Make the core only a source-to-Go transpiler

Rejected because applications would lose one canonical component, render,
asset, and diagnostic contract. Every host plugin could then assign different
semantics to the same `.dreego` source.

### Add a universal Target interface immediately

Rejected as speculative. SSR, build-time static output, and a desktop WebView
have meaningfully different lifecycles. Shared interfaces will be extracted
from working implementations.

## Consequences

- The public `/core` package is expected to be replaced before v1 through an
  explicit migration rather than retained indefinitely through wrappers.
- Rendering must be separated from HTTP before SSG and Wails are implemented.
- Target-specific APIs cannot silently appear in target-neutral component code.
- Documentation must distinguish current SSR behavior from planned target APIs.
- Existing SSR contracts require behavioral tests before internal extraction.
- Plugin capability checks must describe missing requirements at build time.

## Supersedes

- The timeline and SSR-only conclusion in
  [SSG & Wails Integration in V2](ssg-wails-v2.md).
- The post-v1 target restriction in [SSR-First](ssr-first.md).
- Historical target-interface examples in
  [Transpiler Pipeline](transpiler-pipeline.md) where they conflict with the
  capability-first direction.

## Detailed plan

See [`_plan/00-product-architecture.md`](../../_plan/00-product-architecture.md)
and [`_plan/v0.2-render-foundation.md`](../../_plan/v0.2-render-foundation.md).
