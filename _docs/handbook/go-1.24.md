# Go 1.24 Knowledge Patch

Read this page when your reliable Go knowledge ends at 1.23. It contains only
the Go 1.24 delta relevant to maintaining Dreego. Verify exact API details in
the official documentation before editing code.

## Language and build tooling

- Generic type aliases are fully supported without an experiment flag.
- `go.mod` can declare executable dependencies with `tool` directives, and
  `go tool` runs them. Replace `tools.go` blank-import workarounds where present.
- `go build`, `go install`, and `go test` can emit structured JSON diagnostics.
- `go vet` checks malformed test declarations, non-constant single-argument
  `fmt.Printf` calls, invalid point-release build tags, and unsafe lock copies in
  three-clause loops.
- The Go command derives the main module version from version-control metadata
  and adds `+dirty` for uncommitted trees. Preserve Dreego's tag-derived CLI
  version contract deliberately rather than layering accidental version sources.

## Testing and runtime

- Benchmarks should use `for b.Loop() { ... }`. It runs setup once per count and
  prevents the compiler from eliminating benchmark work.
- The runtime uses Swiss Tables for maps and has lower CPU overhead by default.
  Measure Dreego again; do not preserve expectations tied to the old map runtime.
- `testing/synctest` exists only behind `GOEXPERIMENT=synctest` in 1.24. Do not
  establish a stable project contract on its experimental API.

## Filesystem and security

- `os.OpenRoot` returns an `os.Root` whose operations cannot escape the opened
  directory, including through symbolic links. Prefer it for filesystem access
  rooted in projects, plugins, documentation, templates, or extracted archives.
- `crypto/rand.Read` is guaranteed to return a nil error or terminate the
  process. Direct calls no longer need local error branches. Keep error handling
  for injected `io.Reader` values used to test failure paths.
- `crypto/cipher.NewGCMWithRandomNonce` owns nonce generation and prefixes the
  nonce to ciphertext. Adopting it changes encrypted formats, so migrate session
  payloads only with an explicit compatibility decision.
- `encoding/json` supports `omitzero`. Prefer it when zero values, especially
  `time.Time`, should be absent; do not substitute it for `omitempty` when empty
  collections or strings must also be omitted.

## Dreego migration audit

1. Convert every handwritten benchmark from `b.N` to `b.Loop`.
2. Audit path-containment code for replacement by `os.Root`, with traversal and
   symlink regression tests before changing behavior.
3. Simplify direct `crypto/rand.Read` callers while retaining injected-reader
   error tests in session cryptography.
4. Check whether VCS-derived build versions interact with Dreego's linker flag,
   build-info fallback, dirty trees, and `go install module@version` behavior.
5. Run the full race, security, release, and integration suites under Go 1.24.
6. Resolve the new vet diagnostics without compatibility suppressions.

Source: [Go 1.24 release notes](https://go.dev/doc/go1.24).
