#!/bin/sh
set -e
cd "$(dirname "$0")"
dreego init .
rm -rf dreego/routes dreego/components dreego/layouts go.mod go.sum 2>/dev/null
mkdir -p dreego/routes/about dreego/routes/users/about
cat > "dreego/routes/about/get.dreego" << 'DREEGO'
<div><p>about page</p></div>
DREEGO
cat > "dreego/routes/users/about/get.dreego" << 'DREEGO'
<div><p>users about page</p></div>
DREEGO
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
sed -i "s|_ \"gen\"|_ \"t/dreego/gen\"|" main.go
dreego generate
go build -o /dev/null .
echo ok
