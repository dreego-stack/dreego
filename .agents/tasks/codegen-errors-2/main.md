# AKTE: codegen-errors-2 — Absicherungstest + Doku-Update

## Titelseite (für ALLE Subagents)

- **Task-ID:** codegen-errors-2
- **Titel:** Absicherungstest + Doku-Update für Backlog-Eintrag codegen-errors.2
- **Status:** in Arbeit
- **Repo:** dreego (Go Web Framework, Monorepo, go.work)
- **Version:** v0.0.25 (unreleased)

### Kontext (alle Agents lesen)

Der Backlog-Eintrag `codegen-errors.2` in `TODO.md` (Zeile 66) behauptet: `genTemplateNodeComp` (`core/codegen.go:521`) droppe still verschachtelte `{#if}` in nicht-finalen `{#else}`-Zweigen von Komponenten-Templates.

**Analyse-Ergebnis (explore-Agent, Session ses_01ee6079effehaDNTMUGZ9sCSG): Der Bug ist BEREITS BEHOBEN** — in v0.0.22 (Commit `11d33d2`, Feature codegen-errors.1), also VOR dem v0.0.23-Route-Fix, den TODO.md als Muster zitiert.

Beweise:
1. `core/codegen.go:521` existiert nicht mehr — CodeGen wurde in v0.0.22 in Einzeldateien aufgeteilt; Komponenten-Generierung liegt in `core/codegen_component.go`.
2. Chain-Detection (else-if vs. echter else) ist eingebaut: `core/codegen_component.go:48-85` — `chain := true`; Schleife über `n.ElseChildren`; `if ec.Type != NodeIf { chain = false; break }`; `if chain` → `} else if`-Kette, `else` → echter `} else {`-Block mit rekursivem `g.node(ec)`.
3. Regressionstest existiert: `TestGenTemplateNodeCompNestedIfInElseNotDropped` (`core/codegen_template_test.go:99`) — läuft GRÜN.
4. Integrationstest existiert: `_tests/core/Template/component-nested-if-else/test.sh` (nicht unter `_tests/core/Bugs/` wie CHANGELOG behauptet).
5. CHANGELOG.md Zeile 40 (v0.0.22) dokumentiert den Fix explizit.

### Entscheidung (User, 2026-08-08)

**Absicherungs-Test + Doku-Update** (nicht "nur Doku", nicht "nichts tun"):
1. Neuer Negativ-Unit-Test in `core/codegen_component_branches_test.go`: echter else-Branch (gemischte ElseChildren, nicht nur NodeIf) wird NICHT als else-if-Kette missbraucht.
2. `TODO.md` Zeile 66: Backlog-Eintrag als erledigt markieren (Verweis auf v0.0.22, Commit `11d33d2`).
3. `CHANGELOG.md` Zeilen 42, 50: Pfad-Korrektur `_tests/core/Bugs/component-nested-if-else/` → `_tests/core/Template/component-nested-if-else/`.

### Workflow (4 Subagents pro Item)

1. **coder-test** — schreibt Tests (positiv + negativ), verifiziert via smd
2. **coder-implement** — erfüllt Tests / macht Doku-Update, verifiziert via smd (darf selbst `smd` nutzen oder shell starten)
3. **reviewer** — Review + Commit-Gate
4. **shell** — Commit (nur nach reviewer-Freigabe)

### Regeln (alle Agents)

- ALLE Terminal-Befehle via `smd` (Container). Relative Pfade, KEINE Host-Pfade wie /Users/lukas/...
- Repo-Sprache: Englisch (Code, Kommentare, Commit-Messages). Chat mit User: Deutsch.
- Max 300 Zeilen pro Datei, keine unnötigen Kommentare, Go 1.22+, Standard-Library bevorzugt.
- Jeder Agent schreibt sein Ergebnis in sein Kapitel unten (Status + was getan wurde).
- Agentlog via `agentlog`-Tool bei Statuswechseln.

---

## Kapitel: coder-test (nur dieser Agent)

**Aufgabe:** Schreibe einen Negativ-Unit-Test in `core/codegen_component_branches_test.go`.

