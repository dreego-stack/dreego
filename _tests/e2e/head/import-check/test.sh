#!/bin/sh
# What: the Tailwind CDN script declared in a page <head> is served in the
# rendered HTML (black-box counterpart to the head-merging unit tests).
set -e

REPO_DIR="${REPO_DIR:-$(cd "$(dirname "$0")/../../../.." && pwd)}"
DREEGO_BIN="${DREEGO_BIN:-}"

WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

mkdir -p "$WORK/www/routes"
cat > "$WORK/go.mod" <<EOF
module importcheck

go 1.27

require (
	github.com/dreego-stack/dreego/adapter/ssr v0.0.0
	github.com/dreego-stack/dreego/core v0.0.0
)

replace github.com/dreego-stack/dreego => $REPO_DIR
replace github.com/dreego-stack/dreego/core => $REPO_DIR/core
replace github.com/dreego-stack/dreego/adapter/ssr => $REPO_DIR/adapter/ssr
EOF
cp "$REPO_DIR/go.sum" "$WORK/go.sum"

cat > "$WORK/main.go" <<'EOF'
package main

import (
	"log"
	"os"

	"importcheck/www"
	"github.com/dreego-stack/dreego/adapter/ssr"
	dreego "github.com/dreego-stack/dreego/core"
)

func main() {
	app := dreego.New()
	if err := www.Register(app); err != nil {
		log.Fatal(err)
	}
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}
	if err := ssr.Listen(app, addr); err != nil {
		log.Fatal(err)
	}
}
EOF

cat > "$WORK/www/dreego.config.json" <<'EOF'
{
    "logging": {"enabled": false},
    "redirects": [],
    "rewrites": []
}
EOF

cat > "$WORK/www/routes/+page.dreego" <<'EOF'
<head>
    <title>Import check</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>

<body>
    <h1>Import check</h1>
    <p>Verifies the Tailwind CDN script is served in the rendered HTML.</p>
</body>
EOF

if [ -z "$DREEGO_BIN" ]; then
    DREEGO_BIN="$(mktemp -d)/dreego"
    (cd "$REPO_DIR" && go build -o "$DREEGO_BIN" ./cmd/dreego)
fi

(cd "$WORK" && GOWORK=off "$DREEGO_BIN" generate)
(cd "$WORK" && GOWORK=off go build -mod=mod -o server .)

PORT="$(shuf -i 20000-29999 -n 1 2>/dev/null || echo 20000)"
(cd "$WORK" && PORT="$PORT" ./server) &
SERVER_PID=$!
trap 'kill "$SERVER_PID" 2>/dev/null; rm -rf "$WORK"' EXIT

for _ in $(seq 1 50); do
    curl -sf "http://127.0.0.1:$PORT" >/dev/null 2>&1 && break
    sleep 0.1
done

BODY="$(curl -sf "http://127.0.0.1:$PORT")" || {
    echo "FAIL: server did not respond on :$PORT"
    exit 1
}

case "$BODY" in
    *'<script src="https://cdn.tailwindcss.com"></script>'*)
        echo ok
        ;;
    *)
        echo "FAIL: Tailwind CDN script missing from rendered HTML"
        exit 1
        ;;
esac
