---
area: routing
phase: v0.1-blocker
depends_on: pre-v0.1-product-decisions.1
---
# Make file-based routing match its contract

## Goal
Make every route source map deterministically to the documented method and URL.

## Acceptance criteria
- Catch-all directories generate valid Go 1.22 patterns such as `{path...}`.
- The selected optional-segment behavior is implemented and tested by HTTP.
- Only documented method filenames are accepted; unknown names fail with a source path.
- Duplicate generated, plugin, user, and reserved framework routes fail instead of silently overriding.
- Black-box tests send real requests for every documented route form and method.
