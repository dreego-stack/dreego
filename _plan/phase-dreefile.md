# Phase: explicit Dreefile composition

## Release strategy

The Dreefile grammar is an umbrella capability delivered in independent slices.
The first slice is a mandatory correctness fix to head composition that lands
before any grammar keyword. The grammar, layout chaining, directory-resolution
removal, and server import channel follow as separate pull requests. Each slice
is a `version: none` change file on the `stage/*` branch; the stage merge into
`main` is the single deliberate release act. Capability names below are not
version promises.

## Goal

Make composition explicit and free of directory magic. A `.dreego` file declares
its own kind in the file, a page selects its layout by path, and a component
imports other components the way Python does. One rule governs the surface:

**UPPERCASE = keyword, lowercase = value.**

## Grammar

```text
DREEFILE component (title string)          component kind, props on the DREEFILE line
DREEFILE layout                            layout kind
no DREEFILE                                 page (default)
LAYOUT "www/layouts/admin.dreego"          explicit layout selection
COMPONENT "www/components" IMPORT { Card, Button, Card as ProductCard }
GOIMPORT { sync, encoding/json }           explicit Go import channel
```

A component's name comes from its **filename**: `Card.dreego` becomes `<@Card>`.
Props stay on the `DREEFILE` line. `COMPONENT` imports components from a path
with optional aliases (`Card as ProductCard`). `GOIMPORT` is the designated
channel for Go imports; it is parsed and reserved, and codegen consumption
arrives with the server-import slice.

## Locked decisions

These are decided. The slices implement them; they do not reopen them.

1. `DREEFILE` replaces the legacy `Component` header line. The component name
   comes from the filename (`Card.dreego` -> `<@Card>`); props live on the
   `DREEFILE` line. Landed.
2. Local versus module imports keep the existing dot heuristic. `COMPONENT` does
   not rework local-vs-external resolution; the heuristic in
   `internal/transpiler/generate_components.go:141-169` stays until a later,
   separately-proven need.
3. A top-level `import "…"` becomes a **generate error** with a migration
   message pointing to `GOIMPORT` and `<server>`. Because Go imports therefore
   move into `<server>`, this re-scopes
   [`_todo/core/server-stdlib-imports.1.md`](../_todo/core/server-stdlib-imports.1.md)
   from an optional gap into the defined import channel for route and component
   code.
4. The layout directory stays `www/layouts` (plural). Slice 4 removes the
   *directory-driven resolution* (`default.dreego` / `layout.dreego` scope
   cascade and its ambiguity check), not the folder itself; layouts are selected
   only by an explicit `LAYOUT` path.
5. No further keywords. `CACHE` comes later as its own phase and is not implied
   by this grammar.

## Current state

The grammar keywords are parsed, formatted, and enforced: `DREEFILE`,
`LAYOUT`, `COMPONENT ... IMPORT`, and `GOIMPORT` are understood by the lexer
and `dreego fmt`, the legacy `Component` / `import` / `from` forms are generate
errors, and the repository `.dreego` files and public documentation are
migrated. `LAYOUT` and `GOIMPORT` are stored on the file but have no codegen
consumer yet; the anchors below describe the pre-grammar baseline that the
remaining slices change.

- `internal/transpiler/lexer/lexer_header.go` — `ParseFileHeaderStrict` parses
  `DREEFILE`, `LAYOUT`, `COMPONENT ... IMPORT`, and `GOIMPORT`, and returns a
  `file:line:col` error for the legacy `Component` / `import` / `from` forms.
  Done; the remaining slices add codegen consumers, not parsing.
- `internal/transpiler/discovery.go:57-60` — `isComponentsDir`; and
  `internal/transpiler/discovery.go:62-65` — `isLayoutsDir`. Both decide file
  kind by directory name, which is exactly the magic the `DREEFILE` line
  replaces.
- `internal/transpiler/generate_components.go:141-169` — `importedComponentPaths`
  walks only `root/routes`; the local-vs-module split is the dot heuristic at
  `:161` (`strings.Contains(imp.Path, ".")`). Locked decision 2 keeps this.
- `internal/transpiler/fmt.go` — `Format` knows all four header keywords and
  keeps legacy header lines verbatim (so `dreego fmt` does not rewrite or
  corrupt a file that generation will reject). Done.
- `internal/transpiler/html/output/templ.go:80-124` — layout call and head
  merge. The runtime title/meta dedupe at `:94-112` is gated by
  `layout.File.Head != nil`, and `File.Head` is set only for a root-level
  `<head>` (`internal/transpiler/parser/parser.go:107`). Canonical body-level
  layouts therefore never dedupe.
- `internal/transpiler/generate_layout.go:20-88` — directory-based layout
  discovery; `:123-134` — `resolveLayoutForRoute` scope cascade; `:90-111` —
  `detectAmbiguousLayouts`. Slice 4 retires this resolution path.

## Head composition rule

Composition across a layout chain follows one rule, checked against seven
frameworks. **No framework uses outermost-wins.**

- Single-value tags (`title`, `meta name="description"`): innermost wins, with
  outward fallback when the inner scope is silent.
- Many-value tags (for example verification or Open Graph collections): merge
  per key.
- Exactly one `<title>` per document.
- The title template belongs on the **outermost** layout; inner scopes supply
  the page-specific value.
- Ambiguous duplicates that cannot be resolved fail early at generation.

Sources:

- Next.js, `generateMetadata` (evaluation root -> page, shallow merge, later
  duplicate keys replace, parent `title.template` applies to children):
  <https://nextjs.org/docs/app/api-reference/functions/generate-metadata>
