# Phase: internal layering and dreefile restructure

## Goal

Give the repository one consistent internal layer and one clearly named
compiler package, without changing behavior or the public API:

- root `internal/` is the shared implementation layer;
- the compiler is `internal/dreefile/` with the `ir`, `dreecode`, `gogen`,
  `codegen`, `jsoutput` shared packages, the `tokens`/`lexer`/`parser`
  front-end, and the `sections` structure;
- `core/`, `adapter/*`, `cmd/dreego/`, and `dreegotest/` stay slim facades,
  hosts, and binaries over the shared layer.

The target structure and the dependency rule are decided in
[`_docs/decisions/internal-layering-and-dreefile.md`](../_docs/decisions/internal-layering-and-dreefile.md).
This phase implements them; it does not reopen them.

## Target structure

```text
repo-root/
├── internal/
│   ├── render/            target-neutral render contract
│   ├── context/           Context and SSRContext
│   ├── i18n/              runtime locale resolution and catalog access
│   ├── md/                single Markdown implementation (runtime string path + compiler delegation target)
│   ├── server/            App, routing, static, middleware stack
│   ├── middleware/        (existing)
│   ├── session/           (existing)
│   ├── validate/          (existing)
│   ├── gomod/             (existing)
│   ├── templates/         project scaffolds (from cmd/dreego/internal/templates)
│   └── dreefile/          .dreego compiler
│       ├── generate.go    Run, RunCheck, buildPlan, buildRootPlan, buildRootFile
│       ├── generate_check.go  generation plan diff, disk I/O, apply
│       ├── generate_support.go module path, settings load, hashing
│       ├── ir/            AST and types, single source of truth
│       ├── dreecode/      {#if}, {#each}, {{ }}, filter pipes, {#slot}, expression and message scanning
│       ├── gogen/         GoLiteral, ToPascalCase, escaping, source positions
│       ├── codegen/       shared generation state (codegen.State, codegen.Layout)
│       ├── jsoutput/      shared client emission artifact and <script> tag (former js/output)
│       ├── tokens/        token types
│       ├── lexer/         scanner (imports dreecode, tokens, ir)
│       ├── parser/        section parser (imports dreecode, tokens, ir; no lexer in non-test code)
│       ├── sections/
│       │   ├── head/
│       │   ├── style/
│       │   ├── body/html/
│       │   ├── body/md/   Markdown processor home (md-tag scanning moves here)
│       │   └── client/    explicit client orchestrator (former js/process)
│       │       ├── js/
│       │       ├── ts/
│       │       └── lua/
│       ├── route/         route discovery and generation
│       ├── component/     component discovery and generation
│       ├── layout/        layout discovery and generation
│       ├── assets/        static and plugin client assets
│       ├── format/        dreego fmt
│       ├── config/        dreego.config.json
│       ├── check/         accessibility and structural checks
│       └── i18n/          catalog extraction and generation (build-time)
└── ...
```

The tree is the **target**. This phase did not realize all of it: `route/`,
`component/`, `layout/`, and `assets/` stay in the dreefile root (recorded
deviation, see open item 1), as do `format/`, `config/`, and `check/` (already
sanctioned by the Risks). The `header/` and `server/` section processors also
stayed in the root; only `head/`, `style/`, `body/{html,md}`, and `client/`
became `sections/` packages. The realized tree is `ir`, `dreecode`, `gogen`,
`codegen`, `jsoutput`, `i18n`, `tokens`, `lexer`, `parser`, and
`sections/{head,style,body/{html,md},client/{js,ts,lua}}`.

Dependency rule (matches the verified code, not an idealized tree). It states
the maximum set of in-repository imports each package may have; the structural
check inspects non-test files only, because test files add edges the packages
themselves do not have (`dreecode` and `sections/body/md` tests import `lexer`
and `parser`; `parser` tests import `lexer`):

```text
front-end
  lexer                       ->  dreecode, ir, tokens
  parser                      ->  dreecode, ir, tokens   (no lexer in non-test code)
  tokens                      ->  (no in-repo imports)
  ir                          ->  (no in-repo imports)

shared
  dreecode                    ->  ir
  gogen                       ->  ir
  codegen                     ->  ir
  jsoutput                    ->  gogen
  i18n                        ->  (no in-repo imports)

sections (general)
  sections/*                  ->  codegen, dreecode, gogen, ir, jsoutput

sections (explicit composition edges)
  sections/body/html          ->  the general set, plus sections/head, sections/style,
                                  and sections/client (client orchestrator, explicit entry)
  sections/body/md            ->  the general set, plus internal/md (the single shared
                                  Markdown implementation; its only legal compiler consumer)
  sections/client             ->  the general set, plus sections/client/{js,ts,lua}
                                  (parent -> child only)
  sections/client/{js,ts,lua} ->  codegen, ir, jsoutput (never a parent or sibling)

dreefile root (generate.go, codegen_*.go, generate_*.go, parser_facade.go,
lex_facade.go, i18n_export.go)
                              ->  sections (incl. sections/body/md), front-end, codegen,
                                  dreecode, gogen, ir, jsoutput, i18n, internal/gomod
```

