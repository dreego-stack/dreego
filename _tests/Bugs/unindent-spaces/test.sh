#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that go declarations with space indentation are unindented correctly (B18)
set -e

realrepo="$(cd "$(dirname "$0")"/../../.. && pwd)"
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
    type UserForm struct {
        Email string `form:"email" validate:"required,email"`
    }

    func Save(c *core.SSRContext, form UserForm) error {
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
if grep -q "^    type UserForm struct" "$generated"; then
    echo "FAIL: space-indented declaration not unindented (B18)"
    exit 1
fi
if ! grep -q "^type UserForm struct" "$generated"; then
    echo "FAIL: declaration not emitted at all"
    exit 1
fi
go build -o /dev/null .
echo ok
