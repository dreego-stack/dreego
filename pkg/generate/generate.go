package generate

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"codeberg.org/dreego/dreego/pkg/ast"
	"codeberg.org/dreego/dreego/pkg/codegen"
	"codeberg.org/dreego/dreego/pkg/lexer"
	"codeberg.org/dreego/dreego/pkg/parser"
)

var binaryHash string

var methodExt = map[string]string{
	"get":    "GET",
	"post":   "POST",
	"put":    "PUT",
	"delete": "DELETE",
}

func init() {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	f, err := os.Open(exe)
	if err != nil {
		return
	}
	defer f.Close()
	h := sha256.New()
	io.Copy(h, f)
	binaryHash = hex.EncodeToString(h.Sum(nil))
}

func Run(force bool) error {
	start := time.Now()
	var found, generated int

	layout := findLayout()

	err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
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
		if isLayoutDir(path) {
			return nil
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			return nil
		}

		var dreegoFiles []string
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".dreego") {
				dreegoFiles = append(dreegoFiles, filepath.Join(path, e.Name()))
			}
		}
		if len(dreegoFiles) == 0 {
			return nil
		}
		if !isRouteDir(path) {
			return nil
		}

		found++
		outPath := filepath.Join(path, "dree.go")
		pattern := buildPattern(path)
		pageName := buildPageName(path)

		var hashParts []string
		var handlerSources []string
		hashParts = append(hashParts, fmt.Sprintf("bin:%q", binaryHash[:12]))

		for _, fpath := range dreegoFiles {
			data, _ := os.ReadFile(fpath)
			h := sha256.Sum256(data)
			base := strings.TrimSuffix(filepath.Base(fpath), ".dreego")
			hashParts = append(hashParts, fmt.Sprintf("%s:%q", base, hex.EncodeToString(h[:])))

			method := "GET"
			for prefix, m := range methodExt {
				if base == prefix || strings.HasPrefix(base, prefix+"-") {
					method = m
					break
				}
			}

			tokens, err := lexer.Lex(string(data))
			if err != nil {
				return fmt.Errorf("error lexing %s: %w", fpath, err)
			}

			p := parser.NewParser(tokens)
			file, err := p.Parse()
			if err != nil {
				return fmt.Errorf("error parsing %s: %w", fpath, err)
			}

			if len(file.Go) == 0 {
				file.Go = []ast.GoSection{{Method: method}}
			}
			for i := range file.Go {
				file.Go[i].Method = method
			}

			scopeHash := hex.EncodeToString(h[:])[:12]
			pkgName := filepath.Base(path)

			for _, g := range file.Go {
				src, err := codegen.GenerateMethodHandler(file, layout, pkgName, pageName, pattern, g, scopeHash)
				if err != nil {
					return fmt.Errorf("error generating %s: %w", fpath, err)
				}
				handlerSources = append(handlerSources, src)
			}
		}

		hashLine := "// hash:{" + strings.Join(hashParts, ", ") + "}"

		pkgName := filepath.Base(path)
		out := fmt.Sprintf("%s\npackage %s\n\nimport (\n\t\"fmt\"\n\t\"net/http\"\n\t\"strings\"\n\n\t\"codeberg.org/dreego/dreego/pkg/context\"\n\t\"codeberg.org/dreego/dreego/pkg/runtime\"\n)\n\n", hashLine, pkgName)
		out += strings.Join(handlerSources, "")

		if err := os.WriteFile(outPath, []byte(out), 0644); err != nil {
			return fmt.Errorf("error writing %s: %w", outPath, err)
		}
		generated++
		return nil
	})
	if err != nil {
		return err
	}

	elapsed := time.Since(start)
	fmt.Printf("Found %d routes\n", found)
	fmt.Printf("Generated %d dree.go files\n", generated)
	fmt.Printf("in %dns\n", elapsed.Nanoseconds())
	return nil
}

func isRouteDir(path string) bool {
	return strings.Contains(path, "routes/") || strings.HasSuffix(path, "/routes") || filepath.Base(path) == "routes"
}

func buildPageName(path string) string {
	rel := filepath.ToSlash(strings.TrimPrefix(path, "./"))
	if idx := strings.Index(rel, "routes/"); idx >= 0 {
		rel = rel[idx+len("routes/"):]
	} else if strings.HasSuffix(rel, "/routes") || rel == "routes" {
		rel = ""
	}
	parts := []string{}
	for _, seg := range strings.Split(rel, "/") {
		if seg == "" {
			continue
		}
		parts = append(parts, cleanSegment(seg))
	}
	if len(parts) == 0 {
		return "index"
	}
	return strings.Join(parts, "_")
}

func buildPattern(path string) string {
	rel := filepath.ToSlash(strings.TrimPrefix(path, "./"))
	if idx := strings.Index(rel, "routes/"); idx >= 0 {
		rel = rel[idx+len("routes/"):]
	} else if strings.HasSuffix(rel, "/routes") || rel == "routes" {
		rel = ""
	}
	if rel == "" || rel == "." {
		return "/{$}"
	}
	segments := []string{}
	for _, seg := range strings.Split(rel, "/") {
		if seg == "" {
			continue
		}
		segments = append(segments, patternSegment(seg))
	}
	return "/" + strings.Join(segments, "/")
}

func cleanSegment(seg string) string {
	if strings.HasPrefix(seg, "[") && strings.HasSuffix(seg, "]") {
		return seg[1 : len(seg)-1]
	}
	if strings.HasPrefix(seg, "_") && strings.HasSuffix(seg, "_") {
		return seg[1 : len(seg)-1]
	}
	return seg
}

func patternSegment(seg string) string {
	if strings.HasPrefix(seg, "[") && strings.HasSuffix(seg, "]") {
		return "{" + seg[1:len(seg)-1] + "}"
	}
	if strings.HasPrefix(seg, "_") && strings.HasSuffix(seg, "_") {
		return "{" + seg[1:len(seg)-1] + "}"
	}
	return seg
}

func findLayout() *ast.File {
	var layout *ast.File
	filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || layout != nil {
			return nil
		}
		if filepath.Base(path) == "default.dreego" || filepath.Base(path) == "layout.dreego" {
			data, _ := os.ReadFile(path)
			tokens, _ := lexer.Lex(string(data))
			p := parser.NewParser(tokens)
			f, _ := p.Parse()
			if f != nil {
				layout = f
			}
		}
		return nil
	})
	return layout
}

func isLayoutDir(path string) bool {
	return strings.HasSuffix(path, "/layouts") || strings.Contains(path, "/layouts/") ||
		strings.HasSuffix(path, "/components") || strings.Contains(path, "/components/")
}
