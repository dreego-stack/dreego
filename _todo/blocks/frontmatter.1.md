---
id: frontmatter.1
title: Frontmatter Support in .dreego
status: planned
phase: pre-v1.0
requires:
  - transpiler.1
  - context-refactoring.1
created: 2026-07-31
changed: 2026-07-31
---

Parse YAML/TOML frontmatter at the top of `.dreego` files and expose typed metadata via context. Required for documentation sites, blogs, content collections, and SSG. Generates a typed struct or map access through `c.Data("key")`.
