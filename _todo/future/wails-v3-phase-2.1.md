# Wails v3 Phase 2

- Area: target
- Phase: Wails v3 Phase 2
- Goal: harden and expand the desktop host after the upstream Wails v3 API is
  stable and Dreego has operational evidence from Phase 1.
- Entry gate: Wails v3 has a stable upstream release, the Phase 1 suite passes
  against it, the reference application has produced practical compatibility
  evidence, and an architecture review records required migrations.
- Candidate scope: broader typed bridge coverage, packaging and distribution,
  platform compatibility, performance, development tooling, and Wails-specific
  live updates. Recheck native bridge calls and programmatic shutdown under
  headless Linux: Wails v3 beta.20 can block both under Alpine GTK4/Xvfb even
  after the application-started event.
- Depends on: wails-reference-application.1 and an upstream stable Wails v3 release
