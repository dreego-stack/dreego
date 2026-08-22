# Contributing to Dreego

Thanks for contributing. This guide covers the workflow every change goes
through, from local setup to merge. Please also read `AGENTS.md` at the repo
root — it is the authoritative source for conventions.

## Development Setup

- Go 1.22 or newer.
- No external runtime dependencies in `core/` — the standard library only
  (enforced by `_scripts/check-core-deps.sh` in CI).
- Run the full suite with `make test` (Docker-based; builds the `_tests/`
  image and runs the race and integration suites). For a fast subset:
  `go test ./core/... ./cli/dreego/...`.
- All commands may also be run inside the `smd` container. Use
  `smd sh _tests/test.sh` for the full suite; this is optional and never
  required for contributions.
- Never create binaries in the repo — use `/tmp` or `./tmp`.

## Language

Everything in this repository is English: code, comments, commit messages,
tests, and documentation.

## PR Workflow

Every change lands via a pull request. One PR normally implements one todo
item from `_todo/` (delete the item file in the PR that completes it).

1. Create a branch, implement the change, run the relevant tests.
2. Add exactly one uniquely named `.changes/<name>.md` file with YAML
   frontmatter:

   ```yaml
   ---
   version: patch
   ---

   - Bug: fix X
   - Feat: add Y
   ```

3. Open the PR. CI (`pull-request-check.yml`) validates the change file and
   runs the race and full test suites.

After merge, `main-push.yml` combines pending change files into the
changelog and creates the tag. No local tags are required.

### Version bumps

- `version: patch` — `0.0.x` +1.
- `version: none` — changelog lines only, no version bump (e.g. dependency
  updates).
- NEVER use `version: minor` or `version: major` until the v0.1 release —
  a minor bump would tag v0.1.0 prematurely.

## Commit Conventions

- English commit messages.
- One logical change per commit.
- Generated files (`dree.go`) are not committed.
- Only the `.gitignore` entry for `.worktrees/` may appear uncommitted on
  main.

## Bug → Test → Fix Workflow

Every bug gets a permanent regression test.

1. Create `_tests/go/bug_<name>_test.go` that reproduces the bug — it must
   FAIL against the current code.
2. Fix the code until the new test is GREEN.
3. The bug is permanently covered — no regression risk.

`.tmp/<name>/` is only for temporary debugging and exploration, never for
permanent tests.

## Feature Workflow

Every feature follows this cycle:

1. `_tests/` — integration test in `_tests/go/<name>_test.go` using
   `dreegotest` (see `_docs/testing.md` and existing tests for the pattern).
2. `core/` — implementation, one logical thing per file.
3. `_docs/` — update relevant documentation.
4. Test — `go test ./_tests/go/ -run <TestName>` (or `make test`) — GREEN.
5. PR — one `.changes/*.md` file (version bump + changelog lines).
6. `_docs/` — update decision docs in `_docs/decisions/`.

## Style Rules

- Max 300 lines per handwritten file, one logical thing per file. Generated
  fixture output is exempt and must not be manually split.
- No comments unless needed for clarity.
- Go 1.22+, standard library preferred. Core code in `core/` has no external
  dependencies.
- CLI lives in `cli/dreego/` and imports core; plugins live in separate
  repositories under `github.com/dreego-stack/`.
- Build via the `dreego` CLI, not directly `go build`.

## Docs

Public documentation lives in `_docs/`. New docs must be added to
`_docs/sitemap.json` (in the existing entry format) so `dreego docs --list`
shows them. Keep the README documentation table in sync when adding docs.
