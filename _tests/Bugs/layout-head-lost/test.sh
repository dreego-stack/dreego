#!/bin/sh
set -e
cd "$(dirname "$0")"
rm -rf dreego/* 2>/dev/null || true
mkdir -p dreego/layouts dreego/routes
cat > dreego/layouts/default.dreego << 'DREEGO'
<head>
    <title>Layout Title</title>
    <script src="https://cdn.tailwindcss.com"></script>
    {#head}
</head>
<div><main>{#slot}</main></div>
DREEGO
cat > dreego/routes/get.dreego << 'DREEGO'
<head><meta name="desc" content="route meta"></head>
<div><h1>Page</h1></div>
DREEGO
mkdir -p dreego/gen
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.Listen(":0") }
GO
dreego generate 2>&1
grep -q "cdn.tailwindcss.com" dreego/gen/routes.go || { echo "FAIL: layout head dropped — tailwind CDN not in generated code"; exit 1; }
grep -q "Layout Title" dreego/gen/routes.go || { echo "FAIL: layout head dropped — title not in generated code"; exit 1; }
grep -q "route meta" dreego/gen/routes.go || { echo "FAIL: route head dropped — meta not in generated code"; exit 1; }
go build -o /dev/null .
echo ok
