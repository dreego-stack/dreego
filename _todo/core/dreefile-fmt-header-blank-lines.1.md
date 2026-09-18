# `dreego fmt` collapses duplicate blank lines inside the header

- Area: cli
- Phase: Dreefile grammar (slice 3a)
- Goal: pin the header blank-line behaviour of `dreego fmt` and document the change.
- Gap: `Format` collapses the header through `multiBlank.ReplaceAllString` (`internal/transpiler/fmt.go:79-80`), whereas main kept duplicate blank lines in the header verbatim; existing fmt tests and goldens are green, but no test or golden depends on the old behaviour today, so the change is undocumented and unprotected.
- Acceptance: either a test pins the new behaviour or the collapse is reverted to verbatim; the decision is documented either way.
- Depends on: nothing.
