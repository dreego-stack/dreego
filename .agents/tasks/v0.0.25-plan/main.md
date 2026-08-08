---
type: Task
title: v0.0.25 Plan — Plugin Ecosystem + Frontmatter + Dev Server + Embedded Docs
status: in_progress
assign: manager
---

# v0.0.25 Plan

Goal: Ship the frozen plugin contract (monolithic Plugin interface), plugin middleware/route hooks, YAML frontmatter, a full dev server, and embedded docs (no HTTP call). Execute strictly sequentially with TDD: one failing test → one coder fix → review → commit → next.

## Blocks (in execution order)

1. **plugin-interface.1** — Frozen monolithic Plugin interface (v1 contract). Paket-level API (no App object). Import alias `import dreego "github.com/dreego-stack/dreego/core"`. `dreego.UsePlugin(p Plugin)`. Plugins are external Go modules that import core and satisfy the interface; core never imports a plugin.
   - Interface methods (monolithic, "plugins can do everything"): `Name() string`, `RegisterRoutes(...)`, `Middlewares() []func(http.Handler) http.Handler`, `Assets() fs.FS`, `OnStart(ctx) error`, `OnShutdown(ctx) error`.
   - Tests: ≥3 (interface satisfaction, UsePlugin registration, lifecycle).

2. **middleware-hooks.1** — `dreego.UsePlugin` injects plugin middleware into the core stack, FIFO order, order fixated on first Listen. Tests: ≥3.

3. **route-hooks.1** — Plugins register own URL paths programmatically (no filesystem). Generated `dreego/gen/dree.go` collects plugin routes. Tests: ≥3.

4. **docs-extensibility.1** — `dreego docs` reads plugin docs from `plugins/<name>/_docs/` (local) without compile-time dependency on any plugin. Tests: ≥3.

5. **docs-embed.1** (NEW from user) — Embed `_docs/` files into the CLI via `//go:embed` so `dreego docs` works offline without an HTTP call. Currently `cmd/dreego/docs.go` does `http.Get` to codeberg.org. Tests: ≥3 (offline fetch, dump, json).

6. **frontmatter.1** — Parse YAML frontmatter at top of `.dreego` files, expose typed metadata via `c.Data("key")`. YAML only. Tests: ≥3.

7. **dev-server.1** — `dreego dev` with file watcher, auto-regenerate, server restart. Full (watcher + restart). Tests: ≥3.

## Decisions (user-confirmed)

- Plugin interface: monolithic (single interface, all methods).
- API: paket-level `dreego.*` (no App object). Consistent with existing `core.Listen`/`core.Build`.
- Import alias: `import dreego "github.com/dreego-stack/dreego/core"`.
- Frontmatter: YAML.
- Plugin methods: full set (Name, RegisterRoutes, Middlewares, Assets, OnStart, OnShutdown).
- dev-server: full (watcher + restart).
- docs: embedded (no HTTP call).

## Execution mode

- Strictly sequential, one agent at a time.
- Per increment: coder-testN writes failing tests → coder-writeN fixes → reviewer → git commit (no push unless user asks).
- Agent naming: coder-test1..N, coder-write1..N, review1..N, git1..N.

## Status

- [ ] Plan approved by user
- [x] Increment 1: plugin-interface.1 (coder-write1: `core/plugin.go` + wiring)

### Increment 1 — implemented (coder-write1)

- Added `core/plugin.go`: `Plugin` interface, `UsePlugin`, `StartPlugins`, `ShutdownPlugins`, internal slices `plugins`, `pluginMiddlewares`, `pluginAssets`.
- `UsePlugin` calls `RegisterRoutes()`, collects `Middlewares()` (FIFO), stores `Assets()` and the plugin (lifecycle).
- `StartPlugins`/`ShutdownPlugins` iterate plugins and return the first error.
- `Build()` in `core/runtime.go` wraps plugin middleware around the mux (outer layer, FIFO) before redirect/rewrite.
- `Reset()` also clears the three plugin slices (plugins, middleware, assets).
- Side fix: `Register()` is now idempotent per `method+pattern` (replaces handler instead of appending duplicate). Required because the lifecycle test registers `/plugin-lifecycle` across multiple tests while routes intentionally survive `Reset()` (per `TestResetClearsCache`).
- Verification: `go test ./core/...`, `./cmd/dreego/...`, `./dreegotest/...` green; full `_tests/test.sh` 164 passed / 0 failed.

