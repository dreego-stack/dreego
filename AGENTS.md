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

v0.0.36 tagged — Single-source versioning: the latest git tag (`vX.Y.Z`) is the single truth; the CLI derives its version at build time (`-ldflags -X main.version=$(git describe --tags --abbrev=0)`) or from build info (`go install pkg@tag`). Releases are PR-driven: every change lands via a pull request with a `pr.md` (version bump + changelog lines), the tag is created by CI after merge. See `_todo/` for next steps.

## Product Focus

Dreego brings an intuitive, Svelte- and Astro-inspired development experience to Go without hiding or bending Go. Simplicity, type safety, low resource usage, and accessible workflows are core design constraints.

Accessibility is a release quality gate, not a cosmetic enhancement. CLI output, diagnostics, documentation, generated blueprints, and official components must work without relying on sight, color, or pointer input alone. Do not claim that Dreego can make arbitrary user applications automatically accessible.

- Until v1, SSR is the only production target and the core priority.
- Before v0.1, stabilize and harden the SSR core. Do not add an internal client runtime, SSG target, or expanded Wails target.
- The highest-priority pre-v0.1 architecture change is explicit `App` ownership of all runtime state, including explicit generated registration. Do not preserve the current global API through compatibility wrappers.
- HTMX, Alpine.js, and plain JavaScript are the supported progressive-enhancement path before v0.1.
- Between v0.1 and v1, client islands or hydration may be explored outside the stable core. Promote them only after real applications establish the required interfaces.
- Improved Wails support, SSG, and static deployment targets belong after v1.
- Preserve future extension points where inexpensive, but do not add speculative abstractions or delay SSR stability for post-v1 targets.

See `_docs/roadmap.md` for the public, non-binding roadmap.

## Core and Plugin Boundary

- Core contains the SSR framework capabilities needed by a normal Dreego application.
- Optional capabilities, provider integrations, and features with additional dependencies live in separate plugin repositories with their own `go.mod`, releases, tests, and CI.
- Keep optional implementations out of `core/`, even when they currently need only the standard library. SSE and WebSockets are plugins, not core packages.
- Add a provider-neutral interface to core only after at least two real implementations prove the same small contract is necessary.
- Remove the current speculative EventBus, Queue, KVStore, and Storage APIs before v0.1. The session Store remains part of the SSR core.
- Plugins may register route-specific behavior through the owning `App`; they must not weaken unrelated application defaults.

## File Structure

```
repo-root/
├── _todo/                  ← One open item per file; delete the file when its PR completes
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
version: patch
---

- Bug: fix X
- Feat: add Y
```

- `version: none` — no version bump, changelog lines only (e.g. dependabot updates)
- `version: patch` — `0.0.x` +1

NEVER use `version: minor` or `version: major` while in the v0.0.x phase —
a minor bump would tag v0.1.0 before v0.1 stabilization. Only `none` and
`patch` are allowed until the v0.1 release.

## Todo Workflow

Every open work item lives in its own Markdown file under `_todo/`.

- `_todo/core/` — concrete framework, CLI, test, documentation, and hardening work
- `_todo/future/` — long-term architecture ideas without a near-term commitment
- `_todo/plugins/` — external plugin ideas; these do not expand core scope by themselves
- One PR should normally implement one todo item.
- Name files after the stable item ID, for example `_todo/core/server-timeouts.1.md`.
- Concrete items record their area, phase, goal, acceptance criteria, and dependencies when applicable. Long-term idea files may remain concise until promoted into planned work.
- Delete the item file in the PR that completes it. The PR's `pr.md`, changelog, and Git history become the completion record.
- Do not keep completed, rejected, or superseded items in `_todo/`.
- Add newly discovered work as a new item file instead of appending to a shared checklist.

The CI (`pull_request.yml`) validates pr.md and runs `make test`. After approval, `release-prep` (manual, with PR number) applies the changelog + version to the PR branch and removes pr.md. After squash-merge, `release.yml` creates the tag. No local tags.

## Project: dreego

