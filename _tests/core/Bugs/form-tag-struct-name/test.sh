#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: Tests that hasValidateTag/hasFormTag only match the target struct (B4)
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

cat > dreego/routes/get-search.dreego << 'DREEGO'
<go>
type LoginForm struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=8"`
}

type SearchQuery struct {
	Query string `form:"q"`
}

func Search(c *core.SSRContext, form SearchQuery) error {
	return nil
}
</go>
<div>
  <form g-action="Search" method="post">
    <input name="q">
    <button>Search</button>
  </form>
</div>
DREEGO

go run $realrepo/cmd/dreego generate

generated="dreego/gen/routes.go"
if grep -A 30 "func HandleIndexPost" "$generated" | grep -q "ValidateForm"; then
    echo "FAIL: Search handler calls ValidateForm even though SearchQuery has no validate tags (B4)"
    exit 1
fi

go build -o /dev/null .
echo ok
