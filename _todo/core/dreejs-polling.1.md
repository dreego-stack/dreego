# DreeJS polling

- Area: client runtime
- Phase: live updates
- Goal: add bounded polling as the first network update strategy.
- Acceptance: polling pauses or slows while hidden, requests cannot overlap without bounds, cancellation and cleanup stop all work, retry uses capped jittered backoff, stale responses cannot replace newer state, and status changes have accessible announcements.
- Depends on: dreejs-data-islands.1
