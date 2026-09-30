package dreefile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type staticFile struct {
	rel  string
	path string
}

// collectStaticFiles merges an app's static/ directory with the website root's
// shared static/. The returned relative paths are slash-separated; a local file
// shadows a global file with the same relative path.
func collectStaticFiles(appRoot, websiteRoot string) ([]staticFile, error) {
	sources := []string{filepath.Join(appRoot, "static")}
	if websiteRoot != "" {
		sources = append(sources, filepath.Join(websiteRoot, "static"))
	}
	seen := map[string]bool{}
	var out []staticFile
	for _, src := range sources {
		if !hasDir(src) {
			continue
		}
		var found []staticFile
		err := filepath.WalkDir(src, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return fmt.Errorf("error walking %s: %w", path, walkErr)
			}
			if d.IsDir() {
				return nil
			}
			rel, err := filepath.Rel(src, path)
			if err != nil {
				return err
			}
			found = append(found, staticFile{rel: filepath.ToSlash(rel), path: path})
			return nil
		})
		if err != nil {
			return nil, err
		}
		sort.Slice(found, func(i, j int) bool { return found[i].rel < found[j].rel })
		for _, f := range found {
			if seen[f.rel] {
				continue
			}
			seen[f.rel] = true
			out = append(out, f)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].rel < out[j].rel })
	return out, nil
}

// generateStaticAssets registers every static file of an app. The app's static/
// overrides the shared website-root static/ file by file.
func generateStaticAssets(appRoot, websiteRoot string, routePatterns map[string]bool) (src string, count int, err error) {
	files, err := collectStaticFiles(appRoot, websiteRoot)
	if err != nil {
		return "", 0, err
	}

	var buf strings.Builder
	for _, f := range files {
		urlPath := "/" + f.rel
		if routePatterns["GET "+urlPath] {
			return "", 0, fmt.Errorf("static file %q conflicts with existing route %q", f.path, urlPath)
		}

		data, readErr := os.ReadFile(f.path)
		if readErr != nil {
			return "", 0, readErr
		}

		mime := MimeByExt(filepath.Ext(f.path))
		buf.WriteString(registrationStatement(fmt.Sprintf("app.RegisterStatic(%q, %q, %#v)", urlPath, mime, []byte(data))))
		count++
	}

	return buf.String(), count, nil
}
