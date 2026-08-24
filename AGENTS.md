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

The latest `v0.0.x` git tag is the single version source; the CLI derives its version at build time (`-ldflags -X main.version=$(git describe --tags --abbrev=0)`) or from build info (`go install pkg@tag`). Releases are PR-driven: every change lands via a pull request with one unique `.changes/*.md` file, and CI combines pending files into the changelog and creates the tag after merge. `version: none` files are never applied on their own — they stay pending until a `version: patch` file triggers the release. See `_todo/` for next steps.

## Product Focus

Dreego brings an intuitive, Svelte- and Astro-inspired development experience to Go without hiding or bending Go. Simplicity, type safety, low resource usage, and accessible workflows are core design constraints.

Accessibility is a release quality gate, not a cosmetic enhancement. CLI output, diagnostics, documentation, generated blueprints, and official components must work without relying on sight, color, or pointer input alone. Do not claim that Dreego can make arbitrary user applications automatically accessible.

- SSR is the current production target and the v0.1 foundation. Target-neutral
  rendering, SSG, Wails, and DreeJS are planned sequentially for the long v0.x
  line; they are not current behavior until implemented and documented.
- Before v0.1, stabilize and harden SSR and complete the atomic semantic-section
  migration. Do not start SSG, Wails, or DreeJS implementation before the render
  foundation phase.
- Preserve explicit `App` ownership of all runtime state and explicit generated
  registration. Do not reintroduce global compatibility APIs.
- HTMX, Alpine.js, and plain JavaScript are the supported progressive-enhancement path before v0.1.
- After v0.1, extract a target-neutral typed App and render foundation before
  adding the first-party SSG and Wails target packages.
- DreeJS is an optional modular browser layer, not a target or SPA. Local
  presentation state may run in the browser; authoritative business state stays
  in Go or another explicit backend.
- SPA and Wasm remain future investigations outside the planned v0.x line.
- Preserve future extension points only where inexpensive. Do not add a
  universal `Target` or processor interface before implementations prove it.

See `_docs/roadmap.md` for the public, non-binding roadmap.

## Core and Plugin Boundary

- The current `core/` package contains the SSR capabilities needed by a normal
  application. The planned v0.2 migration replaces it with a target-neutral root
  package and explicit `target/ssr`; do not implement that structure piecemeal.
- SSR, SSG, Wails, and DreeJS are first-party monorepo capabilities because they
  share compiler, render, asset, diagnostic, and compatibility contracts.
- Optional capabilities, provider integrations, and features with additional dependencies live in separate plugin repositories with their own `go.mod`, releases, tests, and CI.
- Keep optional implementations out of `core/`, even when they currently need only the standard library. SSE and WebSockets are plugins, not core packages.
- Add a provider-neutral interface to core only after at least two real implementations prove the same small contract is necessary.
- Remove the current speculative EventBus, Queue, KVStore, and Storage APIs before v0.1. The session Store remains part of the SSR core.
- Plugins may register route-specific behavior through the owning `App`; they must not weaken unrelated application defaults.
- Optional language processors such as TypeScript, Markdown, and Lua live in
  external plugin repositories and communicate through a future versioned
  process boundary. Do not add an embedded Lua VM or native Go plugin loader.

## File Structure

```
repo-root/
├── _todo/                  ← One open item per file; delete the file when its PR completes
├── CHANGELOG.md            ← What came in which version
├── README.md               ← Project overview
├── LICENSE                 ← MPL-2.0
├── go.mod                  ← Single root module (one tag per release)
├── .changes/               ← One unique release-note file per pull request
├── _docs/                  ← Public documentation
├── _plan/                  ← Detailed phased architecture and worker guidance
├── _tests/                 ← Integration tests (Docker, `make test`)
│   ├── go/                 ← Go integration tests (bug regressions, transpiler, blackbox, CLI)
│   └── fixtures/           ← Reference apps for integration tests
├── .tmp/                   ← Temporary debug spaces (no permanent tests)
│
├── core/                   ← Runtime framework facade (public API re-exports, no external deps)
├── core/internal/          ← Runtime implementation split into session/, server/, middleware/, context/, validate/
├── internal/transpiler/    ← Transpiler (.dreego → Go), used by CLI and dreegotest
├── cli/dreego/             ← CLI binary
├── .github/workflows/      ← CI: pull-request-check.yml, main-push.yml
│
└── _docs/decisions/        ← Architecture decisions (ADR)
```

## Skills

## Commit Convention

Every change lands via a pull request. The PR must contain exactly one uniquely named `.changes/*.md` file with YAML frontmatter:

```yaml
---
version: patch
---

- Bug: fix X
- Feat: add Y
```

- `version: none` — no version bump; the change file stays pending and is
  applied together with the next `version: patch` release (e.g. dependabot updates)
