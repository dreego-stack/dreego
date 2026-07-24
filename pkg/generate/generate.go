package generate

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"codeberg.org/dreego/dreego/pkg/ast"
	"codeberg.org/dreego/dreego/pkg/codegen"
	"codeberg.org/dreego/dreego/pkg/lexer"
	"codeberg.org/dreego/dreego/pkg/parser"
)

var binaryHash string
var paramRe = regexp.MustCompile(`\[(\w+)\]|_(\w+)_`)

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

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
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
			return fs.SkipDir
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

		found++
		pkgName := filepath.Base(path)
		if pkgName == "routes" {
			pkgName = "routes"
		}

		p := buildPattern(path)
		pageName := buildPageName(path)
		outPath := filepath.Join(path, "dree.go")

		var methods []ast.GoSection

		for _, fpath := range dreegoFiles {
			fname := filepath.Base(fpath)
			method := "GET"
			for prefix, m := range methodExt {
				if strings.HasPrefix(fname, prefix+".dreego") || strings.HasPrefix(fname, prefix+"-") {
					method = m
					break
				}
			}

			input, err := os.ReadFile(fpath)
			if err != nil {
				return fmt.Errorf("error reading %s: %w", fpath, err)
			}

			tokens, err := lexer.Lex(string(input))
			if err != nil {
				return fmt.Errorf("error lexing %s: %w", fpath, err)
			}

			pr := parser.NewParser(tokens)
			file, err := pr.Parse()
			if err != nil {
				return fmt.Errorf("error parsing %s: %w", fpath, err)
			}

			for _, g := range file.Go {
				g.Method = method
				methods = append(methods, g)
			}
			if len(file.Go) == 0 {
				methods = append(methods, ast.GoSection{Method: method})
			}
		}

		if !force {
			if stale, _ := isStale(dreegoFiles, outPath); !stale {
				return nil
			}
		}

		combined := &ast.File{
			Head:     extractHead(dreegoFiles),
			Go:       methods,
			Template: extractTemplate(dreegoFiles),
			Script:   extractScript(dreegoFiles),
			Style:    extractStyle(dreegoFiles),
		}

		combinedHash := hashFiles(dreegoFiles)
		scopeHash := combinedHash[:12]

		handlerCode, err := codegen.GenerateHandler(combined, layout, pkgName, pageName, p, scopeHash)
		if err != nil {
			return fmt.Errorf("error generating %s: %w", path, err)
		}

		out := fmt.Sprintf("// hash:%s/%s\n%s", combinedHash, binaryHash, handlerCode)
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

func extractHead(files []string) *ast.HeadSection {
	for _, f := range files {
		data, _ := os.ReadFile(f)
		tokens, _ := lexer.Lex(string(data))
		p := parser.NewParser(tokens)
		file, _ := p.Parse()
		if file != nil && file.Head != nil {
			return file.Head
		}
	}
	return nil
}

func extractTemplate(files []string) *ast.TemplateSection {
	for _, f := range files {
		data, _ := os.ReadFile(f)
		tokens, _ := lexer.Lex(string(data))
		p := parser.NewParser(tokens)
		file, _ := p.Parse()
		if file != nil && file.Template != nil {
			return file.Template
		}
	}
	return nil
}

func extractScript(files []string) *ast.ScriptSection {
	for _, f := range files {
		data, _ := os.ReadFile(f)
		tokens, _ := lexer.Lex(string(data))
		p := parser.NewParser(tokens)
		file, _ := p.Parse()
		if file != nil && file.Script != nil {
			return file.Script
		}
	}
	return nil
}

func extractStyle(files []string) *ast.StyleSection {
	for _, f := range files {
		data, _ := os.ReadFile(f)
		tokens, _ := lexer.Lex(string(data))
		p := parser.NewParser(tokens)
		file, _ := p.Parse()
		if file != nil && file.Style != nil {
			return file.Style
		}
	}
	return nil
}

func hashFiles(files []string) string {
	h := sha256.New()
	for _, f := range files {
		data, _ := os.ReadFile(f)
		h.Write(data)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func isStale(files []string, outPath string) (bool, error) {
	outInfo, err := os.Stat(outPath)
	if err != nil {
		return true, nil
	}
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			return true, nil
		}
		if info.ModTime().After(outInfo.ModTime()) {
			return true, nil
		}
	}
	return false, nil
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
	filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
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
