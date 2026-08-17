package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func generateStaticAssets(routePatterns map[string]bool) (src string, count int, err error) {
	staticDir := filepath.Join("dreego", "static")
	if _, e := os.Stat(staticDir); os.IsNotExist(e) {
		return "", 0, nil
	}

	var buf strings.Builder

	err = filepath.WalkDir(staticDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil || d.IsDir() {
			return nil
		}

		rel, _ := filepath.Rel(staticDir, path)
		urlPath := "/" + filepath.ToSlash(rel)

		methodPattern := "GET" + " " + urlPath
		if routePatterns[methodPattern] {
			return fmt.Errorf("static file %q conflicts with existing route %q", path, urlPath)
		}

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}

		ext := filepath.Ext(path)
		mime := MimeByExt(ext)

		content := []byte(data)
		buf.WriteString(registrationStatement(fmt.Sprintf("app.RegisterStatic(%q, %q, %#v)", urlPath, mime, content)))
		count++
		return nil
	})

	if err != nil {
		return "", 0, err
	}

	return buf.String(), count, nil
}
