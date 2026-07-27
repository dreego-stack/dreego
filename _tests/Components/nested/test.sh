#!/bin/sh
set -e
cd "$(dirname "$0")"
dreego init .
rm -rf dreego/routes dreego/components dreego/layouts 2>/dev/null
mkdir -p dreego/components dreego/routes
cat > dreego/components/Inner.dreego << 'DREEGO'
Component Inner ()
<div><span>inner</span></div>
DREEGO
cat > dreego/components/Outer.dreego << 'DREEGO'
Component Outer ()
<div><article><@Inner/></article></div>
DREEGO
cat > dreego/routes/get.dreego << 'DREEGO'
<div><@Outer/></div>
DREEGO
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
sed -i "s|_ \"gen\"|_ \"t/dreego/gen\"|" main.go
dreego generate
go build -o /dev/null .
echo ok
