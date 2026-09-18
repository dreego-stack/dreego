# Mixed-shape layout head (root `<head>` plus body `{#head}`)

- Area: compiler
- Phase: layout head merging
- Goal: define and enforce the behaviour of a layout that carries both a root-level `<head>` section and a body-level `{#head}` placeholder.
- Gap: `GenerateLayout` writes `head` at the function top when `file.Head != nil` (`internal/transpiler/codegen_layout.go:15-17`), while the body walk writes it again at the placeholder through `genLayoutNodeState` (`internal/transpiler/codegen_layout.go:105-128`); because the body-level dedupe block is gated on `file.Head == nil` (`internal/transpiler/codegen_layout.go:31`), the merged head is emitted twice.
- Gap: the golden fixture `internal/transpiler/testdata/golden/route_with_layout.golden` is generated from exactly this mixed shape (`internal/transpiler/codegen_golden_test.go:99`), but the golden covers only the route code — no test renders the layout, so the double emission is unobserved.
- Acceptance: decide whether the mixed shape is valid; if valid render the head once, if not fail at generate with a diagnostic; add a rendering test either way.
- Depends on: head-dedupe core fix (slice 1).