The dreefile root orchestrates: it may import the sections, the front-end, the
shared packages, and `internal/gomod`, but nothing imports the root back. The
only compiler package that imports `internal/md` is `sections/body/md`, which
delegates its Markdown parsing and rendering there; `internal/md` never imports
a compiler package, and the runtime string path (`core/markdown.go`) reaches it
directly.

No section imports a sibling section family, with two verified exceptions:
`body/html` composes the leaf processors `sections/head` and `sections/style`,
and it imports `sections/client` as the explicit client orchestrator. The
`client/` language children never import the orchestrator, `body/html`, or each
other. `dreecode` owns the mini-template-language semantics (including message
expression parsing) and is shared by BOTH the front-end (`lexer`/`parser`) and
the sections — not the token/lexer/parser stack and not a section-only package.
Root `internal/` never imports `core` as a Go package; codegen references public
import paths only as emitted strings.

## Release strategy

Every slice lands as its own pull request with exactly one `.changes/*.md` file
and `version: patch`. PR0 (this documentation PR) and PR1 are the first slices;
PR1 is already implemented and green. PR2–PR4 are structural moves. PR2 is
**not** a pure move and rename: it is a move plus real regrouping (introducing
`codegen/` and `jsoutput/` homes, the `sections/client/` orchestrator, the
`body/md/` processor home, and the lexer/parser/token re-homing), so it carries
genuine content and import-graph edits, not just path rewrites. No public API,
generated output, or observable behavior changes. The full `task test` suite
stays green after every slice.

## PR slices

Each slice is a separate PR, in this order. One slice per PR; do not combine.
`PR2a` is a required companion to PR2 rather than a separate position in the
sequence.

### PR0: docs — the ADR, phase plan, and task file (this file's slice)

- Ships `_docs/decisions/internal-layering-and-dreefile.md`,
  `_plan/phase-restructure.md`, and `.agents/tasks/restructure-dreefile/main.md`
  as their own small documentation PR, before any code slice.
- One `.changes/*.md` file, `version: patch`.
- No code changes; `task test` unaffected.

### PR1: remove the three real duplicates — DONE, green

This slice is **already implemented and green** on the branch. It is recorded
here for completeness; no further work is expected beyond review.

- Deleted `internal/transpiler/expr_kind.go` (byte-identical to
  `internal/transpiler/ir/expr_kind.go` modulo naming).
- Deleted `internal/transpiler/source_pos.go` (identical to
  `internal/transpiler/ir/source_pos.go`).
- Removed the pass-through wrapper wall in
  `internal/transpiler/codegen_html.go` that forwarded to `ir`, `output`, and
  `head` without adding behavior.
- Kept the canonical implementations in `ir/`.
- Real call sites that changed to the `ir`/`output`/`head` names:
  - `internal/transpiler/a11y_check.go` (`posToLineCol` -> `ir.PosToLineCol`)
  - `internal/transpiler/layout_chain.go` (`sourceLocation` -> `ir.SourceLocation`)
  - `internal/transpiler/codegen_layout_head_warning.go` (`posToLineCol` -> `ir.PosToLineCol`)
  - `internal/transpiler/generate_route_methods.go` (`setNodeSource`/`setSourceText` -> `ir.*`)
  - `internal/transpiler/generate_components.go` (`setNodeSource` -> `ir.SetNodeSource`)
  - `internal/transpiler/generate_layout.go` (`setNodeSource`/`setSourceText` -> `ir.*`)
  - Tests: `a11y_check_test.go`, `codegen_golden_test.go`,
    `codegen_layout_head_warning_test.go`, `codegen_helpers_test.go`,
    `output_safety_refresh_test.go`.
- `html/output/component_call.go` and `html/output/component_slots.go` already
  used `ir.SourceRef`/`ir.SourceLocation` at the base and required no edit; they
  are listed for completeness so the wrapper-wall scope is unambiguous.
- Behavior unchanged. `task test` green.

### PR2: `internal/transpiler` -> `internal/dreefile`

