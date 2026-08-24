---
type: Decision
title: Semantic sections and external language processors
description: Root sections describe purpose while lang selects a processor
tags: [transpiler, sections, plugins, v0.x]
timestamp: 2026-08-24T00:00:00Z
---
# Semantic sections and external language processors

**Date:** 2026-08-24
**Status:** Accepted direction; migration not yet implemented

## Context

The current root tags combine purpose and implementation language: `<go>` is
server code, `<div>` is the template root, and `<script>` is client JavaScript.
This becomes ambiguous when optional processors add TypeScript, Lua, Markdown,
or another source language.

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

An HTML `<script>` nested within `<body lang="html">` remains ordinary HTML.
Only a root `<client>` section is client source owned by Dreego.

Each component has one body section. Dreego owns `<@Component>`, control flow,
expressions, slots, escaping, source maps, and diagnostics independently of the
body language. A body processor receives protected placeholders and eligible
literal regions; it cannot consume or redefine Dreego constructs.

## Processor registration

A processor registers an exact section, language, and output kind. Supporting
client Lua does not imply supporting server Lua. Initial external processors
run through a versioned subprocess protocol with structured diagnostics,
assets, source maps, capability requirements, cancellation, and limits.

Raw Go, HTML, CSS, and JavaScript remain built in. TypeScript, Markdown, Lua,
and other dependency-bearing processors live in external plugin repositories.
Managed tools require explicit approval, pinned versions, and reproducible lock
information.

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
  and diagnostics must migrate atomically.
- The current syntax remains documented until the migration todo is complete.
- Processor compatibility is checked at generation or build time.
- Markdown and TypeScript serve as the first two protocol proofs.
- TypeScript support must perform real type checking, not only remove types.

## Supersedes

- [5 Sections in .dreego Files](sections-in-dreego.md).
- The TypeScript timing and in-core esbuild direction in
  [TypeScript Deferred to V2](typescript-v2.md).

## Detailed plan

See [`_plan/v0.3-language-processors.md`](../../_plan/v0.3-language-processors.md).
