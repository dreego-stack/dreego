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

Form handler in <go method="post"> block. g-action="login" marks form, handler function in <go> processes POST. Auto-CSRF via csrf.1. Type-safe form data via Go struct tags. Progressive Enhancement: works without JS, HTMX upgrades it. ADR decision form-actions.md exists.