- Move the package to `internal/dreefile/`.
- Introduce the target structure exactly as decided: shared packages `ir/`,
  `dreecode/`, `gogen/`, `codegen/`, `jsoutput/`, `i18n/`; front-end `tokens/`,
  `lexer/`, `parser/`; `sections/head|style|body/{html,md}|client/` with
  `client/{js,ts,lua}`. `route/`, `component/`, `layout/`, `assets/`, `format/`,
  `config/`, `check/`, and the `header`/`server` section processors stay in the
  dreefile root (recorded deviation; see open item 1 and the Risks).
- This is a move **plus real regrouping**: `codegen/` and `jsoutput/` become
  shared homes, `js/process` becomes the `sections/client/` orchestrator, the
  md-tag scanning moves from `parser/` into `sections/body/md/`, and
  lexer/parser/tokens get explicit front-end homes. Expect import-graph edits,
  not only path rewrites.
- Keep central orchestration in `dreefile/generate.go` (`Run`, `RunCheck`,
  `buildPlan`, `buildRootPlan`, `buildRootFile`) plus `generate_check.go` and
  `generate_support.go`.
- Language subfolders only where real variants exist (`body`, `client`).
- Update import paths in `cmd/dreego`, `dreegotest`, `core/markdown.go` (now to
  the shared `internal/md`, see the decided Markdown home), and all internal
  references.
- Behavior unchanged. `task test` green.

### PR2a (required): structural check for the layering rule

- Companion slice to PR2; it does not change the PR0–PR4 numbering.
- Add the dedicated structural check that `_tests/sh/check-core-deps.sh` does
  **not** provide: `check-core-deps.sh` greps out all in-repo imports and only
  verifies external dependencies.
- Assert over **non-test** files: no sibling-section imports except the two
  explicit `body/html` composition edges (`sections/head`, `sections/style`,
  `sections/client`); `dreecode`, `gogen`, and `codegen` import only `ir` and
  `jsoutput` imports only `gogen` (plus standard library); the `sections/client`
  orchestrator is the only parent of `client/{js,ts,lua}`; front-end direction
  `lexer -> dreecode, tokens, ir` and `parser -> dreecode, tokens, ir`
  (`parser` must not import `lexer` in non-test code); `internal/md` is imported
  by `sections/body/md` only and imports no compiler package; root `internal/`
  never imports `core`. Test-file edges (`dreecode` and `sections/body/md` tests
  import `lexer`/`parser`; `parser` tests import `lexer`) are excluded.
- This slice is **required**, not optional; the structural invariants are not
  otherwise checked.
- New file(s) plus test wiring; one `.changes/*.md`, `version: patch`.

### PR3: move `core/internal/*` to root `internal/*`

- Move `core/internal/render|context|i18n|server` to
  `internal/render|context|i18n|server`.
- Slim `core/` to the public facade and re-export layer. Do not rename the
  `render` package contents; only its location changes.
- Keep generated imports (`github.com/dreego-stack/dreego/core`,
  `github.com/dreego-stack/dreego/adapter/ssr`) and the public API stable.
- No forwarding packages inside `core/`; update `core/facade.go`,
  `core/render.go`, `core/i18n.go`, `core/markdown.go`, and internal
  cross-imports to the new paths.
- `core` and `adapter/ssr` still build; `task test` green.

### PR4: `cmd/dreego/internal/templates` -> `internal/templates`

- Move the project scaffold package to root `internal/templates`.
- Update the CLI consumer and the `//go:embed` paths.
- Behavior unchanged. `task test` green.

### Follow-up (non-blocking): re-home `ir/mdtohtml.go`

- `TranslateMdtohtml` is codegen string rewriting, not AST/type code, so it is
  miscategorized in `ir`. Move it to `gogen` or the server section.
- Small, separate change; it must **not** block PR2.

## Acceptance criteria

- `internal/transpiler/` no longer exists; the compiler is `internal/dreefile/`
  with the structure realized above (target minus the recorded scope deviation).
- `core/internal/` no longer exists; `core/` contains only the public facade,
  its tests, and the public helper files.
- `core/internal` is not referenced by any import path in the repository.
- No import cycle: root `internal/` never imports `core`; `internal/dreefile`
  never imports `core` or `adapter/*`.
- The dependency rule holds as stated: no sibling-section imports except the two
  explicit `body/html` composition edges (`head`, `style`, `sections/client`);
  `dreecode`, `gogen`, and `codegen` import only `ir` and `jsoutput` imports
  only `gogen`; the `sections/client` orchestrator is the only parent of
  `client/{js,ts,lua}`; `lexer -> dreecode, tokens, ir` and
  `parser -> dreecode, tokens, ir` with non-test code only.