next: review1
- [ ] Increment 2: middleware-hooks.1

### Increment 2 — tests written (coder-test2)

Defined FIFO semantics (fixated for v1): `UsePlugin(pluginA); UsePlugin(pluginB)` means on request entry **A runs first, then B, then the handler** — the first registered plugin is the outermost middleware. This is the standard Go middleware convention (e.g. Chi). The current `Build()` loop `for _, mw := range pluginMiddlewares { h = mw(h) }` produces LIFO (B first), so the order test is RED.

Tests added to `core/plugin_test.go` (all RED against current code):
1. `TestPluginMiddlewareFIFOOrder` — two order plugins A,B append to a shared log; asserts `[A B]`. Current code yields `[B A]` → RED.
2. `TestPluginMiddlewareOrderFixatedOnFirstBuild` — order is fixated on the first `Build()`; registering B after the first build must not change the order. Currently PASSES (builtHandler cached) — documents the fixated semantics.
3. `TestPluginMiddlewareNilEntryStable` — a nil middleware entry must not panic; stack stays stable. Current code panics in `Build()` → RED.

Verification: `go test ./core/...` → `TestPluginMiddlewareFIFOOrder` FAIL (`[B A]`), `TestPluginMiddlewareNilEntryStable` PANIC, `TestPluginMiddlewareOrderFixatedOnFirstBuild` PASS.

### Increment 2 — implemented (coder-write2)

Fixed `core/runtime.go` `Build()` middleware chain:

1. **FIFO order**: Replaced the forward `for _, mw := range pluginMiddlewares { h = mw(h) }` loop (which produced LIFO, B first) with a reverse iteration `for i := len(pluginMiddlewares)-1; i >= 0; i-- { h = pluginMiddlewares[i](h) }`. Now the first registered plugin (`pluginMiddlewares[0]`) is the outermost middleware and runs first on request entry, matching the standard Go/Chi convention. `TestPluginMiddlewareFIFOOrder` now asserts `[A B]` — GREEN.
2. **nil-guard**: Added `if pluginMiddlewares[i] == nil { continue }` inside the loop so a nil middleware entry is skipped instead of panicking. `TestPluginMiddlewareNilEntryStable` — GREEN.

`TestPluginMiddlewareOrderFixatedOnFirstBuild` still passes (fixated via `builtHandler` cache in `Build()`). `TestUsePluginCollectsMiddleware` (Increment 1) still green — header middleware unaffected by order change.

Verification: `go test ./core/...`, `./cmd/dreego/...`, `./dreegotest/...` all GREEN; full `_tests/test.sh` 164 passed / 0 failed.

next: review2
- [ ] Increment 3: route-hooks.1

### Increment 3 — tests written (coder-test3)

**Analyse:** `UsePlugin` ruft bereits `p.RegisterRoutes()` auf; das Plugin ruft intern `core.Register("GET", "/plugin/...", handler)` → Routen landen automatisch im `routes`-Slice und sind über `ServeMux()` erreichbar. Der Block erfindet also kein neues Verhalten, sondern sichert die bestehende programmatische Registrierung ab.

**Gen/dree.go-Frage geklärt:** `GenerateRouter(routes []RouteInfo)` ist eine reine Codegen-Funktion, die Registration-Code für die *explizit übergebenen* `RouteInfo`-Einträge erzeugt. Es gibt **keine** Plugin-Discovery zur Codegen-Zeit — Plugin-Routen werden zur Laufzeit registriert und können nicht automatisch in `gen/dree.go` gesammelt werden. "Gen/dree.go sammelt Plugin-Routen" heißt: das Tooling muss Plugin-`RouteInfo`-Einträge einsammeln und explizit an `GenerateRouter` übergeben. Das wird als Unit-Test abgesichert (TestGenerateRouterCollectsPluginRoutes).

**Definierte Semantik (Route-Overlap / Last-Wins):**
- `Register()` ist idempotent pro `method+pattern` (Increment 1): die LETZTE Registrierung ersetzt den Handler (last-wins) statt Duplikate anzuhängen.
- Unabhängig von der Reihenfolge (Plugin vor App oder App vor Plugin) gewinnt die zuletzt registrierte Route dasselbe Pattern.
- `routes`-Slice überlebt `Reset()` (bestehende Semantik aus TestResetClearsCache); Plugin-Routen bleiben nach Reset erreichbar.

