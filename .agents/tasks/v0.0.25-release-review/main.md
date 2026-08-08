# AKTE: v0.0.25-release-review — Release-Readiness-Review

## Titelseite (für ALLE Subagents)

- **Task-ID:** v0.0.25-release-review
- **Titel:** Release-Readiness-Review für v0.0.25
- **Status:** in Arbeit
- **Repo:** dreego (Go Web Framework, Monorepo, go.work)
- **Version:** v0.0.25 (unreleased, VERSION-Datei = v0.0.25)

### Kontext (alle Agents lesen)

Der User möchte sicherstellen, dass v0.0.25 korrekt ist und alles enthält, was gewollt ist. Dazu wird ein umfassender Review durchgeführt: CHANGELOG vs. tatsächliche Commits, TODO.md-Items vs. Stand, komplette Test-Suite, Versions-Konsistenz.

### Was v0.0.25 laut CHANGELOG enthalten soll

- **plugin-interface.1** — frozen v1 `core.Plugin` contract + `core.UsePlugin`
- **middleware-hooks.1** — plugin middleware via `app.Use` (FIFO chain)
- **route-hooks.1** — programmatic plugin route registration
- **docs-extensibility.1** — `dreego docs` resolves plugin docs from `plugins/<name>/_docs/`
- **docs-embed.1** — offline embedded docs (`cmd/dreego/embedded/`)
- **frontmatter.1** — YAML frontmatter parsing + typed metadata
- **dev-server.1** — `dreego dev` with file watcher + auto-regenerate
- **head-dedupe.1** — route `<title>`/`<meta description>` wins over layout
- Weitere Commits seit v0.0.24 (siehe git log): tailwind-plugin.1 research, --version/-v flags, ctx vs c naming docs, go-section string fixes, section parsing fixes, scope div before doctype, $loop substitution, prop defaults

### Workflow

1. **reviewer** — umfassender Release-Review (siehe Kapitel)
2. Bei Befunden: zurück an coder-implement → reviewer → shell (Commit)
3. Bei Freigabe: Abschluss

### Regeln (alle Agents)

- ALLE Terminal-Befehle via `smd` (Container). Relative Pfade, KEINE Host-Pfade wie /Users/lukas/...
- Repo-Sprache: Englisch. Chat mit User: Deutsch.
- Jeder Agent schreibt sein Ergebnis in sein Kapitel.
- Agentlog via `agentlog`-Tool bei Statuswechseln.

---

## Kapitel: reviewer (nur dieser Agent)

**Aufgabe:** Umfassender Release-Readiness-Review für v0.0.25. Du änderst NICHTS, nur prüfen und Befunde dokumentieren.

Prüfe systematisch:

1. **CHANGELOG vs. Commits:** Lies `CHANGELOG.md` (v0.0.25-Eintrag, Zeilen 1-15). Vergleiche mit `smd git log --oneline` seit dem letzten v0.0.24-Release. Ist jedes Feature dokumentiert? Fehlen Einträge? Stimmen die Behauptungen (Dateien, Tests, Verhalten)?
2. **TODO.md vs. Stand:** Lies `TODO.md` (v0.0.25-Sektion, Zeilen 7-16). Sind alle 8 Items wirklich umgesetzt? Prüfe die zugehörigen Dateien: `core/plugin.go`, `core/route_hooks_test.go`, `cmd/dreego/docs.go`, `cmd/dreego/embedded/`, `core/frontmatter.go`, `cmd/dreego/dev.go`, `core/codegen_head_dedupe.go`.
3. **Test-Suite:** Führe aus: `smd go test ./core` (Unit) und `smd go test ./cmd/dreego` (CLI-Unit). Optional: `smd _tests/test.sh` (Integration, dauert länger — nur wenn sinnvoll, sonst dokumentieren dass nicht gelaufen).
4. **Versions-Konsistenz:** `VERSION`-Datei = v0.0.25? `go.work` enthält alle Module (core, cmd/dreego, dreegotest, plugins/sample, plugins/sse)? Stimmen die go.mod-Versionen?
5. **Offene Baustellen:** Gibt es bekannte offene Punkte, die v0.0.25 blockieren? (TODO.md "Available Next", Quality Backlog, Rejected)
6. **Git-Zustand:** `smd git status --short` — sind uncommittete Änderungen relevant für den Release? (Test-Artefakte wie main.go, go.mod, dreego-linux sind normal und NICHT Teil des Releases)

**Ergebnis:** Schreibe in dein Kapitel:
- **APPROVE** (Release-ready) ODER **FAIL: <Befunde>** mit Liste
- Pro Befund: Schweregrad (Blocker/Major/Minor), Datei, Beschreibung
- Test-Ergebnisse (welche Suiten grün/rot)
- Empfehlung: Release durchführen oder erst fixen?

