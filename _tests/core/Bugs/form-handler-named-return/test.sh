#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that findFormStruct works with named return values (B19)
set -e

realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"

cat > go.mod << EOF
module t
go 1.22
require codeberg.org/dreego/dreego v0.0.0
replace codeberg.org/dreego/dreego => $realrepo
EOF

cat > main.go << 'GO'
package main
import _ "t/dreego/gen"
func main() {}
GO

mkdir -p dreego/routes

cat > dreego/routes/get.dreego << 'DREEGO'
<go>
type LoginForm struct {
    Email string `form:"email" validate:"required,email"`
}

func Save(c *core.SSRContext, form LoginForm) (err error) {
    return nil
}
</go>
<div>
  <form g-action="Save" method="post">
    <input name="email">
    <button>Save</button>
  </form>
</div>
DREEGO

go run $realrepo/cmd/dreego generate
generated="dreego/gen/routes.go"
if grep -A 20 "func HandleIndexPost" "$generated" | grep -q "var form LoginForm"; then
    :
else
    echo "FAIL: form struct not detected with named return (B19)"
    exit 1
fi
go build -o /dev/null .
echo ok
