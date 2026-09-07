# Phase: live updates

## Goal

Extend the proven DreeJS component lifecycle with explicit network update
strategies: polling, server-sent events, and bidirectional live
connections. Developers choose the cheapest strategy that satisfies the user
experience.

## Strategy order

Implement in increasing operational complexity:

1. `poll`: repeated fetch with visibility pause, cancellation, jitter, and
   backoff.
2. `stream`: unidirectional server updates through an SSE plugin.
3. `live`: bidirectional events and server-rendered component updates through a
   multiplexed WebSocket connection.

Each strategy is independently usable. Installing or using polling must not
ship WebSocket code.

## Server-rendered updates

The server owns business state, handles an event, renders the affected
component region, and returns a structured update. The browser applies the
update through the DreeJS DOM contract.

Messages need:

- protocol version;
- page and component instance identity;
- event or update kind;
- monotonic ordering or revision information;
- HTML or a deliberately chosen patch representation;
- head or asset deltas where supported;
- recoverable error information;
- reconnect and resynchronization behavior.

Do not send arbitrary executable code in server updates. All user data remains
subject to the same contextual escaping as initial rendering.

## Connection ownership

Directives are component-scoped, but physical connections are shared where
possible. One page-level WebSocket can multiplex several live components.
Connection ownership defines:

- authentication and CSRF binding;
- subscription registration and cleanup;
- reconnect with bounded exponential backoff and jitter;
- resume or full resynchronization after missed revisions;
- backpressure and maximum message sizes;
- idle and visibility behavior;
- graceful server shutdown;
- accessible disconnected and retry states.

SSE and WebSocket transports are external plugins. DreeJS owns the client
protocol and compiler integration; core gains a provider-neutral interface only
after at least two real transports prove the same small contract.

## Stateless and stateful server modes

The first implementation should prefer stateless component requests: each
event authenticates, loads authoritative state, renders, and responds. This
scales across ordinary Go instances and survives reconnects naturally.

A later stateful live mode may keep per-connection view state in Go for lower
latency and richer interactions. It requires explicit memory limits, lifecycle,
recovery, deployment guidance, and horizontal scaling behavior. It must never
be enabled merely because a component uses a WebSocket.

## Host behavior

SSR supplies HTTP endpoints for polling and transport plugins. Wails uses an
explicit desktop bridge where supported. Generation fails when a directive has
no compatible endpoint or host capability.

## Distributed operation

Dreego does not synchronize arbitrary in-memory Go state. Horizontal instances
use explicit external services:

- shared session store;
- pub/sub for broadcasts and presence;
- cache or database for authoritative data;
- optional distributed locks or job queues.

Those capabilities belong to external plugins. Core contracts are introduced
only after at least two implementations demonstrate the same need.

## Acceptance criteria

- Each strategy ships only its required runtime module.
- Polling pauses or slows when hidden and cannot overlap unbounded requests.
- Stream and live reconnect behavior is deterministic and tested.
- Multiple live components share a connection without mixing messages.
- Stale or reordered messages cannot overwrite newer component state.
- Authentication, authorization, CSRF, origin, size, and rate limits have
  explicit tests.
- DOM updates preserve focus, form state, and accessibility announcements.
- Multi-instance reference tests document which external shared services are
  required.

## Not in this phase

- Transparent distributed memory.
- A full SPA router.
- Offline-first conflict resolution.
- Collaborative CRDT state.
- Treating every component as permanently connected by default.
