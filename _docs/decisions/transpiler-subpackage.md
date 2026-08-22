---
type: Decision
title: Transpiler as internal subpackage
description: The transpiler lives in internal/transpiler/; core/ is the runtime framework only
tags: [architecture]
timestamp: 2026-08-20T00:00:00Z
---
# Transpiler as internal subpackage

**Date:** 2026-08-20
**Status:** Accepted

## Context

The transpiler (lexer, parser, codegen, generator) historically lived in
`core/`, the package user applications import as the Dreego framework. The
runtime never used the transpiler — the only consumers were the CLI
(`dreego generate`, `dreego fmt`) and the `dreegotest` test helpers. Keeping
both in one package exposed internal pipeline symbols (`Lex`, `Parse`,
`GenerateMethodHandler`, …) as public API without a user.

## Decision

Move the transpiler to `internal/transpiler/` at the repo root. `core/` keeps
the public runtime API (App, SSRContext, middleware, sessions, Safe* helpers,
validation) at the same import path; the CLI and `dreegotest` import
`internal/transpiler` directly. Go's `internal/` rule makes the transpiler
importable only from inside this module.

`internal/` at the root was chosen over `cli/dreego/internal/` because
`dreegotest` is a second consumer; a future watch/build tooling target would
be a third. `internal/` is for shared, unexported implementation; middleware
stays in `core` because users call `app.Use(dreego.Compress())` — it is
public contract, not an implementation detail.

## Consequences

- User applications and generated code keep importing only
  `github.com/dreego-stack/dreego/core`; no API change for applications.
- The transpiler's public surface (`Lex`, `Parse`, AST types, generator) is no
  longer part of the framework API.
- 80+ transpiler source and test files moved; mixed test files
  (`benchmark_test.go`, `error_propagation_test.go`, output-safety tests) were
  split along the runtime/transpiler boundary.
- CI (`_tests/test.sh`, `pull-request-check.yml`) now runs
  `./internal/transpiler/...`; `_scripts/check-core-deps.sh` verifies both
  packages have no external dependencies.
