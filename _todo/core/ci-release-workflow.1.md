---
area: automation
phase: pre-v0.1
---
# Align and serialize the PR release workflow

## Goal
Make the implemented workflows match the PR-driven release process documented in AGENTS.

## Acceptance criteria
- Workflow names and documented release steps agree.
- Release preparation happens on the intended PR commit rather than mutating `main` unexpectedly.
- Concurrent merges cannot race version, changelog, commit, or tag creation.
- A tag is created only for a completely checked commit.
- Workflow contract tests or dry-run validation cover patch, no-version, and failure paths.
