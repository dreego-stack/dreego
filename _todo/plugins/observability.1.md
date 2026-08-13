---
area: plugins
phase: after-v0.1
---
# Metrics and tracing plugins

## Goal
Provide Prometheus metrics and OpenTelemetry spans as separate plugins built on the public plugin interface.

## Acceptance criteria
- Core remains free of external dependencies.
- Metrics and tracing work through separate Go modules.
- Setup and verification are documented.
