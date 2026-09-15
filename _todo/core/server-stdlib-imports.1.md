# Stdlib imports for `<server>` and component code

- Area: compiler
- Phase: code generation
- Goal: let `<server>`/component code use the standard library it names, with the import emitted correctly, while unknown or dynamic imports stay impossible so type safety is preserved.
- Gap: `stdImportsFor` (`internal/transpiler/generate.go:220-231`) recognises only `strings`, `net/http` and `fmt`. Code that needs another stdlib package — e.g. `sync.Mutex` for shared state in a route handler — generates `undefined: sync`, and a route-header `import "sync"` line is ignored.
- Gap: the current workaround is a handwritten sibling `package routes` Go file (precedent: `demo/demo-ssr/blog/routes/news_store.go`, and now `web-app/www/routes/notes_store.go`).
- Acceptance: a `.dreego` server section using a small allow-listed stdlib package compiles; the generated import block is correct; a regression test covers it; the handwritten-sibling workaround is no longer required for the common cases; the allow-list is explicit (no arbitrary import injection).
- Depends on: nothing.
