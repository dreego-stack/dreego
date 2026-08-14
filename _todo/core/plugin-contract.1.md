---
area: api
phase: between-v0.1-and-v1
depends_on: app-runtime.1
---
# Validate the plugin contract

## Goal
Validate whether explicit App-bound registration functions need any shared
capability interfaces after multiple real external plugins exist.

## Acceptance criteria
- Auth, UI, and at least one infrastructure plugin exercise different capabilities.
- Each plugin begins with `Register(app, typedOptions) error` and no required
  central Plugin interface.
- Assets use explicit App APIs only when a real plugin requires them.
- If background work requires startup or shutdown hooks, App integrates them
  safely, including cleanup after partial startup failure.
- The review decides whether any proven common behavior justifies a small
  optional capability interface.
- Registration order, duplicate plugins, late registration, errors, and concurrency have documented semantics and tests.
- The final compatibility promise begins at v1, not before.
