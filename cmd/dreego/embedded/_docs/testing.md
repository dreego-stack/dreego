# Testing Strategy

Every area of the framework needs positive (confirm behavior) and negative (detect errors early) tests.
Tests run as integration tests in `_tests/` via Docker (`make test`). Bugs permanently under `_tests/core/Bugs/`.

## 1. Transpiler

| Test | Type | Description |
|------|-----|-------------|
| basic-page | ✅ pos | `<head>`, `<go>`, `<div>`, `<style>` |
| all-sections | ✅ pos | All 5 sections |
| no-go | ✅ pos | Route without `<go>` |
| unclosed-div | ✅ neg | `<div>` without `</div>` → generate FAIL |
| mismatched-close | ✅ neg | `<div>...</go>` → generate FAIL |
| xss-escaping | ✅ pos | `{var}` escaped `<script>` |
| duplicate-head | ⬜ neg | Two `<head>` sections → generate FAIL |
| duplicate-div | ⬜ neg | Two `<div>` sections → generate FAIL |
| empty-div | ⬜ pos | `<div></div>` — empty template |
| only-go | ⬜ pos | Only `<go>` section without `<div>` |
| large-template | ⬜ pos | 500+ lines HTML in `<div>` |
| unicode | ⬜ pos | UTF-8 characters in template and props |
| comment | ⬜ pos | `{! this is a comment !}` |
| verbatim | ⬜ neg | `{#verbatim}` not yet implemented → error |

## 2. Template Expressions

| Test | Type | Description |
|------|-----|-------------|
| if-true | ✅ pos | `{#if true}` renders |
| if-false | ✅ pos | `{#if false}` does not render |
| each-loop | ✅ pos | `{#each items as item}` with 3 items |
| expression | ✅ pos | `{var}` in `<div>` |
| nested-if | ⬜ pos | `{#if}{#if}{/if}{/if}` |
| if-else | ⬜ neg | `{#else}` not yet implemented |
| each-empty | ⬜ pos | Empty list in `{#each}` → no output |
| each-with-if | ⬜ pos | `{#each}` with `{#if}` inside |
| expression-missing-var | ⬜ neg | `{undefined}` → go build FAIL |
| expression-function | ⬜ pos | `{len(items)}` as expression |

## 3. Layout

| Test | Type | Description |
|------|-----|-------------|
| with-slot | ✅ pos | Layout with `{#slot}` |
| with-head | ✅ pos | Layout with `{#head}` |
| no-layout | ✅ pos | Route without layout file still renders as full fragment (`Layout/no-layout`) |
| layout-not-applied | ✅ bug | Route with layout renders inside layout `{#slot}` (`Bugs/layout-not-applied`) |
| route-head-without-layout | ✅ bug | Route `<head>` appears when no layout exists (`Bugs/route-head-without-layout`) |
| layout-route-head-merge | ✅ bug | Route `<head>` merged into layout `{#head}` (`Bugs/layout-route-head-merge`) |
| nested-slot | ⬜ neg | `{#slot}` in layout, not in route |

## 4. Routing

| Test | Type | Description |
|------|-----|-------------|
| get-method | ✅ pos | GET route |
| post-method | ✅ pos | POST route |
| dynamic | ✅ pos | `[id]` segment |
| catchall | ✅ pos | `[...path]` segment |
| optional | ✅ pos | `[[lang]]` segment |
| groups | ✅ pos | `(group)/` invisible in URL |
| 404-page | ✅ pos | Custom 404 |
| 500-page | ✅ pos | Custom 500 |
| delete-method | ✅ pos | DELETE route |
| put-method | ⬜ pos | PUT route |
| multi-segment | ⬜ pos | `[a]/[b]/get.dreego` |
| bracket-in-name | ⬜ pos | `[id]` in filename with hyphen: `user-[id].dreego` |
| deep-nesting | ⬜ pos | `a/b/c/d/get.dreego` (4 levels) |

## 5. Middleware

| Test | Type | Description |
|------|-----|-------------|
| recovery-panic | ✅ pos | Panic → 500, no crash |
| csrf-token | ✅ pos | Cookie set on GET |
| csrf-post-fail | ⬜ neg | POST without token → 403 |
| csrf-post-pass | ⬜ pos | POST with token → 200 |
| csrf-disabled | ⬜ pos | `SetCSRF(false)` → POST without token OK |

## 6. Session

| Test | Type | Description |
|------|-----|-------------|
| set-get | ✅ pos | Set and read session value |
| delete | ⬜ pos | `DelSessionVal` deletes value |
| destroy | ⬜ pos | `DestroySession` deletes everything |
| no-store | ⬜ pos | Without `SetSessionStore` — SessionVal empty |

## 7. Components

