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
- fenced code blocks;
- pipe tables (GFM-style) with per-column alignment;
- nested lists (indented sub-items);
- images (`![alt](url)` → `<img>`);
- footnotes (references + definitions emitted at the end);
- horizontal rules (`---`);
- raw HTML blocks (see the trust model below).

Indented code blocks are not supported. The `***` and `___` horizontal-rule
forms and setext headings (underlined with `===` or `---`) are not supported
either — use ATX headings and `---` for a horizontal rule.

Raw HTML blocks pass through verbatim. This is a generation-time, trust-based
feature: the Markdown source is developer-authored and compiled at build time,
so raw HTML is trusted exactly like the class passthrough on the inline `<md>`
tag. It is not a runtime sanitization boundary.

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

## Inline `<md>` tag

Inside a regular `<body lang="html">` you can drop a Markdown region anywhere
with the inline `<md>` tag. Dreego processes it with the same Markdown processor
at generation time and wraps the result in a `<div>`:

```html
<body>
<h1>HTML title</h1>
<md class="prose">
# Markdown inside

- a
- b
</md>
<p>After</p>
</body>
```

The `<md>` tag always becomes a `<div>` wrapper — there is no "transparent if no
attributes" magic. The `class` attribute moves onto the `<div>`; every other
attribute passes through verbatim:

```html
<md class="prose" data-x="1"># Hi</md>
<!-- compiles to -->
<div class="prose" data-x="1"><h1>Hi</h1></div>
```

Expressions inside the region are preserved and escaped by Dreego, exactly as in
a `lang="md"` body.

### Error cases

- **Nested `<md>`** inside an `<md>` region is rejected.
- **Control flow** (`{#if}` / `{#each}`) inside an `<md>` region is rejected —
  close the `<md>` tag first.
- **Unclosed `<md>`** (end of body or control flow before `</md>`) is rejected.
- **`<md>` inside `<body lang="md">`** is rejected — Markdown is already the
  body language there.

## Runtime rendering

The same Markdown processor is available at runtime for content stored as
Markdown (for example from a database). Two functions are exported:

- `dreego.MarkdownToHTML(src string) (string, error)` — **safe mode**. Raw HTML
  in the source is escaped, so it is safe for untrusted content (for example
  user-generated Markdown). This is the default and the recommended entry point.
- `dreego.MarkdownToHTMLTrusted(src string) (string, error)` — **trusted mode**.
  Raw HTML passes through verbatim. Use it only for content you fully control;
  it is not a sanitization boundary. On first use it logs a loud warning to
  stderr: `MarkdownToHTMLTrusted: raw HTML passthrough enabled — only use with
  trusted content`.

Both render the Markdown to HTML per call. Control flow and other Dreego
constructs are not interpreted at runtime — they are rendered as literal text,
so use them only for plain Markdown content.

### `dreego.mdtohtml` stdlib syntax

Inside a `<server>` section you can use the `dreego.mdtohtml(...)` stdlib syntax
instead of calling the exported functions directly. The transpiler rewrites it
to the corresponding core call at generation time, so the trust decision is
visible at the call site, per use:

```html
<server>
html, err := dreego.mdtohtml(post.Content)
</server>
```

- `dreego.mdtohtml(x)` → `dreego.MarkdownToHTML(x)` — **safe mode**, the default.
- `dreego.mdtohtml(x, trusted: true)` → `dreego.MarkdownToHTMLTrusted(x)` —
  **trusted mode**; the `trusted: true` argument is stripped and the call is
  rewritten.
- `dreego.mdtohtml(x, trusted: false)` → `dreego.MarkdownToHTML(x)` — safe mode;
  the `trusted: false` argument is stripped.

The `trusted:` argument is a per-call-site decision, not a global setting. If
the call cannot be parsed (for example unbalanced parentheses), the code is left
unchanged and the Go compiler reports the error naturally.

## Roadmap ideas

The following are planned but not supported yet:

- MDX-style imports (rejected for now)
- Syntax highlighting for fenced code blocks
- Automatic table of contents