- **Positiv-Fall** (existiert bereits, NICHT duplizieren): `TestGenTemplateNodeCompNestedIfInElseNotDropped` (`core/codegen_template_test.go:99`) — verschachtelter `{#if}` in `{#else}` wird korrekt generiert. Verifiziere nur, dass er grün ist.
- **Negativ-Fall (NEU):** Ein echter else-Branch mit GEMISCHTEN ElseChildren (z.B. ein `NodeIf` UND ein `NodeText`) darf NICHT als else-if-Kette emittiert werden. Konkret: `NodeIf{Cond:"a"}` mit `ElseChildren: [NodeIf{Cond:"b"}, NodeText{"fallback"}]` → muss `} else {` mit rekursiver Verarbeitung emittieren (beide Kinder), KEIN `} else if b {`. Analog zu `TestCompGenIfElseIfChain` (`codegen_component_branches_test.go:90-119`), aber mit gemischtem Fall.
- **Verifikation:** `smd go test ./core -run "TestCompGen|TestGenTemplateNodeComp" -v` — alle Tests müssen GRÜN sein (der neue Test testet bereits korrektes Verhalten, da der Bug behoben ist).
- **Ergebnis:** Schreibe in dein Kapitel: Testname, was er prüft, Test-Output (PASS/FAIL).

**Ergebnis (coder-test, ses_01ef99687ffeF4uIC1zD3RgU2T):**

- **Neuer Test:** `TestCompGenIfElseMixedChildren` in `core/codegen_component_branches_test.go` (nach `TestCompGenIfElseIfChain`, Zeile ~121).
- **Prüft:** `NodeIf{Cond:"a"}` mit `ElseChildren: [NodeIf{Cond:"b", Children:[NodeText "B"]}, NodeText{"fallback"}]` (gemischte ElseChildren) muss `} else {` mit rekursiver Verarbeitung beider Kinder emittieren (`if b {`, `` `B` ``, `` `fallback` ``) und darf KEIN `} else if b {` erzeugen. Deckt die Chain-Detection in `core/codegen_component.go:48-54` ab (nicht-`NodeIf` in ElseChildren → `chain = false` → echter else-Block).
- **Test-Output:** `TestCompGenIfElseMixedChildren` — PASS. Alle 10 `TestCompGen*`-Tests PASS, alle 5 `TestGenTemplateNodeComp*`-Tests PASS (inkl. Positiv-Fall `TestGenTemplateNodeCompNestedIfInElseNotDropped`). Komplette Suite `smd go test ./core` — ok (0.011s).
- **Hinweis:** `-run "TestCompGen|TestGenTemplateNodeComp"` matcht via smd nicht (Pipe wird als Shell-Pipe interpretiert) — Regexes einzeln ausgeführt.
- **Status:** DONE. next: coder-implement

## Kapitel: coder-implement (nur dieser Agent)

**Aufgabe:** Kein Code-Fix nötig (Bug behoben). Doku-Update:

1. `TODO.md` Zeile 66: `codegen-errors.2` als erledigt markieren — z.B. `- **codegen-errors.2** — ✅ bereits gefixt in v0.0.22 (codegen-errors.1, Commit 11d33d2); Backlog-Eintrag obsolet` (oder entfernen, wenn sauberer).
2. `CHANGELOG.md` Zeilen 42 und 50: Pfad `_tests/core/Bugs/component-nested-if-else/` → `_tests/core/Template/component-nested-if-else/` korrigieren.
3. **Verifikation:** `smd go test ./core` (komplette Unit-Suite) muss grün sein. Optional: `smd _tests/test.sh` (Integration, dauert länger — nur wenn sinnvoll).
4. **Ergebnis:** Schreibe in dein Kapitel: welche Dateien geändert, Test-Ergebnis.

**Ergebnis (coder-implement, ses_01ef99687ffeF4uIC1zD3RgU2T):**

- **Geänderte Dateien:**
  - `TODO.md` (Zeile 66): Backlog-Eintrag `codegen-errors.2` ersetzt durch `- **codegen-errors.2** — ✅ already fixed in v0.0.22 (codegen-errors.1, commit 11d33d2); backlog entry obsolete, covered by TestCompGenIfElseMixedChildren`.
  - `CHANGELOG.md` (Zeile 42, v0.0.22): Pfad `_tests/core/Bugs/component-nested-if-else/` → `_tests/core/Template/component-nested-if-else/` korrigiert.
  - `CHANGELOG.md` (Zeile 50, v0.0.23): Pfad-Korrektur + Hinweis ergänzt, dass die Komponenten-Variante durch `_tests/core/Template/component-nested-if-else/` (v0.0.22) abgedeckt ist.
- **Verifikation:** `smd go test -count=1 ./core` — ok (0.013s), komplette Unit-Suite grün. Integrationstest nicht ausgeführt (nur optional, Doku-Update ohne Code-Änderung).
- **Hinweis:** `cmd/dreego/embedded/CHANGELOG.md` (Zeile 42) enthält denselben falschen Pfad, wird aber von `_scripts/sync-embedded-docs.sh` aus dem Repo-CHANGELOG regeneriert — nicht manuell angefasst.
- **Status:** DONE. next: reviewer

