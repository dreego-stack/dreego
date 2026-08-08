---
type: Task
title: v0.0.24 Plan — Feedback-Driven + Planned Blocks
status: done
assign: manager
---

# v0.0.24 Plan — Full Test-Driven Rollout

Goal: Redesign v0.0.24 around the fresh-start feedback, keep the originally planned blocks that are already testable, and cover every new behavior with ≥ 3 tests each (integration `_tests/core/` and/or unit `*_test.go`). Execute strictly sequentially: one failing test → one coder fix → review → commit → next test. No parallel agents, no premature mass test commits.

## 1. v0.0.24 Block Composition

$%$ Block $%$ Source $%$ What it does $%$ Tests $%$
$%$-------$%$--------$%$--------------$%$-------$%$
$%$ **layout-head.1** $%$ feedback $%$ Layouts apply (`{#slot}`/`{#head}`); route `<head>` works without layout $%$ 4 int + 2 unit $%$
$%$ **scoped-css.2** $%$ feedback $%$ `scopeCSS` preserves declarations in `{}` $%$ 3 int + 3 unit $%$
$%$ **component-attr-props.1** $%$ feedback $%$ `{prop}` substituted in HTML attributes $%$ 4 int + 3 unit $%$
$%$ **scaffold-fix.1** $%$ feedback $%$ `go.sum` in scaffold, `.gitignore` ignores only `dreego/gen/` $%$ 3 int $%$
$%$ **typed-forms.1** $%$ _todo/plan.md (was v0.0.23) $%$ int/bool/slice binding + custom validators + better email $%$ 5 unit + 2 int $%$
$%$ **dreegotest.1** $%$ _todo/plan.md (was v0.0.23) $%$ exported `dreegotest` helper package for route/component unit tests $%$ 3 unit + 2 int $%$
$%$ **golden-tests-core.1** $%$ _todo/plan.md (was v0.0.23) $%$ golden-file assertions for generated `gen/routes.go`/`gen/components.go` $%$ 3 unit $%$

Deferred out of v0.0.24:
- `frontmatter.1`, `dev-server.1`, `docs-extensibility.1` → v0.0.25
- DSGVO asset self-hosting → roadmap / v0.0.30 area
- Cross-build docs → docs backlog for v0.0.25

## 2. Detailed Test List per Block

### layout-head.1

Integration (`_tests/core/`):
1. `Bugs/layout-not-applied` — route with layout renders inside layout `{#slot}`.
2. `Bugs/route-head-without-layout` — `<head>` in route appears when no layout exists.
3. `Bugs/layout-route-head-merge` — route `<head>` merged into layout `{#head}`.
4. `Layout/no-layout` — route without layout file still renders as full fragment.

Unit (`core/*_test.go`):
5. `core/codegen_test.go` or `core/codegen_layout_test.go` — `genTempl` emits layout wrapping code.
6. `core/codegen_head_test.go` — route head is emitted standalone when layout is nil.

### scoped-css.2

Integration:
1. `Bugs/scoped-style-declarations-lost` — `radial-gradient(...)` survives in generated style.
2. `Bugs/scoped-style-comma-parens` — selectors with commas and nested parens preserved.
3. `Bugs/scoped-style-keyframes` — `@keyframes` body is preserved.

Unit:
4. `core/codegen_helpers_test.go` — `scopeCSS` keeps declarations with braces.
5. `core/codegen_helpers_test.go` — `@media` inner selectors still scoped.
6. `core/codegen_helpers_test.go` — pseudo-selectors like `:hover` scoped correctly.

### component-attr-props.1

Integration:
1. `Bugs/component-attr-prop-substitution` — `<a href="{url}">` resolves prop.
2. `Components/prop-expression` — component prop from expression `title={user.Name}` (already listed in _docs/testing.md as planned).
3. `Components/multi-props` — 3+ props passed and rendered.
4. `Components/empty-props` — component without props.

Unit:
5. `core/codegen_component_test.go` — `genComponentCall` emits attribute expression.
6. `core/codegen_helpers_test.go` — `extractAttrValues` handles `{expr}` in quoted value.
7. `core/codegen_template_test.go` — attribute expression escapes value.

### scaffold-fix.1

Integration:
1. `CLI/new-go-sum` — fresh `dreego new` scaffold builds without `go mod tidy`.
2. `CLI/new-gitignore` — `.gitignore` does not ignore `dreego/routes/`, `dreego/components/`, `dreego/config.json`.
3. `CLI/new-layout-exists` — scaffold either contains a working layout or no `layouts/` dir at all.

### typed-forms.1

Unit:
1. `core/validate_test.go` — `BindForm` binds `int` field.
2. `core/validate_test.go` — `BindForm` binds `bool` field.
3. `core/validate_test.go` — `BindForm` binds `[]string` slice.
4. `core/validate_test.go` — custom validator registered and applied.
5. `core/validate_test.go` — `email` rejects missing dot / missing at.

Integration:
6. `FormActions/form-int-binding` — POST binds integer field and validates min.
7. `FormActions/form-bool-binding` — POST binds checkbox to bool.

### dreegotest.1

Unit:
1. `dreegotest/request.go` — simulate GET request and read HTML body.
2. `dreegotest/request.go` — simulate POST form and read validation errors.
3. `dreegotest/render.go` — render a single component with props.

Integration:
4. `_tests/core/dreegotest/basic` — package used in an integration test.
5. `_tests/core/dreegotest/form` — package validates form submission.

### golden-tests-core.1

Unit (`core/codegen_golden_test.go`):
1. Golden for a simple route with head+div.
2. Golden for a component with props and style.
3. Golden for a route with layout `{#slot}`/`{#head}`.

## 3. Additional Test Gaps (from _docs/testing.md and inspection)

The following are not tied to a v0.0.24 feature but increase coverage; add if time remains after the seven blocks above:

- `Routing/put-method` (planned in _docs/testing.md)
- `Routing/delete-method` ✅ already exists — verify coverage
- `Routing/multi-segment`
- `Routing/bracket-in-name`
- `Routing/deep-nesting`
- `Session/delete`, `Session/destroy`, `Session/no-store` (planned)
- `Middleware/csrf-post-fail`, `csrf-post-pass`, `csrf-disabled` (planned)
- `Config/redirect`, `rewrite`, `logging-off`, `invalid-json` (planned)
- `CLI/check-stale`, `no-args`, `unknown-cmd` (planned)
- `Template/nested-if`, `if-else`, `each-empty`, `each-with-if`, `expression-function` (planned)

Decision: only add these **after** all seven v0.0.24 blocks are green and committed, to keep the cycle focused.

## 4. Execution Mode

- Strictly sequential.
- One agent at a time.
- Per increment:
  1. Manager updates this file with current target.
  2. `coder-testN` writes **only** the failing test(s) for one block.
  3. `smd sh _tests/.../test.sh` (or `smd go test ./core/...`) confirms RED.
  4. `coder-writeN` fixes the minimal production code to make that test GREEN.
  5. `smd go test ./core/...` + `smd sh _tests/test.sh` (with filter) confirms GREEN.
  6. `reviewer` reviews the increment.
  7. `git` agent commits (no push).
  8. Manager marks increment done and starts next.

