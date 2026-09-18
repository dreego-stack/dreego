# Dreefile component aliases resolve globally instead of per file

- Area: compiler
- Phase: Dreefile grammar (slice 3a)
- Goal: give `COMPONENT ... IMPORT { X as Y }` clear, file-scoped semantics.
- Gap: `collectComponentAliases` (`internal/transpiler/generate_component_aliases.go:11-42`) walks every `.dreego` file under the website root and registers each alias into the generator-global `State.CompAliases` (`internal/transpiler/codegen/state.go:17`), so an alias declared in one route resolves in every other route, component and layout although `COMPONENT ... IMPORT` is specified as a file-scoped import line.
- Gap: an alias may target a component in another package; `State.Qualify` (`internal/transpiler/codegen/state.go:90-102`) then emits an import for a package the file never asked for.
- Acceptance: a test proves the chosen semantics; either aliases are scoped per file, or the global behaviour is documented as intended.
- Depends on: nothing.
