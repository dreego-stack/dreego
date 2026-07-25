
---
type: Guide
title: Skill: changelog
description: Conventions for maintaining CHANGELOG.md and version tagging workflow
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---
# Skill: changelog

## Purpose

Keep CHANGELOG.md up to date with every meaningful change.

## When to use

- Vor JEDEM Commit: eine Zeile oder kurzer Absatz in CHANGELOG.md
- Neue Version (Tag): Agent schlagt vor (Feature=MINOR, Fix=PATCH, Doc=optional)
- User entscheidet ob neue Version getaggt wird

## Format

```markdown
## v0.0.2 (2026-07-25)

- Context refactoring: map[string]string → Interface + Embedding
- Recovery-Middleware: Panic → 500

## v0.0.1 (2026-07-25)

... (existing entries)
```

- Neueste Version OBEN
- Ein Punkt pro Feature/Fix
- Datum im ISO-Format

## Workflow

1. Code andern
2. CHANGELOG.md updaten (eine Zeile)
3. TODO.md abhaken was erledigt wurde
4. Commit
5. Wenn Feature fertig: `git tag -a vX.Y.Z -m "vX.Y.Z: summary"`
6. Push: `git push origin main --tags`
