---
area: api
phase: v0.1-blocker
---
# Remove the Context.Render stub

## Goal
Remove the unused runtime-style `Context.Render(name string, data any)` API before stabilization. Dynamic URL segments such as `/user/{username}` are unrelated and remain supported.

## Acceptance criteria
- `Render` and `ErrRender` are removed when no remaining caller requires them.
- Generated, type-checked rendering and dynamic route parameters continue to work.
- Tests and documentation contain no reference to the removed API.
- Build-time generator extensions and optional runtime scripting remain separate future concerns.
