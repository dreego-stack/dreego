---
area: hardening
phase: v0.1-blocker
---
# Transpiler fuzzing

## Goal
Exercise the lexer and parser with malformed and unexpected input.

## Acceptance criteria
- Go fuzz targets cover lexer and parser entry points.
- Crashes and hangs produce permanent regression tests.
- Fuzz properties also check deterministic output, bounded resource use, and preservation of valid source rather than crashes alone.
