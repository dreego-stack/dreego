# Go 1.23 Knowledge Patch

Read this page when your reliable Go knowledge ends at 1.22. It contains only
the Go 1.23 delta relevant to maintaining Dreego. Verify exact API details in
the official documentation before editing code.

## Language and collections

- `range` accepts iterator functions with zero, one, or two yielded values.
- The new `iter` package defines `Seq` and `Seq2`.
- `slices` and `maps` provide iterator producers, collectors, sorting helpers,
  and chunking. Prefer them when they remove hand-written collection plumbing;
  ordinary direct loops remain idiomatic.
- Generic type aliases are preview-only in 1.23. Do not make production code
  depend on `GOEXPERIMENT=aliastypeparams`.

## Runtime and tools

- Timer and ticker channels are synchronous for modules declaring Go 1.23 or
  newer. `Stop` and `Reset` no longer leave a stale value to receive. Unused
  timers and tickers can be garbage-collected. Never inspect timer channel
  length or capacity to predict a receive.
- `go vet` includes `stdversion`, which rejects standard-library symbols newer
  than the effective Go version.
- `go mod tidy -diff` checks module-file cleanliness without rewriting files.
- Go telemetry is local-only by default. Upload remains an explicit opt-in.

## Standard library

- `filepath.Localize` converts a valid slash-separated `io/fs` path to an OS
  path and rejects values that cannot be represented safely. Use it at trusted
  `io/fs`-to-OS boundaries instead of assuming `filepath.FromSlash` validates.
- `net/http.Request.Pattern` exposes the matched `ServeMux` pattern.
- Cookies support the `Partitioned` attribute and parsing complete `Cookie` or
  `Set-Cookie` header values.
- `ServeContent`, `ServeFile`, and `ServeFileFS` remove representation headers
  from error responses. Recheck middleware that adds compression or cache
  headers around those helpers.
- `runtime/debug.SetCrashOutput` can duplicate fatal runtime output to a crash
  monitor.
- Go 1.23 requires macOS 11 or newer. It is the final release supporting Linux
  kernels older than 3.2.

## Dreego migration audit

1. Change the module and container baseline together before using new APIs.
2. Run `go vet ./...` and `go mod tidy -diff` inside `smd`.
3. Audit timer reset, stop, and timeout tests for assumptions from Go 1.22.
4. Audit plugin, documentation, component, and static-asset path conversion for
   boundaries where `filepath.Localize` can replace validation plus conversion.
5. Review collection-building loops for clearer iterator APIs without replacing
   straightforward loops mechanically.
6. Verify documented deployment targets against the new OS minimums.

Source: [Go 1.23 release notes](https://go.dev/doc/go1.23).
