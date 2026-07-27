#!/bin/sh
set -e
cd "$(dirname "$0")"
dreego init .
rm -rf dreego/routes dreego/components dreego/layouts 2>/dev/null
mkdir -p dreego/routes
cat > dreego/routes/get.dreego << 'DREEGO'
<head><title>T</title></head>
<div><p>check me</p></div>
DREEGO
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
sed -i "s|_ \"gen\"|_ \"t/dreego/gen\"|" main.go
dreego generate
dreego generate --check 2>&1 | grep -q "up-to-date" || { echo "initial check failed"; exit 1; }
sleep 1
touch dreego/routes/get.dreego
if dreego generate --check 2>/dev/null; then echo "expected stale but got up-to-date"; exit 1; fi
echo ok
