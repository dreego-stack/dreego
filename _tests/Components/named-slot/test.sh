#!/bin/sh
set -e
cd "$(dirname "$0")"
dreego init .
rm -rf dreego/routes dreego/components dreego/layouts 2>/dev/null
mkdir -p dreego/components dreego/routes

printf 'Component Card (title string)\n<div><article>{#slot header}{/slot}<h2>{title}</h2><div>{#slot}</div></article></div>\n' > dreego/components/Card.dreego
printf '<div><@Card title="Hi">{#slot header}<strong>HEADER</strong>{/slot}<p>body</p></@Card></div>\n' > dreego/routes/get.dreego

go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
sed -i "s|_ \"gen\"|_ \"t/dreego/gen\"|" main.go
dreego generate
go build -o /dev/null .
echo ok