---

## Ergebnis: reviewer (ses_01ec7f08affeGVIgNx9Fa67eO0, 2026-08-08)

### VERDICT: FAIL — 1 Blocker, 1 Minor

### Befunde

**🔴 Blocker — `cmd/dreego/embedded/CHANGELOG.md` ist nicht synchron**
- Datei: `cmd/dreego/embedded/CHANGELOG.md` (vs. `CHANGELOG.md`)
- Beschreibung: Die embedded Kopie weicht in den historischen Einträgen v0.0.22/v0.0.23 vom Repo ab (2 Zeilen: `_tests/core/Template/component-nested-if-else/` vs. `_tests/core/Bugs/...`, fehlender Satz im v0.0.23-Eintrag). Der Integrationstest `core/DocsSync` schlägt deshalb fehl: **184 passed, 1 failed**. Der Test ist Teil der Release-Suite und der pre-push-Hook (`make test`) blockiert den Push. Fix: `_scripts/sync-embedded-docs.sh` ausführen und die synchronisierte Datei committen (Commit `91820f4` hat CHANGELOG.md geändert, aber die embedded Kopie nicht nachgezogen).
- Hinweis: Der v0.0.25-Eintrag selbst ist in der embedded Kopie bereits enthalten — nur die historischen Zeilen divergieren.

**🟡 Minor — CHANGELOG v0.0.25 "Full suite: 164 passed" nicht verifizierbar**
- Datei: `CHANGELOG.md:14`
- Beschreibung: Die Behauptung "164 passed, 0 failed" bezieht sich auf die Unit-Suiten. Tatsächlich gemessen: `core` 381 PASS / 0 FAIL (402 RUN), `cmd/dreego` 35 PASS / 0 FAIL, `go vet` sauber, alle Module bauen (`plugins/sse`, `plugins/sample`, `dreegotest`). Die Zahl 164 ist weder die Summe (416) noch eine der Einzelzahlen — vermutlich ein veralteter Stand. Empfehlung: Zahl vor Release korrigieren oder als "Unit-Suiten grün" formulieren.

### Prüfpunkte im Detail

1. **CHANGELOG vs. Commits:** Alle 8 v0.0.25-Features sind im CHANGELOG dokumentiert und durch Commits belegt (9cd8914 plugin-interface.1, 0c01e7b middleware-hooks.1, 00c68e8 route-hooks.1, 87fa4c6 docs-extensibility.1, e424899 docs-embed.1, ddaf242 frontmatter.1, 84cd904 dev-server.1, 91820f4 head-dedupe.1). Zusätzliche Commits (tailwind-plugin.1 research, --version/-v, ctx vs c, go-section strings, section parsing, scope div, $loop, prop defaults) sind im CHANGELOG-Kontext (Zeile 25 der Akte) erwähnt, aber NICHT als eigene CHANGELOG-Einträge — akzeptabel, da Fixes/Research. **Register-Idempotenz-Behauptung stimmt** (`core/runtime.go:58-66` ersetzt Handler bei gleichem method+pattern). **BREAKING-Hinweis korrekt**.
2. **TODO.md vs. Stand:** Alle 8 Items umgesetzt und verifiziert: `core/plugin.go` (frozen v1 Contract, UsePlugin, StartPlugins/ShutdownPlugins), `core/route_hooks_test.go` (5 Tests inkl. last-wins-Overlap), `cmd/dreego/docs.go` (fetchDocLocal mit plugins/-Priorität), `cmd/dreego/embedded/` (existiert, _docs/ + README + CHANGELOG), `core/frontmatter.go` (ParseFrontmatter + normalizeValue), `cmd/dreego/dev.go` (500ms Poll, SIGTERM+reap, Build-Fehler killen Watcher nicht), `core/codegen_head_dedupe.go` (dedupeHeadMerge, route wins). TODO.md listet nur 7 Items (head-dedupe.1 fehlt in der TODO-Liste, steht aber im CHANGELOG + ist umgesetzt) — Minor-Diskrepanz, kein Blocker.
3. **Test-Suite:** `go test -count=1 ./core` ✅ grün (381 PASS/0 FAIL), `go test -count=1 ./cmd/dreego` ✅ grün (35 PASS/0 FAIL), `go vet ./core ./cmd/dreego` ✅, Plugin-Module bauen ✅. **Integration `_tests/test.sh`: ❌ 184 passed, 1 failed (`core/DocsSync`)** — siehe Blocker.
4. **Versions-Konsistenz:** `VERSION` = v0.0.25 ✅. `go.work` enthält alle 5 Module (cmd/dreego, core, dreegotest, plugins/sample, plugins/sse) ✅. `cmd/dreego/go.mod` require core v0.0.25 ✅ (Commit 8d129e3). **Minor:** `plugins/sample/go.mod` require core **v0.0.23** (veraltet, aber sample nutzt nur `dreego.NewSSR` als no-op — baut trotzdem). `dreegotest` + `plugins/sse` nutzen `v0.0.0` + replace (lokal) — konsistent mit Monorepo-Praxis. `_scripts/release.sh` taggt nur core, cmd/dreego, plugins/sample — **plugins/sse wird nicht getaggt** (Minor, sse ist Beispiel-Plugin).
5. **Offene Baustellen:** "Available Next" (observability.1, documentation.1) und Rejected (hot-reload/live-reload/smart-recompile → Air) blockieren v0.0.25 NICHT. Quality Backlog: codegen-errors.2 als erledigt markiert ✅. Nichts blockiert.
6. **Git-Zustand:** Uncommittet: `go.sum` (enthält core v0.0.24-Hashes — Artefakt aus lokalem Test, NICHT release-relevant, aber sollte vor Release nicht committet werden), `.agents/log.md` (normal), Test-Artefakte (main.go, go.mod, dreego-linux, cmd/dreego/dreego, .tmp-Ordner, .DS_Store) — laut Akte normal, NICHT Teil des Releases. Kein uncommitteter Code, der den Release betrifft. HEAD = f01d07b.

