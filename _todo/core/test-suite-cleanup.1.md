---
area: testing
phase: pre-v0.1
---
# Remove the obsolete shell test-header contract

## Goal
Delete the shell-suite header check instead of applying a historical shell
convention to the current Go integration tests.

## Acceptance criteria
- The obsolete standard-header test and shell-header documentation are removed.
- Current Go integration tests remain discoverable through normal Go conventions.
- CI fails when the intended Go test packages are missing or skipped unexpectedly.
