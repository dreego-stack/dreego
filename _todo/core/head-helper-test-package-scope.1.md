# Head-helper uniqueness assertion is package-scoped

- Area: tests
- Phase: head dedupe hardening
- Goal: scope the head-helper uniqueness assertion to one generated package so a body-level layout may emit its own copy.
- Gap: `_tests/go/bug_head_helpers_multi_dir_test.go:43-48` counts `func stripTitleTag(` and `func stripMetaDescriptionTag(` across ALL concatenated generated files; since `layoutNeedsHeadHelpers` (`internal/dreefile/generate.go:180,214`) lets the layouts package emit its own copy, any fixture whose body-level layout has a title or description before `{#head}` would define the helper twice and fail this test for a non-bug.
- Acceptance: the count is scoped per package and a fixture with a body-level layout title passes.
- Depends on: nothing.
