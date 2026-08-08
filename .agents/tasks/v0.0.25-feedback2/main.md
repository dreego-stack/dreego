---
type: Task
title: v0.0.25 Feedback2 — Fix open HashNet feedback items before release
status: in_progress
assign: manager
---

# v0.0.25 Feedback2

Goal: Fix the open items from testspace/feedback.md (HashNet project) before
pushing v0.0.25, so the released version is the best possible. Tailwind
plugin research is done (MIT, standalone binary, feasible) — plugin block
goes to v0.0.26 planning.

## Open items (priority order)

1. **CLI/core version drift** — cmd/dreego/go.mod requires core v0.0.23,
   VERSION is v0.0.25. Fix: bump require to v0.0.25 (or use the local
   replace pattern like new.go does). Test: CLI build uses current core.
2. **init import bug** — blueprints/default/main.go.tmpl has `_ "gen"` but
   generate puts the package in dreego/gen/. Fix: `_ "<modulename>/dreego/gen"`.
   Test: fresh init + generate + build works.
3. **Quoted prop interpolation** — `prop="{var}"` (with quotes) generates a
   literal string; only `prop={var}` (no quotes) works as Go expression.
   Unify: quoted `{var}` should also resolve as expression. Test: component
   with quoted prop compiles and renders.
4. **Props defaults ignored** — `= default` in component signature is parsed
   by parseProps but GenerateComponent only generates required params. Fix:
   generate defaults. Test: component with default prop renders without
   passing it.