Agent naming convention:
- `coder-test1`, `coder-test2`, … for test writing.
- `coder-write1`, `coder-write2`, … for production code.
- `review1`, `review2`, … for reviews.
- `git1`, `git2`, … for commits.

## 5. Order of Increments

1. scaffold-fix.1 tests → fix → commit
2. layout-head.1 tests → fix → commit
3. scoped-css.2 tests → fix → commit
4. component-attr-props.1 tests → fix → commit
5. typed-forms.1 tests → fix → commit
6. dreegotest.1 tests → implementation → commit
7. golden-tests-core.1 tests → implementation → commit
8. Optional additional tests from section 3 if time allows.
9. Final full suite + docs update + version bump to v0.0.24 + final commit.

## 6. Non-Code Deliverables (after all tests green)

- `CHANGELOG.md` v0.0.24 entry.
- `.agents/log.md` update.
- `_docs/layouts.md` explaining `{#slot}`/`{#head}` and route `<head>` behavior.
- `_docs/components.md` update for attribute props + scoped CSS behavior.
- `_docs/deployment.md` cross-build note.
- `_docs/offline.md` or `_docs/cli.md` note on version pinning.
- `TODO.md` and `_todo/plan.md` update to reflect deferred blocks.

## 7. Status Log

- [x] Feedback placed under `.tmp/`
- [x] Initial plan drafted
- [x] Expanded to full v0.0.24 plan with ≥ 3 tests per block
- [x] User approved plan
- [x] Increment 1 started
- [x] scaffold-fix.1 tests written (all RED against v0.0.23)
  - `_tests/core/CLI/new-go-sum/test.sh` — RED (missing `go.sum`)
  - `_tests/core/CLI/new-gitignore/test.sh` — RED (`.gitignore` pattern `/dreego` ignores source dir)
  - `_tests/core/CLI/new-layout-exists/test.sh` — RED (layout exists but generated route has no `<html>`)
- [x] scaffold-fix.1 production code fix (version fallback)
  - `cmd/dreego/version.go` — added `VERSION` file fallback
  - `_tests/core/CLI/new-go-sum/test.sh` — GREEN
  - `_tests/core/CLI/new-gitignore/test.sh` — GREEN (test fixed: case-pattern matched literal heredoc lines as well as file contents; replaced with direct `grep -Eqx` against stripped gitignore lines)
  - `_tests/core/CLI/new-layout-exists/test.sh` — GREEN
  - `go test ./core/... ./cmd/dreego/...` — GREEN
- [x] layout-head.1 tests: DREEGO_BIN fallback added, all four run standalone
  - Added `DREEGO_BIN` build fallback (pattern from `CLI/new-go-sum/test.sh`) after `cd "$workdir"` in all four tests.
  - Fixed `Bugs/layout-not-applied/test.sh` test-construction bug: added missing `main.go` (it previously failed at final `go build .` with "no Go files", not because of layout).
  - **FINDING:** All four layout tests are GREEN against the current code. The layout feature (`{#slot}`, `{#head}`, route `<head>` without layout) is ALREADY implemented in `core/codegen.go` (`genTempl`) and `core/codegen_layout.go` / `core/codegen_head.go`, committed since v0.0.22 (`11d33d2`). No production code change is needed for layout-head.1 — the bug reported in the original feedback is already fixed.
  - Tests verified GREEN standalone:
    - `_tests/core/Bugs/layout-not-applied/test.sh` — ok
    - `_tests/core/Bugs/route-head-without-layout/test.sh` — ok
    - `_tests/core/Bugs/layout-route-head-merge/test.sh` — ok
    - `_tests/core/Layout/no-layout/test.sh` — ok
- [x] layout-head.1 unit tests written (GREEN)
  - `core/codegen_layout_test.go` — `TestGenTemplEmitsLayoutWrapping`, `TestGenTemplLayoutSlotNodeUsesSlot` (genTempl emits `c.Set("slot", pageContent)` + layout renders `c.Get("slot")`).
  - `core/codegen_head_test.go` — `TestGenTemplEmitsHeadStandaloneWithoutLayout`, `TestGenTemplNoHeadWithoutLayout` (route `<head>` emitted standalone when layout is nil).
  - `smd go test ./core/... ./cmd/dreego/...` — GREEN, no flakes.
- [x] scoped-css.2 tests written (increment 3, coder-test3) — all RED against current code
  - Integration `_tests/core/Bugs/scoped-style-declarations-lost/test.sh` — RED (`radial-gradient` dropped)
  - Integration `_tests/core/Bugs/scoped-style-comma-parens/test.sh` — RED (`calc()`/`rgb()` declarations dropped)
  - Integration `_tests/core/Bugs/scoped-style-keyframes/test.sh` — RED (`@keyframes` body dropped)
  - Unit `core/codegen_helpers_test.go` — `TestScopeCSSKeepsDeclarationsWithBraces`, `TestScopeCSSMediaPreservesDeclarationsAndScopesInnerSelectors`, `TestScopeCSSPseudoSelectorKeepsDeclaration` — all RED
  - Root cause: `scopeCSS` in `core/codegen_helpers.go` never copies declaration text between `{` and `}`, so all declarations collapse to `{}`. Only selectors/braces/at-rule names survive.
- [x] layout-head.1 review fixes applied (increment 2, layout-head.1)
  - **Blocker (`sed -i`)**: Removed `sed -i "s/8080/$port/" main.go` from `Bugs/route-head-without-layout`, `Bugs/layout-route-head-merge`, `Layout/no-layout`. Port now written directly into `main.go` via unquoted heredoc with `core.Listen(":$port")` (matches preferred reviewer fix).
  - **Warning (negativcheck)**: `core/codegen_head_test.go:47` — changed `strings.Contains(out, "WriteString(\"<title>")` → backtick pattern `WriteString(\`<title>` to match genHead emission.
  - **`apk add curl || true`**: Removed `|| true` in all three server-based tests so a missing curl fails the test.
  - ⚠️ **VERIFICATION BLOCKED**: Docker/smd daemon unresponsive (all `docker`/`smd` calls time out). `go test` and the four integration tests could NOT be re-run. Fixes are purely test-file edits (no production code) and are logically sound; must be verified when Docker recovers.
