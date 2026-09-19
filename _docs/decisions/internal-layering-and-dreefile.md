---
type: Decision
title: Internal layering and the dreefile package
description: Root internal/ is the shared implementation layer; internal/dreefile/ holds the compiler with ir, dreecode, gogen, and sections
tags: [architecture, compiler, internal, dreefile]
timestamp: 2026-09-18T00:00:00Z
---
# Internal layering and the dreefile package

**Date:** 2026-09-18
**Status:** Accepted

## Context

Two structural problems accumulated as the project grew:

1. The compiler lives in `internal/transpiler/`, a flat package with lexer,
   parser, per-section processors, per-language processors, code generation,
   discovery, formatting, configuration, accessibility checks, and i18n side by
   side. Related concerns (for example the duplicated expression-kind and
   source-position helpers in the root package and in `ir/`) drifted apart.
   A generic `transpiler` name also describes a mechanism rather than the
   artifact: the compiler consumes `.dreego` sources and is the `dreefile`
   compiler.
2. The target-neutral runtime implementation sits under `core/internal/`
   (`core/internal/render|context|i18n|server`) while shared implementation
   (`internal/session`, `internal/validate`, `internal/middleware`,
   `internal/gomod`) already lives at the repository root. The same conceptual
   layer is therefore split across two places, and `core/` carries
   implementation next to its public facade.

The root module and `core` are separate Go modules. Go's internal-visibility
rule is path-prefix based: `github.com/dreego-stack/dreego/internal/...` is
importable from any package whose path begins with
`github.com/dreego-stack/dreego`, which includes the `core` module. This already
works today for `internal/session`, `internal/validate`, and
`internal/middleware`, so moving more shared implementation to the root
introduces no new mechanism.

## Decision

### Root `internal/` is the shared implementation layer

Root `internal/` is the general implementation layer, roughly the role a
project-local standard library plays. `core/`, `adapter/ssr/`,
`adapter/wails/`, `cmd/dreego/`, and `dreegotest/` stay slim
facades, hosts, and binaries over it.

- `core/internal/render|context|i18n|server` move to `internal/render`,
  `internal/context`, `internal/i18n`, `internal/server`.
- `core/` keeps only the public facade and re-export layer.
- `cmd/dreego/internal/templates` moves to `internal/templates`.

No import cycle may be introduced: root `internal/` must never import `core` as
a Go package. Code generation may reference the public import path only as an
emitted string.

### The compiler becomes `internal/dreefile/`

`internal/transpiler/` becomes `internal/dreefile/`. Central orchestration stays
in `dreefile/generate.go` (`Run`, `RunCheck`, `buildPlan`, `buildRootPlan`,
`buildRootFile`) plus `generate_check.go` (plan diff, disk I/O, apply) and
`generate_support.go`.

Shared packages sit directly on the dreefile level:

- `dreefile/ir/` — AST and types, the single source of truth.
- `dreefile/dreecode/` — the mini template language semantics: `{#if}`,
  `{#each}`, `{{ }}`, filter pipes, `{#slot}`, and expression scanning.
  `dreecode` is **not** the lexer/parser/token stack; see the front-end stage
  below.
- `dreefile/gogen/` — Go emission helpers: `GoLiteral`, `ToPascalCase`,
  escaping, and source positions.
- `dreefile/codegen/` — the shared generation state (`codegen.State`,
  `codegen.Layout`) threaded through the sections, discovery, and generation.
  It imports `ir`. It is **not** merged into `gogen`: `gogen` is stateless Go
  emission, while `codegen.State` is the mutable build state (component
  definitions, imports, message uses, Lua features) that nearly every section
  needs.
- `dreefile/jsoutput/` — the shared client emission artifact and `<script>`
  tag generation (former `js/output`: `Artifact`, `GenClient`, `GenClientTo`).
  It imports `gogen`, not `ir`. It is shared because the body/html section and
  every client language emit the same script form.

### Front-end stage

The tokenizing and parsing stack keeps explicit homes at the dreefile level,
parallel to the shared packages:

- `dreefile/tokens/` — token types.
- `dreefile/lexer/` — the scanner; imports `tokens`, `ir`, and `dreecode`.
- `dreefile/parser/` — the section parser; imports `tokens`, `ir`, and
  `dreecode`. It does **not** import `lexer` in non-test code.

