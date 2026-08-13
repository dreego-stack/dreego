---
area: performance
phase: v0.1-blocker
---
# Core benchmarks

## Goal
Measure code generation and runtime performance before making overhead claims.

## Acceptance criteria
- Reproducible Go benchmarks cover code generation and representative requests.
- Documentation reports measured conditions instead of unqualified claims.
- README removes or qualifies claims such as "no hidden overhead" and "zero overhead" until measurements support them.
- Benchmark results are tracked for regressions without enforcing unstable machine-specific absolute timings.
