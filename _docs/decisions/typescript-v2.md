
---
type: Decision
title: TypeScript Deferred to V2
description: TypeScript support in the script block is deferred to V2, V1 uses Vanilla JS
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# TypeScript Deferred to V2

**Date:** 2026-07-28
**Status:** Superseded by
[Semantic sections and external language processors](semantic-sections-and-language-processors.md)

> TypeScript is now planned as an external `<client lang="ts">` processor in
> the v0.x line. Raw JavaScript remains built in. The processor must run real
> type checking and manage pinned tools without adding dependencies to core.

## Context

In `.dreego` files there is a `<client>` block for client-side code. The question was whether this should be TypeScript or Vanilla JavaScript.

## Problem

Building TypeScript into V1 would have caused enormous complexity:

1. Requires JS/TS toolchain (esbuild, Node, Bun, or Deno) on the developer machine
2. Destroys the Go promise: "No external dependencies, only install Go"
3. Massively increases build complexity
4. Scope creep — significantly delays V1

## Lesson from the Past

Elixir's attempt to retroactively introduce types shows how hard that is: Typeifying a dynamic language after 10 years is a massive, painful process. That is not a success story — it is a warning.

**Important for Dreego:** TypeScript in the `<client>` block is uncritical to defer to V2 because:
1. Go itself is already statically typed (server side typed from day 1)
2. The `<client>` block is isolated — it does not influence the core architecture
3. Vanilla JS → TypeScript is an upgrade, not a fundamental redesign

But for OTHER decisions: What we now set in the architecture should be designed so that future extensions are possible without breaking changes (plugin interface, transpiler pipeline, plugin system).

## Decision

**V1: Pure Vanilla JavaScript** in the `<client>` block. No compiler, no bundler, 0 MB extra tooling.

**V2: TypeScript** via esbuild integration (esbuild can be embedded as a Go library).

## Rationale

1. Modern Vanilla JS can already do everything needed: `import/export`, `async/await`, Shadow DOM
2. Dreego only needs to extract JS 1:1 in V1 and embed it in HTML
3. 0 complexity, maximum render speed
4. Focus on the core: transpiler, routing, Go server

## Provision for a future processor

So TypeScript can be added later without breaking changes:

- The transpiler has a clearly defined pipeline: Parse → AST → CodeGen.
- `<client lang="ts">` currently fails with a processor requirement instead of
  silently emitting unchecked TypeScript.
- Types sharing (Go struct → TS interface) is conceptually pre-sketched

## Consequences

- `<client>` block expects pure JavaScript (no `lang="ts"` in V1)
- `lang="ts"` is rejected until an approved processor is installed
- No esbuild dependency in V1
- A future external processor owns type checking, transpilation, and pinned tools
- Types sharing (Go struct → TS interface) also comes only in V2
