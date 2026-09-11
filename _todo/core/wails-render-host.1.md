# Wails render host

- Area: target
- Phase: Wails
- Goal: render a Dreego document in Wails without starting an HTTP listener.
- Acceptance: one component renders through the target-neutral App, the same component remains usable under SSR, unsupported HTTP-only capabilities fail explicitly, and an integration test proves that no TCP socket is opened.
- Depends on: render foundation and Lua client hardening
