---
area: api
phase: pre-v0.1
---
# Remove dead exported API surface

## Goal
Remove the dead exported surface `GenerateRouter`/`RouteInfo`
(`core/router.go`) and `NewHandler` (`core/codegen_ssr.go`) plus their test
references, unless documented as public API in `_docs/compatibility.md` (then
keep and document instead).

## Acceptance criteria
- Check `_docs/compatibility.md` first for whether these are public API.
- If not public: no production references remain.
- Tests are updated or removed accordingly.
- If public: keep and document instead of removing.
