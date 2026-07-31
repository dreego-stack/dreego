
---
type: Decision
title: Name "dreego" / File Extension ".dreego"
description: Naming convention and package convention
tags: [naming, v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---

# Name "dreego" / File Extension ".dreego"

**Date:** 2026-07-28
**Status:** Final (plugin paths superseded by [monorepo-plugin-layout](monorepo-plugin-layout.md) v0.0.21)

## Decision (final)

**Project name:** dreego
**File extension:** `.dreego`
**CLI tool:** `dreego`
**Go Package:** `dreego`

## Package Convention

| Type             | Path                                      |
|------------------|-------------------------------------------|
| Core             | `codeberg.org/dreego/dreego`             |
| Auth Plugin      | `codeberg.org/dreego/dreego/plugins/auth`  |
| Map Plugin       | `codeberg.org/dreego/dreego/plugins/map`   |
| UI Plugin        | `codeberg.org/dreego/dreego/plugins/ui`    |

> Plugin paths superseded by [monorepo-plugin-layout](monorepo-plugin-layout.md) (v0.0.21): official plugins live under `plugins/<name>` in this repository. Community plugins may still use separate repos.

## Consequences

- File extension `.dreego`
- CLI: `dreego` (dreego generate, dreego dev, dreego build)
- Package name: `dreego`
- Website: dreego.dev
- Codeberg Org: codeberg.org/dreego
- GitHub Mirror: github.com/LukasLow/dreego
