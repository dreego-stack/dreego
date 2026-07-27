# Testing Strategy

Jeder Bereich des Frameworks braucht positive (Verhalten bestätigen) und negative (Fehler früh erkennen) Tests. 
Tests laufen als Integration-Tests in `_tests/` via Docker (`make test`). Bugs permanent unter `_tests/Bugs/`.

## 1. Transpiler

| Test | Typ | Beschreibung |
|------|-----|-------------|
| basic-page | ✅ pos | `<head>`, `<go>`, `<div>`, `<style>` |
| all-sections | ✅ pos | Alle 5 Sektionen |
| no-go | ✅ pos | Route ohne `<go>` |
| unclosed-div | ✅ neg | `<div>` ohne `</div>` → generate FAIL |
| mismatched-close | ✅ neg | `<div>...</go>` → generate FAIL |
| xss-escaping | ✅ pos | `{var}` escaped `<script>` |
| duplicate-head | ⬜ neg | Zwei `<head>`-Sektionen → generate FAIL |
| duplicate-div | ⬜ neg | Zwei `<div>`-Sektionen → generate FAIL |
| empty-div | ⬜ pos | `<div></div>` — leeres Template |
| only-go | ⬜ pos | Nur `<go>` Sektion ohne `<div>` |
| large-template | ⬜ pos | 500+ Zeilen HTML im `<div>` |
| unicode | ⬜ pos | UTF-8 Zeichen in Template und Props |
| comment | ⬜ pos | `{! das ist ein Kommentar !}` |
| verbatim | ⬜ neg | `{#verbatim}` noch nicht implementiert → Fehler |

## 2. Template-Expressions

| Test | Typ | Beschreibung |
|------|-----|-------------|
| if-true | ✅ pos | `{#if true}` rendert |
| if-false | ✅ pos | `{#if false}` rendert nicht |
| each-loop | ✅ pos | `{#each items as item}` mit 3 Items |
| expression | ✅ pos | `{var}` im `<div>` |
| nested-if | ⬜ pos | `{#if}{#if}{/if}{/if}` |
| if-else | ⬜ neg | `{#else}` noch nicht implementiert |
| each-empty | ⬜ pos | Leere Liste im `{#each}` → kein Output |
| each-with-if | ⬜ pos | `{#each}` mit `{#if}` drin |
| expression-missing-var | ⬜ neg | `{undefined}` → go build FAIL |
| expression-function | ⬜ pos | `{len(items)}` als Expression |

## 3. Layout

| Test | Typ | Beschreibung |
|------|-----|-------------|
| with-slot | ✅ pos | Layout mit `{#slot}` |
| with-head | ✅ pos | Layout mit `{#head}` |
| no-layout | ⬜ pos | Seite ohne Layout-Datei |
| nested-slot | ⬜ neg | `{#slot}` in Layout, nicht in Route |

## 4. Routing

| Test | Typ | Beschreibung |
|------|-----|-------------|
| get-method | ✅ pos | GET Route |
| post-method | ✅ pos | POST Route |
| dynamic | ✅ pos | `[id]` Segment |
| catchall | ✅ pos | `[...path]` Segment |
| optional | ✅ pos | `[[lang]]` Segment |
| groups | ✅ pos | `(group)/` unsichtbar in URL |
| 404-page | ✅ pos | Custom 404 |
| 500-page | ✅ pos | Custom 500 |
| delete-method | ⬜ pos | DELETE Route |
| put-method | ⬜ pos | PUT Route |
| multi-segment | ⬜ pos | `[a]/[b]/get.dreego` |
| bracket-in-name | ⬜ pos | `[id]` im Dateinamen mit Bindestrich: `user-[id].dreego` |
| deep-nesting | ⬜ pos | `a/b/c/d/get.dreego` (4 Ebenen) |

## 5. Middleware

| Test | Typ | Beschreibung |
|------|-----|-------------|
| recovery-panic | ✅ pos | Panic → 500, kein Crash |
| csrf-token | ✅ pos | Cookie gesetzt auf GET |
| csrf-post-fail | ⬜ neg | POST ohne Token → 403 |
| csrf-post-pass | ⬜ pos | POST mit Token → 200 |
| csrf-disabled | ⬜ pos | `SetCSRF(false)` → POST ohne Token OK |