The front-end stage produces the `ir` tree; sections consume it. `dreecode` is
the mini-template-language semantics that BOTH the front-end (`lexer`, `parser`)
and the sections use for control flow, expressions, filters, and slots — not the
token/lexer/parser stack.

### Sections

`dreefile/sections/` owns the section processors:

```text
dreefile/sections/
├── head/
├── style/
├── body/html/
├── body/md/      Markdown processor home (see below)
├── client/       explicit client orchestrator (former js/process)
│   ├── js/
│   ├── ts/
│   └── lua/
```

The `header/` and `server/` section processors did **not** become own packages
in this phase; they remain in the dreefile root. This is a deliberate,
recorded deviation, not an oversight — see Consequences and the open item in
[`_plan/phase-restructure.md`](../../_plan/phase-restructure.md).

Language subfolders exist **only** where real language variants exist: `body`
and `client`. A section that has one language stays flat.

The `client/` orchestrator (former `js/process`) is the single legal entry point
into the client languages. It imports its children `client/js`, `client/ts`, and
`client/lua`; the children never import the orchestrator or each other. The
orchestrator is a shared emission entry point rather than a section processor,
so `body/html` may import it as an explicit composition edge.

The `body/md/` package is the **Markdown processor home** (former `html/md`).
It owns both the `lang="md"` body transform and the `<md>`-tag region scanning
that previously lived in `parser/parser_md_tag.go`; that scanning moved out of
the parser into `body/md/`. The parser imports no Markdown package.

Message-expression parsing (`ParseMessageExpression`) is
mini-template-language semantics and belongs to `dreecode`. Both the front-end
(`lexer`/`parser`) and the sections use it; `dreecode` is shared semantics, not
a section-only package.

### Discovery and code generation

`dreefile/route/`, `dreefile/component/`, `dreefile/layout/`, and
`dreefile/assets/` were **not** split into their own packages in this phase;
they remain in the dreefile root. Moving them would require exporting root-only
helpers and would create front-end imports the rule does not grant. The
`dreefile/i18n/` catalog-extraction package **was** split out. `format/`,
`config/`, and `check/` staying in the root is already sanctioned by the phase
Risks. See the open item in
[`_plan/phase-restructure.md`](../../_plan/phase-restructure.md).

### Shared Markdown home is `internal/md`

`core/markdown.go` imports the shared Markdown implementation. After the move
it must not import the compiler, so the shared Markdown implementation lives at
**`internal/md`** — *not* `internal/render/markdown`, because `render` is the
target-neutral render contract and must not grow a Markdown dependency.

`internal/md` is shared implementation and therefore imports **no compiler
package** (`ir` included). It owns the runtime string path: it exposes an
HTML-string entry point, the `Mode` values (`ModeSafe`, `ModeTrusted`), the
block parser (`ParseBlocks`), and the inline renderer and shared helpers
(`SafeURL`, `IsHR`, `SafeFenceLanguage`, table alignment). `core/markdown.go`
calls `md.ToHTML` directly. The runtime's exported
`core.MarkdownToHTML`/`MarkdownToHTMLTrusted` keep their current public
signatures.

The compiler keeps its `TemplateNode` path in `dreefile/sections/body/md/`
(`ToNodes`, `TransformNodes`). That package is the compiler-side home of the
Markdown-to-`ir` adapter and imports `ir`. It **delegates inline rendering and
the shared helpers** to `internal/md`; it keeps the `ir.TemplateNode` adapter,
the i18n, control-flow, and `<md>`-tag-scanning logic, and its **own block
parser** inside `TransformNodes`. That block parser stays separate from
`shared.ParseBlocks` because embedded `{#if}`/`{{ }}` boundaries must survive
parsing — a plain parsed string block cannot carry them.

`ToNodes` is the only path to `shared.ParseBlocks`, and it is currently reached
only by tests; the production compiler block path is
`ProcessBody -> TransformNodes`. There is **no parity test** between the
compiler block parser and `shared.ParseBlocks`. That is a flagged open item
(see [`_plan/phase-restructure.md`](../../_plan/phase-restructure.md)), not a
settled equivalence.

The dependency direction is `internal/dreefile/sections/body/md -> internal/md`,
never the reverse: `internal/md` stays compiler-free. The runtime only ever
needs rendered text, so `core` calls `internal/md` directly on the string path
and never touches the node adapter or `ir`.