- [x] scoped-css.2 production code fix (increment 3, coder-write3)  - **Root cause**: `scopeCSS` copied only selectors and `{`/`}`, never the declaration body between braces → all declarations collapsed to `{}`.
  - **Fix** (`core/codegen_helpers.go`): rewrote `scopeCSS` as a recursive brace-tracked parser:
    - `scopeRange` scans for `{`, matches its closing `}` with `matchBrace` (depth-tracking, ignores commas/parens/brackets inside values), and emits `header{ ... }`.
    - Normal selectors (`sel != "" && !@`): body copied verbatim via `scopeSelector` (top-level comma split + `[data-scope=hash] ` prefix per selector, joined with `,\n` so each scoped selector is on its own line).
    - `@media ...` (or any non-`@keyframes` at-rule): header unscoped, inner block recursively `scopeRange`'d → inner selectors get scoped, their declarations preserved.
    - `@keyframes name`: header and full body copied verbatim (unscoped) → `from`/`to`/percent steps stay unscoped, declarations preserved.
    - Declaration values (`radial-gradient(...)`, `calc(...)`, `rgb(1, 2, 3)`) preserved because the entire body between braces is copied byte-for-byte.
  - **Verification** (all GREEN):
    - `smd go test ./core/...` — ok
    - `smd go build ./core/... && go vet ./core/...` — ok
    - `_tests/core/Bugs/scoped-style-declarations-lost/test.sh` — ok (`radial-gradient(circle, #ccfbf1 1px, transparent 1px)` preserved)
    - `_tests/core/Bugs/scoped-style-comma-parens/test.sh` — ok (`calc(100% - 20px)`, `rgb(1, 2, 3)`, both comma selectors scoped → 2 `[data-scope=` lines)
    - `_tests/core/Bugs/scoped-style-keyframes/test.sh` — ok (`@keyframes spin` header + `from`/`to` unscoped, body `rotate(360deg)` preserved, `.spinner` scoped)
    - `_tests/core/Bugs/scoped-css-media/test.sh` — ok (with `DREEGO_BIN` set; the pre-existing test has no local `DREEGO_BIN` fallback but `make test`/`test.sh` provides it)
    - Full suite `_tests/test.sh`: 156 Passed / 1-2 Failed — failures are pre-existing parallel server-port flakes in `Bugs/layout-route-head-merge` / `Bugs/route-head-without-layout` (different test fails on each run, both contain no `<style>`); both pass standalone. No CSS regression.

- [x] component-attr-props.1 production code fix (increment 4, coder-write4)
  - **Fix** `core/codegen_component.go`:
    - `genTemplateNodeComp` `NodeText` case now routes through new `compTextWithAttrs`, which scans the text for `{...}` inside quoted attribute values (`inQuote` + brace-depth tracking) and emits `html.EscapeString(fmt.Sprintf("%v", expr))` for each placeholder, concatenated with `+` literal parts. Text-content `{...}` (unquoted) stays literal as before → route-side `{label}` behavior unchanged.
    - `genComponentCall` self-close path now passes `n.Attrs` via `extractAttrValues` (was raw) → nested component calls resolve `href={url}` / `title={user.Name}` prop expressions.
  - **Fix** `core/codegen_helpers.go`: `attrVal` now also resolves `{expr}` when it appears inside a quoted value (`val[0]=='{' && val[len-1]=='}'` after trimming quotes), so `href="{url}"` becomes the Go expression `url` (route side). Previously only bare `{url}` (unquoted) was resolved.
  - **Verification** (all GREEN):
    - `smd go test ./core/...` — ok
    - `smd go build ./core/... && go vet ./core/...` — ok
    - `smd go test ./cmd/dreego/...` — ok
    - `_tests/core/Bugs/component-attr-prop-substitution/test.sh` — ok (`{url}` resolved + `EscapeString` present, `go build` passes)
    - `_tests/core/Components/prop-expression/test.sh` — ok (route-side `title={user.Name}` still passes via `extractAttrValues`)
    - `_tests/core/Components/multi-props/test.sh` — ok (3 props, incl. `href="mailto:{email}"`)
    - `_tests/core/Components/empty-props/test.sh` — ok
    - Full suite `_tests/test.sh`: **160 Passed / 0 Failed** — no regressions
- [ ] component-attr-props.1 tests written (increment 4, coder-test4) — **RED** against current code
  - **Root cause (confirmed):** `<a href="{url}">` in a component body is lexed as a single `TokenText` node (`scanTag` generic `<` case, `core/lexer.go:245-260`). `genTemplateNodeComp` emits `NodeText` verbatim, so `{url}` in the attribute stays literal. Only text-content `{label}` is split out and substituted. Additionally, nested component calls (`genComponentCall`, `core/codegen_component.go:120-127`) pass `n.Attrs` through raw (not `extractAttrValues`), and `attrVal` (`core/codegen_helpers.go:158-172`) trims quotes from `href="{url}"` returning the literal string `"{url}"`.
  - **RED unit tests:**
    - `core/codegen_component_test.go` — `TestGenComponentCallResolvesAttrExpression` (genComponentCall leaves `{url}` literal)
    - `core/codegen_helpers_test.go` — `TestExtractAttrValuesResolvesExprInQuotedValue` (returns literal `"{url}"`)
    - `core/codegen_template_test.go` — `TestGenTemplateNodeCompAttrExpressionEscapesValue` (emits `WriteString(`<a href="{url}">`)` with no `html.EscapeString`)
  - **RED integration:** `_tests/core/Bugs/component-attr-prop-substitution/test.sh` (`{url}` emitted literally)
  - **GREEN integration (already working):**
    - `_tests/core/Components/prop-expression/test.sh` (route-side `title={user.Name}` already resolved via `extractAttrValues` in `genTemplateNode`)
    - `_tests/core/Components/multi-props/test.sh` (3+ props render)
    - `_tests/core/Components/empty-props/test.sh` (existing, added missing `DREEGO_BIN` standalone fallback; GREEN under runner and standalone)
  - Note: The plan's unit test #7 (attr escaping) was written against the component-body parse path so it reflects the real bug; the route-side attr path already escapes nothing because values are resolved, and the fix must escape the resolved expression.

- [ ] typed-forms.1 tests written (increment 5, coder-test5) — **RED** against current code
  - **Analysis of current state (`core/validate.go`, `core/forms.go`, `core/validate_test.go`):**
    - `BindForm` (`core/validate.go:31-32`) only accepts `reflect.String` fields; any `int`/`bool`/`slice` returns `fmt.Errorf("unsupported field type ...")` → **int/bool/slice binding all RED**.
    - `email` rule (`core/validate.go:78-81`) already rejects missing `@` OR missing `.` → **email tests already GREEN (regression)**.
    - No custom validator registration API exists (no `RegisterRule`/dispatch in `applyRule`) → **custom validator RED (compile error: `undefined: RegisterRule`)**.
    - No FormActions integration tests for int/bool binding exist → **RED**.
  - **Unit tests added (`core/validate_typed_test.go` — new file, validate_test.go was 535 lines):**
    - `TestBindFormIntField` — RED (BindForm errors on int field)
    - `TestBindFormBoolFieldOn` — RED (BindForm errors on bool field)
    - `TestBindFormBoolFieldAbsent` — RED (empty value skipped, bool stays false — no error, but no binding path)
    - `TestBindFormSliceField` — RED (BindForm errors on []string field)
    - `TestRegisterRuleCustom` — RED (compile: undefined RegisterRule)
    - `TestRegisterRuleCustomDoesNotBreakBuiltins` — RED (regression guard, same compile block)
    - `TestEmailRejectsMissingAt` / `TestEmailRejectsMissingDot` / `TestEmailAcceptsValid` — GREEN by inspection (applyRule already correct)
  - **Integration tests added:**
    - `_tests/core/FormActions/form-int-binding/test.sh` — RED (`age=20` returns 200 instead of 303; int binding fails). Uses port schema + `core.Listen(":$port")` directly + `DREEGO_BIN` fallback.
    - `_tests/core/FormActions/form-bool-binding/test.sh` — RED (`subscribe=on` returns 200 instead of 303; bool binding fails). Same port-schema pattern.
  - Note: `go test ./core/...` fails at **compile** time (undefined `RegisterRule`), so the whole core test package is blocked until `RegisterRule` exists. This is expected — the custom-validator test forces the API to exist.
  - **Verification blocked**: `go vet`/`go test ./core/...` cannot compile; the RED unit tests for int/bool are confirmed by code inspection (BindForm returns error for non-string) and the integration tests confirm RED end-to-end.

