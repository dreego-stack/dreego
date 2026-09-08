---
type: Decision
title: Semantic sections and first-party language processors
description: Root sections describe purpose while lang selects a processor
tags: [transpiler, sections, processors, v0.x]
timestamp: 2026-09-01T00:00:00Z
---
# Semantic sections and first-party language processors

**Date:** 2026-09-01
**Status:** Accepted and implemented for the built-in languages; language processors are first-party

> **Current status:** Markdown-to-HTML, TypeScript-to-JavaScript, and the
> dependency-free Browser Lua compiler are implemented. Lua-to-Go is not
> planned.

## Context

The legacy root tags combined purpose and implementation language: `<go>` was
server code, `<div>` was the template root, and root `<script>` was client
JavaScript. This became ambiguous for optional TypeScript, Lua, Markdown, or
other source-language processors.

Arbitrary new root tags such as `<markdown>` or `<lua>` would force the core
parser to infer whether a language produces server code, body nodes, styles, or
client assets. It would also allow processors to conflict with Dreego template
syntax and component calls.

## Decision

Root sections describe semantic purpose:

```html
<server lang="go"></server>
<head lang="html"></head>
<body lang="html"></body>
<style lang="css"></style>
<client lang="js"></client>
```

Default `lang` attributes may be omitted. The breaking migration is:

```text
<go>     -> <server>
<div>    -> <body>
<script> -> <client>
```

An untyped HTML `<script>` nested within `<body lang="html">` remains ordinary
HTML. Inline `<script lang="ts">` and `<script lang="lua">` blocks are explicit
first-party processor inputs and remain at their body position. A root
`<client>` section represents route-level client source.

Each component has one body section. Dreego owns `<@Component>`, control flow,
expressions, slots, escaping, source maps, and diagnostics independently of the
body language. A body processor receives protected placeholders and eligible
literal regions; it cannot consume or redefine Dreego constructs.

## Processor registration

A processor registers an exact section, language, and output kind. Supporting
client Lua does not imply supporting server Lua.

Stable language processors for a small, closed set of source languages are part
of the Dreego monorepo under the relevant transpiler output package:

- Markdown (`md` → HTML IR) under `internal/transpiler/html/md`;
- TypeScript (`ts` → `js`) via Microsoft's pinned native Go compiler for type
  checking and transpilation under `internal/transpiler/js/ts`;
- Lua (`lua` → `js`) through Dreego's own compiler and feature-linked browser
  helpers under `internal/transpiler/js/lua`.

Client Go and Starlark may be evaluated after the Lua patch series, but remain
experimental and outside this stable set. Each experiment must be isolated in
its own processor package, impose no cost when unused, and remain removable
before v1. Its processor name is not reserved until a tested proposal defines
the exact subset and opt-in contract.

Rationale: codegen processors have too much power to run as third-party code,
so their influence must stay reviewable first-party code. The language count is
small (three planned), the VS Code extension can ship the same grammars, and
the processors work out of the box. Raw Go, HTML, CSS, and JavaScript remain
built in.

Runtime plugins (SSE, WebSockets, Tailwind) and provider integrations remain
external plugin repositories using the explicit `Register(app)` model. Managed
tools such as the TypeScript compiler require explicit installation, pinned
versions, and verified release checksums.

## Alternatives considered

### Language-named root tags

Rejected because `<lua>` and `<markdown>` do not say whether their output is
server code, body nodes, styles, or browser code.

### `<$Plugin>` elements

Rejected as a general extension mechanism. UI plugins can expose normal typed
`<@Component>` imports; compiler behavior belongs to a language processor or a
separately specified directive.

### Embedded Lua plugin engine

Rejected because it adds a second implementation ecosystem, VM dependency,
sandbox, debugger, and weakly typed compiler boundary. Lua may still be an
optional client source language.

### Multiple body sections

Rejected because ordering, slots, style scope, control flow, and diagnostics
become ambiguous. Mixed content uses component composition.

## Consequences

- Formatter, parser, generator, scaffolds, fixtures, syntax highlighting, docs,
  and diagnostics use the semantic section model together.
- Legacy root names fail with an actionable migration diagnostic.
- Processor compatibility is checked at generation or build time.
- Markdown, TypeScript, and Lua are the first first-party language processors.
- TypeScript support must perform real type checking, not only remove types.

## Supersedes / Amends

- Replaces the external-processor direction in this document: language
  processors are first-party monorepo code, not external plugin repositories.
- The versioned subprocess protocol for third-party processors is removed from
  the roadmap as speculative.
- Runtime plugins and provider integrations remain external plugin
  repositories.
- [5 Sections in .dreego Files](sections-in-dreego.md).
- The TypeScript timing and in-core esbuild direction in
  [TypeScript Deferred to V2](typescript-v2.md).

## Detailed plan

See [the multi-language Dreego phase](../../_plan/phase-multilanguage-dreego.md).
