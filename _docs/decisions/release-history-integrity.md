---
title: Published release history integrity
description: Keep every published release tag in the ancestry of main
tags: [release, git, ci, pre-v0.1]
---

# Published release history integrity

**Status:** Accepted

## Context

The automated release workflow commits generated changelog updates directly to
`main` and then tags that commit. After v0.0.69, `main` was updated from a
history that did not contain the release commits for v0.0.70 and v0.0.71. The
published tags continued to identify the released source, but later work no
longer descended from them.

This split made the changelog and released features appear absent from `main`.
It also allowed a later release to be calculated from tags that were not part
of the branch being released.

## Decision

Every published release tag through the latest version must be an ancestor of
every pull request and of `main`. Pull-request and main-push CI enforce this
invariant before testing or releasing.

If the histories diverge, restore the published commits with a normal merge
through a reviewed pull request. Published tags must never be deleted, moved,
or recreated to make the check pass.

Repository settings must protect `main` against force pushes and require the
pull-request check before merging. GitHub owns those settings; the workflow
guard is the repository-controlled second line of defense.

## Consequences

- Released source and changelog entries remain reachable from `main`.
- A rewritten or incorrectly based branch fails before it can produce another
  release.
- Intentional history repair stays visible in the merge graph.
- Administrators must preserve branch protection when changing repository
  settings.