5. **$loop in if-conds** — $loop.X is only replaced in each-children, not in
   {#if} conds. Fix: apply replacement in conds too. Test: {#if !$loop.Last}
   compiles.
6. **Duplicate title in {#head} merge** — layout title + route title both
   emitted. Fix: dedupe title/meta in head merge. Test: only one title.
7. **404/500 without layout** — error routes bypass layout, scope div before
   doctype. Fix: don't emit scope div before doctype; optionally wrap error
   pages in layout. Test: 404 page has doctype first.
8. **Parser: text before section** — non-section text at file start makes
   following <go> blocks land as text. Fix: error or handle. Test: file
   starting with <!doctype html> + <go> block works or errors clearly.
9. **< in go strings** — lexer parses <...> inside Go strings/backticks.
   Fix: skip tag parsing inside Go string literals. Test: backtick string
   with < survives.
10. **ctx vs c docs** — document that components use ctx, routes use c.
    Docs only.
11. **--version flag** — help text promises it, only `version` subcommand
    exists. Fix: add --version flag. Test: dreego --version works.

## Tailwind plugin (research done, v0.0.26)

- MIT license, standalone binary (no Node), .dreego files scannable directly.
- Plan: plugins/tailwind calling standalone binary in dreego build step.
- NOT part of this task — create a block file _todo/blocks/tailwind-plugin.1.md
  with the research summary and plan.

## Execution

Sequential TDD per item: coder-test → coder-write → reviewer → commit.
No push.

## Status

- [ ] Task file created
- [x] Item 1: version drift
- [x] Item 2: init import bug (tests ROT, fix pending)
- [ ] Item 3: quoted props (tests done: quoted {var} GREEN, {#if}-in-attr ROT)
- [x] Item 3b: {#if} in attributes (option C: clear error, tests GREEN)
- [ ] Item 4: props defaults (tests ROT, fix pending)
- [ ] Item 5: $loop in conds (tests GREEN — bug already fixed, no prod change)
- [ ] Item 6: duplicate title (tests ROT, fix pending)
- [ ] Item 7: 404/500 layout (tests ROT, fix pending)
- [ ] Item 8: parser text-before-section (tests ROT, fix pending)
- [ ] Item 9: < in go strings (tests ROT, fix pending)
- [ ] Item 10: ctx/c docs
- [x] Item 11: --version flag (tests ROT, fix pending)
- [ ] tailwind-plugin.1 block file created

## Item 1: version drift — tests (ROT) + chosen fix

Tests written (no production code):

- `cmd/dreego/version_test.go` → `TestDreegoVersionMatchesVERSIONFile`: `dreegoVersion()`
  without ldflags must equal the repo-root VERSION file content (v0.0.25), never
  "dev"/"(devel)". GREEN — version resolution itself works.
- `_tests/core/CLI/version-drift/test.sh`: `dreego version` output must contain
  the VERSION-file version, and `cmd/dreego/go.mod` must either require a core
  version matching VERSION or carry a local `replace ... => ../core`. RED —
  go.mod requires core v0.0.23 while VERSION is v0.0.25.

Chosen fix (option a): `cmd/dreego/go.mod` requires `codeberg.org/dreego/dreego/core v0.0.25`
(bumped from v0.0.23, matching VERSION), no `replace` directive. `cmd/dreego/go.sum`
carries the v0.0.25/go.mod hash (core/go.mod unchanged since v0.0.22, hash verified
against the known v0.0.22/v0.0.23 entry). Test `_tests/core/CLI/version-drift/test.sh`
now asserts `require codeberg.org/dreego/dreego/core $VERSION` exactly and `dreego
version` output; no replace path accepted.

RELEASE NOTE: require v0.0.25 set — tags must be pushed before external `go install`
works; in-repo builds use go.work (local core wins). `go mod tidy` in cmd/dreego
fails until the `core/v0.0.25` tag is published remotely (expected, documented).

## Item 2: init import bug — tests (ROT) + chosen fix

Tests written (no production code):

- `_tests/core/CLI/init-import/test.sh`: `dreego init .` in a temp dir with a
  hand-written go.mod (`module t`, core require+replace per how-to-test-sh.md),
  then asserts (a) main.go imports `_ "t/dreego/gen"` and does NOT contain the
  bare `_ "gen"` import, (b) `dreego generate` + `GOWORK=off go build .` succeed.
  RED — main.go still carries `_ "gen"` (init.go copies the template verbatim,
  no placeholder replacement), so the build would fail with "package gen is not
  in std".
- `cmd/dreego/init_test.go` (unit): `TestDefaultBlueprintGenImport` reads
  `blueprints/default/main.go.tmpl` from the embedded FS and requires the
  `_ "§$name$§/dreego/gen"` placeholder import (no bare `_ "gen"`). RED.
  `TestLandingBlueprintGenImport` pins the landing template as the reference
  pattern. GREEN (control test).

Chosen fix (option a, same mechanism as `dreego new`): `cmdInit` replaces the
`§$name$§` placeholder in blueprint files with the module name, exactly like
`cmdNew` does (new.go:52 `strings.ReplaceAll(string(data), "§$name$§", projName)`).
The module name is read from the target directory's existing go.mod
(`go mod edit -json` or a `module ` line parse); `dreego init` does not run
`go mod init` itself, so the go.mod must already exist (or init creates it).
`blueprints/default/main.go.tmpl` changes `_ "gen"` → `_ "§$name$§/dreego/gen"`,
matching the landing blueprint. The integration test's go.mod uses `module t`,
so the expected import is `_ "t/dreego/gen"`.

## Item 3: quoted prop interpolation — tests (GREEN + ROT)

Tests written (no production code):

- `core/codegen_helpers_test.go` (unit, GREEN): `attrVal` behavior pinned:
  - `prop="{var}"` → `var` (expression, NOT literal `"{var}"`) — GREEN, fixed in
    v0.0.24/25 (attrVal trims quotes then resolves the single `{…}` shape).
  - `prop={var}` → `var` — GREEN.
  - `prop="literal"` → `"literal"` (string stays) — GREEN.
  - `prop="{a}-{b}"` → concatenation via concatPlaceholders — GREEN (already
    covered by TestAttrValResolvesMultiplePlaceholdersToConcatenation).
  - `active="false"` → `"false"` (string literal, documented): a bool prop must
    be passed unquoted `active={true}`; quoted `"false"` yields a string and
    fails to compile at the call site. Documented, not a bug.
- `_tests/core/Bugs/component-quoted-prop/test.sh` (integration, GREEN):
  `Component Card (title string, active bool)` + route
  `<@Card title="{myTitle}" active={true}/>` with mixed quotes. `dreego generate`
  + `go build` succeed; generated code contains `Card(myTitle, true)` and no
  literal `{myTitle}`. Fixed in v0.0.24/25 (extractAttrValues/attrVal).
- `_tests/core/Bugs/attr-if-in-attribute/test.sh` (integration, ROT):
  `<a class="nav {#if cond}active{/if}">` — NOT supported, two failure modes:
  - Route path: `{#if}` stays literal in the attribute (`b.WriteString(`<a
    class="nav {#if cond}active{/if}">`)`), so `cond` is never referenced →
    `cond declared and not used` compile error.
  - Component body path (worse): compTextSection treats `{#if` as an expression
    and emits syntactically broken Go: `html.EscapeString(fmt.Sprintf("%v",
    #if cond))`.
  Fix needed: lexer must tokenize `{#if}` inside quoted attribute values (or
  error clearly). Test asserts no `#if` in generated Go + `go build` succeeds.

Summary: quoted `{var}` props are FIXED (GREEN, regression tests added). The
`{#if}`-in-attribute case is OPEN (ROT) — needs a real fix, not just tests.

## Item 3b: {#if} in attributes — chosen option C (clear error)

Chosen option: **C — fail fast with a clear error** (not A/B).

Analysis of the token flow (empirically verified, ROT test reproduced both
failure modes):

- `scanTag` (core/lexer.go:197) scans a whole tag up to `>` as ONE `TokenTagOpen`
  with the raw `Attr` string. `{#if}` inside a quoted attribute value is never
  tokenized as control flow.
- Route path: `parseTemplateNode` → `NodeText` → `genTemplateNode` emits the tag
  literal → `cond` never referenced → `cond declared and not used` compile error.
- Component path: `compGen.node` → `compTextSection` (core/codegen_component.go:150)
  treats `{#if` as an expression → syntactically broken Go
  (`html.EscapeString(fmt.Sprintf("%v", #if cond))`).

Why not A (lexer emits control-flow tokens from inside attributes): the tag is
an atomic token; splitting it would require the parser to interrupt and resume
tag scanning and the codegen to split one `b.WriteString` into multiple
statements — invasive, >100 lines, high regression risk.
Why not B (parse attribute value as template): `compTextSection` concatenates
strings only (no error return, no statement context) and `extractAttrValues`
handles component-call args — both would need a parallel template-parsing path.

Fix (core/parser_section_div.go): `checkAttrControlFlow(attrs, pos)` rejects
`{#if ` / `{#each ` inside quoted attribute values with a clear error
(`"{#if} inside attribute value at position N is not supported; wrap the whole
tag in {#if} instead"`). Called from `parseTemplateNode` for `TokenTagOpen`,
`TokenComponentSelfClose`, `TokenComponentTagOpen` — the single choke point
covering routes, components, layouts, and component calls. Consistent with the
existing fail-fast pattern (`unexpected {#else} outside {#if}`).

Workaround (documented in the error message): wrap the whole tag in `{#if}` —
`{#if cond}<a class="nav active">x</a>{/if}` — which already works.

Tests:
- Unit (core/parser_section_div_test.go): `{#if}`/`{#each}` in route attr and
  component-call attr rejected; control flow around a tag still parses (NodeIf).
- Integration `_tests/core/Bugs/attr-if-in-attribute/test.sh` (rewritten):
  (1) route with `{#if}` in attr → generate fails with "inside attribute value",
  (2) component with `{#if}` in attr → same, (3) workaround (tag wrapped in
  `{#if}`) → generate + `go build` succeed, no `#if` in generated Go.

Verification: bug test GREEN, `go test ./core/... ./cmd/dreego/...` GREEN,
full suite `sh _tests/test.sh` → 177 Passed / 0 Failed.

## Item 4: props defaults — tests (ROT) + chosen design option

Tests written (no production code):

- `core/codegen_component_test.go` (unit, ROT): `TestGenerateComponentAppliesPropDefaults`
  — `Component Badge (title string = "Hi", count int = 5)` parsed via ParseHeader,
  then `GenerateComponent`. Asserts (a) parseProps captures both defaults
  (`"Hi"` / `5`), (b) the generated code applies the string default `"Hi"` and
  the int default `5` to `count`, (c) the output is valid Go. RED — generated
  code is `func Badge(title string, count int) dreego.Component { ... }` with no
  default logic at all; `p.Default` is dropped in codegen.go:324-329.
- `_tests/core/Bugs/component-prop-default/test.sh` (integration, ROT):
  `Component Card (title string = "Default Title")` + route `<@Card/>` without
  the prop. Asserts (a) `dreego generate` succeeds, (b) `"Default Title"` appears
  in `dreego/gen/components.go`, (c) `go build` succeeds (no "not enough
  arguments in call to Card"). RED — generated `func Card(title string)` and the
  call site `Card()` would not compile.

Chosen design option: **a — the default is applied inside the generated
component code** (e.g. `if title == "" { title = "Default Title" }` or a
fallback expression), so `<@Card/>` without the prop works and the call site
stays unchanged. Option b (call site fills missing props) was rejected: the
call site has no access to the component's signature defaults (codegen_template.go
and codegen_component.go only see the attrs string), so the default must live
with the component definition.

Additional finding (documented, not tested as a separate case): `parseProps`
(lexer.go:168-170) only captures `fields[3]` — a quoted default with spaces
(`title string = "Default Title"`) parses as `Default = "\"Default"` (broken).
The integration test uses `"Default Title"` in the .dreego source, but the
fix must also make parseProps handle quoted multi-word defaults (or the fix
must normalize them). The unit test pins the single-word form (`"Hi"`, `5`).

Verification: `go test ./core/...` → only `TestGenerateComponentAppliesPropDefaults`
FAILs, all other core tests GREEN. `sh _tests/core/Bugs/component-prop-default/test.sh`
→ FAIL (default missing from generated code).

## Item 5: $loop in conds — tests (GREEN, bug already fixed)

Tests written (no production code):

- `core/codegen_template_branches_test.go` (unit, GREEN):
  - `TestGenTemplateNodeEachSubstitutesLoopInIfCond` — NodeEach with a NodeIf
    child whose Cond is `!$loop.Last`; asserts generated code contains
    `if !loop.Last {` and no raw `$loop.`.
  - `TestGenTemplateNodeEachSubstitutesLoopInElseIfCond` — `$loop.` in an
    `{#else if}` cond inside an each body is substituted too.
  - `TestGenTemplateNodeEachLoopInIfCondFullParse` — full pipeline
    (lex → parse → codegen) of `{#each items as item}<div>{#if !$loop.Last},
    {/if}{item}</div>{/each}`; asserts `if !loop.Last {`.
- `core/codegen_component_branches_test.go` (unit, GREEN):
  - `TestCompGenEachSubstitutesLoopInIfCond` — same scenario through
    `genTemplateNodeComp` (component body path).
  - `TestCompGenEachLoopInIfCondFullParse` — full pipeline through the
    component codegen.
- `_tests/core/Bugs/each-loop-in-cond/test.sh` (integration, GREEN):
  route with `{#each items as item}<div>{#if !$loop.Last}, {/if}{item}</div>{/each}`;
  `dreego generate` + `go build` succeed; `dreego/gen/routes.go` contains
  `loop.Last` and no raw `$loop`.

Result: **all tests GREEN — the bug is already fixed in v0.0.25.** The
`$loop.` substitution in `genTemplateNode` (codegen_template.go:99) and
`compGen.node` (codegen_component.go:102) is applied to the ENTIRE generated
child code (`code = strings.ReplaceAll(code, "$loop.", "loop.")` after
`genTemplateNode(child)`), not only to direct expressions. A NodeIf child is
generated including its Cond, so `{#if !$loop.Last}` → `if !loop.Last {`
works. The feedback.md report predates this fix (or was based on an older
version).

Chosen design option: **a — the `$loop.` replacement applies to the whole
generated each-body code (children AND nested NodeIf conds), as it already
does.** No production change needed; the tests pin this behavior as a
regression guard. `$loop` remains a per-each-scope alias for `loop` in every
position inside the each body (expressions, if-conds, else-if-conds).

Verification: `go test ./core/... -count=1` → all GREEN (incl. 5 new tests).
`sh _tests/core/Bugs/each-loop-in-cond/test.sh` → ok. Full suite
`sh _tests/test.sh` → 178 Passed / 1 Failed, the single failure being the
pre-existing, documented `run-timer-sigterm` timing flake (isolated run:
ok; unrelated to this change).

## Item 6: duplicate title in {#head} merge — tests (ROT) + chosen semantics

Tests written (no production code):

- `core/codegen_head_dedupe_test.go` (unit, 2 ROT + 1 GREEN control):
  - `TestGenTemplHeadMergeDedupesTitle` — layout head `<title>Site</title>\n{#head}`
    + route head `<title>Page</title>...`; asserts the generated code contains
    NO "Site" (layout title dropped), the route title is emitted at the {#head}
    position, and the route head is still set into the head slot. RED — genTempl
    (codegen_template.go:223-245) writes parts[0] (layout title) then the route
    head verbatim.
  - `TestGenTemplHeadMergeDedupesMetaDescription` — same for
    `<meta name="description">`: layout desc dropped, route desc emitted. RED.
  - `TestGenTemplHeadMergeKeepsLayoutTitleWithoutRouteTitle` — control: route
    head WITHOUT a title → layout title must be kept. GREEN (dedupe must not
    remove layout content the route does not override).
- `_tests/core/Bugs/head-title-dedupe/test.sh` (integration, ROT): layout with
  `<title>Site</title>` + `<meta name="description" content="site desc">` +
  `{#head}`; route with `<title>Page</title>` + route meta desc. `dreego
  generate` + server (DREEGO_PORT) + curl. Asserts the rendered HTML contains
  EXACTLY ONE `<title>` (the route's "Page"), no "Site", exactly one
  `name="description"` (route desc), no "site desc". RED — rendered body has
  2 `<title>` (layout first, route second) and 2 meta descriptions.

Chosen semantics (option a + c combined): **the route head wins.** When the
route head defines a `<title>`, the layout `<title>` is dropped from the merged
output; the route `<title>` stays at the {#head} position. Same for
`<meta name="description">`. The route is more specific than the layout. The
route head is still set into the head slot (`c.Set("head", ...)`) unchanged —
dedupe applies only to the merged output at the {#head} position. Layout head
content the route does NOT override (e.g. charset, scripts, a title when the
route has none) is kept. Dedupe is per-tag-type: `<title>` and
`<meta name="description">` only.

Empirical confirmation (v0.0.25): generated `routes.go` contains
`b.WriteString(\`<title>Site</title>...\`)` followed by
`b.WriteString(\`<title>Page</title>...\`)` — both titles emitted, 3 `<title>`
occurrences in the file (2 in the merge + 1 in the headBuf slot). Rendered HTML
has 2 `<title>`, browser/SEO uses the first (layout "Site") — the reported bug
is real and unfixed.

Verification: `go test ./core/...` → only the 2 new dedupe tests FAIL, all
other core tests GREEN. `sh _tests/core/Bugs/head-title-dedupe/test.sh` → FAIL
(2 titles in body). Full suite `sh _tests/test.sh` → 179 Passed / 2 Failed,
the 2 failures being exactly the new tests (no regressions).

## Item 7: 404/500 without layout — tests (ROT) + chosen semantics

Tests written (no production code):

- `core/codegen_error_test.go` (unit, 1 ROT + 2 GREEN control):
  - `TestGenerateErrorHandlerScopeDivNotBeforeDoctype` — GenerateErrorHandler
    with a template starting `<!doctype html>` must NOT emit `data-scope` at
    all: a scope div before the doctype puts the document into quirks mode.
    RED — codegen_error.go:39 writes the scope div unconditionally before the
    template nodes, so the generated code is
    `b.WriteString("<div data-scope=\"...\">")` followed by the doctype.
  - `TestGenerateErrorHandlerHeadBeforeScopeDiv` (control, GREEN) — a file
    with a `<head>` section and a template WITHOUT doctype: head content is
    emitted before the scope div (pins current ordering, no regression).
  - `TestGenerateErrorHandlerScopeDivKeptWithoutDoctype` (control, GREEN) —
    a plain template without doctype keeps the scope div (the fix must NOT
    remove scoping from non-doctype error pages).
- `_tests/core/Bugs/error-page-doctype/test.sh` (integration, ROT):
  404.dreego with `<!doctype html>` + head + body. `dreego generate` + server
  (DREEGO_PORT) + curl on a non-existent route. Asserts (a) HTTP 404 status,
  (b) the body starts with `<!doctype html>` (NOT `<div data-scope=`),
  (c) no `data-scope` anywhere in the body, (d) head + body content present.
  RED — rendered body is `<div data-scope="..."><!doctype html>...`.
- `_tests/core/Bugs/error-page-layout/test.sh` (integration, ROT): layout
  (`<title>Layout Site</title>` + `{#head}` + `{#slot}`) AND 404.dreego with
  doctype + head (incl. CSS link) + body + `<style>` + `<script>` sections.
  Asserts (a) HTTP 404, (b) body starts with `<!doctype html>`, (c) no
  `data-scope`, (d) NO layout content ("Layout Site" must NOT appear — the
  error page is self-contained, not wrapped in the layout), (e) the error
  page's own head CSS link, style section (scoped), and script section are
  all rendered. RED — scope div before doctype; layout title check cannot
  even be reached.

Chosen semantics (option a, pragmatic v0.0.25 scope): **the scope div is
suppressed when the error template starts with `<!` (doctype/comment); the
layout is NOT applied to error pages.** Error pages are self-contained
documents: the template author must include the full `<!doctype html><html>…
</html>` skeleton (including head, CSS links, and scripts) themselves.
Layout wrapping of error pages is a documented limitation / future feature
(v0.0.26 candidate) — the test pins that layout content must NOT leak into
the error page today, so adding layout support later is a deliberate,
observable change. The `<head>` section (HeadSection) is still rendered
BEFORE the template body — no scope div can precede a doctype either way.
Normal (non-error) routes are unaffected: genTempl keeps its existing scope
div behavior (isGET branch).

Fix sketch for the coder (production change, NOT done here):
codegen_error.go — before emitting `b.WriteString("<div data-scope=…>")`,
inspect the first template node: if `file.Template.Nodes[0]` is NodeText
with content starting `<!` (or any node whose generated output starts with
`<!`), skip the scope div AND the closing `</div>`. The head section
(genHead) stays where it is (before the scope div).

Verification: `go test ./core/...` → only
`TestGenerateErrorHandlerScopeDivNotBeforeDoctype` FAILs, all other core
tests GREEN (incl. the 2 new control tests). Both integration tests FAIL at
the body-prefix check (scope div before doctype); HTTP status 404 already
GREEN. Filtered runner `DREEGO_FILTER=error-page sh _tests/test.sh` → 0
Passed / 3 Failed, the 3 failures being exactly the new tests (no
regressions).

## Item 8: parser text-before-section — tests (ROT) + chosen option b (fail fast)

Tests written (no production code):

- `core/parser_test.go` (unit, 1 GREEN + 2 ROT + 1 GREEN control):
  - `TestParseTextBeforeGoSectionSwallowed` (GREEN, pins the CURRENT bug):
    `<!doctype html>\n<go>msg := "hi"</go>\n<div><p>x</p></div>` parses
    WITHOUT error, `file.Go` is empty, and the `<go>` body lands as
    `NodeText` in the template. Documents the silent-corrupt behavior that
    must be replaced.
  - `TestParseTextBeforeGoSectionFailsFast` (ROT): the same input must
    produce an error containing "must come before template content".
  - `TestParseTextBeforeDivSectionFailsFast` (ROT): `plain text\n<div>…`
    must fail the same way (any section after leading text).
  - `TestParseGoSectionFirstStillWorks` (GREEN control): `<go>` at the file
    start still parses as Go code (`file.Go[0].Code == "msg := \"hi\""`).
- `_tests/core/Bugs/text-before-section/test.sh` (integration, ROT):
  (1) route starting `<!doctype html>` + `<go>` + `<div>` → `dreego generate`
  must FAIL with "must come before template content", (2) control: `<go>`
  first → generate + `go build` succeed and `dreego/gen/routes.go` contains
  NO `msg := "hi"` as template text.

Root cause (empirically verified): `<!doctype html>` lexes as `TokenTagOpen`
with tag `!doctype` (lexer.go:247-262), which is not a section tag, so the
lexer does not push the section stack. In `Parse` (parser.go:32-42) the
first token is a `TokenTagOpen` with an unknown tag → the `default` branch
(parser.go:89-98) calls `parsePlainTemplate()` because `file.Template ==
nil`, and `parsePlainTemplate` (parser.go:123-136) consumes ALL remaining
tokens as template nodes — `<go>`, `</go>`, `<div>`, `</div>` all become
`NodeText`. `file.Go` stays empty; the Go code is emitted via
`b.WriteString("msg := \"hi\"")` into the HTML. Silent corrupt, no error.

Chosen option: **b — fail fast with a clear error** (not a). Option a
(recognize sections after leading text) would require `parsePlainTemplate`
to stop at section tags and hand control back to the main loop, plus a new
AST field for the pre-section template text (there is no "template prefix"
slot today) — invasive, touches lexer/parser/AST/codegen, high regression
risk for v0.0.25. Option b is a small change in the `default` branch of
`Parse`: when the first non-section token is template text and a later
token is a section tag (`go`/`div`/`head`/`script`/`style`), return
`<go> section must come before template content at position N` (or the
matching tag name). Consistent with the existing fail-fast pattern
(`expected section tag, got Text`, `unknown section <x>`).

Fix sketch for the coder (production change, NOT done here): in
`parsePlainTemplate` (or in the `default` branch before calling it), scan
the remaining tokens for a `TokenTagOpen` whose tag is a section tag; if
found, return an error naming the section tag and its position. The error
message must contain "must come before template content" (pinned by the
tests).

Verification: `go test ./core/ -count=1` → only the 2 new FailsFast tests
FAIL, all other core tests GREEN (incl. the 2 new GREEN tests). Filtered
runner `DREEGO_FILTER=text-before-section sh _tests/test.sh` → 0 Passed /
2 Failed, the 2 failures being exactly the new tests (no regressions).

## Item 9: < in go strings — tests (ROT) + chosen option b (raw reconstruction)

Tests written (no production code):

- `core/parser_test.go` (unit, 2 ROT + 1 GREEN control):
  - `TestParseGoSectionKeepsLtInQuotedString` (ROT): `<go>msg := "TO:
    <HASH>"</go>` must parse to `file.Go[0].Code == "msg := \"TO:
    <HASH>\""`. RED — currently `msg := "TO: "` (`<HASH>` silently dropped).
  - `TestParseGoSectionKeepsLtInBacktickString` (ROT): `<go>svg :=
    `<svg viewBox="0 0 24 24"></svg>`</go>` must parse to the full backtick
    literal. RED — currently `svg := `` ` (the `<svg>...</svg>` content is
    dropped, empty string).
  - `TestLexGoSectionLtLexesAsTags` (GREEN control, documents lexer
    behavior): the lexer has NO Go-string awareness — `"TO: <HASH>"` inside
    a `<go>` section lexes as `TagOpen(go), TagOpen(HASH), TagClose(go)`.
    Pins that the fix lives in the parser, not the lexer (option b).
- `_tests/core/Bugs/go-string-lt/test.sh` (integration, ROT): route with
  `<go>` containing `msg := "TO: <HASH>"` and
  `` svg := `<svg viewBox="0 0 24 24"><path d="M12 2"/></svg>` ``.
  `dreego generate` + `go build` must succeed; `dreego/gen/routes.go` must
  contain `TO: <HASH>` and `<svg viewBox="0 0 24 24">` (splitGoSections puts
  non-declaration code inline in routes.go); the go body must NOT leak as
  template text into `dreego/gen/dree.go`. RED — `<HASH>` missing from
  routes.go.

Root cause (empirically verified): the lexer (`Lex`, core/lexer.go:14-75)
has no Go-string awareness — every `<` inside a `<go>` section is handed to
`scanTag`. `"TO: <HASH>"` lexes as `TokenText("TO: ")` + `TokenTagOpen(HASH)`
+ `TokenText("\"")`; the backtick SVG lexes as `TokenText("svg := `")` +
`TokenTagOpen(svg)` + `TokenText(" viewBox=...")` + `TokenTagClose(svg)` +
`TokenText("`")`. `parseGoSection` (parser_section_go.go:8-26) only keeps
`TokenText` and silently DISCARDS `TokenTagOpen`/`TokenTagClose` → the `<...>`
content vanishes without any error. Silent data corruption — the generated
Go still compiles (it just holds the wrong string), which is why the bug
survived: `\u003c` escapes and backtick-string avoidance were the only
workarounds.

Chosen design option: **b — the go-section scanner reconstructs the raw
content up to `</go>` (like parseNonDivSection), the lexer stays as-is.**
Rationale: (1) `parseNonDivSection` (parser_section_go.go:28-63) already
does exactly this for head/script/style — it rebuilds tags from tokens
(`"<" + tok.Tag + attrs + ">"` and `"</" + tok.Tag + ">"`) — so option b is
the established codebase pattern, ~10 lines in `parseGoSection`. (2) Option
a (Go-string-aware lexer) would require the lexer to track quote/backtick
state across token boundaries inside `<go>` sections and would change token
streams for head/script/style too — invasive, high regression risk for
v0.0.25. (3) Option b also fixes `<` in raw Go code that is not a string
(e.g. a generic-less comparison in inline code), because `</go>` is the only
thing the scanner needs to recognize — and `</go>` cannot appear inside a Go
string without breaking the Go code anyway. Known edge (documented, not a
regression): a Go string containing a literal `</go>` cannot be expressed —
it would close the section; the fix does not change that.

Fix sketch for the coder (production change, NOT done here): in
`parseGoSection`, mirror `parseNonDivSection` — on `TokenTagOpen` write
`"<" + tag + (attrs) + ">"`, on `TokenTagClose` write `"</" + tag + ">"`,
keep `TokenText` as-is, stop at `TokenTagClose` with `Tag == "go"`. The
`</go>` token is consumed (like the current code). No lexer change.

Verification: `go test ./core/ -count=1` → only the 2 new ROT tests FAIL,
all other core tests GREEN (incl. the new GREEN lexer test). `sh
_tests/core/Bugs/go-string-lt/test.sh` → FAIL at the `<HASH>` check (missing
from routes.go; no regressions elsewhere).

## Item 11: --version flag — tests (ROT) + chosen option a (--version AND -v)

Tests written (no production code):

- `_tests/core/CLI/version-flag/test.sh` (integration, ROT): runs the CLI
  with `--version`, `-v`, and the `version` subcommand. Asserts for each:
  (a) exit 0, (b) non-empty output, (c) output is NOT the build-info
  placeholder `(devel)`. Then asserts `--version` and `-v` produce EXACTLY
  the same output as the `version` subcommand (all three must be
  interchangeable). RED — `--version` falls into the `default` branch of the
  main() switch (cmd/dreego/main.go:45-48): `unknown command: --version` +
  help + exit 1. Verified: `sh _tests/core/CLI/version-flag/test.sh` → FAIL
  at the first check ("dreego --version exited non-zero").

Chosen design option: **a — `--version` AND `-v` both supported, same
output as the `version` subcommand.** Consistent with the existing
`help`/`--help`/`-h` trio (main.go:43) which already accepts both the bare
word and the flag forms. Fix sketch for the coder (production change, NOT
done here): extend the switch — `case "version", "--version", "-v":
fmt.Println(dreegoVersion())`. No other change needed; the help text already
documents `version` (the flag forms are conventional aliases, no help-text
change required).

Verification: `sh _tests/core/CLI/version-flag/test.sh` → FAIL (ROT, as
expected). After the fix: all three invocations must print the same version
string and exit 0.
