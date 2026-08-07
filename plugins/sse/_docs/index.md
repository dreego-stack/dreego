# SSE Plugin

Server-Sent Events for dreego. The plugin registers a `GET /sse` route that
streams `text/event-stream` responses with a heartbeat, and exposes a
`Broadcast` method to push messages to all connected subscribers.

## Usage

```go
import (
    dreego "codeberg.org/dreego/dreego/core"
    "codeberg.org/dreego/dreego/plugins/sse"
)

func main() {
    p := sse.New()
    dreego.UsePlugin(p)
    dreego.Listen(":8080")
}
```

## Route

- `GET /sse` — streams events as `data: <message>` lines. A comment line
  (`: ping`) is sent every 15 seconds to keep the connection alive.

## Broadcast

```go
p.Broadcast("hello subscribers")
```

Every connected client receives `data: hello subscribers`.

### Backpressure

`Broadcast` sends non-blocking: each subscriber has a buffered channel
(capacity 16). When a subscriber's channel is full, the message is **dropped
for that subscriber** instead of blocking the broadcaster. This is a deliberate
backpressure decision — a slow client must never stall the broadcast loop.
Slow consumers may miss messages; reconnect or heartbeat-based resync is the
caller's responsibility.

## Streaming notes

- The SSE stream works behind the core `Compress()` middleware: the gzip
  response writer implements `http.Flusher`, so `Accept-Encoding: gzip`
  requests stream correctly.
- The request log line for a long-lived SSE connection is only written when
  the connection closes (the logging middleware wraps the handler), so
  `duration` reflects the full connection lifetime, not the time to first
  byte.

## Assets

The plugin exposes `assets/sse.js`, a small helper that connects an
`EventSource` to the `/sse` route:

```html
<script src="/assets/sse.js"></script>
<script>
  connectSSE("/sse", (data) => console.log(data))
</script>
```
