# `Must` helpers for App registration

- Area: core API
- Phase: v0.x
- Goal: Remove the repeated `if err := app.Register(...); err != nil { log.Fatal(err) }`
  blocks in generated `Register(app)` functions and handwritten endpoint files.
- Proposal: add a small `Must`-style helper (for example
  `app.MustRegister(method, pattern, handler)`) that panics on a registration
  error, so generated registration code stays one line per route while the
  checked `Register` stays available.
- Acceptance: generated registration uses the helper; `Register` keeps returning
  an error; a conflicting route still fails loudly at build/start time.
- Dependencies: none.
