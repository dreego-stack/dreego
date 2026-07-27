#!/bin/sh
set -e
cd "$(dirname "$0")"
dreego init .
rm -rf dreego/routes dreego/components
mkdir -p dreego/components dreego/routes
cat > dreego/components/Box.dreego << 'DREEGO'
Component Box ()
<div><div class="box"><p>scoped</p></div></div>
<style>.box{border:1px solid red}</style>
DREEGO
cat > dreego/routes/get.dreego << 'DREEGO'
<head><title>T</title></head>
<div><@Box/><p class="box">unscoped</p></div>
<style>.box{color:blue}</style>
DREEGO
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
sed -i "s|_ \"gen\"|_ \"t/dreego/gen\"|" main.go
dreego generate
go build -o /dev/null .
echo ok
