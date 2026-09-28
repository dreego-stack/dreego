# Dependencies

Generated with `go1.27.1` from the `go.work` workspace.
For every module the standard-library, external and workspace-internal
packages **actually imported** are listed (imports, test imports and
external test imports; deduplicated and sorted).

Command per module: `go list -f '{{range .Imports}}{{.}}\n{{end}}{{range .TestImports}}{{.}}\n{{end}}{{range .XTestImports}}{{.}}\n{{end}}' ./...`

Classification:
- **Standard library**: first path segment contains no dot.
- **Workspace-internal modules**: package belongs to a module listed in `go.work`.
- **External dependencies**: first path segment contains a dot, excluding workspace-internal packages.

## github.com/dreego-stack/dreego

### Standard library

archive/tar
bufio
bytes
compress/gzip
context
crypto/aes
crypto/cipher
crypto/hmac
crypto/rand
crypto/sha256
crypto/subtle
crypto/tls
encoding/base64
encoding/hex
encoding/json
errors
flag
fmt
go/ast
go/parser
go/scanner
go/token
html
io
log/slog
maps
net
net/http
net/http/httptest
net/url
os
os/exec
path
path/filepath
reflect
regexp
runtime
runtime/debug
slices
sort
strconv
strings
sync
testing
time
unicode
unicode/utf8

### External dependencies

golang.org/x/mod/modfile
golang.org/x/text/language

### Internal/workspace modules

github.com/dreego-stack/dreego/internal/dreefile/codegen
github.com/dreego-stack/dreego/internal/dreefile/dreecode
github.com/dreego-stack/dreego/internal/dreefile/gogen
github.com/dreego-stack/dreego/internal/dreefile/i18n
github.com/dreego-stack/dreego/internal/dreefile/ir
github.com/dreego-stack/dreego/internal/dreefile/jsoutput
github.com/dreego-stack/dreego/internal/dreefile/lexer
github.com/dreego-stack/dreego/internal/dreefile/parser
github.com/dreego-stack/dreego/internal/dreefile/sections/body/html
github.com/dreego-stack/dreego/internal/dreefile/sections/body/md
github.com/dreego-stack/dreego/internal/dreefile/sections/client
github.com/dreego-stack/dreego/internal/dreefile/sections/client/js
github.com/dreego-stack/dreego/internal/dreefile/sections/client/lua
github.com/dreego-stack/dreego/internal/dreefile/sections/client/ts
github.com/dreego-stack/dreego/internal/dreefile/sections/head
github.com/dreego-stack/dreego/internal/dreefile/sections/style
github.com/dreego-stack/dreego/internal/dreefile/tokens
github.com/dreego-stack/dreego/internal/gomod
github.com/dreego-stack/dreego/internal/md
github.com/dreego-stack/dreego/internal/session

## github.com/dreego-stack/dreego/core

### Standard library

bytes
context
crypto/tls
encoding/json
encoding/xml
errors
fmt
html
io
log/slog
maps
math
mime
net
net/http
net/http/httptest
net/url
os
path
reflect
regexp
slices
strconv
strings
sync
sync/atomic
testing
time
unicode

### External dependencies

golang.org/x/text/currency
golang.org/x/text/feature/plural
golang.org/x/text/language
golang.org/x/text/message

### Internal/workspace modules

github.com/dreego-stack/dreego/core/internal/context
github.com/dreego-stack/dreego/core/internal/i18n
github.com/dreego-stack/dreego/core/internal/render
github.com/dreego-stack/dreego/core/internal/server
github.com/dreego-stack/dreego/internal/md
github.com/dreego-stack/dreego/internal/middleware
github.com/dreego-stack/dreego/internal/session
github.com/dreego-stack/dreego/internal/validate

## github.com/dreego-stack/dreego/cmd/dreego

### Standard library

bytes
embed
encoding/json
errors
fmt
io
io/fs
maps
os
os/exec
os/signal
path
path/filepath
reflect
regexp
runtime
runtime/debug
slices
sort
strings
syscall
testing
testing/fstest
time

### External dependencies

_none_

### Internal/workspace modules

github.com/dreego-stack/dreego/cmd/dreego/internal/templates
github.com/dreego-stack/dreego/internal/dreefile
github.com/dreego-stack/dreego/internal/dreefile/sections/client/ts
github.com/dreego-stack/dreego/internal/gomod

## github.com/dreego-stack/dreego/adapter/ssr

### Standard library

context
errors
io
net
net/http
os
os/signal
reflect
runtime
sync
syscall
testing
time

### External dependencies

_none_

### Internal/workspace modules

github.com/dreego-stack/dreego/core
github.com/dreego-stack/dreego/internal/middleware
github.com/dreego-stack/dreego/internal/session
github.com/dreego-stack/dreego/internal/validate

## github.com/dreego-stack/dreego/adapter/wails

### Standard library

errors
log/slog
net/http
net/http/httptest
path
strings
testing

### External dependencies

_none_

### Internal/workspace modules

github.com/dreego-stack/dreego/core

## github.com/dreego-stack/dreego/dreegotest

### Standard library

crypto/sha256
encoding/hex
fmt
io
io/fs
net
net/http
net/http/httptest
net/url
os
os/exec
path/filepath
reflect
strings
sync
testing
time

### External dependencies

_none_

### Internal/workspace modules

github.com/dreego-stack/dreego/core
github.com/dreego-stack/dreego/dreegotest
github.com/dreego-stack/dreego/internal/dreefile

## github.com/dreego-stack/dreego/_tests/go

### Standard library

bufio
bytes
encoding/json
errors
fmt
io
io/fs
net
net/http
net/http/httptest
net/url
os
os/exec
path/filepath
regexp
sort
strconv
strings
sync
syscall
testing
time
unicode

### External dependencies

github.com/wailsapp/wails/v3/pkg/application

### Internal/workspace modules

github.com/dreego-stack/dreego/adapter/ssr
github.com/dreego-stack/dreego/adapter/wails
github.com/dreego-stack/dreego/core
github.com/dreego-stack/dreego/dreegotest

## demo

### Standard library

encoding/json
errors
fmt
html
log
net/http
net/http/httptest
net/url
os
strings
sync
testing

### External dependencies

_none_

### Internal/workspace modules

demo/blog
demo/blog/components
demo/blog/layouts
demo/blog/routes
demo/lua
demo/lua/routes
demo/saas
demo/saas/routes
demo/www
demo/www/layouts
demo/www/routes
github.com/dreego-stack/dreego/adapter/ssr
github.com/dreego-stack/dreego/core

## demo-wailsv3

### Standard library

context
log
net/http
strings
sync
testing
time

### External dependencies

github.com/wailsapp/wails/v3/pkg/application

### Internal/workspace modules

demo-wailsv3/app
demo-wailsv3/app/routes
github.com/dreego-stack/dreego/adapter/wails
github.com/dreego-stack/dreego/core
