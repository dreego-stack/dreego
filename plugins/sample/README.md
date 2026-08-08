# dreego/plugins/sample

Example plugin showing the monorepo layout. It imports `github.com/dreego-stack/dreego/core` via a `replace` directive in `plugins/sample/go.mod`, and `go.work` links the modules for local development.

Core must never import this package or any other plugin package.
