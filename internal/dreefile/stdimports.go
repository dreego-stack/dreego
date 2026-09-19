package dreefile

import (
	"fmt"
	"sort"
	"strings"
)

var stdlibAllowList = map[string]bool{
	"bytes":           true,
	"context":         true,
	"encoding/base64": true,
	"encoding/hex":    true,
	"encoding/json":   true,
	"errors":          true,
	"fmt":             true,
	"html":            true,
	"io":              true,
	"log":             true,
	"maps":            true,
	"math":            true,
	"net/http":        true,
	"net/url":         true,
	"path":            true,
	"path/filepath":   true,
	"regexp":          true,
	"slices":          true,
	"sort":            true,
	"strconv":         true,
	"strings":         true,
	"sync":            true,
	"time":            true,
	"unicode":         true,
	"unicode/utf8":    true,
}

func allowedStdlibImport(path string) bool {
	return stdlibAllowList[path]
}

func allowedStdlibImports() string {
	paths := make([]string, 0, len(stdlibAllowList))
	for path := range stdlibAllowList {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return strings.Join(paths, ", ")
}

func registerGoImports(gen *Generator, pkg, source string, paths []string) error {
	for _, path := range paths {
		if path == "" {
			return fmt.Errorf("%s: GOIMPORT contains an empty package path", source)
		}
		if !allowedStdlibImport(path) {
			return fmt.Errorf("%s: GOIMPORT %q is not allowed; supported packages: %s", source, path, allowedStdlibImports())
		}
		gen.AddGoImportPath(pkg, path)
	}
	return nil
}

func stdImportsFor(gen *Generator, pkg, src string) string {
	used := map[string]bool{}
	if strings.Contains(src, "strings.") {
		used["strings"] = true
	}
	if strings.Contains(src, "http.") {
		used["net/http"] = true
	}
	if strings.Contains(src, "fmt.") {
		used["fmt"] = true
	}
	for _, path := range gen.GoImportPaths[pkg] {
		used[path] = true
	}
	paths := make([]string, 0, len(used))
	for path := range used {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	lines := make([]string, 0, len(paths))
	for _, path := range paths {
		lines = append(lines, fmt.Sprintf("%q", path))
	}
	return strings.Join(lines, "\n\t")
}