### Empfehlung

**Erst fixen, dann releasen:**
1. `_scripts/sync-embedded-docs.sh` ausführen → `cmd/dreego/embedded/CHANGELOG.md` synchronisieren → committen (behebt den einzigen roten Integrationstest).
2. Danach `_tests/test.sh` erneut laufen lassen — muss 185/185 grün sein.
3. Optional (Minor): CHANGELOG-Zahl 164 korrigieren, `plugins/sample/go.mod` auf v0.0.25 heben, sse-Tag in release.sh ergänzen.
4. Erst dann `_scripts/release.sh` + Tags pushen.

**Kein Release in diesem Zustand** — der DocsSync-Fail bricht den pre-push-Hook (`make test`) und die embedded Docs wären beim Release nachweislich veraltet.

---

## Kapitel: coder (nur dieser Agent)

**Aufgabe:** Blocker-Fix — embedded Docs synchronisieren, DocsSync-Test verifizieren.

### Ergebnis: coder (ses_01ebe0c1affemgZt18mgtSrHie, 2026-08-08)

**Blocker behoben.** Der einzige rote Integrationstest (`core/DocsSync`) ist jetzt grün.

1. **Sync ausgeführt:** `sh _scripts/sync-embedded-docs.sh` (direkter Aufruf scheiterte mit "permission denied" — Skript ist nicht executable, daher via `sh`). Skript spiegelt `_docs/`, `README.md`, `CHANGELOG.md` nach `cmd/dreego/embedded/`.
2. **Diff verifiziert:** `git diff --stat cmd/dreego/embedded/` → nur **1 Datei** geändert: `cmd/dreego/embedded/CHANGELOG.md` (2 insertions, 2 deletions). README und `_docs/` waren bereits synchron. Die 2 Zeilen entsprechen exakt dem Reviewer-Befund: v0.0.22-Eintrag (`_tests/core/Template/component-nested-if-else/` vs. `_tests/core/Bugs/...`) und v0.0.23-Eintrag (fehlender Satz zur component-path-Variante). Der v0.0.25-Eintrag war bereits enthalten.
3. **Test verifiziert:** Runner unterstützt Filter via `DREEGO_FILTER` (Regex auf Testname, `_tests/test.sh:51`). Aufruf: `env DREEGO_FILTER=DocsSync sh _tests/test.sh` → **PASS, 1 Passed, 0 Failed**. Danach komplette Suite: `sh _tests/test.sh` → **185 Passed, 0 Failed** (vorher 184/1). GO Tests PASS.
4. **Offen für shell:** `cmd/dreego/embedded/CHANGELOG.md` ist uncommittet (M) und muss committet werden. Hinweis: `_scripts/sync-embedded-docs.sh` ist nicht executable (chmod +x fehlt) — Minor, optional fixen.

**Empfehlung:** Blocker ist behoben, Release kann nach Commit + optionalen Minor-Fixes (CHANGELOG-Zahl 164, `plugins/sample/go.mod` v0.0.23→v0.0.25) durchgeführt werden. next: reviewer → shell (Commit).

## Re-Review: reviewer (ses_01ebb5a14ffehae6eslbzMvYC3, 2026-08-08)