## Kapitel: reviewer (nur dieser Agent)

**Aufgabe:** Review des gesamten Items (Test + Doku-Änderungen).

- Prüfe: Deckt der Negativ-Test den gemischten else-Branch-Fall wirklich ab? Ist die Doku-Änderung korrekt und vollständig? Keine Regressionen?
- Commit-Gate: Nur freigeben, wenn würdig. Bei Befund: `FAIL: <Befund> next: coder-implement` in dein Kapitel schreiben.
- **Ergebnis:** Schreibe in dein Kapitel: Freigabe (APPROVE) oder Befund (FAIL) mit Begründung.

## Kapitel: shell (nur dieser Agent)

**Aufgabe:** Commit erstellen — NUR nach reviewer-Freigabe (APPROVE in reviewer-Kapitel).

- `smd git add` der geänderten Dateien (core/codegen_component_branches_test.go, TODO.md, CHANGELOG.md)
- Commit-Message nach Repo-Konvention (siehe `.agents/guides/changelog.md`): z.B. `fix(codegen): add negative else-branch guard test for component templates` + Doku-Update
- **Ergebnis:** Schreibe in dein Kapitel: Commit-Hash + Message.

**Ergebnis (reviewer, ses_01ed0ef62ffeueD5qAlsHlcdjS):**

**APPROVE** — Item ist würdig, freigegeben für Commit.

Befund im Detail:

1. **Test `TestCompGenIfElseMixedChildren`** (`core/codegen_component_branches_test.go:124-154`):
   - Deckt den gemischten Fall real ab: `ElseChildren: [NodeIf{b}, NodeText{"fallback"}]` → Chain-Detection (`core/codegen_component.go:48-54`) setzt `chain=false` → echter `} else {` mit rekursivem `g.node(ec)` für BEIDE Kinder. Positive Asserts (`if a {`, `} else {`, `if b {`, `` `A` ``, `` `B` ``, `` `fallback` ``) + negatives Assert (`} else if b {` darf NICHT vorkommen) — der Test ist diskriminierend: im Chain-Pfad würde `} else {` fehlen (NodeIf b hat keine ElseChildren) und `} else if b {` auftauchen.
   - Stil konsistent mit `TestCompGenIfElseIfChain` (gleiche Struktur: TemplateNode → genTemplateNodeComp → want-Schleife → Negativ-Assert, gleicher Kommentarstil).
   - Keine Duplikation mit `TestGenTemplateNodeCompNestedIfInElseNotDropped` (`codegen_template_test.go:99`): jener ist ein Full-Pipeline-Test (lex→parse→codegen) ohne Negativ-Assert; der neue Test pinnt die exakte Emissionsform als Unit-Test. Komplementär, nicht dupliziert.
2. **TODO.md:66** — korrekt als erledigt markiert (✅, Verweis v0.0.22/Commit 11d33d2 — Commit existiert, per `git log` verifiziert; CHANGELOG v0.0.22:40 bestätigt den Fix).
3. **CHANGELOG.md:42,50** — Pfad-Korrektur korrekt: `_tests/core/Template/component-nested-if-else/test.sh` existiert (verifiziert), alter Pfad `_tests/core/Bugs/component-nested-if-else/` existiert NICHT (glob: 0 Treffer). Zeile 50-Zusatz (Komponenten-Variante via v0.0.22-Test abgedeckt) ist akkurat.
4. **Verifikation:** `smd go test -count=1 ./core` — ok (0.015s). `TestCompGenIfElseMixedChildren` PASS, `TestCompGenIfElseIfChain` PASS. Keine Regressionen.

Hinweise für shell (keine Blocker):
- 🟡 `go.sum` (Root) wurde als Test-Nebeneffekt modifiziert (Hash für `codeberg.org/dreego/dreego/core v0.0.24` aus Root-`go.mod` "gzipflush" — die Root-`go.mod` selbst ist UNTRACKED). Nicht Teil des Tasks → NICHT committen, nur die 3 Task-Dateien.
- 🟡 `.agents/log.md` enthält die Agentlog-Einträge des Tasks (282 Zeilen) — Task-Spec listet nur 3 Dateien; log.md-Änderung ist Workflow-Konvention, Entscheidung bei shell/manager.
- 🟡 `cmd/dreego/embedded/CHANGELOG.md:42` enthält noch den alten Pfad — korrekt NICHT manuell angefasst (wird von `_scripts/sync-embedded-docs.sh` regeneriert, verifiziert).
- 🟡 Viele untracked Artefakte im Repo (.DS_Store, dreego-linux, main.go, .tmp-review-cli, cmd/dreego/dreego-Binary) — `git add` nur die 3 Task-Dateien.

next: shell