- `version: patch` — `0.0.x` +1, applies all pending files at once

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
- Delete the item file in the PR that completes it. The PR's change file, changelog, and Git history become the completion record.
- Do not keep completed, rejected, or superseded items in `_todo/`.
- Add newly discovered work as a new item file instead of appending to a shared checklist.

The CI (`pull-request-check.yml`) validates the change file and runs the race and full test suites. After merge, serialized `main-push.yml` refreshes to the latest `main`, reruns the suite, combines every pending change file, pushes the changelog commit with retry protection, and creates the tag. `version: none` files are never applied on their own: they stay pending until a `version: patch` file triggers the release. No local tags.

## Project: dreego

- **Name:** dreego | **Package:** `dreego` | **Module:** `github.com/dreego-stack/dreego`
- **Org:** github.com/dreego-stack | **Main repo:** github.com/dreego-stack/dreego
- **Approach:** Compile-Time Transpiler (.dreego → Go-Code) + net/http + HTMX/Alpine.js

## Note: smd

All development commands run inside `smd` (Docker container). Never run `make test`, `go build`, or any dev command directly on the host. The committed root `smd.toml` uses `golang:1.22-alpine` and includes the tools required by the test and release scripts. Run the full suite inside the container with `smd sh _tests/test.sh`; `make test` remains the Docker-based host and CI entry point.

The `smd.toml` configuration exists ONLY at the repo root. Never create `smd.toml` in subdirectories (e.g. `core/`, `demo/`, worktrees copy the root file when a container image is needed).

## Git Operations

All git operations run on the HOST via the shell subagent — never inside the
smd container. Worktree `.git` files point to
host paths that do not exist in the container.

- Worktree setup: `git worktree add -b <branch> .worktrees/<name>` (shell agent)
- Commits: `git -C <repo>/.worktrees/<name> add/commit` (shell agent, host)
- Push and PR creation: `git -C <repo>/.worktrees/<name> push origin <branch>` and `gh pr create` (shell agent)
- Never run `git checkout`, `git reset`, or `git stash` on main.
- Main must stay clean; only the `.gitignore` entry for `.worktrees/` may appear uncommitted.

## Coding Rules

- Max 300 lines per handwritten file, one logical thing per file. Generated fixture output is exempt and must not be manually split.
- No code comments (except where needed for clarity)
- Go 1.22+, prefer standard library
- Single root module `github.com/dreego-stack/dreego` (one `go.mod` at repo root, one tag per release)
- Core code in `core/` (facade re-exports) and `core/internal/` (session, server, middleware, context, validate — no external deps, enforced by `_scripts/check-core-deps.sh` in CI)
- Transpiler in `internal/transpiler/` (no external deps, same CI check; importable only from within this repo: CLI, dreegotest)
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
2. **Code** — Implement in `core/internal/` or `internal/transpiler/`; public API lives in `core/` (facade, one logical thing per file, max 300 lines)
3. **`_docs/`** — Update relevant documentation
4. **Test** — `go test ./_tests/go/ -run <TestName>` (or `make test`) — must be GREEN
5. **PR** — Create a PR with one `.changes/*.md` file (version bump + changelog lines); CI validates it
6. **Docs** — Update `_docs/` + relevant decision docs in `_docs/decisions/`

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

SSR is the current production target and v0.1 foundation. Planned v0.x work
extracts a target-neutral typed App and renderer, then adds explicit first-party
SSR, SSG, and Wails target packages. DreeJS is the optional browser layer.
There is no universal `Target` interface until working implementations prove a
small shared contract. See
[target-neutral-application-and-first-party-targets](_docs/decisions/target-neutral-application-and-first-party-targets.md)
and `_plan/`.

### 1. Server Section: No hard `*http.Request`
The current `<go>` section uses `dreego.Context`. It will be renamed to
`<server>` before v0.1; non-HTTP rendering must use explicit capabilities rather
than nil HTTP fields. → [decisions/context-design](_docs/decisions/context-design.md)

### 2. Plugin Contracts Stay Provisional Until v1
Real external plugins between v0.1 and v1 must validate the contract before a stability promise. → [plugin-contract.1](_todo/core/plugin-contract.1.md)

### 3. File-based Routing
Filename-based routing is the released pre-v0.1 implementation; the accepted
v0.1 target is one route file per URL (`+page.dreego` and method sections).
→ [decisions/routing-and-components](_docs/decisions/routing-and-components.md)

### 4. Asset System: Dual-Mode (Embedded + Disk)

### 5. Template Rendering without HTTP Server

The v0.2 render foundation must preserve typed generated inputs and expose
non-HTTP rendering before SSG or Wails is implemented. Do not make
`map[string]any` the primary page-data contract.
