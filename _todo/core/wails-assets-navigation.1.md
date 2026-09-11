# Wails assets and navigation

- Area: target
- Phase: Wails
- Goal: provide deterministic embedded assets and explicit desktop route navigation.
- Acceptance: head entries, scoped styles, scripts, and static assets load without HTTP; literal routes navigate with documented history semantics; unsupported dynamic behavior reports a build diagnostic; and traversal attempts cannot leave the embedded asset root.
- Depends on: wails-render-host.1
