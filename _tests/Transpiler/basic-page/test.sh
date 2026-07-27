#!/bin/sh
set -e
cd "$(dirname "$0")"

dreego init .

go mod init basic-test >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
sed -i 's|_ "gen"|_ "basic-test/dreego/gen"|' main.go

dreego generate

go build -o /dev/null .

echo "ok"
