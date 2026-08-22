---
type: Decision
title: Per-File Line Limit Raised to 300
description: Hard file size limit changed from 120 to 300 lines to avoid artificial splits of tightly coupled logic
tags: [v0.0.22]
timestamp: 2026-08-03T22:30:00Z
---
# Per-File Line Limit Raised to 300

**Date:** 2026-08-03
**Status:** Accepted

## Context

The project had a hard rule of *max 120 lines per file*. Over time, two areas repeatedly hit the ceiling without a clean split point:

- Code generation logic (`core/codegen*.go`) — branches for parser, emitter, and error propagation share types and helper functions; splitting them creates package-internal APIs that are harder to follow than a single file.
- Session management (`core/session.go`) — initialization, encryption/decryption, and request-scoped helpers are tightly coupled around the same state type.

Artificial splits made files *smaller* but the system *harder* to read, because related state transitions were scattered across files whose only reason for existence was the 120-line rule.

## Decision

Raise the hard per-file line limit from **120 to 300 lines**.

The companion rule — *one logical thing per file* — stays in effect. A file should still contain a single concept (e.g., one type and its methods, one phase of the transpiler, one feature). The limit is a guardrail, not a goal.

## Consequences

- Generated and hand-written core files can stay coherent when the logic genuinely exceeds 120 lines.
- Agents must still enforce one logical thing per file during review.
- Existing files do not need bulk reformatting unless a future change naturally touches them.
- Documentation (`AGENTS.md`, `coding-standards.md`, task files) updated to reflect 300 lines.