**Tests (alle GREEN gegen aktuellen Code — Verhalten existiert schon):**
1. `TestPluginRegistersMultipleRoutes` — Plugin registriert GET, POST und dynamisches Pattern (`/plugin/multi/{id}`); alle über `ServeMux()` erreichbar (httptest, Statuscodes + Body). GREEN.
2. `TestPluginRoutesSurviveReset` — Plugin-Route bleibt nach `Reset()` erreichbar (dokumentiert bestehende Routen-Survival-Semantik). GREEN.
3. `TestPluginRouteLastWinsOverridesAppRoute` — App-Route zuerst, Plugin-Override danach → Plugin-Handler gewinnt. GREEN.
4. `TestPluginRouteLastWinsAppWinsOverPlugin` — Plugin-Route zuerst, App-Override danach → App-Handler gewinnt. GREEN.
5. `TestGenerateRouterCollectsPluginRoutes` — `GenerateRouter` erzeugt korrekte `mux.HandleFunc`-Zeilen für Plugin-RouteInfo (GET /admin, POST /api/auth/login). GREEN.

Dateien: `core/route_hooks_test.go` (130 Zeilen), `core/route_hooks_test_helpers.go` (57 Zeilen, Plugins `multiRoutePlugin`, `overlapRoutePlugin`). Kein Produktivcode geändert.

Verification: `go test ./core/ -count=1` → ok (alle Tests PASS).
- [ ] Increment 4: docs-extensibility.1

### Increment 4 — tests written (coder-test4)

**Analyse:** `cmd/dreego/docs.go` liest aktuell ausschließlich per HTTP (`fetchDoc` → `http.Get` auf codeberg.org). Offizielle Plugins leben unter `plugins/<name>/` mit `_docs/` als primärer Quelle. Es darf KEINE Compile-Zeit-Abhängigkeit auf ein konkretes Plugin geben → Discovery per Dateisystem, nicht per Plugin-Import.

**Definierte API (in `cmd/dreego/docs.go`, noch nicht implementiert):**

```go
// pluginDocsRoot ist der Wurzelordner, unter dem Plugins liegen (default "plugins").
// In Tests überschreibbar (t.TempDir()).
var pluginDocsRoot = "plugins"

// fetchDocLocal sucht zuerst lokale Plugin-Docs unter
// plugins/<name>/_docs/<path>, dann den Fallback (embedded/remote).
// fromLocal=true wenn die lokale Plugin-Docs geliefert wurde.
func fetchDocLocal(path string) ([]byte, fromLocal bool, err error)

// fetchDocFallback ist der Fallback-Loader (aktuell fetchDoc/HTTP,
// Increment 5 stellt auf embedded um). In Tests überschreibbar.
var fetchDocFallback = fetchDoc
```

**Prioritätsreihenfolge (fixiert):** lokal/Plugin-Docs → embedded/remote. D.h. `plugins/<name>/_docs/<path>` gewinnt immer vor dem Fallback. `fetchDocLocal` liefert `fromLocal=true` nur, wenn die lokale Plugin-Docs existiert und gelesen wurde; sonst `fromLocal=false` und der Fallback liefert den Body.

**Tests (alle RED — Feature fehlt, Kompilierung schlägt fehl):**
1. `TestFetchDocLocalPluginDiscovery` — legt `plugins/sample/_docs/index.md` in `t.TempDir()` an, setzt `pluginDocsRoot`, ruft `fetchDocLocal("/plugins/sample/_docs/index.md")` → erwartet `fromLocal=true` + exakter Body. RED (Funktion fehlt).
2. `TestFetchDocLocalReadsFilesystem` — liest `plugins/auth/_docs/guide.md` per Dateisystem (kein Plugin-Import); fehlendes Plugin (`plugins/missing/_docs/...`) → `fromLocal=false`, kein Fehler. RED.
3. `TestFetchDocLocalFallback` — leeres `pluginDocsRoot` + überschriebener `fetchDocFallback` → `fromLocal=false`, Body vom Fallback. RED.
4. `TestFetchDocLocalPriorityLocalWins` — lokale Docs existieren UND Fallback gesetzt → lokale Docs gewinnen (`fromLocal=true`, lokaler Body). RED.

Datei: `cmd/dreego/docs_test.go` (Unit-Tests, `t.TempDir()`). Kein Produktivcode geändert.

