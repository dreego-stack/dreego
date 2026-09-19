# WebSocket presence and broadcast hook

- Area: plugin (plugin-websocket)
- Phase: plugin
- Goal: Remove the manual `Broadcast` call after every state change and the
  client-side "Online" approximation.
- Proposal: expose an exported client count and a small presence/broadcast hook
  that an application can register once, so state changes fan out without each
  handler remembering to broadcast.
- Note: This lives in the plugin repository, not in `core/`.
- Acceptance: presence count is server-authoritative; an app registers a
  broadcast hook instead of calling `Broadcast` after each write.
- Dependencies: plugin-websocket contract.
