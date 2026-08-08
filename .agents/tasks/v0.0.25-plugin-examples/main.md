---
type: Task
title: v0.0.25 Plugin Examples — SSE (in-repo) + external plugin-example
status: in_progress
assign: manager
---

# v0.0.25 Plugin Examples

Goal: Ship two working plugin examples that exercise the frozen v1 Plugin
interface (core/plugin.go), one in-repo and one external, so the plugin
ecosystem and the upcoming module-docs discovery have real references.

## Deliverables

1. **In-repo plugin: `plugins/sse/`** — Server-Sent Events plugin.
   - Own module `codeberg.org/dreego/dreego/plugins/sse` with its own go.mod
     (requires core, replace to ../core for local dev).
   - Implements the full Plugin interface: Name, RegisterRoutes, Middlewares,
     Assets, OnStart, OnShutdown.
   - SSE functionality: a route that streams events (text/event-stream),
     e.g. `/sse` with a simple event loop; a helper to broadcast messages.
   - `_docs/` folder with `dreego-sitemap.xml` (manifest listing doc files)
     and at least one markdown doc — this is the seed for the future
     module-docs discovery.
   - Tests: interface satisfaction (`var _ dreego.Plugin = ...`), route
     registration via UsePlugin, SSE stream works (httptest).

2. **External plugin: `/Users/lukas/home/proj/dreego/plugin-example`** —
   a separate directory OUTSIDE the dreego repo (sibling folder), acting as
   an external/community plugin example.
   - Own go.mod (module e.g. `example.com/plugin-example` or
     `codeberg.org/dreego/plugin-example`), requires core with a replace to
     the local core path (../dreego/core) for local dev.
   - Implements the full Plugin interface (e.g. a "hello" plugin: registers
     `/hello` route, a middleware that sets a header, assets with a small
     file, lifecycle logs).
   - `_docs/` with `dreego-sitemap.xml` + one markdown doc.
   - Tests: interface satisfaction + route works via UsePlugin.

## Rules

- Core must never import a plugin; plugins import core.
- Plugins with external deps get their own go.mod; these two use stdlib only.
- Max 300 lines per file, English, no comments unless needed.
- The external plugin dir is NOT part of the dreego git repo (it lives in
  /Users/lukas/home/proj/dreego/plugin-example) — do not commit it into the
  dreego repo.

## Execution

1. coder creates both plugins + tests.
2. Verify: `smd go test ./plugins/sse/...` (in-repo), and for the external
   plugin run its tests from its own dir (go test ./... with replace).
3. reviewer reviews.
4. git commits only the in-repo plugin (plugins/sse/). The external plugin
   stays outside the repo.

## Status

- [x] Task file created
- [x] plugins/sse implemented + tested
- [ ] external plugin-example implemented + tested
- [ ] review + commit (in-repo only)
