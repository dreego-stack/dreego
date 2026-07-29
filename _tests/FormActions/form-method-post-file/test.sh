#!/bin/sh
set -e
cd "$(dirname "$0")"
rm -rf dreego/* 2>/dev/null || true
mkdir -p dreego/routes
cat > dreego/routes/post-login.dreego << 'DREEGO'
<go>
    type LoginForm struct {
        Email string
    }
    func Login(c core.Context, form LoginForm) error {
        return nil
    }
</go>
<form g-action="Login" method="post">
    <input name="email">
    <button>Login</button>
</form>
DREEGO
cat > main.go << 'GO'
package main
import (_ "t/dreego/gen"; core "codeberg.org/dreego/dreego/core")
func main() { core.Listen(":0") }
GO
go mod init t >/dev/null 2>&1
go mod edit -replace codeberg.org/dreego/dreego=../../..
go mod edit -require codeberg.org/dreego/dreego@v0.0.0
dreego generate 2>&1
grep -q "HandleIndexGet" dreego/gen/routes.go || { echo "FAIL: GET handler missing"; exit 1; }
grep -q "HandleIndexPost" dreego/gen/routes.go || { echo "FAIL: POST handler missing"; exit 1; }
go build -o /dev/null .
echo ok