- [x] typed-forms.1 production code fix (increment 5, coder-write5) — all GREEN
  - **Fix `core/validate.go`:**
    - Added `RegisterRule(name string, fn func(string) string)` + `customRules map[string]validatorFunc`. `applyRule` dispatches registered rules in a new `default` case (built-in required/email/min/max untouched). `validatorFunc` type alias for the signature.
    - Extended `BindForm` beyond `reflect.String`:
      - `int` → `strconv.Atoi`, error on invalid input.
      - `bool` → `SetBool(len(values)>0 && values[0]=="on")`; absent/empty → false, no error.
      - `[]string` → collects all form values, skips empties.
      - String path keeps the `val==""` skip (existing behavior preserved).
      - Any other kind (map, struct, etc.) still returns `unsupported field type` error.
    - `ValidateForm` now uses `fmt.Sprint(v.Field(i).Interface())` instead of `v.Field(i).String()` so `min`/`max`/`required` rules work on bound `int` values (e.g. `Age int validate:"min=2"`).
  - **Test updates (obsolete behavior that directly contradicted the new feature):**
    - `core/validate_test.go` `TestBindFormNonStringFieldReturnsError`: struct field changed from `Count int` → `Count map[string]string` (still exercises the unsupported-type error path; int is now supported).
    - `_tests/core/Bugs/bindform-non-string/test.sh`: `Profile.Age int`/`Profile.Admin bool` previously asserted error; now asserts int/bool bind correctly AND map field still returns `unsupported field type` error. `admin=true` → `admin=on` (checkbox convention).
  - **Verification (all GREEN):**
    - `smd go test ./core/...` — ok (incl. new validate_typed_test.go)
    - `smd go vet ./core/...` — ok
    - `smd go test ./cmd/dreego/...` — ok
    - `_tests/core/FormActions/form-int-binding/test.sh` — ok (303 on valid, /adult redirect, 200 re-render when min fails)
    - `_tests/core/FormActions/form-bool-binding/test.sh` — ok (/subscribed when checked, /skipped when absent)
    - `_tests/core/Bugs/bindform-non-string/test.sh` — ok
    - Full suite `_tests/test.sh`: **162 Passed / 0 Failed** (two consecutive runs). Earlier single `run-timer-sigterm` failure was a transient flake — SIGTERM/server test unrelated to form changes; passes in subsequent full runs.
  - Note: `core/validate.go` grew 134→178 lines (under 300 limit); `forms.go`/`splitGoSections` untouched.

- [x] dreegotest.1 production code implemented (increment 6, coder-write6) — all GREEN
  - **New module `dreegotest/`** with `request.go` + `render.go`, implementing the public API defined by coder-test6:
    - `Get(t, path) *Response` — `httptest.NewRequest(GET)` → `core.ServeMux()` → `httptest.NewRecorder`; returns `*Response{StatusCode, Body}`. Tests pass status 200/body and 404.
    - `PostForm(t, path, form url.Values) *Response` — `httptest.NewRequest(POST, body=form.Encode())` with `Content-Type: application/x-www-form-urlencoded` → `core.ServeMux()`. `r.FormValue` binding + Content-Type both verified by tests.
    - `RenderComponent(t, fn core.ComponentFunc, props ...any) string` — `core.NewSSR(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))`; alternating key/value props written via `ctx.Set(key, html.EscapeString(fmt.Sprintf("%v", value)))`. Tests verify rendered HTML and XSS escaping.
  - **Analysis (no test changes needed):** `ServeMux()` calls `Build()` which is idempotent (early-returns if `builtHandler != nil`). `core.Reset()` nils `builtHandler` for a deterministic rebuild. The tests call `core.Reset()` + `core.Register(...)` directly, so no `dreego generate`/generated routes are required — this is a genuine unit test of the routing pipeline with no file-system dependency.
  - **Design note:** XSS escaping happens at prop-injection time (`ctx.Set`), not on the whole component output, so static markup (`<h1>Welcome</h1>`) stays untouched while `ctx.Get("name")` returns an already-escaped value → `<script>` → `&lt;script&gt;`. This matches the XSS test contract.
  - **Verification (all GREEN):**
    - `smd go test ./dreegotest/...` — ok (6 tests)
    - `smd go build ./core/... ./cmd/dreego/... ./plugins/sample/... ./dreegotest/...` — ok
    - `smd go vet ./dreegotest/...` — ok
    - `smd go test ./core/... ./cmd/dreego/... ./plugins/sample/... ./dreegotest/...` — ok (core + cmd + dreegotest pass; plugins/sample has no tests). No regression.
  - Note: `go build ./...` from repo root is not supported by go.work (`directory prefix . does not contain modules`); modules are built individually as above.

- [x] golden-tests-core.1 tests written (increment 7, coder-test7) — **GREEN** (golden fixtures created once via `-update`, then compared)
  - **New `core/codegen_golden_test.go`** — golden-file assertions against `core/testdata/golden/*.golden`:
    - `TestGoldenSimpleRoute` — route `<head>` + `<div>` + expression, no layout → standalone head + scoped wrapper div.
    - `TestGoldenComponentWithStyle` — component with props + `<style>` → scoped CSS, `html.EscapeString` on props.
    - `TestGoldenRouteWithLayout` — route with layout `{#slot}`/`{#head}` → `c.Set("slot"...)`/`c.Set("head"...)` wrapping.
    - `TestGoldenRouter` — `GenerateRouter` registration code.
  - **Mechanik**: helper `parseFile` (Lex+NewParser), `scopeHashFor` (sha256 first 12 hex, mirrors `generate.go:126`/`scanComponents:455`), `assertGolden` with `-update` flag (write) vs compare (`t.Errorf` with full diff). No file-system access to `dreego/routes` — uses hand-built/parsed `File` objects + direct `GenerateMethodHandler`/`GenerateComponent`/`GenerateRouter` calls.
  - **Fixtures**: `core/testdata/golden/{simple_route,component_with_style,route_with_layout,router}.golden`.
  - **Verification**:
    - Fixtures created: `smd go test ./core/... -run TestGolden -update` → 4 PASS.
    - Compare mode: `smd go test ./core/... -run TestGolden` → 4 PASS.
    - Regression check: mutating a fixture → `t.Errorf` FAIL with diff (regression protection confirmed).
    - Full suite: `smd go test ./core/...` → ok; `smd go vet ./core/...` → ok.
  - Note: golden-tests produce fixtures (no classic RED start) — value is regression protection. No production code touched.
