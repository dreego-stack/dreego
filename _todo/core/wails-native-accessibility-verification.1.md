# Verify Wails reference application accessibility on a desktop host

- Area: accessibility
- Phase: Wails v3
- Goal: verify the interactive behavior of the Wails reference application on a real desktop host.
- Gap: the automated gate in `_tests/go/wails_timer_demo_test.go` only asserts markup strings in the generated Go. Tab order, Enter/Space activation, focus visibility, and screen-reader output are documented as a manual check in `demo/demo-wailsv3/README.md` but have never been performed.
- Acceptance: the reference application is checked with a desktop screen reader and keyboard on at least one supported platform, and the result is recorded.
- Depends on: nothing.
