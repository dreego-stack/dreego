# Gate the native Wails integration tests in CI

- Area: tests
- Phase: Wails v3
- Goal: run the native Wails integration tests in a CI pipeline and keep them green.
- Gap: `_tests/go/wails_native_no_tcp_linux_test.go`, `_tests/go/wails_native_navigation_linux_test.go`, and `_tests/go/wails_native_output_linux_test.go` are guarded by `//go:build linux && !race`. Every CI test invocation uses `-race`, so no pipeline ever compiles or runs those three tests.
- Gap: a CI step was attempted and reverted. On GitHub Actions (ubuntu-24.04, amd64, GTK 4.14.5, WebKitGTK 2.52.6) the Wails app starts and logs Build/Platform/AssetServer info, but no window is ever mapped, so both native tests hit their 15s deadline. The same tests pass in the smd container (Alpine, aarch64, GTK 4.22.4, WebKitGTK 2.48.1).
- Gap: a diagnosis harness exists on the local branch `feat/native-wails-ci-gating` (commit `a8c48c4`, NOT merged): a shared `_tests/go/wails_native_harness_linux_test.go` with a per-test free display instead of a hardcoded `:97`, an `xdpyinfo`-based readiness probe instead of a socket-existence check, and a `DIAG:` failure block (xwininfo tree, helper liveness, sandbox markers, Xvfb output).
- Acceptance: a CI step runs the native Wails tests and is green on the GitHub runner, or the tests are explicitly excluded from the pipeline with a documented reason; the shared harness is used instead of duplicated per-test setup.
- Depends on: nothing.
