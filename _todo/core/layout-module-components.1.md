# Layout components imported from an external module

- Area: compiler
- Phase: layout and component discovery
- Goal: discover and compile a component that an external module provides when a layout, not a route, is the only file importing it.
- Gap: `importedComponentPaths` (`internal/transpiler/generate_components.go:141-169`) walks only `root/routes`, so a layout header import such as `import "example.com/ui/components"` is invisible during discovery and the component remains unregistered.
- Gap: the layout parser sets `f.Imports = imports` (`internal/transpiler/generate_layout.go:60`) but no consumer reads `File.Imports`; the state mirrors route parsing and is currently dead.
- Acceptance: a layout using an external-module component is discovered and the generated app builds; a regression test covers the layout-only import; the follow-up decides whether `File.Imports` becomes the discovery input or is removed/documented, and `gen.Src` retaining the last layout after the loop (`internal/transpiler/generate_layout.go:178`) is either consumed, reset, or documented as harmless.
- Depends on: layout component generation fix (component imports in generated layouts, layout header import parsing).
