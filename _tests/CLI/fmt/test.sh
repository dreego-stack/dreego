#!/bin/sh
set -e
cd "$(dirname "$0")"

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

dreego fmt --stdout messy.dreego > actual.dreego

if diff expected.dreego actual.dreego > /dev/null 2>&1; then
    echo "ok: basic formatting"
else
    echo "FAIL: basic formatting"
    diff expected.dreego actual.dreego
    exit 1
fi

dreego fmt --stdout actual.dreego > formatted_twice.dreego
if diff actual.dreego formatted_twice.dreego > /dev/null 2>&1; then
    echo "ok: idempotent"
else
    echo "FAIL: not idempotent"
    diff actual.dreego formatted_twice.dreego
    exit 1
fi

echo "ok: all fmt tests passed"
