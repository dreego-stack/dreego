# Go 1.26 Knowledge Patch

Read this page when your reliable Go knowledge ends at 1.25. It contains only
the Go 1.26 delta relevant to maintaining Dreego. Verify exact API details in
the official documentation before editing code.

## Language and modernization

- `new` accepts an expression and returns a pointer initialized with its value.
  Prefer `new(value)` over local pointer-helper functions or temporary variables.
- Generic constraints may refer recursively to the generic type being
  constrained.
- `go fix` is now an analyzer-based suite of modernizers for current Go idioms
  and library APIs. Run it as a reviewed migration input: inspect every diff,
  retain behavior tests, and do not preserve superseded forms after approval.
- `go mod init` intentionally writes the previous supported Go version rather
  than the running toolchain version. Generated Dreego applications must set
  their intended baseline explicitly.

## Runtime and diagnostics

- The Green Tea garbage collector is enabled by default. Re-run allocation,
  rendering, request, and transpiler benchmarks without carrying the temporary
  Go 1.25 experiment flag.
- Heap base addresses are randomized on 64-bit systems.
- Goroutine leak profiling exists behind
  `GOEXPERIMENT=goroutineleakprofile`. Use it for migration diagnostics and CI
  experiments, but do not make it part of Dreego's stable tooling contract yet.

## Standard library

- `errors.AsType[T]` is a type-safe generic alternative to most `errors.As`
  calls.
- `testing.ArtifactDir` provides a managed location for retained test output.
- `testing.B.Loop` no longer prevents inlining, removing the remaining reason
  to keep legacy `b.N` benchmark loops.
- `testing/cryptotest.SetGlobalRandom` supplies deterministic cryptographic
  randomness for tests. Several crypto APIs now ignore caller-provided random
  readers and use secure global randomness instead; injected readers used by
  Dreego's own code remain ordinary `io.Reader` contracts.
- `net/http.ServeMux` trailing-slash redirects use status 307 instead of 301.
- `net/http/httputil.ReverseProxy.Director` is deprecated because hop-by-hop
  headers make it unsafe; use `Rewrite`.
- `net/url.Parse` rejects malformed hosts containing ambiguous colons.
- `os/signal.NotifyContext` records the triggering signal as its cancellation
  cause.
- Reflection exposes iterator methods for fields, methods, inputs, and outputs.
- Go 1.26 is the last release supporting macOS 12.

## Dreego migration audit

1. Run `go fix ./...`, review each transformation, and add tests for any change
   whose behavioral equivalence is not obvious.
2. Find pointer helpers and temporary address variables replaceable by
   `new(value)`.
3. Replace remaining `errors.As` and `b.N` forms where the new APIs apply.
4. Recheck redirect status expectations, URL rejection, dev proxies, signal
   shutdown behavior, and generated application `go.mod` versions.
5. Compare benchmark and race results under the new garbage collector.
6. Exercise goroutine-leak profiling against server lifecycle, dev watcher,
   middleware, and Markdown concurrency tests without enabling it by default.

Source: [Go 1.26 release notes](https://go.dev/doc/go1.26).