Verification: `go test ./cmd/dreego/...` → build failed (undefined: pluginDocsRoot, fetchDocLocal, fetchDocFallback). `go build ./cmd/dreego` → ok (Produktivcode kompiliert).

### Increment 4 — implemented (coder-write4)

Implemented the local plugin-docs discovery in `cmd/dreego/docs.go`:

1. **`var pluginDocsRoot = "plugins"`** — package-level var, overridable in tests (`t.TempDir()`).
2. **`var fetchDocFallback = fetchDoc`** — fallback loader (HTTP, stays until Increment 5).
3. **`fetchDocLocal(path string) ([]byte, bool, error)`** — pure filesystem discovery, no compile-time dep on any plugin:
   - Paths not under `plugins/` → delegate to `fetchDocFallback`, `fromLocal=false`.
   - `plugins/<name>/_docs/<path>` → `os.ReadFile(pluginDocsRoot/<name>/_docs/<path>)`.
   - Found → body + `fromLocal=true`. Not found (`os.IsNotExist`) → `nil, false, nil` (caller uses fallback). Real I/O errors propagate.
4. **`cmdDocs`** uses `fetchDocLocal`; when `fromLocal=false` it calls `fetchDocFallback` for the body.
5. **`cmdDump`** uses the same local-first flow (consistent).

HTTP remains the fallback (Increment 5 switches to embedded). File: `cmd/dreego/docs.go` (221 lines, <300).

Verification: `go test ./cmd/dreego/... ./core/... ./dreegotest/...` all GREEN (incl. 4 new tests); `go build ./cmd/dreego` ok; full `_tests/test.sh` 164 passed / 0 failed (one flaky `run-timer-sigterm` timing failure on first run, green on re-run — unrelated to this change).

### Increment 4 — blocker fix (coder-write4b)

**Blocker (Reviewer):** `fetchDocLocal` rief für Nicht-Plugin-Pfade intern `fetchDocFallback` auf, lieferte aber `fromLocal=false` zurück. Der Aufrufer (`cmdDocs`/`cmdDump`) interpretierte `fromLocal=false` als "Body fehlt" und rief den Fallback ERNEUT auf → `dreego docs /_docs/index.md` machte 2 HTTP-Requests. Verletzte die Semantik "Fallback liefert den Body → der Aufrufer".

**Fix 1 — `cmd/dreego/docs.go`:** `fetchDocLocal` liefert für Pfade, die NICHT unter `plugins/` liegen, jetzt `nil, false, nil` — OHNE intern den Fallback aufzurufen. Der Aufrufer (`cmdDocs`/`cmdDump`) macht den Fallback. `fromLocal=true` nur wenn die lokale Plugin-Docs geliefert wurde.

**Fix 2 — `cmd/dreego/docs_test.go`:** `TestFetchDocLocalFallback` auf korrekte Semantik umgestellt: erwartet `(nil, false, nil)` und prüft per Fallback-Zähler, dass `fetchDocLocal` den Fallback NICHT intern aufruft (0 Calls). Zusätzlich `TestCmdDocsFallbackCalledOnce` (Integrationstest): ruft `cmdDocs` mit überschriebenem Fallback auf und verifiziert, dass der Fallback GENAU 1× aufgerufen wird (Fallback-Zähler).

Verification: `go test ./cmd/dreego/... ./core/... ./dreegotest/...` all GREEN; full `_tests/test.sh` 164 passed / 0 failed.

next: review4b
- [ ] Increment 5: docs-embed.1

### Increment 5 — tests written (coder-test5)

**Strukturelle Entscheidung (Embedding-Lösung):** `cmd/dreego` ist ein eigenes Modul; `_docs/` liegt bei `../_docs`. `//go:embed` kann NICHT über `..` hinausgehen. Gewählter pragmatischer Weg = **Variante (b)**: Ein Verzeichnis `cmd/dreego/embedded/` wird per `//go:embed` eingebettet; ein Build/Pre-Commit-Schritt kopiert `_docs/`, `README.md`, `CHANGELOG.md` dorthin. Damit bleibt das Embedding innerhalb des Moduls (kein `..`), und die eingebetteten Dateien sind Teil des CLI-Binaries.

**Definierte API (in `cmd/dreego/docs.go`, noch nicht implementiert):**

