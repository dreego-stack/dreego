# Go 1.25 Knowledge Patch

Read this page when your reliable Go knowledge ends at 1.24. It contains only
the Go 1.25 delta relevant to maintaining Dreego. Verify exact API details in
the official documentation before editing code.

## Language and tools

- Go 1.25 adds no language change that affects Go programs.
- `go vet` adds `waitgroup`, which detects `WaitGroup.Add` inside the launched
  goroutine, and `hostport`, which prefers `net.JoinHostPort` for IPv6-safe
  network addresses.
- The `ignore` directive in `go.mod` excludes directories from package-pattern
  matching. It does not replace Dreego's own project-file discovery rules.
- Updating a `go` directive no longer adds a `toolchain` directive automatically.

## Runtime and testing

- On Linux, default `GOMAXPROCS` observes cgroup CPU limits and updates when the
  limit changes. Avoid redundant container CPU tuning unless measurements prove
  an application-specific requirement.
- `testing/synctest` is stable. Use `synctest.Test` and `synctest.Wait` for
  deterministic concurrent tests whose correctness depends on timers or
  goroutines instead of real sleeps and polling deadlines.
- `WaitGroup.Go` starts and accounts for a goroutine together. Prefer it over a
  separate `Add` plus `go` statement when the surrounding error and panic
  semantics fit.
- The compiler now performs nil checks correctly when a result is used before
  its accompanying error is checked. Always check errors immediately before
  using other returned values.
- DWARF 5 is the default debug format and may change binary sizes and tooling
  behavior. Verify release and debugging workflows rather than disabling it.

## Standard library

- `net/http.CrossOriginProtection` rejects unsafe cross-origin browser requests
  using Fetch Metadata and origin information. Evaluate it as defense in depth;
  do not silently replace Dreego's token-based CSRF contract.
- `os.Root` gains write, rename, link, symlink, recursive removal, and metadata
  operations, making it suitable for more complete rooted workflows.
- `go/parser.ParseDir` and legacy `go/ast` package-merging APIs are deprecated.
- `encoding/json/v2` and `jsontext` are experimental behind
  `GOEXPERIMENT=jsonv2`. Do not use them in Dreego's stable contracts yet.
- The Green Tea garbage collector remains experimental. Benchmark it separately;
  do not require it in builds or tests.
- Go 1.25 requires macOS 12 or newer.

## Dreego migration audit

1. Resolve every new `go vet` finding instead of suppressing the analyzers.
2. Replace time-based concurrency tests with stable `testing/synctest` where the
   test models blocking, cancellation, shutdown, or timeout behavior.
3. Audit all multi-result calls for use-before-error-check ordering.
4. Revisit rooted project and plugin file operations using the expanded
   `os.Root` API.
5. Benchmark server behavior inside CPU-limited containers and update expected
   concurrency only from measured results.
6. Keep experimental JSON and garbage-collector modes out of the supported
   production baseline.

Source: [Go 1.25 release notes](https://go.dev/doc/go1.25).