- Unhead, handling duplicates (most recent wins; many-value names preserved;
  title singleton):
  <https://unhead.unjs.io/docs/head/guides/core-concepts/handling-duplicates>
- Unhead, titles (`titleTemplate` prefix/suffix):
  <https://unhead.unjs.io/docs/head/guides/core-concepts/titles>
- SvelteKit SEO (`<title>` and `<meta description>` per page):
  <https://svelte.dev/docs/kit/seo>
- SvelteKit issue 10089 (title deduped, meta duplicates slip through):
  <https://github.com/sveltejs/kit/issues/10089>
- Astro, transferring slots (nested layouts transfer the head slot):
  <https://docs.astro.build/en/basics/astro-components/#transferring-slots>
- Phoenix `live_title` (`prefix` / `suffix` / `default` around `page_title`):
  <https://hexdocs.pm/phoenix_live_view/Phoenix.Component.html#live_title/1>
- Rails layouts and rendering (`yield :head`, `content_for`, nested layout
  inheritance):
  <https://guides.rubyonrails.org/layouts_and_rendering.html>
- Django template inheritance (`{% block %}` override, parent content as
  fallback):
  <https://docs.djangoproject.com/en/5.0/ref/templates/language/#template-inheritance>

## Implementation slices

One pull request per slice, in this order. Slice 1 is mandatory first.

1. **Head dedupe fix.** Make the head merge produce exactly one `<title>` and one
   meta description for both root-level and body-level layout shapes. This is the
   bug in `_todo/core/layout-head-title-dedupe.1.md`
   (`templ.go:94-112` gated by `layout.File.Head != nil`; `parser.go:107`).
   Dependency-free, and a prerequisite for trustworthy chained-head tests.
2. **Grammar and fmt support.** Parse and emit `DREEFILE` (component / layout /
   page), `LAYOUT`, `COMPONENT ... IMPORT { ... }`, and `GOIMPORT` in
   `lexer/lexer_header.go:9-50`; teach `internal/transpiler/fmt.go:24-46` the
   same keywords; turn `Component`, `import "…"`, and `from "…" import` into
   generate errors with a migration message. Component name comes from the
   filename. **Landed.**
3. **Layout chaining with cycle detection.** Allow a layout to declare
   `LAYOUT` and be selected by another layout; reject cycles and missing targets
   with a source diagnostic. This makes the head-composition rule from the
   section above observable across a chain.
4. **Directory-resolution removal.** Retire `isLayoutsDir`
   (`discovery.go:62-65`), the `default.dreego` / `layout.dreego` scope cascade
   (`generate_layout.go:123-134`), and `detectAmbiguousLayouts`
   (`generate_layout.go:90-111`). Layout files stay under `www/layouts`
   (decision 4) and are referenced only by explicit `LAYOUT`.
5. **Server import channel.** Route and component stdlib imports move into
   `<server>` and compile through `GOIMPORT` plus the allow-list from
   `_todo/core/server-stdlib-imports.1.md`. The handwritten sibling-`package`
   workaround stops being necessary for the common cases.

## Smallest decisive fixture

One reference fixture proves the slice 2-3 codegen contract:

```text
base.dreego            layout, no LAYOUT
admin.dreego           layout, LAYOUT base
routes/dashboard/+page.dreego   page, LAYOUT admin
```

The fixture asserts:

- nested render: base wraps admin wraps the route body;
- pairwise title dedupe: exactly one `<title>` across the pair at each level;
- cascade coexistence: during slices 2-3 the explicit `LAYOUT` chain and the
  existing directory cascade both work until slice 4 removes the latter;
- the 404 diagnostic: a `LAYOUT` path that does not exist fails generation with
  a diagnostic naming file and line;
- `dreego fmt` round-trip: formatting is idempotent and preserves all four
  keywords.

## Phase gate

Adapted from [`_plan/README.md`](README.md):

- public behavior has black-box integration coverage;
- generated code compiles through the supported Dreego CLI workflow;
- documentation and migration guidance match released behavior;
- accessibility, security, race, and dependency checks pass;
- a reference application exercises the capability;
- unresolved design questions are answered or explicitly deferred.

Phase-specific additions:

- every layout chain yields exactly one `<title>` and never a duplicate meta
  description, in both layout shapes;
- the `import "…"` generate error names the migration path (`GOIMPORT` /
  `<server>`);
- `dreego fmt` round-trips the new grammar without moving keywords into the
  body;
- a cycle or missing `LAYOUT` target fails at generation, not at runtime.

## Not in this phase

- **No YAML frontmatter.** A YAML dependency is not allowed in the transpiler;
  `---` collides with the header terminator; and fence arithmetic in
  `internal/transpiler/source_pos.go:27` would shift every `file:line:col`
  diagnostic. The grammar is line-keyword based on purpose.
- **Do not remove `dreego generate`.** The transpiler is `internal/`,
  `_tests/sh/check-core-deps.sh` enforces the dependency boundary, the accepted
  ADR `cmd/dreego/_docs/decisions/transpiler-vs-runtime.md` requires compile-time
  safety, and `README.md` states that Dreego is a compile-time transpiler.
- **No universal interfaces.** No `Target`, processor, or cache interface before
  real implementations prove a small shared contract.
- `CACHE` keyword and cache policy (later phase).
- Replacing the existing component dot heuristic (locked decision 2).

## Superseded draft

A previously drafted version of this file on the discarded branch
`docs/dreefile-phase` (base `c7eeadf`) used `LAYOUT` in PascalCase and a `BASE`
keyword. That draft is superseded by the all-UPPERCASE grammar above; no part of
it is carried forward.
