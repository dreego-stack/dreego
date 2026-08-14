---
area: documentation
phase: pre-v0.1
---
# Consolidate normative architecture sources

## Goal
Make AGENTS, the decision index, and accepted ADR statuses point to one current
architecture without deleting useful history.

## Acceptance criteria
- Superseded Chi, validator, Tailwind-core, monorepo-plugin, target, and runtime decisions are marked historical.
- The decision index distinguishes current, provisional, research, and superseded material.
- AGENTS and current ADRs agree on the dependency-free Core and external-plugin boundary.
- A future agent can identify the current v0.1 contract without reading chronological logs.
