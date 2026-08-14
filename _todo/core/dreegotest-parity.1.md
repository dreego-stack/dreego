---
area: testing
phase: pre-v0.1
---
# Dreegotest generation parity

## Goal
Prevent divergence between CLI `generate` output and the generation path used by `dreegotest`.

## Acceptance criteria
- A parity test compares both outputs for the same fixture.
- Any output divergence fails CI with a useful diff.
