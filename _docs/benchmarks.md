# Benchmarks

Reproducible Go benchmarks measure code generation and representative request
handling. They exist to support performance claims with measured conditions and
to track regressions — they do not enforce machine-specific absolute timings.

## Running

```bash
go test ./internal/transpiler/ -bench=. -benchtime=1s -run=^$ -count=3
go test ./core/ -bench=. -benchtime=1s -run=^$ -count=3
```

Benchmarks live in `internal/transpiler/benchmark_test.go` (code generation)
and `core/benchmark_test.go` (request pipeline) and run without external
dependencies.

## Benchmarks

| Benchmark | Measures |
|-----------|----------|
| `BenchmarkGeneratePage` | Full transpile of a representative page (head + go + div with if/each/component call) |
| `BenchmarkGenerateComponent` | Full transpile of a component with props |
| `BenchmarkRequestPage` | Full request through the built handler chain: SSR page render with escaping, if, each, component call |
| `BenchmarkRequestJSON` | Full request through the built handler chain: JSON API response |
| `BenchmarkRequestSimple` | Full request through the built handler chain: minimal handler |

`BenchmarkRequestPage` uses a hand-written render function that mirrors the
shape of the generated code, not the actual transpiler output.

## Reference Results

Measured 2026-08-15 in the smd container (`golang:1.22-alpine`, linux/arm64,
`-benchtime=200x`). These numbers are environment-specific and will vary by
machine; treat them as a baseline, not a guarantee.

| Benchmark | ns/op | B/op | allocs/op |
|-----------|-------|------|-----------|
| BenchmarkGeneratePage | 46255 | 37221 | 238 |
| BenchmarkGenerateComponent | 5259 | 6607 | 54 |
| BenchmarkRequestPage | 3035 | 2384 | 37 |
| BenchmarkRequestJSON | 2751 | 2192 | 30 |
| BenchmarkRequestSimple | 2650 | 1616 | 23 |

## Regression Tracking

- Run the benchmarks before and after a change that touches the transpiler or
  the request pipeline.
- Compare relative changes (e.g. "page generation got 20% slower") rather than
  absolute ns values, which are not comparable across machines.
- Do not add CI gates on absolute timings; they are unstable across runners.
