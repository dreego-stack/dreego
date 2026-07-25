
---
type: Guide
title: Skill: knowledge-base
description: Guide for maintaining the .agents/ Obsidian-style knowledge vault
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---
# Skill: knowledge-base

## Purpose

Maintain the `.agents/` Obsidian-style knowledge vault for the Dreego project.

## When to use

- After every architectural decision → write ADR in `decisions/`
- After researching a topic → write reference in `KB/`
- When a concept is fully formed → write in `concepts/`
- When project conventions change → update `guides/`
- Always update `_index.md` when adding new files

## File conventions

- Every note: `# Title` as first line
- Cross-reference with `[text](../path/file.md)`
- Facts only — no prose, no fluff
- Group related files in subdirectories

## Priority order when updating

1. `_index.md` — TOC, always first
2. `decisions/` — ADRs for architecture
3. `concepts/` — fleshed-out feature designs
4. `KB/` — external research, reference material
5. `guides/` — coding standards, architecture guide
