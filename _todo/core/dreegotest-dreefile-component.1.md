# `dreegotest` returns an empty component for DREEFILE components

- Area: tests
- Phase: Dreefile grammar (slice 3a)
- Goal: let `dreegotest` generate components declared with the new `DREEFILE component` grammar.
- Gap: the component name now comes from the filename, so `ParseHeader` leaves `ComponentDef.Name` empty for DREEFILE components; `generateComponent` bails out early on `comp == nil || comp.Name == ""` (`dreegotest/generate.go:124-126`) and silently returns an empty string, so the public helper cannot build such components.
- Acceptance: `dreegotest` generates a DREEFILE component correctly; a test covers it.
- Depends on: nothing.
