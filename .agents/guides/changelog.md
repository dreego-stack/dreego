
---
type: Guide
title: Skill: changelog
description: Conventions for maintaining CHANGELOG.md and version tagging workflow
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Skill: changelog

## Purpose

Keep CHANGELOG.md up to date with every meaningful change.

## When to use

- Every change lands via a pull request with a `pr.md` (see AGENTS.md Commit Convention)
- The `pr.md` frontmatter declares the version bump: `none | patch` — `minor` and `major` are blocked until the v0.1 release
- The changelog lines in `pr.md` become the CHANGELOG entry

## Format

```markdown
## v0.0.2 (2026-07-25)

- Context refactoring: map[string]string → Interface + Embedding
- Recovery-Middleware: Panic → 500

## v0.0.1 (2026-07-25)

... (existing entries)
```

- Newest version on TOP
- One bullet per feature/fix
- Date in ISO format
- `version: none` prepends only the lines, no version header

## Workflow

1. Change code
2. Create a PR with `pr.md` (copy `pr.md.example`): `version:` + changelog lines
3. CI (`pull_request.yml`) validates pr.md and runs `make test`
4. After approval, run `release-prep` (manual, with PR number) — it applies the changelog + version to the PR branch and removes pr.md
5. Squash-merge the PR
6. `release.yml` creates the tag after merge

No local tags, no `git tag -a`, no `git push --tags`.
