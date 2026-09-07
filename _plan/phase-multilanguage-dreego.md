# Current phase: multi-language Dreego

## Goal

Extend the already implemented purpose-based root sections with first-party
processors for a small, closed set of source languages. Dreego continues to own
document structure, component composition, template control flow, escaping,
diagnostics, and generated Go integration.

Language processors live in the Dreego monorepo. Their conceptual matrix uses
the first name for the normalized output and the second for the accepted input:

```text
html/html  html/head  html/css  html/md  html/go
js/js      js/ts      js/lua
```

The physical package layout may group closely related processors when that is
clearer than mirroring every matrix cell. Third-party processor support is not
part of this phase.

## Implemented section foundation

The pre-v0.1 breaking migration was:

```text
<go>     -> <server>
<div>    -> <body>
<script> -> <client>
```

`<head>` and `<style>` retain their names. Defaults are:

```html
<server lang="go"></server>
<head lang="html"></head>
<body lang="html"></body>
<style lang="css"></style>
<client lang="js"></client>
```

The `lang` attributes may be omitted for defaults. Unknown attributes,
languages, duplicate singleton sections, invalid ordering, and unsupported
section/language pairs are generation errors.

An HTML `<script>` nested inside an HTML body is ordinary markup. It is not
treated as a root client section:

```html
<body>
    <script type="application/ld+json">{"name":"Dreego"}</script>
</body>
```

## Ownership of template syntax

Body processors do not receive an unstructured source string that can consume
Dreego syntax. The outer parser first recognizes protected Dreego constructs:

- imported component calls such as `<@Card>`;
- control flow such as `{#if}` and `{#each}`;
- slots;
- escaped and explicitly safe expressions;
- source ranges and diagnostic boundaries.

The body processor receives only eligible literal regions plus structured
placeholders. It returns body nodes or safe fragments that are reassembled by
Dreego. A Markdown body can therefore contain Dreego conditions and components
without redefining their meaning.

```html
<body lang="md">
# Account

{#if user.IsAdmin}
<@AdminPanel user={user} />
{/if}
</body>
```

Exact interpolation rules inside language-owned literal regions require tests
with Markdown punctuation, code fences, escaped braces, and raw HTML.

## Processor categories

A processor registers an exact section and language pair with an output kind:

```text
section=body   language=md   output=body-nodes
section=client language=ts   output=javascript
section=client language=lua  output=javascript
```

Support for one pair never implies support for another. For example, a client
Lua processor does not authorize `<server lang="lua">`.

The initial stable categories should be no broader than:

- body source to validated body nodes;
- client source to JavaScript plus source map and assets;
- style source to CSS plus assets.

Head processors and arbitrary parser passes are deferred until real use proves
their contract.

## Process boundary

First-party language processors run inside the Dreego transpiler. Do not use
Go's native `plugin` package, reflection-based module loading, or an embedded
Lua VM.

The versioned subprocess protocol for third-party processors is not planned. It
was removed as speculative API: codegen processors have too much power
to run as third-party code, and the language set is small and closed (Markdown,
TypeScript, Lua-later). The VS Code extension can ship the same grammars.

TypeScript is the one exception: it still needs `node` as an external tool
behind an explicit, pinned, approval-based flow. npm remains opt-in. The
processor records its compiler version, supports reproducible CI, reports
diagnostics against the `.dreego` file, and never installs an unpinned `latest`
version during a build.

## Discovery and installation

Runtime Go plugins continue to use explicit typed `Register` functions.
First-party compiler processors are built into the transpiler and need no
discovery or installation. They are not scanned from arbitrary `go.mod`
dependencies.

The closed language set is Markdown (`md` → `html`), TypeScript (`ts` → `js`),
and Lua (`lua` → `js`) later. Markdown uses stdlib-first parsing. TypeScript
uses a pinned, approved `node` toolchain; npm remains opt-in. First use of the
TypeScript toolchain requires an explicit warning and approval. CI uses a
checked-in lock and allowlist. Offline builds work after approved tools are
cached.

## Reference processors

The first-party processor set is implemented in this order:

1. Markdown body processor — **DONE** (shipped in v0.3, `html`/`md`,
   stdlib-first). It exercises structured body output and protected Dreego
   template placeholders, and preserves protected Dreego constructs.
2. TypeScript client processor — **NEXT**; it will exercise external tooling,
   diagnostics, source maps, JavaScript assets, and type checking. It remains
   the second proof.

3. Lua client processor — **AFTER TYPESCRIPT**; it produces JavaScript through
   the same normalized JavaScript output stage. Lua-to-Go is not planned.

## TypeScript requirements

The TypeScript processor must run a real TypeScript type checker. A transformer
that only removes type annotations is insufficient. It may manage a pinned
native TypeScript compiler through npm or another verified distribution, but
developers interact only with Dreego commands.

The processor records its compiler version, supports reproducible CI, reports
diagnostics against the `.dreego` file, and never installs an unpinned `latest`
version during a build. `node` is an external tool behind an explicit, pinned,
approval-based flow; npm remains opt-in.

## Acceptance criteria

- The section rename is implemented atomically with formatter, scaffolds,
  diagnostics, fixtures, LSP/syntax metadata, and migration documentation.
- Nested HTML scripts remain ordinary body markup.
- Unknown language pairs fail with an installation hint.
- Markdown cannot reinterpret Dreego components or control flow.
- TypeScript type errors fail generation with correct source positions.
- Lua client input produces deterministic JavaScript through the shared output
  stage without embedding a runtime in Dreego.
- A processor crash cannot corrupt existing generated output.
- Identical locked inputs produce identical generated files and assets.
- No processor dependency is added to Dreego core or the transpiler module.

## Not in this phase

- Arbitrary AST mutation across the whole file.
- In-process execution of untrusted plugin code.
- A Dreego-owned Lua interpreter.
- Lua-to-Go compilation.
- Automatic installation without approval.
- Third-party processor support (removed as speculative API).
- Claiming a stable ecosystem protocol after only one processor.
