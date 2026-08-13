---
area: api
phase: v0.1-blocker
depends_on:
  - app-runtime.1
  - component-correctness.1
  - context-render.1
  - error-propagation.1
  - reference-apps.1
  - remove-optional-infrastructure.1
  - routing-correctness.1
  - server-timeouts.1
  - session-security.2
---
# Pre-v0.1 API review

## Goal
Review the public API and define the compatibility policy before v0.1.

## Acceptance criteria
- Exported contracts are reviewed against real applications and plugins.
- The breaking-change policy is documented.
- Prematurely frozen interfaces are revised before the compatibility promise.
- The plugin contract remains explicitly provisional until v1 and is excluded from any v0.1 stability promise.
