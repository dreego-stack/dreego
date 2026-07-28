
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

- Before EVERY commit: one line or short paragraph in CHANGELOG.md
- New version (Tag): Agent proposes (Feature=MINOR, Fix=PATCH, Doc=optional)
- User decides whether to tag a new version

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

## Workflow

1. Change code
2. Update CHANGELOG.md (one line)
3. Check off completed items in TODO.md
4. Commit
5. When feature is complete: `git tag -a vX.Y.Z -m "vX.Y.Z: summary"`
6. Push: `git push origin main --tags`
