#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: POST binds an integer field and validates min on the bound value
set -e
realrepo="$(cd "$(dirname "$0")"/../../../.. && pwd)"
workdir="$(mktemp -d)"
trap "rm -rf $workdir" EXIT

cd "$workdir"
if [ -z "$DREEGO_BIN" ]; then
    DREEGO_BIN="$workdir/.dreego-bin"
    (cd "$realrepo" && go build -o "$DREEGO_BIN" ./cmd/dreego) || { echo "FAIL: could not build dreego CLI"; exit 1; }
    export DREEGO_BIN
fi
apk add --no-cache curl >/dev/null 2>&1 || true

port="${DREEGO_PORT:-$(( ( $(od -An -N2 -i /dev/urandom | tr -d ' ') % 50000 ) + 10000 ))}"

cat > go.mod << EOF
module t
go 1.22
require github.com/dreego-stack/dreego v0.0.0
replace github.com/dreego-stack/dreego => $realrepo
EOF

mkdir -p dreego/routes
BT='`'
cat > dreego/routes/get-age.dreego << DREEGO
<go>
    type AgeForm struct {
        Age int ${BT}validate:"min=2"${BT}
    }
    func SubmitAge(c dreego.Context, form AgeForm) error {
        if form.Age == 20 {
            return c.Redirect("/adult", 303)
        }
        return c.Redirect("/other", 303)
    }
</go>
<form g-action="SubmitAge" method="post">
    <input name="age" type="number">
    <button type="submit">Send</button>
</form>
DREEGO
cat > main.go << GO
package main
import (_ "t/dreego/gen"; dreego "github.com/dreego-stack/dreego/core")
func main() { dreego.SetCSRF(false); dreego.Listen(":$port") }
GO
$DREEGO_BIN generate 2>&1
go build -o $workdir/srv .
$workdir/srv &
PID=$!
trap "kill $PID 2>/dev/null; rm -rf $workdir" EXIT
for i in $(seq 1 30); do curl -s -o /dev/null http://localhost:$port/ && break; sleep 0.1; done

# Valid bound int (min passes) -> redirect
CODE=$(curl -s -o /dev/null -w "%{http_code}" -H "Content-Type: application/x-www-form-urlencoded" -d "age=20" http://localhost:$port/)
[ "$CODE" = "303" ] || { echo "FAIL: expected 303 redirect for valid age, got $CODE"; exit 1; }

# Bound int used in handler decision -> /adult for age>=18
LOC=$(curl -s -o /dev/null -w "%{redirect_url}" -H "Content-Type: application/x-www-form-urlencoded" -d "age=20" http://localhost:$port/)
case "$LOC" in
    *adult*) : ;;
    *) echo "FAIL: expected redirect to /adult for age=20, got $LOC"; exit 1 ;;
esac

# min validation fails on bound int (length 1 / numeric < 2) -> re-render
CODE=$(curl -s -o /dev/null -w "%{http_code}" -H "Content-Type: application/x-www-form-urlencoded" -d "age=1" http://localhost:$port/)
[ "$CODE" = "200" ] || { echo "FAIL: expected 200 re-render when min fails, got $CODE"; exit 1; }

echo ok
