---
id: typed-forms.1
title: Typed Form Binding and Validation
status: 41
phase: v0.0.24
requires:
  - form-actions.1
  - transpiler.1
created: 2026-07-31
changed: 2026-08-08
---

Extend `BindForm` and `ValidateForm` to support `int`, `bool`, `time.Time`, slices, and improve the built-in `email` validator. The `email` rule is part of the core form tag validator set so that any form using `validate:"email"` gets a proper check.

Done in v0.0.24: `BindForm` binds int/bool/[]string, `core.RegisterRule(name, fn)` custom validators.
