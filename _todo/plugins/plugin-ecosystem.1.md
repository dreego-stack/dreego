---
area: plugins
phase: after-v0.1
---
# Plugin ecosystem foundations

## Goal
Validate discovery, documentation, compatibility, and distribution through real auth, UI, admin, and database plugins.

## Acceptance criteria
- Plugins remain separate Go modules.
- Core changes require evidence from multiple implementations.
- Every optional plugin uses a separate repository and `go.mod`, including dependency-free plugins.