| Test | Type | Description |
|------|-----|-------------|
| basic | ✅ pos | Prop + rendering |
| self-closing | ✅ pos | `<@Icon/>` without body |
| not-found | ✅ pos | `<@Missing/>` → go build FAIL |
| scoped-style | ✅ pos | Component CSS does not leak |
| with-go | ✅ pos | `<go>` in component |
| with-slot | ✅ pos | Default slot with children |
| nested-component | ✅ pos | Component calls another component (`Components/nested`) |
| empty-props | ✅ pos | Component without props (`Components/empty-props`) |
| multi-props | ✅ pos | Component with 3+ props (`Components/multi-props`) |
| prop-default | ⬜ pos | Prop with default value |
| prop-expression | ✅ pos | Prop value from expression: `title={user.Name}` (`Components/prop-expression`, `prop-expr`) |
| attr-prop-substitution | ✅ bug | `{prop}` substituted in HTML attributes: `<a href="{url}">` (`Bugs/component-attr-prop-substitution`) |
| slot-missing | ⬜ pos | Component with `{#slot}`, call without body |
| slot-named | ⬜ v0.0.7 | `{#slot header}` |
| recursive | ⬜ neg | Component calls itself → error or warning |
| import-alias | ⬜ v0.0.7 | `import C "path"` → `<@C/>` |

## 8. Imports

| Test | Type | Description |
|------|-----|-------------|
| basic | ✅ pos | `import "dreego/components/Card"` |
| multi-file | ✅ pos | `import "dreego/components/button"` |
| missing | ✅ pos | Import path does not exist → no crash |
| subdir | ⬜ v0.0.7 | `import "dreego/components/button"` → `<@Login/>` |
| alias | ⬜ v0.0.7 | `import Btn "path"` → `<@Btn/>` |
| duplicate-import | ⬜ pos | Same import twice → no error |

## 9. CLI

| Test | Type | Description |
|------|-----|-------------|
| init | ✅ pos | `dreego init .` creates files |
| check | ✅ pos | `dreego generate --check` |
| check-stale | ⬜ pos | After .dreego change → `--check` FAIL |
| no-args | ⬜ pos | `dreego` without args → help |
| unknown-cmd | ⬜ pos | `dreego invalid` → error |

## 10. Config

| Test | Type | Description |
|------|-----|-------------|
| redirect | ⬜ pos | `dreego/config.json` with redirect |
| rewrite | ⬜ pos | `dreego/config.json` with rewrite |
| logging-off | ⬜ pos | `Logging.Enabled: false` |
| invalid-json | ⬜ neg | Broken JSON in config |

## 11. Bugs (Regression)

| Test | Type | Description |
|------|-----|-------------|
| component-close-tag | ✅ bug | `</@Card>` lexer fix |
| component-quoted-attrs | ✅ bug | `title="Hello World"` with spaces |
| component-attr-prop-substitution | ✅ bug | `{prop}` resolved inside HTML attributes (`Bugs/component-attr-prop-substitution`) |
| scoped-style-declarations-lost | ✅ bug | Declarations in `{}` preserved (`radial-gradient`) (`Bugs/scoped-style-declarations-lost`) |
| scoped-style-comma-parens | ✅ bug | Selectors with commas + nested parens (`calc()`, `rgb()`) (`Bugs/scoped-style-comma-parens`) |
| scoped-style-keyframes | ✅ bug | `@keyframes` body preserved (`Bugs/scoped-style-keyframes`) |
| scoped-css-media | ✅ bug | `@media` inner selectors scoped (`Bugs/scoped-css-media`) |
| layout-not-applied | ✅ bug | Layout applies (`Bugs/layout-not-applied`) |
| route-head-without-layout | ✅ bug | `<head>` works without layout (`Bugs/route-head-without-layout`) |
| layout-route-head-merge | ✅ bug | Route `<head>` merges into layout `{#head}` (`Bugs/layout-route-head-merge`) |
| bindform-typed | ✅ bug | int/bool/slice binding + unsupported map type error (`Bugs/bindform-non-string`) |
| run-timer-sigterm | ✅ bug | SIGTERM graceful shutdown (`Bugs/run-timer-sigterm`) |
| div-in-slot | ⬜ bug | `<@Card><div>hi</div></@Card>` — HTML in children |

## 12. Form Actions

| Test | Type | Description |
|------|-----|-------------|
| g-action-basic | ✅ pos | `<form g-action>` submits + redirect |
| form-int-binding | ✅ pos | POST binds integer field + validates min (`FormActions/form-int-binding`) |
| form-bool-binding | ✅ pos | POST binds checkbox to bool (`FormActions/form-bool-binding`) |

## 13. dreegotest

| Test | Type | Description |
|------|-----|-------------|
| dreegotest-basic | ✅ unit | `dreegotest.Get` simulates GET + reads body |
| dreegotest-form | ✅ unit | `dreegotest.PostForm` validates form submission |
| dreegotest-render | ✅ unit | `dreegotest.RenderComponent` renders with props + XSS escape |

## Summary

| Status | Count |
|--------|-------|
| ✅ Implemented | 36 |
| ⬜ Planned (v0.0.7) | 30 |
| ⬜ Named Slots (v0.0.8) | 4 |

**Test Philosophy:** Every behavior, every edge case, every bug becomes a test.
Positive tests confirm correct behavior. Negative tests ensure errors are detected EARLY
(at `dreego generate`, not only at `go build` or at runtime).
