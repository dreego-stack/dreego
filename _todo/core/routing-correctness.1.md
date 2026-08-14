---
area: routing
phase: v0.1-blocker
---
# Make file-based routing match its contract

## Goal
Make every route source map deterministically to the documented method and URL.

## Acceptance criteria
- Catch-all directories generate valid Go 1.22 patterns such as `{path...}`.
- Double-bracket optional segments are rejected with a source-aware diagnostic; developers define each route explicitly.
- Flat route files and `+page.dreego` map deterministically to one URL;
  `index.dreego` and conflicting flat/directory definitions fail with a source path.
- `<go>` and `<div>` default to GET; explicit method attributes register only
  their declared methods and must form a valid response pair.
- Successful method logic renders the matching method-specific `<div>`;
  redirects are explicit and errors enter the App error path.
- Duplicate generated, plugin, user, and reserved framework routes fail instead of silently overriding.
- Black-box tests send real requests for every documented route form and method.
