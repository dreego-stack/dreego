---
area: api
phase: between-v0.1-and-v1
depends_on: app-runtime.1
---
# Validate the plugin contract

## Goal
Replace the premature frozen-v1 promise with a contract proven by multiple real external plugins before v1.

## Acceptance criteria
- Auth, UI, and at least one infrastructure plugin exercise different capabilities.
- Assets are either served through a defined lifecycle or removed from the contract.
- Server startup and shutdown integrate plugin lifecycle safely, including cleanup after partial startup failure.
- The review decides whether routes, middleware, assets, and lifecycle become small optional capability interfaces.
- Registration order, duplicate plugins, late registration, errors, and concurrency have documented semantics and tests.
- The final compatibility promise begins at v1, not before.
