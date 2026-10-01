# Dreego Agent Skill

How to work in a Dreego repository without guessing. Read this before you change
a `.dreego` file or answer a question about how Dreego behaves. Print it again
at any time with `dreego docs skill`.

## What Dreego is

Dreego is a compile-time transpiler: `.dreego` files become Go source and are
hosted on `net/http` (SSR). There is no runtime template parser, no embedded
template engine, and no hidden localhost server. Optional client languages
(TypeScript, Lua) and Markdown compile to JavaScript and HTML at generate time.
Dreego is built on the Go standard library; keep it boring and idiomatic.

## Rule 1: read the docs, not the implementation

Do not open the framework source under `go/pkg`, `$GOMODCACHE`, or
`vendor/github.com/dreego-stack/...` to learn how Dreego works. The bundled
documentation is the contract and it is one command away:

```bash
dreego docs                     # the index
dreego docs /_docs/security.md  # one page
dreego docs --list              # every page, plus plugin docs
dreego docs --dump              # every page in one output (best for agents)
dreego docs --json              # structured JSON: headings, code blocks, links
dreego docs -p <plugin>         # an external plugin's docs
dreego docs skill               # this page
```

If a question is not answered there, say so and propose a docs change instead of
inferring behavior from the implementation. Never tell the user that a feature
exists until you have seen it in `_docs/`.

## Rule 2: a `.dreego` file is sections, not free text

A `.dreego` file has an optional header followed by root sections. The surface
rule is: **uppercase is a keyword, lowercase is a value.**

Header directives (before the first section, at column zero):

- `DREEFILE page|layout|component (props)` declares the file kind.
- `LAYOUT "path/to/layout.dreego"` picks an explicit parent layout.
- `COMPONENT "path" IMPORT { Card, Card as ProductCard }` imports components.
- `GOIMPORT { strings, alias "path" }` declares Go imports for the file.
- `PROFILE "name"` binds the route folder to a registered profile.

The five root sections, `lang` selects the input language:

- `<server>` (Go) — request-time code, handlers, and typed responses.
- `<head>` (HTML) — document metadata, merged with layouts.
- `<body>` (HTML or `lang="md"`) — the page markup or a Markdown body.
- `<style>` (CSS) — route or component CSS.
- `<client>` (JS, TS, or Lua) — browser code.

Header directives are the only content allowed outside these sections.

## Rule 3: let `dreego fmt` own indentation

Write `.dreego` files and run `dreego fmt`. Do not hand-align or fight the
formatter; the canonical shape is:

- Root section tags and header directives at column zero.
- Markup indented by nesting depth, four spaces per level; a closing tag aligns
  with its opening tag; `{#if}`/`{#each}` indent their children.
- Code sections (`server`, `style`, `client`) shifted one level under their tag,
  with the code's own relative indentation preserved.

The formatter never rewrites code semantics: it preserves `<pre>`,
`{#verbatim}`, multi-line string literals, and unindented `lang="md"` bodies.
Run `dreego fmt --check` in CI.

## Rule 4: the pipeline is generate → build

`dreego generate` transpiles `.dreego` files to Go (`dree.go`), one package per
route folder, plus one registrar per app. Generated files carry a marker line
and are never committed. `dreego build` runs `generate` and `go build` into
`build/bin/<name>`. Do not run `go build` directly for a Dreego app; use the CLI.

## Rule 5: security is on by default

Context-aware output escaping is chosen at generate time from the position of
`{{ expression }}`. CSRF, security headers, session cookies with secure
defaults, request limits, recovery, and request IDs are enabled by default. Do
not claim a capability is missing before checking `_docs/security.md` and
`_docs/middleware.md`. `|raw` is the only escaping opt-out and it is explicit.

## Rule 6: report bugs through the CLI

If something is wrong, reproduce it and report it; do not patch the framework
from a consumer project. Open the issue form with:

```bash
dreego feedback
```

## Language rules for this repository

- Chat with the user in German.
- Everything committed to the repository is English: code, comments, docs,
  commit messages, tests, and configuration.

## Keep reading

Start at `dreego docs` (the index). The most useful pages:

- `/_docs/getting-started.md` — quick start
- `/_docs/file-anatomy.md` — section model and header directives
- `/_docs/routing.md`, `/_docs/layouts.md`, `/_docs/components.md`
- `/_docs/server-section.md`, `/_docs/body-html.md`, `/_docs/template-logic.md`
- `/_docs/security.md`, `/_docs/middleware.md`, `/_docs/runtime.md`
- `/_docs/testing.md`, `/_docs/forms.md`, `/_docs/i18n.md`