```go
// embeddedDocs ist das eingebettete Docs-FS (//go:embed embedded/*).
// In Tests überschreibbar (t.TempDir()-FS).
var embeddedDocs fs.FS

// fetchDocEmbedded liest eine Docs-Datei aus embeddedDocs.
// Der führende Slash des URL-Pfads wird auf die FS-Wurzel gemappt.
func fetchDocEmbedded(path string) ([]byte, error)

// fetchDocFallback zeigt auf fetchDocEmbedded (statt auf fetchDoc/HTTP).
var fetchDocFallback = fetchDocEmbedded
```

**Tests (alle RED — `embeddedDocs`/`fetchDocEmbedded` fehlen, Kompilierung schlägt fehl):**
1. `TestFetchDocEmbeddedReadsEmbeddedFS` — `fetchDocEmbedded("/_docs/index.md")` liefert den eingebetteten Body (kein Netzwerk). RED.
2. `TestFetchDocEmbeddedMissingFile` — fehlende eingebettete Datei → Fehler. RED.
3. `TestFetchDocFallbackPointsToEmbedded` — `fetchDocFallback` zeigt per Funktionszeiger auf `fetchDocEmbedded`, NICHT auf `fetchDoc` (HTTP). RED.
4. `TestFetchDocFallbackUsesEmbedded` — Fallback liefert eingebetteten Body (beweist: kein `http.Get`). RED.
5. `TestFetchDocEmbeddedCentralFiles` — `_docs/index.md`, `README.md`, `CHANGELOG.md` offline erreichbar. RED.

Datei: `cmd/dreego/docs_embed_test.go` (Unit-Tests, `t.TempDir()`-FS via `newTestEmbeddedFS`-Helper). Kein Produktivcode geändert.

Verification: `go test ./cmd/dreego/...` → build failed (undefined: embeddedDocs, fetchDocEmbedded).

### Increment 5 — implemented (coder-write5)

Implemented embedded docs (offline `dreego docs`, no HTTP call):

1. **`cmd/dreego/embed.go`** (new, 34 lines): `//go:embed all:embedded` into `embeddedFS`; `var embeddedDocs fs.FS = mustSubFS(embeddedFS, "embedded")` (root mirrors repo layout: `_docs/`, `README.md`, `CHANGELOG.md`); `fetchDocEmbedded(path)` trims the leading slash and reads via `fs.ReadFile(embeddedDocs, rel)`.
2. **`cmd/dreego/docs.go`**: `var fetchDocFallback = fetchDocEmbedded` (was `fetchDoc`/HTTP). `fetchDoc` (HTTP) stays as unused code; `http` import still used by `openBrowser`/`fetchDoc`. Local-first flow from Increment 4 intact (plugin docs → embedded fallback).
3. **`_scripts/sync-embedded-docs.sh`** (new): `rm -rf cmd/dreego/embedded`, copies `_docs/`, `README.md`, `CHANGELOG.md` from repo root. Run before committing after doc changes. Executed once to create the committed mirror.
4. **`cmd/dreego/embedded/`**: created via the sync script (16 `_docs/*.md` + `README.md` + `CHANGELOG.md`).

Verification: `go test ./cmd/dreego/...` (5 embed tests + 4 extensibility tests) GREEN; `./core/... ./dreegotest/...` GREEN; `go build ./cmd/dreego` ok; full `_tests/test.sh` 164 passed / 0 failed; `sh _scripts/sync-embedded-docs.sh` runs clean; built binary serves `docs /_docs/index.md`, `--json`, and `--dump` fully offline.

next: review5
- [ ] Increment 6: frontmatter.1

### Increment 6 — tests written (coder-test6)

**Architektur-Entscheidung — gewählte Option C:** `core/` muss dependency-frei bleiben (nur Standardbibliothek). Ein voller YAML-Parser (gopkg.in/yaml.v3) wäre eine externe Dependency und ist ausgeschlossen. Stattdessen ein **minimaler, dependency-freier Frontmatter-Parser in core**, der einfache YAML-ähnliche `key: value`-Zeilen und einfache Listen unterstützt, aber NICHT das volle YAML-Spezifikum (keine verschachtelten Maps, keine Anchors, kein Multi-Document-YAML). Das ist für Docs/Blogs-Metadaten (title, date, description, tags) ausreichend. Option B (Parsing in cmd/Plugin) verworfen, weil `c.Data` in core lebt und der Parser das Ergebnis als map an core liefern müsste — unnötige Komplexität und Verschleppung der Abhängigkeit.

