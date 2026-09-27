---
version: patch
---

- Test: rework `_tests/` into one recursive tree (`_tests/e2e/<group>/<test>/`) with a single orchestrator `_tests/test.sh`; Go and shell tests are discovered and run the same way, each test owns its folder and input files
- Test: remove the `_tests/fixtures` reference-app scaffolding and `dreegotest.Fixture`/`ServeFixture`; tests write their own project or use `dreego new`
- Test: move the repo-wide gates (layering, core-deps, no-binaries, coverage, fuzz) and the release-prep contract into `_tests/e2e/invariants/` and `_tests/e2e/release/` as ordinary tests
- Test: `task test` runs the unit tests once and the recursive tree once; the duplicate `go test`, coverage and fuzz steps in `pull-request-check.yml` are gone and `main-push.yml` runs the suite once, cutting the CI test time roughly in half
