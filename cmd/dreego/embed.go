package main

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed all:embedded
var embeddedFS embed.FS

// embeddedDocs is the embedded docs filesystem. Its root mirrors the repo
// layout: _docs/, README.md, CHANGELOG.md. Overridable in tests.
var embeddedDocs fs.FS = mustSubFS(embeddedFS, "embedded")

func mustSubFS(fsys fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		panic(err)
	}
	return sub
}

// fetchDocEmbedded reads a docs file from embeddedDocs. The leading slash of
// the URL path is mapped onto the embedded FS root.
func fetchDocEmbedded(path string) ([]byte, error) {
	rel := strings.TrimPrefix(path, "/")
	return fs.ReadFile(embeddedDocs, rel)
}
