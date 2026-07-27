#!/bin/sh
set -e
cd "$(dirname "$0")"
dreego init .
rm -rf dreego/routes dreego/components dreego/layouts 2>/dev/null
mkdir -p dreego/components dreego/routes
cat > dreego/components/Greet.dreego << 'DREEGO'
Component Greet ()
<div><p>{#slot header}{/slot} {#slot}</p></div>
DREEGO
cat > dreego/routes/get.dreego << 'DREEGO'
<go>name := "World"</go>
<div><@Greet>{#slot header}<strong>{name}</strong>{/slot}!</@Greet></div>
DREEGO
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
sed -i "s|_ \"gen\"|_ \"t/dreego/gen\"|" main.go
dreego generate
go build -o /dev/null .
echo ok
