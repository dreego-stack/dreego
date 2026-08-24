
---
type: Decision
title: 5 Sections in .dreego Files
description: Dreego files have 5 clearly separated sections for server and client code
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# 5 Sections in .dreego Files

**Date:** 2026-07-28
**Status:** Superseded by
[Semantic sections and external language processors](semantic-sections-and-language-processors.md)

> This document records the historical section model. Its original root names
> were `<go>`, `<div>`, and `<script>`; the implemented semantic model uses
> `<server>`, `<body>`, and `<client>`.

## Context

A `.dreego` file must clearly separate what runs on the server and what runs in the browser. Many frameworks fail at this separation.

## Decision

A `.dreego` file is divided into **5 clearly separated sections**:

1. **`<head>`** — Component-specific meta tags, scripts, CSS links
2. **`<server>`** — Server-side Go code (data fetching, logic)
3. **`<body>`** — The one template root containing HTML and component calls
4. **`<client>`** — Client-side JavaScript (V1: Vanilla JS)
5. **`<style>`** — Scoped CSS (automatically with hashes)

Only these five section tags may appear at the file root. `Component` and
`import` header directives may appear before them. Free text, HTML elements,
and `<@Component>` calls outside `<body>` are generation errors.

Escaped output uses `{{ expression }}`. Control flow keeps its distinct
`{#if}`, `{#each}`, and slot syntax. Typed component props use unquoted Go
expressions such as `<@Card count={count} />`.

## Rationale

1. **No confusion:** It's always clear which code runs where
2. **Full Go power on the server:** `<server>` has DB access, request context, etc.
3. **Real JavaScript for the browser:** `<client>` is sent 1:1 to the client
4. **Component-based assets:** `<head>` loads scripts only when the component is rendered
5. **Scoped CSS:** `<style>` doesn't pollute the global namespace

## The `<head>` Innovation

The `<head>` tag is a core innovation for plugins and performance:

- `dreego-map` declares Mapbox scripts only in its `<head>`
- The Dreego transpiler injects these only when the component is actually rendered
- No global loading of heavy libraries for all pages

## Consequences

- The transpiler must be able to parse and separate the 5 sections
- Each section is processed differently:
  - `<server>` → Go code (server)
  - `<body>` → Go code (HTML generation)
  - `<style>` → Collected, scoped, into CSS file
  - `<client>` → Extracted, embedded in HTML
  - `<head>` → Dynamically injected into the final HTML head
