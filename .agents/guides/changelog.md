
---
type: Guide
title: Skill: changelog
description: Conventions for maintaining CHANGELOG.md and version tagging workflow
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Skill: changelog

## Rule: Agents NEVER edit CHANGELOG.md directly

Agents NEVER edit CHANGELOG.md directly. New changelog lines belong only in the PR's unique `.changes/*.md` file. `release-prep` applies all pending files to CHANGELOG.md after merge. If a vX.Y.Z header is missing in CHANGELOG.md, that is normal until release-prep runs.

## Purpose

Keep CHANGELOG.md up to date with every meaningful change.

## When to use

- Every change lands via a pull request with one `.changes/*.md` file
- The change-file frontmatter declares `none | patch`; larger bumps are blocked before v0.1
- Concurrent pending files are combined into one release entry

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
2. Add a uniquely named `.changes/*.md` file with `version:` and changelog lines
3. CI (`pull-request-check.yml`) validates the file and runs the test suites
4. Squash-merge the PR
5. Serialized `main-push.yml` tests the latest main, combines pending files, updates the changelog, and creates the tag

No local tags, no `git tag -a`, no `git push --tags`.
