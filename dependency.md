# Dependencies

Generated with `go1.27.1` from the `go.work` workspace.

Per module, only the packages that leave this repo are listed:
the **standard library** and the **external dependencies**.
All imports internal to this repo (the workspace modules and their
packages) are intentionally omitted.

Command per module: `go list -f '{{range .Imports}}{{.}}\n{{end}}{{range .TestImports}}{{.}}\n{{end}}{{range .XTestImports}}{{.}}\n{{end}}' ./...`

Classification:
- **Standard library**: first path segment contains no dot and is not part of this repo.
- **External dependencies**: first path segment contains a dot and is not part of this repo.

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
