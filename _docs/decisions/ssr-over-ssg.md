---
type: Decision
title: Prefer dynamic SSR and explicit caching over static site generation
description: Remove SSG from the planned targets and keep Wails independent from static output
tags: [architecture, ssr, cache, wails]
timestamp: 2026-09-07T00:00:00Z
---
# Prefer dynamic SSR and explicit caching over static site generation

**Date:** 2026-09-07
**Status:** Accepted direction; cache contracts remain provisional

## Context

Dreego previously planned a first-party SSG target between the render
foundation and Wails. That target required a separate route-enumeration model,
build-time data ownership, output paths, atomic publication, canonical URL and
sitemap behavior, static-host capability checks, and a client strategy for
dynamic regions.

Production Dreego applications are expected to remain dynamic. Public output
can be accelerated through response and CDN caches while private regions load
through explicit DreeJS data islands. Wails needs render functions, embedded
assets, navigation, and a typed host bridge; it does not need an SSG pipeline.

## Decision

Dreego does not plan a first-party SSG target. SSR remains the production web
host, with explicit cache policy and invalidation developed after DreeJS proves
the public-page and private-island boundary. Wails builds directly on the
target-neutral renderer and JavaScript processor output.

Cache implementations remain provider plugins until at least two real
implementations prove a shared contract. SSE and WebSocket transports likewise
remain plugins; DreeJS may own a transport-neutral update model only after the
implementations establish common semantics.

The renderer remains callable without an HTTP server. This is required for
tests and Wails and must not be presented as implicit SSG support.

## Consequences

- No `target/ssg`, static route enumeration, or static deployment CLI is
  planned.
- The roadmap moves from version-numbered milestones to capability phases.
- Multi-language client output precedes Wails: Markdown-to-HTML is shipped,
  followed by TypeScript-to-JavaScript and Lua-to-JavaScript.
- Wails does not wait for DreeJS, but benefits from the completed JavaScript
  processor pipeline.
- DreeJS introduces private data islands before shared page caching is
  stabilized.
- CDN, Redis, and in-memory caches remain provider-owned proof points.
- Static marketing sites may use another tool without expanding Dreego's core
  product scope.

## Alternatives considered

### Keep SSG as a later target

Rejected because it retains architecture and compatibility obligations without
a committed product need.

### Model Wails as static generation

Rejected because a desktop WebView has a runtime host, lifecycle, native bridge,
and navigation semantics. Embedded frontend assets do not make Wails an SSG
target.

### Put cache providers in core

Rejected because provider behavior differs and no common contract has yet been
proven by multiple implementations.

## Supersedes

- The SSG target direction in
  [Target-neutral application and first-party targets](target-neutral-application-and-first-party-targets.md).
- The already superseded proposal in
  [SSG and Wails Integration in V2](ssg-wails-v2.md).
