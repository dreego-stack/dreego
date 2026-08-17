package core

import (
	"path/filepath"
	"strings"
)

func isSkippedDir(name string) bool {
	switch name {
	case "vendor", "node_modules", ".git", ".worktrees", ".tmp":
		return true
	}
	return strings.HasPrefix(name, ".") && name != "." && name != ".."
}

func isInDreegoRoot(path string) bool {
	rel := filepath.ToSlash(strings.TrimPrefix(path, "./"))
	if rel == "dreego" || strings.HasPrefix(rel, "dreego/") {
		return true
	}
	return false
}

func isGeneratedDir(path string) bool {
	rel := filepath.ToSlash(strings.TrimPrefix(path, "./"))
	return rel == "dreego/gen" || strings.HasPrefix(rel, "dreego/gen/")
}

func isDreegoRoutesDir(path string) bool {
	rel := filepath.ToSlash(strings.TrimPrefix(path, "./"))
	if !isInDreegoRoot(rel) {
		return false
	}
	parts := strings.Split(rel, "/")
	for i, p := range parts {
		if p == "dreego" && i+1 < len(parts) && parts[i+1] == "routes" {
			rest := strings.Join(parts[i+2:], "/")
			if rest == "" {
				return true
			}
			for _, seg := range strings.Split(rest, "/") {
				if seg == "components" || seg == "layouts" {
					return false
				}
			}
			return true
		}
	}
	return false
}

func isDreegoComponentsDir(path string) bool {
	rel := filepath.ToSlash(strings.TrimPrefix(path, "./"))
	if !isInDreegoRoot(rel) {
		return false
	}
	parts := strings.Split(rel, "/")
	for i, p := range parts {
		if p == "dreego" && i+1 < len(parts) && parts[i+1] == "components" {
			rest := strings.Join(parts[i+2:], "/")
			if rest == "" {
				return true
			}
			for _, seg := range strings.Split(rest, "/") {
				if seg == "routes" || seg == "layouts" || seg == "static" || seg == "gen" {
					return false
				}
			}
			return true
		}
	}
	return false
}

func isDreegoLayoutsDir(path string) bool {
	rel := filepath.ToSlash(strings.TrimPrefix(path, "./"))
	if !isInDreegoRoot(rel) {
		return false
	}
	base := filepath.Base(rel)
	if base != "layouts" {
		return false
	}
	parts := strings.Split(rel, "/")
	for i, p := range parts {
		if p == "dreego" && i+1 < len(parts) {
			rest := strings.Join(parts[i+1:], "/")
			if rest == "layouts" || strings.Contains(rest, "/layouts") {
				return true
			}
		}
	}
	return false
}
