---
id: frontmatter.1
title: Frontmatter Support in .dreego
status: 48
phase: v0.0.25
requires:
  - transpiler.1
  - context-refactoring.1
created: 2026-07-31
changed: 2026-08-08
---

Parse YAML/TOML frontmatter at the top of `.dreego` files and expose typed metadata via context. Required for documentation sites, blogs, content collections, and SSG. Generates a typed struct or map access through `c.Data("key")`.

Done in v0.0.25: `core.ParseFrontmatter(src)` splits leading `---` block, exposes key:value pairs.
