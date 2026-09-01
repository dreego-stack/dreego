---
type: Concept
title: "Markdown Body Processor"
description: "Write `<body lang=\"md\">` bodies in Markdown, rendered to HTML at build time"
tags: [sections, body, markdown, processors, v0.3]
timestamp: 2026-09-01T00:00:00Z
---
# Markdown Body Processor

Set the body section language to `md` to author the page or component body in
Markdown. Dreego renders the Markdown to validated HTML at generation time,
before the rest of the template is compiled.

```html
<body lang="md">
# Account

Welcome back, **Dreego** developer.
</body>
```

## Syntax

A `<body>` section whose `lang` attribute is `md` contains Markdown instead of
raw HTML. Markdown paragraphs, headings, lists, links, emphasis, blockquotes,
code, and other constructs are converted to their HTML equivalents. The default
`<body lang="html">` remains raw HTML.

## Supported constructs

- ATX headings, lists, blockquotes, and paragraphs;
- emphasis, strong, inline code, and links;
- fenced code blocks.

Tables, raw HTML, and indented code blocks are not supported.

The processor is a small hand-written Markdown parser with no external
dependencies. It lives at `internal/transpiler/html/md` in the transpiler
matrix and requires no external tooling.

## Protected Dreego constructs

Dreego owns component calls, template control flow, expressions, slots,
escaping, and source positions. The body processor receives only the eligible
literal regions plus structured placeholders, so Markdown cannot reinterpret or
consume Dreego syntax:

```html
<body lang="md">
# Account

{#if user.IsAdmin}
<@AdminPanel user={user} />
{/if}
</body>
```

This example mixes Markdown with a `{#if}` condition and an `<@AdminPanel>`
component call; the Markdown processor preserves them verbatim for Dreego to
compile.

## Example

```html
<server>post := c.Data("post")</server>

<body lang="md">
# {{ post.Title }}

Published on **{{ post.Date }}**.

{{ post.Body }}
</body>
```

The Markdown processor emits the surrounding structure; Dreego keeps control of
expressions and escaping for the dynamic values.
