# DreeJS component lifecycle

- Area: client runtime
- Phase: DreeJS
- Goal: add an optional component-scoped browser lifecycle without creating a SPA runtime.
- Acceptance: inactive components emit zero runtime bytes, activated instances have deterministic identity and cleanup, local events and presentation state remain scoped, CSP-compatible assets contain only used modules, and a countdown reference preserves keyboard and screen-reader behavior.
- Depends on: the Wails v3 Phase 1 reference application being released (`adapter/wails/_docs/phase-1.md`)
