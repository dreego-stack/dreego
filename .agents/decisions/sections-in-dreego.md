
---
type: Decision
title: 5 Sections in .dreego Files
description: Dreego files have 5 clearly separated sections for server and client code
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# 5 Sections in .dreego Files

**Date:** 2026-07-28
**Status:** Accepted

## Context

A `.dreego` file must clearly separate what runs on the server and what runs in the browser. Many frameworks fail at this separation.

## Decision

A `.dreego` file is divided into **5 clearly separated sections**:

1. **`<head>`** — Component-specific meta tags, scripts, CSS links
2. **`<go>`** — Server-side Go code (data fetching, logic)
3. **Template (HTML)** — The markup with Dreego template syntax
4. **`<script>`** — Client-side JavaScript (V1: Vanilla JS)
5. **`<style>`** — Scoped CSS (automatically with hashes)

## Rationale

1. **No confusion:** It's always clear which code runs where
2. **Full Go power on the server:** `<go>` has DB access, request context, etc.
3. **Real JavaScript for the browser:** `<script>` is sent 1:1 to the client
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
  - `<go>` → Go code (server)
  - Template → Go code (HTML generation)
  - `<style>` → Collected, scoped, into CSS file
  - `<script>` → Extracted, embedded in HTML
  - `<head>` → Dynamically injected into the final HTML head