- `internal/md` is the single Markdown implementation and the shared Markdown
  home; neither `core` nor `internal/render` imports `internal/dreefile`. The
  only compiler consumer of `internal/md` is `sections/body/md`, which delegates
  its parsing/rendering there; `internal/md` imports no compiler package.
- The public API is unchanged: generated code and applications compile without
  edits.
- Generated output is byte-identical before and after each slice (no golden
  file changes; import-path strings emitted into generated code do not change).
- `task test` (race suite, coverage gate, dependency check) is green after each
  slice.
- The required structural check (PR2a) passes and is wired into the suite.
- The 300-line rule holds for every moved file, or the violating file is split
  in the same slice.
- `AGENTS.md` and `_docs/dreego-architecture.md` describe the new layout.

## Verification

Run through `smd` only:

- `task test` after each slice.
- `sh _tests/sh/check-core-deps.sh` after PR2 and PR3. Note: this script only
  checks **external** dependencies and does not catch layering violations.
- The required structural check (PR2a) after PR2.
- A repository-wide search confirms no `internal/transpiler`, `core/internal`,
  or `cmd/dreego/internal/templates` import path remains.
- `go build ./...` through the workspace, or the CLI build via `task build`.
- Golden-code tests (`internal/dreefile/...` testdata) must pass unchanged.

## Risks

- **Large structural diff.** The compiler tree has 160+ files including tests
  (194 Go files: 109 non-test, 85 test). PR2 is a move plus real regrouping
  (`codegen/`/`jsoutput/` homes, `sections/client/` orchestrator, `body/md/`
  processor placement, lexer/parser/token re-homing), so it is larger than a
  path rewrite. Mitigation: one slice per PR, PR2 keeps behavior byte-identical,
  and the structural invariants are asserted by the required check.
- **Accidental import cycle.** `core` already imports root `internal/*`, but a
  reverse import would fail to build. Mitigation: the required structural check
  plus `go build ./...` in CI; codegen only ever emits import strings.
- **Layering gaps go unnoticed.** `_tests/sh/check-core-deps.sh` only checks
  external dependencies and does not catch section-to-section or front-end
  direction violations. Mitigation: the required structural check (PR2a).
- **Two modules, one internal tree.** The root and `core` modules share
  `internal/` by path prefix. This works today but is easy to break by adding a
  third module or a `replace` that changes path prefixes. Mitigation: document
  the rule and keep the shared packages free of module-specific logic.
- **Naming proximity: `internal/i18n` vs `internal/dreefile/i18n`.** Two
  packages named `i18n` with different roles invite wrong imports. See open
  item 1; resolve before PR3.
- **`format`, `check`, `config` as single-purpose packages.** Splitting these
  out may create tiny packages and extra indirection. Keep them only where the
  dependency direction justifies a package boundary; otherwise they stay in the
  dreefile root.
- **Test-file boundaries.** Mixed tests that use unexported internals must move
  with their package, as in the earlier `core/internal` split. Mitigation: keep
  each test next to the code it tests during the move.
- **Anchor staleness.** `_plan/phase-dreefile.md` and other docs cite
  `internal/transpiler/file.go:line` anchors. Mitigation: re-point them in PR2
  (documentation slice of that PR).

## Not in this phase

- No behavior change, no new feature, no public API change.
- No rename of the `render` package or its contents; only its location moves.
- No compatibility or forwarding packages for the old import paths.
- No universal `Target`, processor, cache, or live interface.
- No new module, no change to which modules participate in a coordinated
  release, no tag or version-strategy change.
- No `_plan` phase renumbering or roadmap change.
- No deletion of the `AGENTS.md` follow-up items listed under open items.

## Decided (resolved; no longer open)

- **Shared Markdown home is `internal/md`.** Chosen over
  `internal/render/markdown`: `render` is the target-neutral render contract
  and must not grow a Markdown dependency. `internal/md` owns the runtime
  string path and the shared inline-rendering/helper surface (`ParseBlocks`,
  `ToHTML`, `SafeURL`, `IsHR`, table alignment). `core/markdown.go` imports it
  directly (string path); the compiler section `sections/body/md/`
  (`ToNodes`, `TransformNodes`) delegates inline rendering and the shared
  helpers to it and keeps the `ir.TemplateNode` adapter, the i18n,
  control-flow, and `<md>`-tag-scanning logic, and its own block parser. The
  block parser stays separate because embedded `{#if}`/`{{ }}` boundaries must
  survive parsing. The dependency direction is
  `sections/body/md -> internal/md`, never shared -> compiler. This resolves
  the former open item before PR2; the missing block-parser parity test is now
  open item 6.
