# WebSocket hub registration, broadcast, and README alignment

- Area: plugin (plugin-websocket)
- Phase: plugin
- Goal: Make the first-party WebSocket plugin fit the multi-client use case it
  is chosen for, and align its README with the released API.
- Gaps (reported by an external test user on v0.10.3):
  - The connection handler is echo-only: it writes each payload back to the
    sender and never fans out, so a chat or any multi-client feature does not
    work.
  - The `Hub` exposes `Broadcast` and a client count, but connections are never
    registered with the hub, so `Broadcast` has no recipients.
  - The README "Quick Start" uses APIs that do not exist in v0.10.3
    (`app.Listen(":8080")` and a `.../core` import without `adapter/ssr`).
- Proposal: register each accepted connection with the hub; implement a real
  `Broadcast` over the registered set; document the connection lifecycle and
  shutdown; rewrite the README against the released API, ideally as a compiled
  example covered by a test.
- Note: This lives in the plugin repository, not in `core/`.
- Acceptance: two clients exchanging messages through the plugin receive each
  other's messages; `Broadcast` reaches every registered client; the README
  example compiles and matches the released API.
- Dependencies: plugin-websocket contract; related to
  `_todo/plugins/websocket-presence-hook.1.md`.
