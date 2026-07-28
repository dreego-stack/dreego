#!/bin/sh
set -e
cd "$(dirname "$0")"
dreego init .
rm -rf dreego/routes dreego/components dreego/layouts go.mod go.sum 2>/dev/null
mkdir -p dreego/layouts dreego/routes
cat > "dreego/layouts/default.dreego" << 'DREEGO'
<head>
    <title>Layout Title</title>
</head>
<div>{#slot}</div>
DREEGO
cat > "dreego/routes/get.dreego" << 'DREEGO'
<head><script src="route-script.js"></script></head>
<div><p>hello</p></div>
DREEGO
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
sed -i "s|_ \"gen\"|_ \"t/dreego/gen\"|" main.go
dreego generate 2>&1
grep -q "Layout Title" dreego/gen/routes.go || { echo "FAIL: layout head dropped — title not in generated code"; exit 1; }
grep -q "route-script.js" dreego/gen/routes.go || { echo "FAIL: route head dropped — script not in generated code"; exit 1; }
go build -o /dev/null .
echo ok
