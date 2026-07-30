# How to write `test.sh`

Every integration test in `_tests/` follows the same pattern.

## Template

```sh
#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
set -e

realrepo="$(cd "$(dirname "$0")"/../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

cat > go.mod << EOF
module t
go 1.22
require codeberg.org/dreego/dreego v0.0.0
replace codeberg.org/dreego/dreego => $realrepo
EOF

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<div><p>hello</p></div>
DREEGO

go run codeberg.org/dreego/dreego/cmd/dreego generate

grep -q 'hello' dreego/gen/dree.go
echo "ok"
```

## Rules

0. **First line** after `#!/bin/sh` must be `# Using standard: _tests/how-to-test-sh.md` — makes non-compliant files discoverable via `head -1`
1. **`realrepo`** — absolute path to repo root, always `../../..` from `_tests/<Group>/<name>/` (3 levels up)
2. **`workdir`** — always `mktemp -d`, never create files inside `_tests/`
3. **`trap "rm -rf $workdir" EXIT`** — mandatory cleanup on success *and* failure
4. **CLI** — always `go run "$realrepo/cmd/dreego"`, never a pre-built binary
5. **`go.mod`** — always fresh with `require` + `replace` to `$realrepo` (use `cat > go.mod`, never `go mod init`)
6. **`mkdir -p dreego/routes`** — scaffold minimal project structure as needed
7. **No files left behind** — test does all I/O inside `$workdir`
8. **Random port for server tests** — if the test starts an HTTP server, use a random port to avoid conflicts:
   - Add `port=$(awk 'BEGIN{srand();print int(rand()*50000)+10000}')` after `cd "$workdir"`
   - Write `main.go` with `:8080` inside the heredoc
   - Add `sed -i "s/8080/$port/" main.go` after the `main.go` heredoc (before `go run`)
   - Use `localhost:$port` in all `curl` commands, never `localhost:8080`

## Why

| Concern | Old way | New way |
|---------|---------|---------|
| Cleanup after failure | ❌ leftovers pollute repo | ✅ `trap` deletes on exit |
| Parallel safety | ❌ shared state | ✅ each test gets unique `mktemp -d` |
| CLI freshness | ❌ pre-built binary may be stale | ✅ `go run` compiles from current source |
| Portability | ❌ relative paths break in Docker | ✅ `$realrepo` is absolute |

## Example

See [Static/subdir/test.sh](./Static/subdir/test.sh) for a complete working example.
