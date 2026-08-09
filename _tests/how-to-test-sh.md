# How to write `test.sh`

Every integration test in `_tests/` follows the same pattern. The test runner (`_tests/test.sh`) pre-compiles the dreego CLI and exports it as `$DREEGO_BIN`. Tests use `$DREEGO_BIN` instead of `go run`.

## Template

```sh
#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: <one-line summary of what this test verifies>
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

cat > go.mod << EOF
module t
go 1.22
require github.com/dreego-stack/dreego v0.0.0
replace github.com/dreego-stack/dreego => $realrepo
EOF

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<div><p>hello</p></div>
DREEGO

$DREEGO_BIN generate

grep -q 'hello' dreego/gen/dree.go
echo "ok"
```

## Rules

0. **First line** after `#!/bin/sh` must be `# Using standard: _tests/how-to-test-sh.md` — makes non-compliant files discoverable via `head -1`
0a. **Second line** must be `# What: <summary>` — one-line description of what this test verifies
1. **`realrepo`** — absolute path to repo root, always `../../../..` from `_tests/core/<Group>/<name>/` (4 levels up). For `_tests/plugins/<name>/` adjust depth accordingly.
2. **`$DREEGO_BIN`** — pre-compiled dreego CLI, exported by the test runner. Use `$DREEGO_BIN generate` / `$DREEGO_BIN run` etc. Never `go run`. Exception: tests that spawn a real subprocess and signal it (e.g. `run-timer-sigterm`) use `go run "$realrepo/cli/dreego"` instead — the runner's pre-compiled binary would be reused across tests and the signal would hit the wrong process.
3. **`workdir`** — always `mktemp -d`, never create files inside `_tests/`
4. **`trap "rm -rf $workdir" EXIT`** — mandatory cleanup on success *and* failure
5. **`go.mod`** — always fresh with `require` + `replace` for `github.com/dreego-stack/dreego` (use `cat > go.mod`, never `go mod init`)
6. **`mkdir -p dreego/routes`** — scaffold minimal project structure as needed
7. **No files left behind** — test does all I/O inside `$workdir`
8. **Port for server tests** — the test runner (`_tests/test.sh`) exports `DREEGO_PORT` for every test, deterministically ascending from `DREEGO_PORT_BASE` (default `20000`), so ports never overlap in a parallel run:
   - Read the port after `cd "$workdir"`:
     ```sh
     port="${DREEGO_PORT:-$(( ( $(od -An -N2 -i /dev/urandom | tr -d ' ') % 50000 ) + 10000 ))}"
     ```
     The `od` fallback only applies when the test runs standalone without the runner.
   - Write `core.Listen(":$port")` **directly** into the `main.go` heredoc. The heredoc must be **unquoted** (`<< GO`, not `<< 'GO'`) so `$port` expands.
   - Do **not** use `sed -i "s/8080/$port/" main.go` anymore — it is non-portable (BSD-sed) and flakes.
   - Use `localhost:$port` in all `curl` commands, never `localhost:8080`.
   - `curl` is installed once, sequentially, by the runner before any test runs. In a test it is enough to add `apk add --no-cache curl >/dev/null 2>&1 || true`.

## Why

| Concern | Old way | New way |
|---------|---------|---------|
| Cleanup after failure | leftovers pollute repo | `trap` deletes on exit |
| Parallel safety | shared state | each test gets unique `mktemp -d` |
| CLI freshness | `go run` recompiles each time | `$DREEGO_BIN` pre-compiled once |
| Portability | relative paths break in Docker | `$realrepo` is absolute |

## Example

See [Static/subdir/test.sh](./core/Static/subdir/test.sh) for a complete working example.