This replaces the earlier open item; it is decided, not deferred.

### Structural check is required

`_tests/sh/check-core-deps.sh` only checks **external** dependencies (it greps
out every `github.com/dreego-stack/dreego` and `golang.org/x/` import). It does
not catch layering or section-to-section violations. A dedicated structural
check for the rules above is therefore **required**, not optional, and lands as
its own slice before the structural moves can be trusted.

### Dependency rule

The rule must match the real code, not an idealized tree. It states the maximum
set of in-repository imports each package may have; a package may use a subset.
The required structural check inspects **non-test** files only, because test
files add edges the packages themselves do not have: `dreecode` and
`sections/body/md` tests import `lexer` and `parser`, and `parser` tests import
`lexer`.

`dreecode` is shared semantics used by BOTH the front-end and the sections — it
is **not** a section-only package. The `lexer` and `parser` stages use
`dreecode` for message-expression and mini-template-language scanning, and the
sections use the same semantics for control flow, expressions, filters, and
slots.

```text
front-end
  lexer                  -> dreecode, ir, tokens
  parser                 -> dreecode, ir, tokens
  tokens                 -> (no in-repo imports)
  ir                     -> (no in-repo imports)

shared
  dreecode               -> ir
  gogen                  -> ir
  codegen                -> ir
  jsoutput               -> gogen
  i18n                   -> (no in-repo imports)

sections (general)
  sections/*             -> codegen, dreecode, gogen, ir, jsoutput (nothing else
                            in-repo, except the explicit composition edges below)

sections (explicit composition edges)
  sections/body/html     -> the general set, plus sections/head, sections/style,
                            and sections/client (the explicit client orchestrator)
  sections/body/md       -> the general set, plus internal/md (the single shared
                            Markdown implementation; it is the only legal
                            compiler consumer of internal/md)
  sections/client        -> the general set, plus sections/client/{js,ts,lua}
                            (parent -> child only)
  sections/client/{js,ts,lua} -> codegen, ir, jsoutput (never a parent or sibling)

dreefile root (generate.go, codegen_*.go, generate_*.go, parser_facade.go,
lex_facade.go, i18n_export.go)
                         -> sections (incl. sections/body/md), front-end, codegen,
                            dreecode, gogen, ir, jsoutput, i18n, internal/gomod
```

The dreefile root orchestrates: it may import the sections, the front-end, the
shared compiler packages, and the root shared implementation (`internal/gomod`),
but nothing imports the root back. The only compiler package that imports
`internal/md` is `sections/body/md`, which delegates its Markdown parsing and
rendering there; `internal/md` never imports a compiler package, and the runtime
string path (`core/markdown.go`) reaches it directly.

No section imports a sibling section **family**, with two verified exceptions:
`body/html` composes the leaf processors `sections/head` and `sections/style`,
and it imports `sections/client` as the explicit client orchestrator. The leaf
sections (`head`, `style`, `body/md`) and the language children
(`client/{js,ts,lua}`) import no sibling and no parent. A section that needs
shared behavior otherwise uses `codegen` (generation state), `dreecode`
(language semantics), `gogen` (Go emission), `jsoutput` (client script
emission), `ir` (types), or — for `body/md` only — the shared implementation
`internal/md`.

`dreecode` owns the mini-template-language semantics, including message
expression parsing; the parser stack is a consumer of `dreecode`, not a
dependency of the sections.

### Layer summary

| Layer | Packages | May import | Must not import |
|---|---|---|---|
| Shared implementation | `internal/{render,context,i18n,server,middleware,session,validate,gomod,templates,md}` | standard library, `golang.org/x/*`, each other per direction | `core`, `adapter/*`, `cmd/*`, `dreegotest`, `internal/dreefile/...` as Go packages |
| Compiler shared | `internal/dreefile/{ir,dreecode,gogen,codegen,jsoutput,i18n}` | standard library, `ir`; `jsoutput` may import `gogen` | other compiler packages outside the rule, shared implementation, `core`, `adapter/*` |
| Compiler front-end | `internal/dreefile/{tokens,lexer,parser}` | `tokens`, `ir`, and `dreecode`; `parser` must not import `lexer` in non-test code | sections, `core`, `adapter/*` |
| Compiler sections | `internal/dreefile/sections/...` | `codegen`, `dreecode`, `gogen`, `ir`, `jsoutput`; `body/html` also `head`, `style`, `sections/client` (orchestrator only); `client` also `client/{js,ts,lua}`; `body/md` also `internal/md` (the single Markdown implementation; its only legal compiler consumer) | other sibling sections, parents or siblings from `client/{js,ts,lua}`, `core`, `adapter/*` |
| Compiler root | `internal/dreefile` | sections (incl. `sections/body/md`), front-end, `codegen`, `dreecode`, `gogen`, `ir`, `jsoutput`, `i18n`, `internal/gomod` | `core`, `adapter/*`, `cmd/*`, `dreegotest`, `internal/md` |
| Facades and hosts | `core`, `adapter/ssr`, `adapter/wails`, `cmd/dreego`, `dreegotest` | shared implementation, public facades | — |