- **Name:** dreego | **Package:** `dreego` | **Module:** `github.com/dreego-stack/dreego`
- **Org:** github.com/dreego-stack | **Main repo:** github.com/dreego-stack/dreego
- **Approach:** Compile-Time Transpiler (.dreego → Go-Code) + net/http + HTMX/Alpine.js

## Note: smd

All commands run inside `smd` (Docker container). Never run `make test`, `go build`, or any dev command directly on the host. The smd container uses `golang:1.22-alpine`. Install curl once per session: `smd apk add --no-cache curl`.

## Git Operations

All git operations run on the HOST via the shell subagent — never inside the
smd container. The smd image has no git, and worktree `.git` files point to
host paths that do not exist in the container.

- Worktree setup: `git worktree add -b <branch> .worktrees/<name>` (shell agent)
- Commits: `git -C <repo>/.worktrees/<name> add/commit` (shell agent, host)
- Push and PR creation: `git -C <repo>/.worktrees/<name> push origin <branch>` and `gh pr create` (shell agent)
- Never run `git checkout`, `git reset`, or `git stash` on main.
- Main must stay clean; only the `.gitignore` entry for `.worktrees/` may appear uncommitted.

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

Every bug gets a permanent test in `_tests/go/bug_<name>_test.go`. Workflow:
1. Bug found → create `_tests/go/bug_<name>_test.go` that reproduces the bug (must FAIL)
2. Fix code until `make test` shows the new test GREEN
3. Bug is permanently covered — no regression risk

`.tmp/<name>/` is ONLY for temporary debugging/exploration — never for permanent tests.

## Feature Workflow

Every feature follows this cycle:

1. **`_tests/`** — Create integration test in `_tests/go/<name>_test.go` using `dreegotest` (see `_docs/testing.md` and existing `_tests/go/*_test.go` for the pattern)
2. **Code** — Implement in `core/` (one logical thing per file, max 300 lines)
3. **`_docs/`** — Update relevant documentation
4. **Test** — `go test ./_tests/go/ -run <TestName>` (or `make test`) — must be GREEN
5. **PR** — Create a PR with `pr.md` (version bump + changelog lines); CI validates it
6. **KB** — Update `.agents/log.md` + relevant concept/decision docs

For multi-step features, repeat the cycle for each step. Commit after each step.

## Type Safety

1. Build time (`dreego generate`, `go build`)
2. Start time (`server.Listen()`)
3. Runtime (per request, as local as possible)

Generated application APIs, component props, route contracts, and primary
application data are strongly typed. Unknown, duplicate, and missing generated
fields fail as early as possible.

Dynamic strings remain valid at boundaries whose schemas Dreego does not own,
including HTTP headers, URL and form values, sessions, configuration, and
request-local extension state. Missing values and conversions must be explicit
at those boundaries. Convert boundary data into typed application structures
before using it as domain data. Do not claim that Core contains no string keys.

---

## Architecture Guarantees

Until v1, SSR is the only production target and the core priority. SSG,
expanded Wails support, and static deployment targets belong after v1; the
former V2 preparation (Target interface, reserved CLI flags) is no longer
required — extension points are preserved only where inexpensive, without
speculative abstractions. See [decisions/ssg-wails-v2](.agents/decisions/ssg-wails-v2.md)
(superseded) and the Product Focus section above.

### 1. `<go>` Block: No hard `*http.Request`
Solution: `dreego.Context` Interface. → [decisions/context-design](.agents/decisions/context-design.md)

### 2. Plugin Contracts Stay Provisional Until v1
Real external plugins between v0.1 and v1 must validate the contract before a stability promise. → [plugin-contract.1](_todo/core/plugin-contract.1.md)

### 3. File-based Routing
Filename-based routing is the released pre-v0.1 implementation; the accepted
v0.1 target is one route file per URL (`+page.dreego` and method sections).
→ [decisions/routing-and-components](.agents/decisions/routing-and-components.md)

### 4. Asset System: Dual-Mode (Embedded + Disk)

### 5. Template Rendering without HTTP Server