- [ ] dreegotest.1 tests written (increment 6, coder-test6) — **RED** (package does not exist yet)
  - **Package location**: `dreegotest/` at repo root. New Go module `github.com/dreego-stack/dreego/dreegotest` (own `dreegotest/go.mod`, `replace` → `../core`, same pattern as `demo/`). Added `./dreegotest` to `go.work`. This is required because `dreegotest` imports `core` and must live outside `core/`.
  - **Defined public API (unit tests `dreegotest/dreegotest_test.go`, package `dreegotest_test`):**
    - `dreegotest.Get(t *testing.T, path string) *Response` — simulates a GET through `core.ServeMux()`; returns `*Response{StatusCode int; Body string}`. Tests: status 200 + body, non-OK status (404).
    - `dreegotest.PostForm(t *testing.T, path string, form url.Values) *Response` — simulates an `application/x-www-form-urlencoded` POST; returns `*Response`. Tests: binds values (`r.FormValue("name")`), sets `Content-Type: application/x-www-form-urlencoded`.
    - `dreegotest.RenderComponent(t *testing.T, fn core.ComponentFunc, props ...any) string` — renders a single component; props are alternating key/value pairs written into the SSRContext via `ctx.Set(key, value)`. Tests: returns rendered HTML (`<h1>Welcome</h1>`), escapes XSS (`<script>` → `&lt;script&gt;`).
  - **Verification (expected RED):**
    - `smd go test ./dreegotest/...` — build failed: `no non-test Go files in /app/dreegotest` (the `dreegotest` package/import does not exist yet — desired RED).
    - `smd go vet ./dreegotest/...` — same RED build failure.
    - `smd go build ./core/... ./cmd/dreego/... ./plugins/sample/...` — GREEN (no regression from the new `go.work` entry).

## Increment 8 — component-attr-props.1 edge cases (coder, coder-write-edge)

Task `v0.0.24-attr-edge`. Writes ONLY failing tests for the two review edge cases of
`compTextWithAttrs`/`attrVal`. No production code. All four tests RED against current code.

### Edge case 1: `<script>`/`<style>` bodies must stay literal
`compTextWithAttrs` (`core/codegen_component.go:129`) is applied to **every** `NodeText`
in a component body. Inside `<script>const s = "{x}";</script>` the lexer emits the script
content as `TokenText` → `NodeText`, so `{x}` gets resolved to
`html.EscapeString(fmt.Sprintf("%v", x))` → `go build` breaks if `x` is undefined. The
lexer treats script/style content as raw text blocks where `{…}` is NOT a Go expression
(`Lex` only scans `{` when `!inSection`, `core/lexer.go:26`), so `compTextWithAttrs` must
leave those sections untouched.

### Edge case 2: multi-placeholder in one attribute value
`attrVal` (`core/codegen_helpers.go:158`) only handles the exact `{...}` (single
placeholder) shape. `href="{a}-{b}"` falls through to `return val[1:len(val)-1]` producing
the broken Go expression `a}-{b`. Must be split into `+`-concatenation of both expressions.
Affected path: `extractAttrValues` used for **component call arguments** (route-side
`genTemplateNode` NodeComponentCall and nested `genComponentCall`). The component **body**
path (`compTextWithAttrs`) already handles multi-placeholder correctly.

### Tests written (all RED)

| Test | Location | RED evidence |
|------|----------|--------------|
| `TestCompTextWithAttrsLeavesScriptStyleBodiesUntouched` | `core/codegen_component_test.go` | `{x}` in `<script>`/`<style>` body resolved to `html.EscapeString(fmt.Sprintf("%v", x))` |
| `TestAttrValResolvesMultiplePlaceholdersToConcatenation` | `core/codegen_helpers_test.go` | `attrVal("href=\"{a}-{b}\"")` → `a}-{b` (broken) |
| `Bugs/component-script-body-literal/test.sh` | `_tests/` | generated `b.WriteString(\`const s = "literal \` + html.EscapeString(fmt.Sprintf("%v", x)) + \`";\`)` |
| `Bugs/component-multi-placeholder-attr/test.sh` | `_tests/` | `Card(a}-{b).Render(c)` — broken expression, `go build` fails |

### Verification
- `smd go test ./core/...` — FAIL: only the 2 new unit tests RED; all other core tests PASS (no regression).
- `smd sh _tests/core/Bugs/component-script-body-literal/test.sh` — RED (script body `{x}` resolved as expr).
- `smd sh _tests/core/Bugs/component-multi-placeholder-attr/test.sh` — RED (`Card(a}-{b)` broken, build fails).
- **Note:** The multi-placeholder integration test initially used the component **body** path
  (`<a href="{a}-{b}">` inside the component) and passed (compTextWithAttrs already handles it).
  Rewritten to exercise the actual broken path — a **component call** passing
  `href="{a}-{b}"` → now correctly RED.

## Increment 8 fix — component-attr-props.1 edge cases (coder-write-edge, all GREEN)

Task `v0.0.24-attr-edge`. Production-code fix for both review edge cases.

### Edge case 1 fix — `<script>`/`<style>` bodies stay literal

**Root cause:** `GenerateComponent` (`core/codegen.go:346`) called `genTemplateNodeComp(n)` fresh
per `NodeText`. `<script>`/`<style>` content is lexed as raw text where `{…}` is NOT an expression
(`Lex` only scans `{` when `!inSection`, `core/lexer.go:26`), but `compTextWithAttrs` resolved
`{x}` inside those sections → generated `html.EscapeString(fmt.Sprintf("%v", x))`, breaking
`go build` when `x` is undefined.

**Fix (`core/codegen_component.go`, rewritten):**
- Introduced a stateful generator `compGen{inSection bool}`. `genTemplateNodeComp(n)` now
  constructs one `compGen` per top-level call; `GenerateComponent` (via `g.node(n)`) reuses the
  **same** generator across all sibling `NodeText` nodes so section state carries over (the parser
  splits `<script>` open tag, body, and `</script>` close tag into **three separate** `NodeText`
  nodes — confirmed via debug run `[0:"<script>"][0:"const s = ..."][0:"</script>"]`).
- `compTextSection(content, inSection) (string, bool)`: when not in a section, on encountering
  `<script`/`<style` it stops placeholder scanning, emits the segment up to there, and flips into
  section mode; while in a section it emits everything verbatim (goLiteral) until `</script>`/`</style>`
  (`sectionCloseLen`), then flips back. Outside sections the existing quoted-attribute `{…}`
  resolution is unchanged → Increment-4 attr props still work.
- `compTextWithAttrs(s)` kept as a thin wrapper (single-shot, no cross-node state) so the unit
  test `TestCompTextWithAttrsLeavesScriptStyleBodiesUntouched` (which passes a single self-contained
  string) stays valid.

### Edge case 2 fix — multi-placeholder in one quoted attribute value

**Root cause:** `attrVal` (`core/codegen_helpers.go:158`) had two raw-passthrough branches that
returned `val[1:len(val)-1]` when the value started and ended with braces. `href="{a}-{b}"`
matched the single-placeholder shape and produced the broken Go expression `a}-{b`.