**Definierte API (in `core/frontmatter.go`, noch nicht implementiert — RED):**

```go
// ParseFrontmatter splits a leading YAML-like frontmatter block off src.
// The opening "---\n" and the closing "---" delimit the block (only at the
// very top of src). Each "key: value" line becomes a map entry; the first ':'
// splits key from value, so a ':' inside the value is preserved. A list value
// "tags: [go, web]" is normalized to the string "go, web". src without a
// leading frontmatter block yields an empty map and the whole src as body.
func ParseFrontmatter(src string) (frontmatter map[string]string, body string)
```

**Unterstützte Frontmatter-Features:**
- Einzelne `key: value`-Zeilen, getrennt durch den ersten `:`. Wert = alles danach (Doppelpunkt im Wert bleibt erhalten, z.B. `description: "a: b"` → `"a: b"`).
- Mehrzeilige key:value-Blöcke (title, date, description, ...).
- Einfache Listenwerte `tags: [go, web]` → normalisiert zu String `"go, web"`.
- Kein volles YAML: keine verschachtelten Maps, keine Block-Skalare, keine Anchors/Tags.
- Body = kompletter Rest nach dem schließenden `---`.

**Tests (alle RED — Kompilierung schlägt fehl, `ParseFrontmatter` fehlt):**
1. `TestParseFrontmatterExtractsKeys` — `---\ntitle: X\n---\n<body>` → map{title:X}, Body=`<body>`.
2. `TestParseFrontmatterNoFrontmatter` — kein `---` → leere map, Body=kompletter Text.
3. `TestParseFrontmatterMultiLineValues` — title/date/description mehrzeilig → 3 Keys.
4. `TestParseFrontmatterListValue` — `tags: [go, web]` → `"go, web"`.
5. `TestParseFrontmatterColonInValue` — `description: "a: b"` → Wert `"a: b"`.
6. `TestParseFrontmatterIntegratesWithData` — via `c.Set`/`c.Data("title")` Integration mit SSRContext.

Datei: `core/frontmatter_test.go` (106 Zeilen). Kein Produktivcode geändert.

Verification: `go test ./core/ -run TestParseFrontmatter` → build failed (`undefined: ParseFrontmatter`). `go build ./core` → ok (Produktivcode kompiliert).

### Increment 6 — implemented (coder-write6)

Implemented `ParseFrontmatter` in `core/frontmatter.go` (new, stdlib only, no external deps):

1. **`ParseFrontmatter(src) (map[string]string, string)`** — returns `(nil, src)` unless src starts with `"---\n"`. Finds the closing `"\n---"` line; everything between is the block, everything after (minus one leading newline) is the body. No closing `---` → `(nil, src)`.
2. **`parseBlock(block)`** — per line: skip blank and `#`-comment lines; split on the first `:`; trim both sides; strip surrounding double quotes from the value (e.g. `description: "a: b"` → `a: b`); key must be non-empty.
3. **`normalizeValue(v)`** — a list value `[go, web]` becomes `go, web` (brackets removed, comma-joined with `, `); non-list values returned unchanged.
4. **Side fix `core/context.go`** — `NewSSR` is now nil-request-safe (uses `context.Background()` when `r == nil`). Required by `TestParseFrontmatterIntegratesWithData` which calls `NewSSR(nil, nil)`.

File: `core/frontmatter.go` (66 lines, <300). No comments except the doc comment on `ParseFrontmatter`.

Verification: `go test ./core/ -run TestParseFrontmatter` → 6/6 PASS. `go test ./core/... ./cmd/dreego/... ./dreegotest/...` all GREEN. Full `_tests/test.sh` 164 passed / 0 failed.

next: review6
- [ ] Increment 7: dev-server.1

### Increment 7 — tests written (coder-test7)

**Testbarer Kern (Manager-Vorgabe):** Ein Daemon-Feature ist schwer end-to-end zu testen. Deshalb definiert dieser Increment einen deterministisch testbaren Kern in `cmd/dreego/` — reine Funktionen ohne Hintergrund-Daemon:

