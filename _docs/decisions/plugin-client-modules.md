---
type: Decision
title: Declarative Plugin Client Modules
description: Installed plugins expose selectable JavaScript modules as data
tags: [plugins, client, security]
timestamp: 2026-09-09T00:00:00Z
---
# Declarative Plugin Client Modules

**Date:** 2026-09-09
**Status:** Accepted for v0.5

## Context

An authentication plugin needs small browser helpers for passwords, passkeys,
TOTP, and one-time codes. Embedding every helper in the plugin's Go package
bypasses Dreego's client pipeline, while a shell build hook executes code and
cannot omit disabled features deterministically.

## Decision

`dreego-plugin.json` may declare one client bundle with a public URL and named
JavaScript modules. A module can be required or depend on other module IDs. The
application selects optional IDs in `dreego.config.json`; generation resolves
the dependency graph and embeds only the resulting files.

Plugin client files are treated as untrusted data. Paths must remain within the
resolved plugin directory after symlink evaluation, only `.js` files are
accepted, and per-file and bundle size limits apply. Reading client modules
does not require build-hook approval because no command is executed.

The generated bundle is deterministic and starts with a SHA-256 content marker.
Its manifest URL remains stable so SSR layouts can reference it without knowing
the generated hash. Future minification can operate on the assembled bundle
without changing the manifest contract.

## Consequences

- Plugins can ship independent browser features without forcing all of them
  into every application.
- Runtime Go options and build-time client selection remain separate explicit
  configuration and must be kept consistent by the application.
- TypeScript in plugins is compiled before publication until a public compiler
  extension contract is proven.
- `RegisterStatic` remains supported and is used by generated registration.
