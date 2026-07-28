---
id: form-actions.1
title: Form Actions (g-action / g-submit)
status: planned
phase: v0.0.16
requires:
  - context-refactoring.1
  - routing.1
  - csrf.1
  - session.1
created: 2026-07-26
changed: 2026-07-29
---

Form handler in `<go method="post">` block with `g-action="Login"` on `<form>`. Generated handler: CSRF check, r.ParseForm(), struct mapping via form tags, validation via go-playground/validator, then calls handler function `func Login(c dreego.Context, form LoginForm) error`. c.Errors("email"), c.Old("email") for template access. Progressive Enhancement: works without JS, HTMX upgrades it. ADR decision form-actions.md exists. Blocks without g-action remain plain forms.

