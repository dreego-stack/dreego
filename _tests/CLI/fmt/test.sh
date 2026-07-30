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

cat > messy.dreego << 'DREEGO'
<head>
    <title>test</title>
</head>


<go>
    msg   :=    "hello"
</go>

<div>
  <p>{  msg  }</p>
    {#if  show}
        <span>visible</span>
    {/if}
</div>
DREEGO

cat > expected.dreego << 'DREEGO'
<head>
    <title>test</title>
</head>

<go>
    msg   :=    "hello"
</go>

<div>
  <p>{msg}</p>
    {#if show}
        <span>visible</span>
    {/if}
</div>
DREEGO

go run codeberg.org/dreego/dreego/cmd/dreego fmt --stdout messy.dreego > actual.dreego

if diff expected.dreego actual.dreego > /dev/null 2>&1; then
    echo "ok: basic formatting"
else
    echo "FAIL: basic formatting"
    diff expected.dreego actual.dreego
    exit 1
fi

go run codeberg.org/dreego/dreego/cmd/dreego fmt --stdout actual.dreego > formatted_twice.dreego
if diff actual.dreego formatted_twice.dreego > /dev/null 2>&1; then
    echo "ok: idempotent"
else
    echo "FAIL: not idempotent"
    diff actual.dreego formatted_twice.dreego
    exit 1
fi

echo "ok: all fmt tests passed"
