#!/bin/sh
set -e
cd "$(dirname "$0")"
dreego init .
rm -rf dreego/routes dreego/components dreego/layouts go.mod go.sum 2>/dev/null
mkdir -p dreego/routes dreego/static/images
cat > "dreego/routes/get.dreego" << 'DREEGO'
<div><p>hello</p></div>
DREEGO
printf '<svg></svg>' > dreego/static/images/logo.svg
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
sed -i "s|_ \"gen\"|_ \"t/dreego/gen\"|" main.go
dreego generate
grep -q 'RegisterStatic("/images/logo.svg"' dreego/gen/routes.go
go build -o /dev/null .
echo ok
