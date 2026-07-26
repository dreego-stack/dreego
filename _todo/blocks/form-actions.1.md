---
id: form-actions.1
title: Form Actions (g-action / g-submit)
status: planned
phase: v0.0.3
requires:
  - context-refactoring.1
  - routing.1
  - csrf.1
created: 2026-07-26
changed: 2026-07-26
---

Form-Handler im <go method="post">-Block. g-action="login" markiert Form, Handler-Funktion im <go> verarbeitet POST. Auto-CSRF via csrf.1. Typensichere Form-Daten via Go-Struct-Tags. Progressive Enhancement: funktioniert ohne JS, HTMX upgraded es. ADR-Entscheidung form-actions.md liegt vor.
