---
area: cli
phase: v0.1-blocker
---
# Make generate check content-based

## Goal
Make `dreego generate --check` verify generated content rather than file modification times.

## Acceptance criteria
- Check mode generates expected output without modifying the working tree.
- It compares every generated file affected by routes, components, layouts, static assets, and configuration.
- Missing, extra, or stale generated files fail with a concise path-level diff.
- Checkout, rebase, copied files, and manipulated timestamps cannot produce a false pass.
- The same generation implementation serves normal and check modes to prevent drift.
