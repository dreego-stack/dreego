---
area: generator
phase: v0.1-blocker
depends_on: app-runtime.1
---
# Make source discovery and layout resolution deterministic

## Goal
Restrict generation to the declared Dreego project roots and resolve layouts by
the documented route hierarchy.

## Acceptance criteria
- Unrelated directories named `routes` or `components` outside the project root are ignored.
- Layout lookup is route-local and cascades through documented parent directories.
- Ambiguous components, layouts, error pages, and project roots fail with useful diagnostics.
- Add, move, and delete tests verify deterministic discovery on every supported platform.
