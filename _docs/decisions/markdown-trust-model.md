---
type: Decision
title: Markdown trust model — generation-time vs runtime
description: Two trust domains for markdown rendering — author-trusted generation-time passthrough and safe-by-default runtime rendering
tags: [v0.3.1]
timestamp: 2026-09-05T00:00:00Z
---
# Markdown trust model: generation-time vs runtime

**Date:** 2026-09-05
**Status:** Accepted (v0.3.1)

## Context

The markdown processor was designed for author-trusted `.dreego` files: raw
HTML passes through unchanged at generation time. v0.3 added runtime rendering
for database/CMS content via `dreego.MarkdownToHTML`. That content is not
author-trusted — it can be stored by end users, which introduces a stored XSS
risk if raw HTML and unsafe URLs are emitted verbatim.

## Decision

Two trust domains, chosen per call site:

1. **Generation-time** (`TransformNodes`, `ModeTrusted`): raw HTML passthrough,
   author-trusted, unchanged. This is the existing behavior for `.dreego`
   bodies and is not altered.
2. **Runtime** (`MarkdownToHTML`, `ModeSafe`): raw HTML is escaped, URLs are
   scheme-validated (`http`/`https`/`mailto`/relative; `data:image` raster only
   for `img`), fenced-code language attributes are escaped, and each call uses a
   per-call renderer (no global state, no data race).

`dreego.MarkdownToHTMLTrusted` exists for fully-controlled content and emits a
start-time warning. The trust decision is per call site — `dreego.mdtohtml(...,
trusted: true)` in `.dreego` files — never global, never in config files. The
runtime never reads `dreego.config.json`; global policies would violate explicit
`App` ownership of runtime state.

## Consequences

- Safe mode is the default for all runtime content.
- An XSS test matrix (15 hostile inputs) and a fuzz target for the safe runtime
  path are a release gate.
- Trusted mode is grep-auditable at call sites.

## Supersedes

The v0.3 "runtime rendering" docs section, updated in `_docs/markdown.md`.
