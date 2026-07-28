---
id: each-else.1
title: {#each else} — Empty List Fallback
status: 25
phase: v0.0.9
requires:
  - transpiler.1
  - each-loop.1
created: 2026-07-28
changed: 2026-07-28
---

{#each items as item}...{#each else}...{/each} — else branch renders on empty list. TokenEachElse, parseEachElseNodes(). Codegen: `if len(items) > 0 { for ... } else { ... }`.
