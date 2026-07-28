---
id: if-else.1
title: {#else} in {#if}-Block
status: 24
phase: v0.0.9
requires:
  - transpiler.1
created: 2026-07-28
changed: 2026-07-28
---

{#if cond}...{#else}...{/if} — else branch in if-block. TokenElse in lexer, parseElseNodes() in parser, ElseChildren in AST. Codegen: if/else Go block.
