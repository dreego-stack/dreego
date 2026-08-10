# Agent Instructions for Dreego

- Don't create binaries here — only in /tmp or ./tmp

## Language Rule

- **Chat with user : German**
- **Everything in this repository: English**
  - All `.go` files, comments, variable names
  - All `.md` documentation (docs, agents, todos)
  - All commit messages
  - All test files and scripts
  - All configuration files

## Model Strategy (Pro + Flash)

Two models share the work with strict role separation:

| Role | Model | Strengths | Cost |
|------|-------|-----------|------|
| **Architect** | `deepseek-v4-pro` | Decisions, architecture, spec design, code review, quality gates | $$$ |
| **Worker** | `deepseek-v4-flash` | Executes defined tasks: read files, write code, run tests, return diffs | $ |

**Rules:**

1. **Pro (Architect):**
   - Human-facing: all chat messages go through Pro
   - Decides *what* to build and *how* — design, architecture, specs, task breakdown
   - Reviews Flash output before it reaches the user
   - Never delegates ambiguous, creative, or decision-heavy tasks to Flash
   - Verifies Flash-generated code against specs, coding rules, and test results

2. **Flash (Worker):**
   - Receives fully-specified, self-contained tasks from Pro
   - Reads files, searches codebase, returns structured summaries
   - Writes code and test files according to exact specifications
   - Returns diffs and test results — no interpretation, no decisions
   - No direct user interaction; all output passes through Pro

3. **Token Efficiency:**
   - Batch independent Flash tasks in parallel (e.g. read 3 files at once)
   - Flash reads large files and returns summaries — Pro sees summary, not full content
   - Flash writes boilerplate/tests/diffs — Pro reviews the result once
   - Pro makes a decision → Flash executes → Pro verifies → User sees final result

4. **Quality Gate:**
   - After Flash writes code, Pro must verify: compilation (`go build`), test pass (`make test`), line count (max 300), coding rules, no comments unless needed
   - If Flash output violates any rule, Pro fixes or re-tasks Flash with corrective instructions

## Current Phase: pre v0.1

v0.0.27 tagged — Single-source versioning: the latest git tag (`vX.Y.Z`) is the single truth; the CLI derives its version at build time (`-ldflags -X main.version=$(git describe --tags --abbrev=0)`) or from build info (`go install pkg@tag`). Releases are PR-driven: every change lands via a pull request with a `pr.md` (version bump + changelog lines), the tag is created by CI after merge. See TODO.md for next steps.

## File Structure

```
repo-root/
├── TODO.md                 ← NEXT code changes (short, prioritized)
├── TODO-Future.md          ← Long-term ideas (SSG, Wails, plugin ideas)
├── CHANGELOG.md            ← What came in which version
├── README.md               ← Project overview
├── LICENSE                 ← MPL-2.0
├── go.mod                  ← Single root module (one tag per release)
├── pr.md.example           ← PR metadata template (version + changelog lines)
├── _docs/                  ← Public documentation
├── _tests/                 ← Integration tests (Docker, `make test`)
│   └── core/<Category>/    ← Core/framework test suites
├── .tmp/                   ← Temporary debug spaces (no permanent tests)
│
├── core/                   ← Core package (single package, no external deps)
├── cli/dreego/             ← CLI binary
├── .github/workflows/      ← CI: pull_request.yml, release-prep.yml, release.yml
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

Every change lands via a pull request. The PR must contain a `pr.md` (copy `pr.md.example`) with YAML frontmatter:

```yaml
---
version: patch        # none | patch | minor | major
---

- Bug: fix X
- Feat: add Y
```

- `version: none` — no version bump, changelog lines only (e.g. dependabot updates)
- `version: patch` — `0.0.x` +1
- `version: minor` — `0.x.0` +1
- `version: major` — `x.0.0` +1

The CI (`pull_request.yml`) validates pr.md and runs `make test`. After approval, `release-prep` (manual, with PR number) applies the changelog + version to the PR branch and removes pr.md. After squash-merge, `release.yml` creates the tag. No local tags.

## Project: dreego

- **Name:** dreego | **Package:** `dreego` | **Module:** `github.com/dreego-stack/dreego`
- **Org:** github.com/dreego-stack | **Main repo:** github.com/dreego-stack/dreego
- **Approach:** Compile-Time Transpiler (.dreego → Go-Code) + net/http + HTMX/Alpine.js

## Note: smd

All commands run inside `smd` (Docker container). Never run `make test`, `go build`, or any dev command directly on the host. The smd container uses `golang:1.22-alpine`. Install curl once per session: `smd apk add --no-cache curl`.

## Coding Rules

- Max 300 lines per file, one logical thing per file
- No code comments (except where needed for clarity)
- Go 1.22+, prefer standard library
- Single root module `github.com/dreego-stack/dreego` (one `go.mod` at repo root, one tag per release)
- Core code in `core/` (no external deps — enforced by `_scripts/check-core-deps.sh` in CI)
- CLI in `cli/dreego/` (imports core)
- Plugins live in separate repos under `github.com/dreego-stack/` (each with own `go.mod`)
- Build via `dreego` CLI, not directly `go build`
- Generated `dree.go` not committed

## Bug → Test → Fix Workflow

Every bug gets a permanent test in `_tests/core/Bugs/<name>/`. Workflow:
1. Bug found → create `_tests/core/Bugs/<name>/` that reproduces the bug (must FAIL)
2. Fix code until `make test` shows the new test GREEN
3. Bug is permanently covered — no regression risk

`.tmp/<name>/` is ONLY for temporary debugging/exploration — never for permanent tests.

## Feature Workflow

Every feature follows this cycle:

1. **`_tests/`** — Create integration test in `_tests/core/<FeatureGroup>/<name>/test.sh`
2. **Code** — Implement in `core/` (one logical thing per file, max 300 lines)
3. **`_docs/`** — Update relevant documentation
4. **Test** — `DREEGO_FILTER=<name> make test` — must be GREEN
5. **PR** — Create a PR with `pr.md` (version bump + changelog lines); CI validates it
6. **KB** — Update `.agents/log.md` + relevant concept/decision docs

For multi-step features, repeat the cycle for each step. Commit after each step.

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
→ [concepts/plugin-ecosystem](.agents/concepts/plugin-ecosystem.md)

### 5. File-based Routing: Crawlable for SSG
→ [decisions/routing-and-components](.agents/decisions/routing-and-components.md)

### 6. Asset System: Dual-Mode (Embedded + Disk)

### 7. Template Rendering without HTTP Server

### 8. CLI Interface: Reserved Flags
`dreego build --static | --wails | --mobile`
