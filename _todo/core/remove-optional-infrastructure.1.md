---
area: architecture
phase: v0.1-blocker
priority: 2
---
# Remove speculative infrastructure APIs from core

## Goal
Remove EventBus, Queue, KVStore, and Storage contracts and implementations from core before v0.1. They are optional infrastructure and their contracts have not been proven by multiple real implementations.

The session Store remains in core because sessions are part of the normal SSR web foundation.

## Scope
- Remove `EventBus`, `Subscription`, and the in-memory event bus implementation.
- Remove `Queue`, `Job`, `JobHandler`, and `JobMiddleware`.
- Remove `KVStore`.
- Remove `Storage`.
- Remove their core tests and public core documentation.
- Preserve the product ideas as separate future plugin items whose APIs initially live in their own repositories.

## Acceptance criteria
- A normal Dreego SSR application loses no routing, rendering, form, session, middleware, or static-asset capability.
- Core exposes no optional event bus, queue, key-value, or blob-storage API.
- No compatibility shims or deprecated aliases remain.
- Future plugins own their APIs until at least two implementations prove a shared provider-neutral contract.
- Core dependency and full test checks pass after removal.
