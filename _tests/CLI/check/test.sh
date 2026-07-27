#!/bin/sh
set -e
cd "$(dirname "$0")"
dreego init .
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
sed -i "s|_ \"gen\"|_ \"t/dreego/gen\"|" main.go
dreego generate
dreego generate --check 2>&1 | grep -q "up-to-date" || { echo "check failed"; exit 1; }
echo ok