### VERDICT: APPROVE — Release-ready

### 1. Fix verifiziert ✅
- `git diff cmd/dreego/embedded/CHANGELOG.md`: **exakt 2 Zeilen** geändert (v0.0.22-Eintrag: `_tests/core/Template/component-nested-if-else/`; v0.0.23-Eintrag: ergänzter Satz zur component-path-Variante). Stimmt mit coder-Dokumentation und meinem ursprünglichen Befund überein.
- Sync-Richtung korrekt: `diff CHANGELOG.md cmd/dreego/embedded/CHANGELOG.md` → **IDENTICAL** (embedded == Repo).
- `env DREEGO_FILTER=DocsSync sh _tests/test.sh` → **PASS, 1 Passed, 0 Failed** ✅

### 2. Komplette Suite ✅ (nach 2. Lauf)
- Lauf 1: `184 Passed, 1 Failed` (`core/Bugs/run-timer-sigterm`) — **Flakiness, kein Regression**: Test isoliert (`DREEGO_FILTER=run-timer-sigterm`) → PASS 1/0. Ursache: 15s-SIGTERM-Timeout unter paralleler Last (test.sh:65-68, JOBS=nproc). 
- Lauf 2: `sh _tests/test.sh` → **185 Passed, 0 Failed** ✅
- GO Tests: PASS (test.sh:14-20, `./core/...` + `./cmd/dreego/...`).

### 3. Minor-Punkte (alle NICHT release-blockierend)
$%$ Punkt $%$ Bewertung $%$ Empfehlung $%$
$%$---$%$---$%$---$%$
$%$ `_scripts/sync-embedded-docs.sh` nicht executable (`-rw-r--r--`, root) $%$ **Minor, kein Blocker** — Skript ist Tooling, wird nur via `sh` aufgerufen; Release-Pfad und pre-push-Hook hängen nicht davon ab. $%$ Optional als eigener Mini-Commit: `chmod +x` + Commit. Nicht Release-blockierend. $%$
$%$ `plugins/sample/go.mod` require core **v0.0.23** $%$ **Minor, nicht release-relevant** — Sample nutzt nur no-op `dreego.NewSSR`, baut grün, ist Beispiel-Code. $%$ Nach Release auf v0.0.25 heben (kein Eiltempo). $%$
$%$ `_scripts/release.sh` taggt `plugins/sse` nicht $%$ **Minor, nicht release-relevant** — sse nutzt v0.0.0+replace (lokal, konsistent), wurde nie getaggt (auch v0.0.22/v0.0.23 nicht). $%$ Bewusst so lassen; ggf. als TODO für später. $%$

### 4. Commit-Gate: APPROVE
**Zu committen (1 Datei):** `cmd/dreego/embedded/CHANGELOG.md` — der Blocker-Fix.
**Optional im selben oder eigenem Commit:** chmod +x für `_scripts/sync-embedded-docs.sh` (trivial, sauberer Hygiene-Fix).
**NICHT committen:** `go.sum` (enthält core v0.0.24-Hashes — Test-Artefakt), `.agents/log.md` (normal), alle Test-Artefakte (main.go, go.mod, dreego-linux, cmd/dreego/dreego, .dreego-review-bin, .DS_Store, _tests/core/CLI/check-stale/dreego/gen/, smd.toml.example).

### Empfehlung
**Release durchführen.** Blocker behoben, Suite komplett grün (185/185), pre-push-Hook (`make test`) läuft durch. Minor-Punkte als Follow-up, nicht als Blocker behandeln.

---

## Kapitel: shell (nur dieser Agent)

**Aufgabe:** Commit des Reviewer-freigegebenen Blocker-Fixes.

### Ergebnis: shell (2026-08-08)

- Commit: `0e9c8c7` — `docs(embed): sync embedded CHANGELOG with repo`
- Body: "Re-run _scripts/sync-embedded-docs.sh after the component-nested-if-else path correction; fixes the DocsSync integration test."
- Enthalten: `cmd/dreego/embedded/CHANGELOG.md` (2+/2-, Blocker-Fix) + `_scripts/sync-embedded-docs.sh` (chmod +x Hygiene-Fix, 100644→100755).
- NICHT committet: `go.sum`, `.agents/log.md`, `.agents/tasks/`, Test-Artefakte (main.go, go.mod, dreego-linux, cmd/dreego/dreego, .dreego-review-bin, .tmp-review-cli, smd.toml.example, _tests/core/CLI/check-stale/dreego/gen/).
- Branch: main. Commit-Gate: APPROVE durch reviewer (ses_01ebb5a14ffehae6eslbzMvYC3).
