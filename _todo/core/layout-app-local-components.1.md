# Layout cannot reference an app-local component

- Area: compiler
- Phase: components / layouts
- Goal: a layout may use a component from the app-local `components/` tree of the app that uses it.
- Gap: `buildRootPlan` generated layouts before the per-app `scanComponentsAt` (`internal/dreefile/generate.go`), so an app-local component referenced by a layout failed with `unknown component`.
- Fix: register every app's local components before `generateLayouts`, so layout codegen can resolve them; each app pass still regenerates its own component files.
- Acceptance: a layout referencing `<app>/components/…` generates and calls the app-local component; covered by `internal/dreefile/bug_layout_app_local_component_test.go` and `_tests/go/bug_layout_app_local_component_test.go`.
