package generate

import (
	"bufio"
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
var paramRe = regexp.MustCompile(`\[(\w+)\]`)

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
	var found, generated, skipped int

	if err := checkConflicts(); err != nil {
		return err
	}

	layout := findLayout()

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			base := filepath.Base(path)
			if base != "." && (strings.HasPrefix(base, ".") || strings.HasPrefix(base, "_")) {
				return fs.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".dreego" {
			return nil
		}
		if isLayout(path) {
			return nil
		}

		found++

		input, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading %s: %w", path, err)
		}

		srcHash := hashContent(input)
		outPath := cleanPath(strings.TrimSuffix(path, ".dreego")) + "_dreego.go"

		if !force {
			if cSrc, cBin, ok := readHashes(outPath); ok && cSrc == srcHash && cBin == binaryHash {
				skipped++
				return nil
			}
		}

		pkgName := filepath.Base(filepath.Dir(path))
		pageName := buildPageName(path)
		pattern := buildPattern(path)

		tokens, err := lexer.Lex(string(input))
		if err != nil {
			return fmt.Errorf("error lexing %s: %w", path, err)
		}

		p := parser.NewParser(tokens)
		file, err := p.Parse()
		if err != nil {
			return fmt.Errorf("error parsing %s: %w", path, err)
		}

		handlerCode, err := codegen.GenerateHandler(file, layout, pkgName, pageName, pattern, srcHash[:12])
		if err != nil {
			return fmt.Errorf("error generating code for %s: %w", path, err)
		}

		out := fmt.Sprintf("// hash:%s/%s\n%s", srcHash, binaryHash, handlerCode)
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
	fmt.Printf("Found %d .dreego files\n", found)
	fmt.Printf("Generated %d _dreego.go files\n", generated)
	if skipped > 0 {
		fmt.Printf("Skipped %d (unchanged)\n", skipped)
	}
	fmt.Printf("in %dns\n", elapsed.Nanoseconds())
	return nil
}

func buildPageName(path string) string {
	rel := filepath.ToSlash(strings.TrimPrefix(path, "./"))
	idx := strings.Index(rel, "routes/")
	if idx >= 0 {
		rel = rel[idx+len("routes/"):]
	}
	dir, file := filepath.Split(rel)
	baseName := strings.TrimSuffix(file, ".dreego")

	parts := []string{}
	if dir != "" && dir != "./" && dir != "." {
		for _, seg := range strings.Split(strings.Trim(dir, "/"), "/") {
			parts = append(parts, paramRe.ReplaceAllString(seg, "${1}"))
		}
	}

	if baseName != "index" {
		parts = append(parts, paramRe.ReplaceAllString(baseName, "${1}"))
	}

	if len(parts) == 0 {
		return "index"
	}
	return strings.Join(parts, "_")
}

func buildPattern(path string) string {
	rel := filepath.ToSlash(strings.TrimPrefix(path, "./"))
	idx := strings.Index(rel, "routes/")
	if idx >= 0 {
		rel = rel[idx+len("routes/"):]
	}

	dir, file := filepath.Split(rel)
	baseName := strings.TrimSuffix(file, ".dreego")

	segments := []string{}
	if dir != "" && dir != "./" && dir != "." {
		for _, seg := range strings.Split(strings.Trim(dir, "/"), "/") {
			seg = paramRe.ReplaceAllString(seg, "{$1}")
			segments = append(segments, seg)
		}
	}

	if baseName == "index" {
		if len(segments) == 0 {
			return "/{$}"
		}
		return "/" + strings.Join(segments, "/") + "/"
	}

	baseName = paramRe.ReplaceAllString(baseName, "{$1}")
	segments = append(segments, baseName)

	return "/" + strings.Join(segments, "/")
}

func checkConflicts() error {
	routes := map[string]string{}

	return filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			base := filepath.Base(path)
			if base != "." && (strings.HasPrefix(base, ".") || strings.HasPrefix(base, "_")) {
				return fs.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".dreego" {
			return nil
		}
		if isLayout(path) {
			return nil
		}

		p := buildPattern(path)
		if existing, ok := routes[p]; ok {
			return fmt.Errorf("route conflict: %s and %s both resolve to %s", existing, path, p)
		}
		routes[p] = path

		alt := strings.TrimRight(p, "/")
		if alt != p {
			if existing, ok := routes[alt]; ok {
				return fmt.Errorf("route conflict: %s -> %s and %s -> %s", existing, alt, path, p)
			}
		}
		alt = p + "/"
		if alt != p {
			if existing, ok := routes[alt]; ok {
				return fmt.Errorf("route conflict: %s -> %s and %s -> %s", existing, alt, path, p)
			}
		}
		return nil
	})
}

func readHashes(path string) (src, bin string, ok bool) {
	f, err := os.Open(path)
	if err != nil {
		return "", "", false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return "", "", false
	}
	line := scanner.Text()
	if !strings.HasPrefix(line, "// hash:") {
		return "", "", false
	}
	parts := strings.SplitN(strings.TrimPrefix(line, "// hash:"), "/", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func findLayout() *ast.File {
	var layout *ast.File
	filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || layout != nil {
			return nil
		}
		if filepath.Base(path) == "default.dreego" || filepath.Base(path) == "layout.dreego" {
			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			tokens, err := lexer.Lex(string(data))
			if err != nil {
				return nil
			}
			p := parser.NewParser(tokens)
			f, err := p.Parse()
			if err == nil {
				layout = f
			}
		}
		return nil
	})
	return layout
}

func cleanPath(path string) string {
	path = strings.ReplaceAll(path, "[", "p_")
	path = strings.ReplaceAll(path, "]", "")
	return path
}

func isLayout(path string) bool {
	return filepath.Base(path) == "default.dreego" || filepath.Base(path) == "layout.dreego"
}

func hashContent(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}
