#!/bin/sh
set -e
cd "$(dirname "$0")"
dreego init .
rm -rf dreego/routes dreego/components dreego/layouts go.mod go.sum 2>/dev/null
mkdir -p dreego/routes dreego/components/button
printf 'Component FlatButton (label string)\n\n<div><button>{label}</button></div>' > dreego/components/button.dreego
printf 'Component NestedButton (label string)\n\n<div><button class="nested">{label}</button></div>' > dreego/components/button/button.dreego
cat > "dreego/routes/get.dreego" << 'DREEGO'
<div><@FlatButton label="Click"/><@NestedButton label="Go"/></div>
DREEGO
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
sed -i "s|_ \"gen\"|_ \"t/dreego/gen\"|" main.go
dreego generate
go build -o /dev/null .
echo ok
