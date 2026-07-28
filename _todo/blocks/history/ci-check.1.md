---
id: ci-check.1
title: dreego generate --check (CI Mode)
status: 16
phase: v0.0.6
requires:
  - cli.1
created: 2026-07-27
changed: 2026-07-28
---

Templ-inspired: `dreego generate --check` validates that all generated files are up-to-date. Exit non-zero if stale. Enables CI pipeline: "Fail if generated code is stale". No re-generation, only validation.
