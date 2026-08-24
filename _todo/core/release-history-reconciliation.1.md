---
area: release
phase: pre-v0.1
---
# Reconcile main with published release tags

## Goal

Restore one coherent, non-rewritten release history before v0.1. At the time
this item was recorded, `origin/main` at `d2882ca` descends from `v0.0.69`, but
published `v0.0.70` and `v0.0.71` are not ancestors of main.

## Acceptance criteria

- Audit the code and documentation changes reachable from `v0.0.70` and
  `v0.0.71` but absent from `origin/main`.
- Restore every intended released change through a normal pull request without
  moving or recreating published tags.
- Resolve conflicts against the current roadmap deliberately; do not select a
  side by commit date alone.
- Verify that every published release tag through the latest `v0.0.x` is an
  ancestor of the resulting main branch.
- Verify changelog entries, pending `.changes` files, generated version output,
  and GitHub release metadata against the repaired graph.
- Run the full test, race, release-prep, and dependency suites before merge.
- Document the cause and add branch-protection or workflow safeguards that make
  an accidental main rewrite harder.

## Out of scope

- Deleting, moving, or force-recreating published tags.
- Tagging v0.1 before the graph and released behavior agree.