```go
// detectChanges scans dir for .dreego files and compares their modtimes
// against the previous map. It returns the changed files (relative paths)
// plus the updated mtime map. A file is "changed" when it is new, its
// modtime moved, or it disappeared from the previous map.
func detectChanges(dir string, mtimes map[string]time.Time) (changed []string, updated map[string]time.Time)

// shouldRestart reports whether a server restart is required. Any .dreego
// change needs codegen + rebuild, so it is true whenever changed is non-empty.
func shouldRestart(changed []string) bool
```

**Semantik (fixiert):**
- `detectChanges` betrachtet NUR `.dreego`-Dateien (Suffix `.dreego`); andere Dateien (z.B. `.md`) werden ignoriert.
- Neu = in `mtimes` nicht vorhanden → als geändert gemeldet, in `updated` aufgenommen.
- Modifiziert = ModTime weicht von `mtimes` ab → als geändert gemeldet, `updated` bekommt die neue ModTime.
- Entfernt = in `mtimes` vorhanden, aber nicht mehr auf der Platte → als geändert gemeldet, aus `updated` entfernt.
- `shouldRestart` = `len(changed) > 0`.

**Tests (alle RED — Kompilierung schlägt fehl, `detectChanges`/`shouldRestart` fehlen):**
1. `TestDetectChangesNewFile` — neue `.dreego`-Datei wird erkannt.
2. `TestDetectChangesNoChange` — unveränderte Dateien → keine Änderung, Map unverändert.
3. `TestDetectChangesModified` — ModTime-Änderung (via `os.Chtimes`) wird erkannt, `updated` bekommt neue ModTime.
4. `TestDetectChangesRemoved` — gelöschte Datei (in previous, nicht mehr da) wird erkannt, aus `updated` entfernt.
5. `TestDetectChangesIgnoresNonDreego` — `.md`-Datei wird ignoriert.
6. `TestShouldRestartTrueOnChange` — `.dreego`-Änderung → Restart nötig.
7. `TestShouldRestartFalseOnNoChange` — `nil`/leere Liste → kein Restart.

Datei: `cmd/dreego/dev_test.go` (Unit-Tests, `t.TempDir()`). Kein Produktivcode geändert.

Verification: `go test ./cmd/dreego/...` → build failed (undefined: detectChanges, shouldRestart). `go build ./cmd/dreego` → ok (Produktivcode kompiliert).

next: coder-write7

### Increment 7 — implemented (coder-write7)

Implemented the `dreego dev` dev-server in `cmd/dreego/`:

1. **`cmd/dreego/dev.go`** (new):
   - `detectChanges(dir, mtimes) ([]string, map[string]time.Time)` — walks `dir`, considers only `.dreego` files (suffix filter), ignores everything else. New file (absent from `mtimes`) or modtime mismatch → reported as changed + added/updated in `updated` map. Files present in `mtimes` but no longer on disk (removed) → reported as changed + absent from `updated`.
   - `shouldRestart(changed) bool` — `len(changed) > 0`.
   - `cmdDev(args)` — daemon loop: `core.Run(false)` + `cmdBuild(nil)` once, then `startServer(bin)`. A `time.Ticker` (500ms) applies `detectChanges` to the working dir each tick; on change it regenerates (`cmdBuild`) and restarts the running binary via `restartServer` (SIGTERM + `Wait`, then `startServer`). SIGINT/SIGTERM via `signal.Notify` stops the server and exits cleanly.
   - `startServer`/`restartServer` helpers mirror `cmdRun`'s binary start pattern (`exec.Command(bin)`, stdout/stderr wired to the CLI).
2. **`cmd/dreego/main.go`** — added `case "dev": cmdDev(os.Args[2:])` to the main switch; registered `dev` in `printHelp` (commands + example).

Rules respected: only `dev.go` (new) + `main.go` (switch/help) changed; stdlib only (`os`, `os/exec`, `os/signal`, `path/filepath`, `strings`, `syscall`, `time`) + existing `core` import; `dev.go` 177 lines <300.

Verification: `go test ./cmd/dreego/...` GREEN (7 dev tests + existing); `go test ./core/... ./dreegotest/...` GREEN; `go build ./cmd/dreego` ok; `go vet ./cmd/dreego/...` ok; full `_tests/test.sh` 164 passed / 0 failed.

next: review7
- [ ] Final: version bump v0.0.25, docs, commit

### Increment 7 — review fixes (coder-write7b)

Fixed the three reviewer blockers in `cmd/dreego/`:

