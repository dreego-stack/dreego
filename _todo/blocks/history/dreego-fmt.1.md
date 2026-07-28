---
id: dreego-fmt.1
title: dreego fmt (Formatter)
status: 27
phase: v0.0.12
requires:
  - transpiler.1
created: 2026-07-27
changed: 2026-07-28
---

Templ-inspired: `dreego fmt .` formats `.dreego` files in-place. Uses Prettier for `<script>`/`<style>` blocks when available. `--check` flag for CI enforcement. Stdout mode via `--stdout`.
