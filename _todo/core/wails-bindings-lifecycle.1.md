# Wails bindings and lifecycle

- Area: target
- Phase: Wails
- Goal: generate typed, least-privilege Go-to-client bindings and own desktop startup and shutdown.
- Acceptance: one explicit Go contract generates client declarations, unsupported values fail before runtime, filesystem and shell access are absent unless registered, shutdown releases target resources, and development reload preserves deterministic output.
- Depends on: wails-assets-navigation.1
