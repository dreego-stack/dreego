
---
type: Decision
title: Name "dreego" / File Extension ".dreego"
description: Naming convention and package convention
tags: [naming, v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---

# Name "dreego" / File Extension ".dreego"

**Date:** 2026-07-28
**Status:** Final for naming; plugin paths now use separate repositories and modules

## Decision (final)

**Project name:** dreego
**File extension:** `.dreego`
**CLI tool:** `dreego`
**Go Package:** `dreego`

## Package Convention

| Type             | Path                                      |
|------------------|-------------------------------------------|
| Core             | `github.com/dreego-stack/dreego`             |
| Auth Plugin      | `github.com/dreego-stack/plugin-auth`          |
| Map Plugin       | `github.com/dreego-stack/plugin-map`           |
| UI Plugin        | `github.com/dreego-stack/dreego/plugins/ui`    |

> Plugin paths superseded (v0.0.21): official plugins initially lived under
> `plugins/<name>` in this repository, then moved to separate repositories
> under `github.com/dreego-stack/`. See AGENTS.md "Core and Plugin Boundary".

## Consequences

- File extension `.dreego`
- CLI: `dreego` (dreego generate, dreego dev, dreego build)
- Package name: `dreego`
- Website: dreego.dev
- Codeberg Org: codeberg.org/dreego
- GitHub Mirror: github.com/LukasLow/dreego
