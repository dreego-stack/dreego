# Go 1.27 Knowledge Patch

Read this page when your reliable Go knowledge ends at 1.26. It contains only
the Go 1.27 delta relevant to maintaining Dreego. Verify exact API details in
the official documentation before editing code.

## Language and tools

- Methods may declare type parameters. Interface methods still cannot be
  generic, and a generic method cannot implement an interface method.
- Struct literal keys may use valid promoted field selectors.
- Generic function inference also applies when assigning or converting a
  function to a matching function type.
- `go test` runs the `stdversion` analyzer by default.
- `go fix` adds modernizers for typed atomics, embedded literals, backward slice
  iteration, and unsafe functions.
- `go mod tidy` consolidates requirements into at most one direct and one
  indirect block while preserving dependency comments.

## Runtime and diagnostics

- Timer channels are always synchronous; the old `asynctimerchan` compatibility
  mode has been removed.
- Goroutine leak profiling is generally available through `runtime/pprof` and
  `/debug/pprof/goroutineleak`.
- Tracebacks include goroutine labels by default. Treat labels as potentially
  sensitive data before attaching user, route, or request information.
- Small allocations are faster but add roughly 60 KB to binaries. Re-baseline
  performance and binary-size assertions rather than preserving old numbers.

## Standard library

- `encoding/json/v2` and `encoding/json/jsontext` are available without an
  experiment. V2 rejects invalid UTF-8 and duplicate object names by default.
  The v1 API remains supported and uses the new implementation with compatible
  behavior, although exact error text may change.
- The new `uuid` package provides standard UUID generation and parsing. Do not
  change Dreego request-ID format unless its public contract is migrated.
- `net/http.Server.MaxHeaderValueCount` adds a direct request-hardening limit.
- `net/http/httptest.NewTestServer` uses an in-memory fake network compatible
  with `testing/synctest`.
- `URL.Clone`, `Values.Clone`, `strings.CutLast`, and `bytes.CutLast` replace
  common hand-written copying and last-separator logic.
- `testing/synctest.Sleep` combines virtual sleep with synchronization.
- Unicode tables advance from Unicode 15 to Unicode 17.
- Go 1.27 requires macOS 13 or newer.

## Dreego migration audit

1. Evaluate a deliberate move of configuration, catalog, manifest, and build
   metadata parsing to `encoding/json/v2`; add duplicate-key, invalid-UTF-8,
   unknown-field, and error-contract tests before selecting strict defaults.
2. Run all `go fix` modernizers and remove superseded code after review.
3. Add goroutine-leak checks to lifecycle and concurrency verification.
4. Evaluate `MaxHeaderValueCount` alongside existing request-size and header
   limits.
5. Replace real-time HTTP tests with `httptest.NewTestServer` and `synctest`
   where this removes ports, polling, or sleeps without weakening integration
   coverage.
6. Re-run lexer, identifier, locale, case-folding, and pseudolocale tests against
   Unicode 17.
7. Verify redirect, compression, TLS, JSON error, macOS deployment, binary-size,
   and benchmark contracts under the new defaults.

Source: [Go 1.27 release notes](https://go.dev/doc/go1.27).
