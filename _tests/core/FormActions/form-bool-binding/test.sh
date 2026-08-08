#!/bin/sh
# Using standard: _tests/how-to-test-sh.md
# What: POST binds a checkbox value to a bool struct field
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
require github.com/dreego-stack/dreego/core v0.0.0
replace github.com/dreego-stack/dreego/core => $realrepo/core
EOF

mkdir -p dreego/routes
BT='`'
cat > dreego/routes/get-news.dreego << DREEGO
<go>
    type NewsForm struct {
        Email    string ${BT}validate:"required"${BT}
        Subscribe bool
    }
    func SubmitNews(c dreego.Context, form NewsForm) error {
        if form.Subscribe {
            return c.Redirect("/subscribed", 303)
        }
        return c.Redirect("/skipped", 303)
    }
</go>
<form g-action="SubmitNews" method="post">
    <input name="email" type="email">
    <input name="subscribe" type="checkbox">
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

# Checkbox checked -> Subscribe=true -> /subscribed
LOC=$(curl -s -o /dev/null -w "%{redirect_url}" -H "Content-Type: application/x-www-form-urlencoded" -d "email=a@b.c&subscribe=on" http://localhost:$port/)
case "$LOC" in
    *subscribed*) : ;;
    *) echo "FAIL: expected /subscribed when checkbox checked, got $LOC"; exit 1 ;;
esac

# Checkbox unchecked/absent -> Subscribe=false -> /skipped
LOC=$(curl -s -o /dev/null -w "%{redirect_url}" -H "Content-Type: application/x-www-form-urlencoded" -d "email=a@b.c" http://localhost:$port/)
case "$LOC" in
    *skipped*) : ;;
    *) echo "FAIL: expected /skipped when checkbox absent, got $LOC"; exit 1 ;;
esac

echo ok
