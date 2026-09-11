# DreeJS private data islands

- Area: client runtime
- Phase: DreeJS
- Goal: load personalized fragments without placing private data in shared SSR cache entries.
- Acceptance: loading, empty, error, cancellation, and fallback states are accessible; personalized responses default to `private, no-store`; shared page output contains no private data; authorization is revalidated per request; and a black-box test proves isolation between users.
- Depends on: dreejs-typed-props.1 and cache-interface.1