1. **Blocker 1 — unprimed `mtimes`:** `cmdDev` now primes the baseline before the loop with `_, mtimes := detectChanges(wd, map[string]time.Time{})` so the first tick diffs against a real baseline instead of treating every `.dreego` file as new.
2. **Blocker 2 — build error killed watcher:** extracted `cmdBuildE(args) error` from `cmdBuild` in `main.go` (generation + `go build`, returns `error`, no `os.Exit`). `cmdBuild` wraps it (`if err != nil { print; os.Exit(1) }`). `cmdDev` calls `cmdBuildE(nil)`, on error prints to stderr and `continue` — the old server keeps running, rebuild only on the next valid change. `cmdDev` startup build also uses `cmdBuildE`.
3. **Blocker 3 — double signal + missing Wait in SIGINT path:** replaced `Signal(TERM)`+`Kill()` with `Signal(SIGTERM)`+`Process.Wait()`. `restartServer` documented (TERM + `Wait`; a misbehaving server could hang the dev tool, acceptable).

**Tests (dev_test.go, +2):**
- `TestDetectChangesInitialPriming` — first scan with empty map establishes the baseline; second scan unchanged reports no changes (Blocker 1).
- `TestCmdBuildEErrorsInsteadOfExit` — `cmdBuildE(nil)` in a `t.TempDir()` without a valid module returns an error instead of `os.Exit` (Blocker 2).

Files changed: `cmd/dreego/dev.go` (136 lines), `cmd/dreego/main.go` (294 lines, cmdBuild→cmdBuildE refactor), `cmd/dreego/dev_test.go` (177 lines). All <300.

Verification: `go test ./cmd/dreego/... ./core/... ./dreegotest/...` GREEN (incl. 2 new tests); `go build ./cmd/dreego` ok; `go vet ./cmd/dreego/...` ok; full `_tests/test.sh` 164 passed / 0 failed.

next: review7b

## Plugin Interface (frozen v1 contract)

Monolithic, paket-level. Plugins are external Go modules that import core and satisfy this interface. Core never imports a plugin.

```go
// core/plugin.go
package core

import (
    "context"
    "io/fs"
    "net/http"
)

type Plugin interface {
    Name() string
    RegisterRoutes() // plugin calls dreego.Register(...) internally
    Middlewares() []func(http.Handler) http.Handler
    Assets() fs.FS
    OnStart(ctx context.Context) error
    OnShutdown(ctx context.Context) error
}

func UsePlugin(p Plugin) // registers routes, collects middleware, assets, lifecycle
```

## API defined by tests (increment 1, coder-test1)

`core/plugin_test.go` (package core) fixes the following contract. It is RED:
compilation fails because `Plugin`, `UsePlugin`, `StartPlugins`, `ShutdownPlugins` do not exist yet.

- `UsePlugin(p Plugin)`:
  - calls `p.RegisterRoutes()` (the plugin calls `Register(...)` internally, like an external module importing core);
  - collects `p.Middlewares()` for injection into the `Build()` stack (FIFO);
  - stores `p.Assets()` (not asserted in this increment);
  - registers `p.OnStart`/`p.OnShutdown` for later lifecycle invocation.
- Lifecycle (fixated for v1): `StartPlugins(ctx)` calls `OnStart` on every registered plugin; `ShutdownPlugins(ctx)` calls `OnShutdown`. Both return an error if any plugin method returns one, and propagate the first error to the caller.
- Middleware sampling: collected plugin middleware must run inside `Build()`/`ServeMux()` (test asserts a response header set by plugin middleware).

Tests (4 + 1 edge):
1. `TestPluginInterfaceSatisfaction` — `var _ Plugin = (*testPlugin)(nil)` compile-time check.
2. `TestUsePluginRegistersRoutes` — plugin route reachable via `ServeMux()` after `UsePlugin` + `Build()` (httptest).
3. `TestUsePluginCollectsMiddleware` — plugin middleware header applied in the stack.
4. `TestUsePluginLifecycle` — `StartPlugins`/`ShutdownPlugins` call OnStart/OnShutdown on every plugin.
5. `TestUsePluginLifecycleErrors` — lifecycle errors propagate to the caller.

Usage (user-confirmed import alias):
```go
import dreego "github.com/dreego-stack/dreego/core"
import auth "github.com/example/dreego-auth"

func main() {
    dreego.UsePlugin(auth.New("secret"))
    dreego.Listen(":8080")
}
```