- **`codegen.State` home is `dreefile/codegen/`.** Kept as a shared package
  rather than merged into `gogen`: `gogen` is stateless Go emission, while
  `codegen.State` is the mutable build state (component definitions, imports,
  message uses, Lua features) that nearly every section imports. The rule
  becomes `sections/* -> codegen, dreecode, gogen, ir, jsoutput`.
- **Client script emitter home is `dreefile/jsoutput/`** (former `js/output`).
  It holds `Artifact`, `GenClient`, `GenClientTo`, imports `gogen`, and is
  importable by body/html and every client language.
- **Client orchestrator is `dreefile/sections/client/`** (former `js/process`).
  It imports its children `client/js`, `client/ts`, `client/lua`
  (parent -> child only); `body/html` imports the orchestrator as the explicit
  client entry. The children import no parent and no sibling.
- **Front-end homes are `dreefile/tokens/`, `dreefile/lexer/`, and
  `dreefile/parser/`.** `dreecode` is the mini-template-language semantics
  (including message expression parsing), shared by the front-end and the
  sections — not the token/lexer/parser stack.
- **`body/md/` is the Markdown processor home.** The `<md>`-tag region scanning
  moved out of the parser into `sections/body/md/`; the parser imports no
  Markdown package.
- **The structural check is required, not optional** (PR2a). It is the only
  check that asserts the layering rule; `check-core-deps.sh` only covers
  external dependencies.

## Open items (flagged, not silently decided)

1. **`route/`, `component/`, `layout/`, `assets/` did not leave the dreefile
   root.** This phase deliberately deviates from the target structure above:
   these four concerns were **not** split into their own packages and remain in
   `internal/dreefile/`. Reason: moving them would require exporting root-only
   helpers and would create front-end imports the dependency rule does not
   grant, so a split would either widen the rule or duplicate code. `format/`,
   `config/`, and `check/` staying in the root is already sanctioned by the
   phase Risks (single-purpose packages stay in the root unless the dependency
   direction justifies a boundary). `i18n/` *was* split out. Follow-up work, if
   pursued, must first establish the needed helper boundaries and an explicit
   rule extension; until then the root keeps these files. This is a recorded
   deviation, not an unfinished silent step.
2. **`internal/i18n` (runtime) vs `internal/dreefile/i18n` (catalog).** Both
   packages are called `i18n` and both may be imported from a file that touches
   localization. Options: keep the names and use a named import alias; rename
   the compiler package to `internal/dreefile/catalog` or
   `internal/dreefile/messages`; rename the runtime package. Decide before PR3,
   because the runtime move and the compiler move meet there.
3. **`core/i18n_selection.go` is an HTTP handler in the public facade.** It
   implements `LocaleSelectionHandler` on `net/http` and belongs next to the SSR
   host or the runtime i18n HTTP glue, not in the target-neutral facade.
   Moving it is a public-API question (`core` currently exports
   `LocaleSelectionHandler`), so it is flagged for a follow-up slice rather than
   folded into the mechanical PR3 move.
4. **`AGENTS.md` follow-up prose.** The "File Structure" block, the
   `internal/transpiler/` Coding-Rules bullet, and the Feature Workflow path are
   re-pointed to `internal/dreefile/`. The remaining "Core and Plugin Boundary"
   prose still uses "transpiler" as the mechanism name, which stays valid; any
   further prose re-pointing is a separate follow-up if it becomes ambiguous.
5. **`ir/mdtohtml.go` re-homing.** `TranslateMdtohtml` is codegen string
   rewriting, miscategorized in `ir`; move it to `gogen` or the server section
   as a small non-blocking follow-up. It must not block PR2.
6. **No parity test between the two Markdown block parsers.** The compiler's
   production block path is `ProcessBody -> TransformNodes`
   (`sections/body/md/transform.go`), which contains its own block parser and
   cannot call `shared.ParseBlocks` because embedded `{#if}`/`{{ }}` boundaries
   must survive parsing. `shared.ParseBlocks` is reached only through
   `sections/body/md/markdown.go` `ToNodes`, which today is called only from
   tests (`markdown_test.go`, `edge_cases_test.go`, `security_test.go`). Nothing
   asserts that the two block parsers agree. Open: add a parity test over the
   shared corpus (headings, lists, fences, tables, HR, HTML blocks) or state
   explicitly which constructs are allowed to diverge. Flagged, not silently
   assumed equivalent.

The former open items "core/markdown.go imports compiler internals" and
"structural-invariant test" are now decided above.
