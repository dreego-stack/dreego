---
id: golden-tests.1
title: Golden File Tests for Generator
status: planned
phase: v0.0.5
requires:
  - transpiler.1
  - dreegotest.1
created: 2026-07-27
changed: 2026-07-27
---

Templ-inspired: ~50 test subdirectories with `.dreego` input + expected `_dreego.go` output. Golden file pattern: Parse → Generate → Compare. Every template construct ({#if}, {#each}, section, layout) gets own test fixtures. Replaces manual CI checks.
