---
version: patch
---

- Fix: move CLI from `cmd/dreego/` to `cli/dreego/` — `go install github.com/dreego-stack/dreego/cli/dreego@latest` now works and tracks the single root tag
- Fix: `go install github.com/dreego-stack/dreego/cmd/dreego@latest` was broken because the old `cmd/dreego` path was cached by the Go module proxy as a separate submodule (its v0.0.22–v0.0.26 tags declared a stale `codeberg.org/dreego/dreego` module path). The proxy cache is immutable, so that path can never resolve to the root module again — renaming to `cli/dreego/` gives a fresh, never-seen path that resolves via the single root tag.
- Fix: update all build scripts, Makefile, Dockerfiles, tests, and docs that referenced the old `cmd/dreego` path.
