# DreeJS stream transport

- Area: client runtime
- Phase: live updates
- Goal: integrate server-sent updates through the external SSE plugin without moving connection ownership into core.
- Acceptance: component subscriptions clean up deterministically, reconnect and resynchronization are bounded and tested, authentication and message limits remain explicit, multiple components cannot mix updates, and only the stream runtime module is emitted.
- Depends on: dreejs-polling.1 and a released compatible SSE plugin
