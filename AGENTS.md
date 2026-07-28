# Agent Instructions for Dreego

## Language Rule

- **Chat with user (Lukas): German**
- **Everything in this repository: English**
  - All `.go` files, comments, variable names
  - All `.md` documentation (docs, agents, todos)
  - All commit messages
  - All test files and scripts
  - All configuration files

## Current Phase: pre v0.1

Blueprint scaffolding + integration tests working. v0.0.10 tagged. See TODO.md for next steps.

## File Structure

```
repo-root/
├── TODO.md                 ← NEXT code changes (short, prioritized)
├── ROADMAP.md              ← Release pipeline (high-level)
├── CHANGELOG.md            ← What came in which version
├── README.md               ← Project overview
├── LICENSE                 ← MPL-2.0
├── _docs/                  ← Public documentation
├── _tests/                 ← Integration tests (Docker, `make test`)
├── .tmp/                   ← Temporary debug spaces (no permanent tests)
│
.agents/                    ← Knowledge Base (OKF format)
├── index.md                 ← Start here (OKF TOC)
├── log.md                   ← Change history
├── tips.md                  ← 50 tips + checklist
├── KB/                      ← Reference material
├── decisions/               ← Architecture decisions (ADR)
├── concepts/                ← Worked-out concepts
└── guides/                  ← Coding standards, skills, OKF conventions
```

## Skills

- [Knowledge Base](.agents/guides/knowledge-base.md) — How the knowledge base is maintained
- [Changelog](.agents/guides/changelog.md) — How CHANGELOG.md and versioning works
- [Open Knowledge Format](.agents/guides/open-knowledge-format.md) — OKF conventions (YAML frontmatter, types, links)

## Commit Convention

See [Changelog Guide](.agents/guides/changelog.md) for the full workflow.

## Project: dreego

- **Name:** dreego | **Package:** `dreego` | **Module:** `codeberg.org/dreego/dreego`
- **Org:** codeberg.org/dreego | **Mirror:** github.com/LukasLow/dreego
- **Approach:** Compile-Time Transpiler (.dreego → Go-Code) + net/http + HTMX/Alpine.js

## Coding Rules

- Max 120 lines per file, one logical thing per file
- No code comments (except where needed for clarity)
- Go 1.22+, prefer standard library
- Core code in `dreego-core/` (single package), plugins in `dreego-plugin/`
- Build via `dreego` CLI, not directly `go build`
- Generated `dree.go` not committed

## Bug → Test → Fix Workflow

Every bug gets a permanent test in `_tests/Bugs/<name>/`. Workflow:
1. Bug found → create `_tests/Bugs/<name>/` that reproduces the bug (must FAIL)
2. Fix code until `make test` shows the new test GREEN
3. Bug is permanently covered — no regression risk

`.tmp/<name>/` is ONLY for temporary debugging/exploration — never for permanent tests.

## Type Safety

1. Build time (`dreego generate`, `go build`)
2. Start time (`server.Listen()`)
3. Runtime (per request, as local as possible)

No `map[string]string`, no `interface{}` cast, no string key in core.

---

## Architecture Guarantees for V2

### 1. Target-Agnostic Transpiler Pipeline
V1: `TargetSSR`, V2: `TargetSSG`, `TargetWails`. → [decisions/ssg-wails-v2](.agents/decisions/ssg-wails-v2.md)

### 2. `<go>` Block: No hard `*http.Request`
Solution: `dreego.Context` Interface. → [decisions/context-design](.agents/decisions/context-design.md)

### 3. Transpiler Pipeline with Extension Points
→ [decisions/typescript-v2](.agents/decisions/typescript-v2.md)

### 4. Plugin Interface: First Release = Final Contract
→ [concepts/addon-ecosystem](.agents/concepts/addon-ecosystem.md)

### 5. File-based Routing: Crawlable for SSG
→ [decisions/routing-and-components](.agents/decisions/routing-and-components.md)

### 6. Asset System: Dual-Mode (Embedded + Disk)

### 7. Template Rendering without HTTP Server

### 8. CLI Interface: Reserved Flags
`dreego build --static | --wails | --mobile`
