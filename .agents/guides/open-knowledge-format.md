
---
type: Guide
title: Skill: Open Knowledge Format (OKF)
description: Conventions for maintaining agent-readable knowledge in the .agents/ directory
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---

## Purpose

Maintain the `.agents/` directory as an OKF (Open Knowledge Format) knowledge bundle. OKF is a Google Cloud standard (v0.1, June 2026) for agent-readable knowledge: markdown files with YAML frontmatter, typed concepts, and standard markdown links forming a knowledge graph.

## When to use

- When creating a new `.md` file in `.agents/`
- When updating an existing file
- After architectural decisions
- When the knowledge graph needs new edges

## File format

Every file starts with YAML frontmatter:

```yaml
---
type: Decision
title: Decision: Name "dreego"
description: Naming: dreego / .dreego
tags: [naming]
timestamp: 2026-07-23T00:00:00Z
---
```

Required: `type`. Recommended: `title`, `description`, `tags`, `timestamp`.

## Type values

| Type | Directory | Example |
|------|-----------|---------|
| `Decision` | `decisions/` | ADR-format architecture decisions |
| `Concept` | `concepts/` | Feature designs, architecture overviews |
| `Reference` | `KB/` | External research, framework comparisons |
| `Guide` | `guides/` | Coding standards, workflows |
| `Plan` | root | TODO.md / TODO-Future.md |
| `Index` | root | index.md (TOC) |
| `Log` | root | log.md (changelog) |

## Links

Always use standard markdown: `[text](../path/to/file.md)` — NEVER `[[wiki-links]]`.

From `decisions/name-dreego.md` to `index.md`: `[index](../index.md)`.
From `KB/dreego-concept.md` to `decisions/` file: `[Decision](../decisions/name-dreego.md)`.

## log.md

Record every change in `log.md`:

```markdown
## 2026-07-25

- Converted entire knowledge base to OKF format
- Added YAML frontmatter with `type` field to all files
- Replaced `[[wiki-links]]` with standard markdown links
```

## index.md

Lists all concepts grouped by directory:

```markdown
## Decisions
- [Name "dreego"](decisions/name-dreego.md) — Naming
- [Technology Stack](decisions/technology-stack.md) — Tech Stack

## Concepts
- [Dreego Architecture](concepts/dreego-architecture.md) — Architecture Overview
```
