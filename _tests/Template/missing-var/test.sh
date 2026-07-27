#!/bin/sh
set -e
cd "$(dirname "$0")"
dreego init .
rm -rf dreego/routes dreego/components dreego/layouts 2>/dev/null
mkdir -p dreego/routes
cat > dreego/routes/get.dreego << 'DREEGO'
<div><p>{undefined}</p></div>
DREEGO
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
sed -i "s|_ \"gen\"|_ \"t/dreego/gen\"|" main.go
dreego generate
if go build -o /dev/null . 2>/dev/null; then echo "expected build failure but succeeded"; exit 1; fi
echo ok
