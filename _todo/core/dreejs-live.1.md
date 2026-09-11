# DreeJS bidirectional live transport

- Area: client runtime
- Phase: live updates
- Goal: support multiplexed bidirectional component updates through the external WebSocket plugin.
- Acceptance: messages carry protocol version, component identity, ordering, and resynchronization data; stale messages cannot overwrite newer state; authorization, origin, CSRF, size, rate, and backpressure limits are tested; focus and form state survive updates; and multiple components safely share one connection.
- Depends on: dreejs-stream.1 and a released compatible WebSocket plugin