**Fix (`core/codegen_helpers.go`):**
- Both single-placeholder raw branches now guard with
  `strings.Count(val,"{")==1 && strings.Count(val,"}")==1` so only an exact `{...}` is passed
  through raw (`title={user.Name}` → `user.Name`, `href={url}` → `url`).
- Anything else containing `{` falls through to new `concatPlaceholders(val)`, which splits the
  value into `+`-concatenated parts: each `{expr}` → `html.EscapeString(fmt.Sprintf("%v", expr))`
  (consistent with `compTextWithAttrs`), literal segments → `strconv.Quote` (double-quoted, per
  test contract). `href="{a}-{b}"` → ``html.EscapeString(fmt.Sprintf("%v", a)) + "-" + html.EscapeString(fmt.Sprintf("%v", b))``.

### Test-scaffolding fixes (unrelated to the edge cases, required for `go build` to pass)

- `_tests/core/Bugs/component-script-body-literal/test.sh`: `Component Snippet (x int)` + route
  `x=42` caused `cannot use "42" (untyped string constant) as int`. Changed param to `x string`
  (behavior under test — literal script body — unchanged).
- `_tests/core/Bugs/component-multi-placeholder-attr/test.sh`: `<go>a := "x"; b := "y"</go>`
  collided with the generated `var b strings.Builder` → `cannot use "y" ... as strings.Builder`.
  Renamed locals to `left`/`right` (behavior under test — multi-placeholder splitting — unchanged).

### Verification (all GREEN)

- `smd go test ./core/...` — ok (incl. both new unit tests + all golden tests `-count=1`)
- `smd go test ./cmd/dreego/...` — ok
- `smd go test ./dreegotest/...` — ok
- `smd go build ./core/... ./cmd/dreego/... ./dreegotest/...` — ok
- `smd go vet ./core/...` — ok
- `_tests/core/Bugs/component-script-body-literal/test.sh` — ok (script body kept literal `{x}`, `go build` passes)
- `_tests/core/Bugs/component-multi-placeholder-attr/test.sh` — ok (`left}-{right` split, `go build` passes)
- `_tests/core/Bugs/component-attr-prop-substitution/test.sh` — ok (no regression)
- `_tests/core/Components/prop-expression/test.sh` — ok (no regression)
- Full suite `smd sh _tests/test.sh` — **164 Passed / 0 Failed** (no regressions)

Line counts: `core/codegen_component.go` 220, `core/codegen_helpers.go` 211 — both under 300.

next: review-edge2

## Increment 8 follow-up — attr double-escaping fix (coder, v0.0.24-attr-edge)

Task `v0.0.24-attr-edge` (follow-up to reviewer-edge2 finding).

### Reviewer finding: double-escaping of multi-placeholder component-call attrs

`concatPlaceholders` (`core/codegen_helpers.go`) wrapped each `{expr}` in
`html.EscapeString(fmt.Sprintf("%v", expr))`. These concatenations are used as
**component-call arguments** (`Card(extractAttrValues(...)).Render(ctx)` in
`genComponentCall` and the route-side `genTemplateNode` NodeComponentCall path,
`core/codegen_template.go:138-143`). The component body then escapes **again** at the
prop-injection point (`{url}` → `html.EscapeString(fmt.Sprintf("%v", url))` in
`compTextSection`/`NodeExpression`). The single-placeholder raw path (`attrVal`
guarded branches) returns the raw unescaped `expr`, so escaping happens once at
injection. Net result: `href="{a}-{b}"` passed as a call arg was **double-escaped**
(`&amp;lt;` etc.) whereas `href="{url}"` was escaped once — inconsistent.

### Fix (1 line, `core/codegen_helpers.go`)

Removed `html.EscapeString(...)` from `concatPlaceholders`' expression branch —
each `{expr}` now emits only `fmt.Sprintf("%v", expr)` (literal segments stay
`strconv.Quote`). Both paths now consistently defer escaping to the prop-injection
point, so escaping happens exactly once.

### Regression test

- `core/codegen_helpers_test.go` — `TestConcatPlaceholdersDoesNotEscape`: asserts
  the output contains **no** `html.EscapeString` and still emits
  `fmt.Sprintf("%v", a)`/`fmt.Sprintf("%v", b)` plus the literal separator `"-"`.
  Guards against reintroducing double-escaping (fails if a future change re-wraps
  expressions in EscapeString).

### Verification (all GREEN)

- `smd go test ./core/... ./cmd/dreego/... ./dreegotest/...` — ok (new unit test included)
- `smd go vet ./core/...` — ok
- `_tests/core/Bugs/component-multi-placeholder-attr/test.sh` — ok
- `_tests/core/Bugs/component-attr-prop-substitution/test.sh` — ok
- `_tests/core/Components/prop-expression/test.sh` — ok (no regression)
- `_tests/core/Bugs/component-script-body-literal/test.sh` — ok (no regression)
- `_tests/core/Components/multi-props/test.sh` — ok (no regression)
- Full suite `smd sh _tests/test.sh` — **164 Passed / 0 Failed**

Line counts: `codegen_helpers.go` 211, `codegen_helpers_test.go` 106 — both under 300.

next: review-edge3

## Increment 8 regression test — GenerateComponent stateful generator (coder, v0.0.24-attr-edge)

Reviewer-Auflage: `GenerateComponent` (`core/codegen.go:315`) was switched to the stateful
`compGen` (`g := &compGen{}` reusing one generator across all `file.Template.Nodes`), but the
existing `TestCompTextWithAttrsLeavesScriptStyleBodiesUntouched` only exercises the wrapper
`compTextWithAttrs`, not the real production path. Added a direct unit test that drives
`GenerateComponent`.

### Test added (tests only, no production code)

- `core/codegen_component_test.go` — `TestGenerateComponentStatefulGenerator`:
  - Builds a `File` for a component `Card (x string, url string)` whose template contains
    `<script>const s = "literal {x}";</script>` (a raw-text section body) AND
    `<a href="{url}">go</a>` (a quoted attribute placeholder).
  - Uses existing helpers `parseFile` + `scopeHashFor` (from `codegen_golden_test.go`) and
    `ParseHeader` to strip the component header before parsing the body.
  - Asserts:
    1. `literal {x}` stays in the output and `fmt.Sprintf("%v", x)` does NOT appear (script body
       is literal, cross-node section state carried correctly across the three split `NodeText`
       nodes `<script>` / body / `</script>`).
    2. `fmt.Sprintf("%v", url)` appears and `href="{url}"` does not (attribute placeholder still
       resolved to a Go expression).
    3. The emitted component parses as valid Go via stdlib `go/parser` + `go/token`
       (prefixed `package comp\n` since `parser.ParseFile` needs a package clause).

### Verification (all GREEN)

- `smd go test ./core/... -run TestGenerateComponent -count=1` — ok
- `smd go test ./core/... -count=1` — ok (full core suite)
- `smd sh _tests/core/Bugs/component-script-body-literal/test.sh` — ok (no regression)
- `smd sh _tests/core/Bugs/component-attr-prop-substitution/test.sh` — ok (no regression)

Line count: `core/codegen_component_test.go` 94 — under 300. Stdlib only, no new comments beyond
the existing explanatory style.

