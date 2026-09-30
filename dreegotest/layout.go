package dreegotest

import (
	"sort"
	"strings"
)

// TestAppDir is the app subdirectory the test helpers use when a test writes
// the historical `<root>/routes/…` shape. Test files may also name the app dir
// explicitly (`<root>/<app>/routes/…`).
const TestAppDir = "app"

// websiteRootOf returns the website root implied by a test file path: for a
// routes file it is the parent of the app directory; for a shared
// components/layouts tree it is the containing directory (unless the path sits
// under a routes/ tree, in which case it is route-local and has no root).
func websiteRootOf(path string) (string, bool) {
	parts := strings.Split(path, "/")
	seenRoutes := false
	for i, p := range parts {
		switch p {
		case "routes":
			if i >= 2 {
				return strings.Join(parts[:i-1], "/"), true
			}
			seenRoutes = true
		case "components", "layouts":
			if i >= 1 && !seenRoutes {
				return strings.Join(parts[:i], "/"), true
			}
		}
	}
	return "", false
}

// appPackagePath returns the module-relative import path of the app package for
// the files ("www/app" for a root named www with an app dir named app). It
// prefers a path that contains routes/ and falls back to any website root. When
// several apps exist, the lexicographically smallest route-bearing app wins so
// the choice is deterministic.
func appPackagePath(files map[string]string) string {
	paths := make([]string, 0, len(files))
	fallback := ""
	for path := range files {
		if _, ok := websiteRootOf(path); !ok {
			continue
		}
		parts := strings.Split(path, "/")
		for i, p := range parts {
			if p == "routes" && i >= 2 {
				paths = append(paths, strings.Join(parts[:i], "/"))
				break
			}
		}
		if fallback == "" {
			root, _ := websiteRootOf(path)
			if root == "" {
				fallback = TestAppDir
			} else {
				fallback = root + "/" + TestAppDir
			}
		}
	}
	if len(paths) > 0 {
		sort.Strings(paths)
		return paths[0]
	}
	return fallback
}

// appPackageName returns the Go package identifier generated for the app
// directory: the sanitized base name of the app path.
func appPackageName(files map[string]string) string {
	pkgPath := appPackagePath(files)
	base := pkgPath
	if i := strings.LastIndex(pkgPath, "/"); i >= 0 {
		base = pkgPath[i+1:]
	}
	return sanitizeTestPkg(base)
}

func sanitizeTestPkg(name string) string {
	var b strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-', r == '.', r == ' ':
			b.WriteRune('_')
		}
	}
	s := b.String()
	if s == "" {
		return TestAppDir
	}
	if s[0] >= '0' && s[0] <= '9' {
		s = "pkg" + s
	}
	return s
}
