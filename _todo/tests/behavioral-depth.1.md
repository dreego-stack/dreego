---
area: tests
phase: after-v0.1
---
# Upgrade compile-only integration tests to behavioral assertions

## Goal
Upgrade compile-only integration tests (~60 files) to behavioral assertions
where feasible, starting with `middleware_*`, `routing_*`, `template_*`.

## Acceptance criteria
- A sample batch is converted to behavioral assertions.
- A list of remaining candidates is recorded.
