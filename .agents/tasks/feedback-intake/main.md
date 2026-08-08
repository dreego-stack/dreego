# Task: feedback-intake

## Ziel
Die zwei echten Issues aus `.tmp/dreego-feedback.md` (PeerNet-Doku-Erfahrung, dreego v0.0.22) ins Projekt einbringen: reproduzieren, fixen, permanent testen.

## Kontext
- Feedback-Datei: `.tmp/dreego-feedback.md` (Punkt 3 = reiner Umgebungs-Kontext, kein Issue).
- Issue A: Verschachtelte `{#if}` im Route-Template werden still verworfen → `dreego generate` ohne Fehler, leeres Template, Folgefehler `declared and not used`. Workaround: flache Blöcke. Wunsch: Fehler statt Silent-Drop ODER korrekte Generierung.
- Issue B: Ausdrücke in `<head>`-Sektion einer Route werden nicht aufgelöst → `{doc.Title}` erscheint wörtlich. `core/parser_section_head.go` liest Head als rohen String (`readUntilClose`). Wunsch: Head-Ausdrücke auflösen ODER als "statisch" dokumentieren.
- TODO.md Quality Backlog enthält bereits `codegen-errors.1` ("Replace silent fails in CodeGen with explicit errors") — Issue A ist ein konkreter Fall davon.

## Relevante Dateien
- `core/parser_template.go` — parseIfNodes/parseElseNodes (verschachtelte ifs)
- `core/codegen_template.go` — NodeIf-Codegen (Zeile 33-62, Problemzone: `return ""` bei 47)
- `core/parser_section_head.go` — Head als roher String
- `core/codegen_ssr.go` — SSR-Rendering (head → c.Set)
- Tests: `core/codegen_template_test.go`, `_tests/core/<Group>/` Integration

## Status
- [ ] A1: Reproduktionstest schreiben (nested {#if} im else-Zweig) → muss FAIL
- [ ] A2: Fix im Parser/Codegen (nested {#if} korrekt generieren ODER expliziter Fehler)
- [ ] B1: Reproduktionstest schreiben (head-Ausdruck) → muss FAIL
- [ ] B2: Fix (head-Ausdrücke auflösen ODER Doku-Entscheidung)
- [ ] make test GREEN
- [ ] CHANGELOG.md + KB Update

## Entscheidungen (Lukas, 2026-08-03)
- Issue A: KORREKTE GENERIERUNG (nested {#if} im {#else}-Zweig unterstützen). NICHT nur Fehler werfen.
- Issue B: HEAD-AUSDRÜCKE AUFLÖSEN (wie Template-Body). Nicht nur dokumentieren.
- Umsetzung: Parallel Chains. A: parser_template/codegen_template. B: parser_section_head/codegen_ssr.

## Status (2026-08-03)
- [x] A1: Test (core Unit + _tests/core/Bugs/nested-if-in-else/) — FAIL vor Fix
- [x] A2: Fix codegen_template.go NodeIf (chain vs. echter else-Zweig, `return ""` entfernt)
- [x] B1: Test _tests/core/Bugs/head-expression-raw/ — FAIL vor Fix
- [x] B2: Fix neue Datei core/codegen_head.go (genHead) + 4 Emissionsstellen in core/codegen.go
- [x] make test: 144 Passed, 0 Failed (via smd, DiD-Ersatz)
- [ ] CHANGELOG.md + KB-Update
- [ ] Commit
- [ ] Report an User

## Follow-up-Befund (offen, neuer Task nötig)
coder: core/codegen.go genTemplateNodeComp (Zeile ~505-534) enthält DENSELBEN Bug (`return ""` bei Zeile ~519) für verschachtelte {#if} in KOMPONENTEN. Außerhalb Task-Scope gelassen. → Eigener Bug-Fix-Task, gleiches Muster wie Issue A.

## Status final (2026-08-03)
- [x] CHANGELOG.md + KB-Update (docs: v0.0.23, codegen-errors.2 Backlog)
- [x] QS durch git-Agent (9 Dateien, keine Blocker, <=120 Zeilen, Imports ok)
- [x] Commit f78c1db "fix(codegen): support nested {#if} in else branch and resolve head expressions" (10 Dateien, +269/-13)
- [x] Working Tree sauber (nur .agents/tasks, .agents/chains untracked)
