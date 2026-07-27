#!/bin/sh
cd "$(dirname "$0")"
dreego init .
rm -rf dreego/routes dreego/components dreego/layouts 2>/dev/null
mkdir -p dreego/routes
cat > dreego/routes/get.dreego << 'DREEGO'
<div><p>first</p></div>
<div><p>second</p></div>
DREEGO
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
if dreego generate 2>/dev/null; then echo "expected failure but succeeded"; exit 1; fi
echo ok