## 6. Session

| Test | Typ | Beschreibung |
|------|-----|-------------|
| set-get | ✅ pos | Session-Wert setzen und lesen |
| delete | ⬜ pos | `DelSessionVal` löscht Wert |
| destroy | ⬜ pos | `DestroySession` löscht alles |
| no-store | ⬜ pos | Ohne `SetSessionStore` — SessionVal leer |

## 7. Components

| Test | Typ | Beschreibung |
|------|-----|-------------|
| basic | ✅ pos | Prop + Rendering |
| self-closing | ✅ pos | `<@Icon/>` ohne Body |
| not-found | ✅ pos | `<@Missing/>` → go build FAIL |
| scoped-style | ✅ pos | Component-CSS leaked nicht |
| with-go | ✅ pos | `<go>` im Component |
| with-slot | ✅ pos | Default-Slot mit Children |
| nested-component | ⬜ pos | Component ruft anderen Component auf |
| empty-props | ⬜ pos | Component ohne Props |
| multi-props | ⬜ pos | Component mit 3+ Props |
| prop-default | ⬜ pos | Prop mit Default-Wert |
| prop-expression | ⬜ pos | Prop-Wert aus Expression: `title={user.Name}` |
| slot-missing | ⬜ pos | Component mit `{#slot}`, Aufruf ohne Body |
| slot-named | ⬜ v0.0.7 | `{#slot header}` |
| recursive | ⬜ neg | Component ruft sich selbst → Fehler oder Warnung |
| import-alias | ⬜ v0.0.7 | `import C "path"` → `<@C/>` |

## 8. Imports

| Test | Typ | Beschreibung |
|------|-----|-------------|
| basic | ✅ pos | `import "dreego/components/Card"` |
| multi-file | ✅ pos | `import "dreego/components/button"` |
| missing | ✅ pos | Import-Pfad existiert nicht → kein Crash |
| subdir | ⬜ v0.0.7 | `import "dreego/components/button"` → `<@Login/>` |
| alias | ⬜ v0.0.7 | `import Btn "path"` → `<@Btn/>` |
| duplicate-import | ⬜ pos | Gleicher Import zweimal → kein Fehler |

## 9. CLI

| Test | Typ | Beschreibung |
|------|-----|-------------|
| init | ✅ pos | `dreego init .` erstellt Dateien |
| check | ✅ pos | `dreego generate --check` |
| check-stale | ⬜ pos | Nach .dreego-Änderung → `--check` FAIL |
| no-args | ⬜ pos | `dreego` ohne Args → Hilfe |
| unknown-cmd | ⬜ pos | `dreego invalid` → Fehler |

## 10. Config

| Test | Typ | Beschreibung |
|------|-----|-------------|
| redirect | ⬜ pos | `dreego/config.json` mit Redirect |
| rewrite | ⬜ pos | `dreego/config.json` mit Rewrite |
| logging-off | ⬜ pos | `Logging.Enabled: false` |
| invalid-json | ⬜ neg | Kaputtes JSON in config |

## 11. Bugs (Regression)

| Test | Typ | Beschreibung |
|------|-----|-------------|
| component-close-tag | ✅ bug | `</@Card>` Lexer-Fix |
| component-quoted-attrs | ✅ bug | `title="Hello World"` mit Spaces |
| div-in-slot | ⬜ bug | `<@Card><div>hi</div></@Card>` — HTML in Children |

## Zusammenfassung

| Status | Count |
|--------|-------|
| ✅ Implementiert | 36 |
| ⬜ Geplant (v0.0.7) | 30 |
| ⬜ Named Slots (v0.0.8) | 4 |

**Test-Philosophie:** Jedes Verhalten, jeder Edge-Case, jeder Bug wird ein Test. 
Positiv-Tests bestätigen korrektes Verhalten. Negativ-Tests stellen sicher, dass Fehler FRÜH erkannt werden 
(beim `dreego generate`, nicht erst beim `go build` oder zur Laufzeit).