next: review-edge3

## Files Changed

- `dreegotest/request.go` (new, dreegotest.1: Get + PostForm via core.ServeMux())
- `dreegotest/render.go` (new, dreegotest.1: RenderComponent with prop escaping)
- `dreegotest/dreegotest_test.go` (new, dreegotest.1 unit tests defining public API)
- `core/codegen_golden_test.go` (new, golden-tests-core.1)
- `core/testdata/golden/*.golden` (new, golden-tests-core.1 fixtures)
- `dreegotest/go.mod` (new, module github.com/dreego-stack/dreego/dreegotest, replace core)
- `go.work` (added `./dreegotest`)
- `core/validate_typed_test.go` (new, typed-forms.1 unit tests)
- `core/validate.go` (typed-forms.1: RegisterRule + int/bool/slice BindForm + ValidateForm Sprint)
- `core/validate_test.go` (BindFormNonStringFieldReturnsError → map type)
- `_tests/core/Bugs/bindform-non-string/test.sh` (updated for typed binding + unsupported-map error)
- `_tests/core/FormActions/form-int-binding/test.sh` (new)
- `_tests/core/FormActions/form-bool-binding/test.sh` (new)
- `core/codegen_component.go` (attr-prop fix: compTextWithAttrs + genComponentCall extractAttrValues)
- `core/codegen_helpers.go` (attrVal resolves {expr} inside quoted value)
- `core/codegen_component_test.go` (new)
- `core/codegen_helpers.go`
- `cmd/dreego/version.go`
- `core/codegen_layout_test.go`
- `core/codegen_head_test.go`
- `core/codegen_helpers_test.go`
- `_tests/core/Bugs/scoped-style-declarations-lost/test.sh`
- `_tests/core/Bugs/scoped-style-comma-parens/test.sh`
- `_tests/core/Bugs/scoped-style-keyframes/test.sh`
- `_tests/core/Bugs/route-head-without-layout/test.sh` (review fix)
- `_tests/core/Bugs/layout-route-head-merge/test.sh` (review fix)
- `_tests/core/Layout/no-layout/test.sh` (review fix)
- `_tests/test.sh` (Port-Schema + einmaliges curl-Install)
- 28 Server-Tests (DREEGO_PORT-Fallback, sed entfernt, `:$port` direkt) — siehe Abschnitt 8
- `.agents/tasks/v0.0.24-feedback-intake/main.md`

## Finalization (coder, ses_0309a3419ffe2p0l9LOGU6IHpg)

- [x] `_tests/core/Bugs/run-timer-sigterm/test.sh` — added `DREEGO_BIN` fallback (pattern from `CLI/new-go-sum/test.sh`) after `cd "$workdir"`. Verified standalone: `smd sh _tests/core/Bugs/run-timer-sigterm/test.sh` → **ok** (no "generate: not found").
- [x] `VERSION` bumped `v0.0.23` → `v0.0.24`.
- [x] `CHANGELOG.md` — added v0.0.24 entry summarizing all 7 blocks + port-schema + DREEGO_BIN fallback.
- [x] `_docs/testing.md` — marked newly added Layout/Components tests as done, added Bugs (scoped-css, layout, typed-forms, run-timer-sigterm), new `## 12. Form Actions` and `## 13. dreegotest` sections.
- [x] `_docs/components.md` — added `## Attribute Props` section (`{prop}` in attributes, escaped) + `## Scoped CSS` section (declarations preserved, `@media`, `@keyframes`) + rule #7.
- [x] `_docs/layouts.md` (new) — `{#slot}`/`{#head}` syntax + route `<head>` with/without layout behavior + generated Go.
- [x] `_docs/deployment.md` (new) — build, `dreego build --target <os>/<arch>` cross-compile (GOOS/GOARCH), container Dockerfile, runtime notes.
- [x] `_docs/index.md` — linked `layouts.md` and `deployment.md`.
- [x] `TODO.md` — v0.0.24 done block, deferred `frontmatter.1`/`dev-server.1`/`docs-extensibility.1` → v0.0.25, cleaned done backlog items.
- [x] `_todo/plan.md` — v0.0.24 in Done, v0.0.25 now holds the 3 deferred blocks + plugin blocks.
- [x] `.agents/log.md` — written via `agentlog` tool (not manually).
- [x] Full suite verification: `smd go test ./core/... ./cmd/dreego/... ./dreegotest/...` and `smd sh _tests/test.sh` (see below).
- No commit made (git agent does that).

**Suite verification result:**
- `smd go test ./core/... ./cmd/dreego/... ./dreegotest/...` — ok (cached, all pass)
- `smd sh _tests/test.sh` — **161 Passed / 1 Failed** (`core/CLI/new-go-sum`)
  - **Root cause (expected, release-time):** The `VERSION` bump to `v0.0.24` makes `dreego new` write `require github.com/dreego-stack/dreego/core v0.0.24`. `go mod tidy` inside `new` fails with `unknown revision core/v0.0.24` because the git tag `core/v0.0.24` does not exist yet (only `core/v0.0.22` and `core/v0.0.23` are present). This is the normal release-order dependency: the tag is created/pushed by the release process (`_scripts/release.sh`) after committing. Once `core/v0.0.24` (and `cmd/dreego/v0.0.24`, `plugins/sample/v0.0.24`) tags exist, `new-go-sum` turns green. **Not a code bug** — no other test regressed.

## Summary

v0.0.24 finalization complete. Added the missing `DREEGO_BIN` fallback to the pre-existing `run-timer-sigterm` flake test so it no longer fails standalone with "generate: not found" (SIGTERM timing flake remains pre-existing and unrelated). Bumped `VERSION` to v0.0.24 and documented the release across CHANGELOG and the `_docs/` tree (new `layouts.md` + `deployment.md`, components attribute-props/scoped-CSS, testing table additions, index links). Marked v0.0.24 done in TODO.md and `_todo/plan.md` and moved the three deferred blocks (`frontmatter.1`, `dev-server.1`, `docs-extensibility.1`) to v0.0.25. Full suite: go tests pass; integration suite 161/162 green. The single `core/CLI/new-go-sum` failure is the expected release-time dependency on the `core/v0.0.24` git tag (not yet created — created by `_scripts/release.sh` at release time), not a code regression.

## Files to Monitor

Increment scaffold-fix.1 (Nachzügler) completed. `dreegoVersion()` now reads the repo `VERSION` file as a fallback, keeping local dev builds compatible with `dreego new` and `go mod tidy`. All three scaffold CLI integration tests and the Go unit tests pass.
## Files to Monitor

- `core/codegen.go`, `core/codegen_template.go`, `core/codegen_component.go`, `core/codegen_helpers.go`, `core/codegen_layout.go`, `core/codegen_head.go`
- `core/generate.go` (`findLayout`, `isRouteDir`)
- `core/validate.go`, `core/forms.go`
- `cmd/dreego/main.go`, `cmd/dreego/blueprints/landing/`
- `_tests/core/CLI/new/test.sh`

## 8. Port-Flake-Beseitigung (deterministisches Port-Schema) — Teil der v0.0.24