The compiler dependency rule is a review invariant and is enforced by the
required structural check described above; generated import strings
(`github.com/dreego-stack/dreego/core`,
`github.com/dreego-stack/dreego/adapter/ssr`) are output, not imports of the
compiler package.

## Consequences

- The public API is unchanged. Generated code and applications keep importing
  only `github.com/dreego-stack/dreego/core` and, for HTTP hosting,
  `github.com/dreego-stack/dreego/adapter/ssr`.
- `internal/dreefile/` is importable only from inside the repository (CLI,
  `dreegotest`), exactly like `internal/transpiler/` was.
- The split packages in `core/internal/` keep their package names and contents;
  only their location changes. No re-export or forwarding package is added
  inside `core/`.
- The 300-line rule stays enforced per file inside the split.
- `_tests/sh/check-core-deps.sh` keeps the **external** dependency boundary
  checked after the moves. It does **not** check the layering rule, so the
  required structural check is a separate, necessary artifact.
- Structural invariants (no sibling-section imports except the two explicit
  `body/html` composition edges; `dreecode`, `gogen`, and `codegen` import only
  `ir` and `jsoutput` imports only `gogen`; the client orchestrator is the only
  parent of the language children; `sections/body/md` is the only compiler
  consumer of `internal/md`, and `internal/md` imports no compiler package; root
  `internal/` never imports `core`) are asserted by the required dedicated check
  over non-test files.
- The `internal/transpiler` name disappears; documentation, plans, and todos
  that cite `internal/transpiler/...` anchors must be re-pointed.
- **Recorded scope deviation.** `route/`, `component/`, `layout/`, and
  `assets/` did **not** move into their own packages in this phase; they remain
  in the `dreefile` root. Splitting them would have required exporting
  root-only helpers and would have created front-end imports the dependency
  rule does not grant, so the move was deliberately deferred rather than hidden.
  `format/`, `config/`, and `check/` staying in the root was already sanctioned
  by the phase Risks. The deviation is tracked as an explicit open item in
  [`_plan/phase-restructure.md`](../../_plan/phase-restructure.md).
- `ir/mdtohtml.go` (`TranslateMdtohtml`) is codegen **string rewriting**, not
  AST/type code, so it is miscategorized in `ir`. Moving it to `gogen` or the
  server section is a small, non-blocking follow-up and is not part of the
  structural move.

## Supersedes

- The **location parts** of
  [Core runtime split into internal subpackages](core-internal-subpackages.md).
  Its runtime dependency direction (`context -> session, middleware`;
  `middleware -> session`; `server -> session, middleware, validate`) is carried
  over to the new root `internal/*` locations; the `core/internal/*` paths are
  superseded.
- The `core/internal/*` package map in
  [`_plan/phase-render-foundation.md`](../../_plan/phase-render-foundation.md)
  (`core/internal/app|render|compiler|assets|validate`). The capability
  conclusions of that completed phase remain valid; only the location map is
  superseded.
- The package **name and location** in
  [Transpiler as internal subpackage](../../cmd/dreego/_docs/decisions/transpiler-subpackage.md);
  its reasoning (a shared root `internal/` for multiple in-repo consumers) still
  holds, the compiler is now `internal/dreefile/`.

## See also

- [Phase: internal layering and dreefile restructure](../../_plan/phase-restructure.md)
- [Target-neutral application and first-party adapters](target-neutral-application-and-first-party-targets.md)
- [Dreego Architecture](../dreego-architecture.md)
