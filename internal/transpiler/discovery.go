package transpiler

import (
	"os"
	"path/filepath"
	"strings"
)

const configFileName = "dreego.config.json"

func isSkippedDir(name string) bool {
	switch name {
	case "vendor", "node_modules", ".git", ".worktrees", ".tmp":
		return true
	}
	return strings.HasPrefix(name, ".") && name != "." && name != ".."
}

func isWebsiteRoot(path string) bool {
	info, err := os.Stat(filepath.Join(path, configFileName))
	return err == nil && !info.IsDir()
}

func findWebsiteRoots() ([]string, error) {
	var roots []string
	err := filepath.WalkDir(".", func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !d.IsDir() {
			return nil
		}
		base := filepath.Base(path)
		if base == "." {
			return nil
		}
		if strings.HasPrefix(base, ".") {
			return filepath.SkipDir
		}
		if isSkippedDir(base) {
			return filepath.SkipDir
		}
		if isWebsiteRoot(path) {
			roots = append(roots, path)
			return filepath.SkipDir
		}
		return nil
	})
	return roots, err
}

func isRoutesDir(root, path string) bool {
	rel := relToRoot(root, path)
	return rel == "routes" || strings.HasPrefix(rel, "routes/")
}

func isComponentsDir(root, path string) bool {
	rel := relToRoot(root, path)
	return rel == "components" || strings.HasPrefix(rel, "components/")
}

func isLayoutsDir(root, path string) bool {
	rel := relToRoot(root, path)
	return rel == "layouts" || strings.HasPrefix(rel, "layouts/") || strings.HasSuffix(rel, "/layouts")
}

func isStaticDir(root, path string) bool {
	rel := relToRoot(root, path)
	return rel == "static" || strings.HasPrefix(rel, "static/")
}

func relToRoot(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return ""
	}
	return filepath.ToSlash(rel)
}

func routeDirRel(root, path string) string {
	rel := relToRoot(root, path)
	if rel == "routes" {
		return ""
	}
	return strings.TrimPrefix(rel, "routes/")
}

func sanitizePkgName(name string) string {
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
		return "app"
	}
	if s[0] >= '0' && s[0] <= '9' {
		s = "pkg" + s
	}
	return s
}