Problem: Server-basierte Tests verwenden `port=$(od -An -N2 -i /dev/urandom ...)` + `sed -i "s/8080/$port/" main.go`. Zwei Probleme:
1. `sed -i` ist nicht portabel (BSD-sed im Container) → Port wird teils nicht ersetzt → Flakes.
2. Zufallsports können sich im Parallel-Lauf überlappen → Flakes.

Lösung (deterministisch, überlappungsfrei):
- `_tests/test.sh` (Runner) iteriert SEQUENTIELL über Test-Dirs und vergibt vor jedem `(...)&`-Start einen deterministisch aufsteigenden Port aus `DREEGO_PORT_BASE` (Default 20000). Exportiert `DREEGO_PORT` in die Sub-Shell. Basis 20000 + ~160 Tests bleibt unter 65535.
- Jeder Server-Test liest `port=${DREEGO_PORT:-$(od ...)}` (Fallback für Standalone) und schreibt den Port DIREKT ins `main.go`-Heredoc. `sed -i` entfällt komplett.
- Der Runner muss für Nicht-Server-Tests keine Ports vergeben, aber ein einfacher Weg: Port immer hochzählen und exportieren (kostet nichts).

Betroffene Tests (~28, Server-Ports):
- Alle mit `port=$(od -An -N2 -i /dev/urandom` UND `sed -i "s/8080/$port/"`.
- Enthält auch die 3 neuen Layout-Tests (route-head-without-layout, layout-route-head-merge, no-layout), die schon `:$port` direkt nutzen → nur DREEGO_PORT-Fallback ergänzen.

Transformationsregel pro Test:
1. Ersetze
   ```sh
   port=$(od -An -N2 -i /dev/urandom $%$ tr -d ' ')
   port=$((port % 50000 + 10000))
   ```
   durch
   ```sh
   port="${DREEGO_PORT:-$(od -An -N2 -i /dev/urandom $%$ tr -d ' ')}"
   port=$((port % 50000 + 10000))
   ```
   (Oder einfacher: `port="${DREEGO_PORT:-$(( ( $(od -An -N2 -i /dev/urandom $%$ tr -d ' ') % 50000 ) + 10000 ))}"`)
2. Entferne die `sed -i "s/8080/$port/" main.go`-Zeile; schreibe `core.Listen(":$port")` direkt ins Heredoc (falls nicht schon so).
3. Keine weiteren Änderungen.

Status:
- [x] Port-Schema im Runner (_tests/test.sh)
- [x] Alle 28 Server-Tests umgestellt
- [x] Verifikation: mehrere Läufe `smd sh _tests/test.sh` grün, keine Port-Flakes

### Umsetzung (v0.0.24-feedback-intake, Teil: Port-Flake-Beseitigung)

- **Runner `_tests/test.sh`:** `DREEGO_PORT_BASE` (Default 20000) + `port_counter`, vor jedem `(...)&`-Start `export DREEGO_PORT=$port_counter; port_counter=$((port_counter+1))`. Sequentieller aufsteigender Port pro Test → keine Überlappung im Parallel-Lauf. `DREEGO_PORT_BASE` und `DREEGO_PORT` exportiert.
- **Alle 28 Server-Tests:** `port=${DREEGO_PORT:-$(( ( $(od -An -N2 -i /dev/urandom | tr -d ' ') % 50000 ) + 10000 ))}`; `sed -i "s/8080/$port/"` entfernt; `core.Listen(":$port")` direkt ins `main.go`-Heredoc (unquoted) geschrieben. Enthält die 3 Layout-Tests (route-head-without-layout, layout-route-head-merge, no-layout) + run-timer-sigterm, die schon `:$port` direkt nutzten.
- **Zusätzlicher Flake gefunden (Root-Cause):** Parallel-Läufe flakten NICHT wegen Ports, sondern wegen **apk-Datenbanklock-Race**: alle 27 Server-Tests starten `apk add --no-cache curl` gleichzeitig → `ERROR: Unable to lock database`. Die 3 Layout-Tests ohne `|| true` brachen daraufhin hart ab. Fix: curl wird EINMAL sequentiell im Runner vor der Test-Schleife installiert (deterministisch), und in den 3 Layout-Tests das `apk add ... || true` (curl ist durch den Runner garantiert vorhanden).
- **Verifikation:** 3× kompletter Lauf `smd sh _tests/test.sh` grün (157 Passed / 0 Failed), keine Port-Flakes. 5× gezielte Läufe je Test `layout-route-head-merge`, `route-head-without-layout`, `no-layout` alle grün. `CLI/check-stale` ist eine vorbestehende mtime-Granularity-Flake (unabhängig vom Port-Schema), standalone grün.

## Finalization — new.go offline replace (echter Fix, coder)

- [x] **`cmd/dreego/new.go`** — `dreego new` fügt jetzt selbst eine `replace`-Direktive auf das lokale `core`-Verzeichnis ins Scaffold ein, wenn es als repo-lokaler Build läuft:
  - `findLocalCore()` (neu): löst `<repo>/core` über `runtime.Caller(0)` relativ zum Quellverzeichnis auf (Quellpfad `<repo>/cmd/dreego/new.go` → `../../core`), verifiziert `core/go.mod` existiert. Muster analog `versionFromSourceRoot` in `version.go`.
  - Nach `go mod edit -require ...`: falls `findLocalCore()` nicht leer → `go mod edit -replace=github.com/dreego-stack/dreego/core=<abs coreDir>`.
  - `go mod tidy` läuft jetzt mit `GOWORK=off`, damit die Auflösung gegen das lokale replace offline/deterministisch erfolgt.
  - **Release-sicher:** Schlägt NICHT fehl, wenn kein lokales core existiert (release-installiertes Binary hat keins) → tidy löst dann den gepushten Tag remote. Damit bricht das Release-Verhalten nicht.
- [x] **`_tests/core/CLI/new-go-sum/test.sh`** — vereinfacht auf reales Verhalten ohne manuellen replace-Workaround:
  - `$DREEGO_BIN new testapp`, `cd testapp`, prüft dass `go.mod` eine `^replace github.com/dreego-stack/dreego/core => `-Zeile enthält (von `dreego new` geschrieben), `GOWORK=off go mod tidy`, `$DREEGO_BIN generate`, `GOWORK=off go build .`.
- [x] **`cmd/dreego/go.mod`** — unnötiger Whitespace-Diff (vom früheren Agenten) revertiert; `git diff cmd/dreego/go.mod` ist leer. Nicht committet.
- [x] Verifikation:
  - `smd sh _tests/core/CLI/new-go-sum/test.sh` — grün (exit 0), offline.
  - `smd sh _tests/core/CLI/new-gitignore/test.sh` — grün.
  - `smd sh _tests/core/CLI/new-layout-exists/test.sh` — grün.
  - `smd go test ./core/... ./cmd/dreego/... ./dreegotest/...` — grün (clean `-count=1`).
  - `smd sh _tests/test.sh` — **162 Passed / 0 Failed** (keine Regression).
  - `smd go vet ./cmd/dreego/...` — grün.
- Kein Commit (git-Agent übernimmt das).

next: review-final3
