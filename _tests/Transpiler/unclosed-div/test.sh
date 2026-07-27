#!/bin/sh
cd "$(dirname "$0")"
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
if dreego generate 2>/dev/null; then echo "expected failure but succeeded"; exit 1; fi
echo ok
