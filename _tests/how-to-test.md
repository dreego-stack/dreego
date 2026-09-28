# _tests layout

`_tests/` holds black-box and end-to-end tests. One test lives in one folder.

```
_tests/
  test.sh                 recursive orchestrator (the only entry point)
  run-one.sh              runs a single test folder (used by test.sh)
  how-to-test.md          the contract every test.sh / <name>_test.go follows
  e2e/
    <group>/<test-name>/
      test.sh             a shell test
      <name>_test.go      a Go test (package tests)
      www/…               real input files the test needs, if any
```

## Rules

- One test = one folder. A folder contains `test.sh`, or a Go `_test.go` file, or both (two tests).
- Every test is standalone: it writes its own project into a temp dir, or uses
  `dreego new` (which loads its template from the CLI binary). There is no
  shared fixture or template system.
- No cross-folder helpers. A helper belongs to the test that uses it.
- Input files (`www/…`) live next to the test and are the test's own data, not a
  reusable template.

## Writing a test

### Shell (`test.sh`)

```sh
#!/bin/sh
# What: one line describing what this test proves.
set -eu
cd "$(dirname "$0")"
work="$(mktemp -d)"; trap 'rm -rf "$work"' EXIT
"$DREEGO_BIN" new app -t web-minimal      # template comes from the binary
cd "$work/app" && "$DREEGO_BIN" generate && go build -o /dev/null .
echo ok
```

### Go (`test.go`, package `tests`)

```go
package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestThing(t *testing.T) {
	dreegotest.MustBuild(t, map[string]string{
		"www/routes/+page.dreego": `<body><p>hi</p></body>`,
	})
}
```

## Environment (provided by test.sh)

- `REPO_DIR` — repository root (absolute).
- `DREEGO_BIN` — the freshly built `dreego` CLI.
- `DREEGO_LOCAL_REPO` — same as `REPO_DIR`; read by `dreegotest`.
- Success = exit code 0 and a final `ok` line.
