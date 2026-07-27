#!/bin/sh
set -e
cd "$(dirname "$0")"
dreego init .
rm -rf dreego/routes dreego/components
mkdir -p dreego/components dreego/routes
cat > dreego/components/Card.dreego << 'DREEGO'
Component Card (title string)
<div><article><h2>{title}</h2></article></div>
DREEGO
cat > dreego/routes/get.dreego << 'DREEGO'
import "dreego/components/Card"
<div><@Card title="Imported!"/></div>
DREEGO
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
sed -i "s|_ \"gen\"|_ \"t/dreego/gen\"|" main.go
dreego generate
go build -o /dev/null .
echo ok
