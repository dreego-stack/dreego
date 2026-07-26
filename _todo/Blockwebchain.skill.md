---
id: Blockwebchain
title: Blockwebchain — Dependency-Graph TODO System
type: Skill
created: 2026-07-26
---

# Blockwebchain

Non-linear TODO system. Features/Ideas are Markdown blocks with YAML frontmatter. A Python script computes the dependency graph and shows what's next.

## Directory Structure

```
_todo/
├── blocks/                    # active blocks (planned, in-progress)
│   ├── form-actions.1.md
│   └── csrf-protection.2.md
├── blocks/history/            # completed blocks (done, rejected)
│   ├── transpiler.1.md
│   └── context-refactoring.3.md
├── process.py                 # dependency graph engine
└── index.md                   # generated: current state
```

## Block Frontmatter

```yaml
---
id: form-actions.1               # filename without .md
title: Form Actions              # human-readable
status: planned                  # draft | planned | in-progress | done | rejected
phase: v0.0.3                    # target version
requires:                        # block IDs that must be done first
  - csrf-protection.2
  - session-interface.4
created: 2026-07-23
changed: 2026-07-26
---
```

## Creating a New Block

1. Choose a unique `name.integer` filename (check no duplicates across blocks/ and history/)
2. Write frontmatter with all required fields
3. Every `requires` entry must reference an existing block ID
4. Run `python process.py` to verify the graph is valid
5. The script will show if the block is "available next" (all requires are done)

## Moving to History

When a block is done:
1. Change `status: done`
2. Move file from `blocks/` to `blocks/history/`
3. Run `python process.py` — it will detect newly unblocked dependents

## Validation Rules (built into process.py)

- No circular dependencies (A requires B, B requires A)
- Every `requires` resolves to an existing block ID (across both dirs)
- No self-requirement
- No duplicate filenames across blocks/ and history/

## Status System — Chain vs Web

- **Integer status** (1, 2, 3...) = CHAIN. Done, in history. Sequential, no gaps, no duplicates.
- **String status** (planned, in-progress, draft) = WEB. Still open. Dependency graph.
- Status 0 is invalid (error).
- process.py validates: chain starts at 1, no gaps, no duplicate integers.
- Current code = max chain integer + 1. Next completed block gets this code.

## process.py Output

```
BLOCKWEBCHAIN
────────────────────────────────────────
Chain: 1–13  |  Next code: 14

CHAIN (History):
  01  transpiler.1  — Transpiler-Pipeline
  ...

AVAILABLE NEXT:
  session.1  — Session-Interface
  static-assets.1  — Static Assets

BLOCKED:
  csrf.1  — CSRF-Schutz  (missing: session.1)
```
