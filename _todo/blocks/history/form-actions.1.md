---
id: form-actions.1
title: Form Actions (g-action / g-submit)
status: 33
phase: v0.0.16
requires:
  - context-refactoring.1
  - routing.1
  - csrf.1
  - session.1
created: 2026-07-26
changed: 2026-07-29
---

`<form g-action="Login">` generates full handler pipeline: BindForm → ValidateForm → Handler → Redirect. `validate:"required,email,min,max"` struct tags with no external deps. `c.Errors(field)` and `c.Old(field)` for template re-render. `c.Redirect(url, code)` for PRG pattern. 112 tests total.

