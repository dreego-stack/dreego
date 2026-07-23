package generate

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"codeberg.org/dreego/dreego/pkg/ast"
	"codeberg.org/dreego/dreego/pkg/codegen"
	"codeberg.org/dreego/dreego/pkg/lexer"
	"codeberg.org/dreego/dreego/pkg/parser"
)

func Run() error {
	start := time.Now()
	var found, generated, skipped int

	layout := findLayout()

	type job struct {
		path     string
		pkgName  string
		pageName string
	}
	var jobs []job

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
		pkgName := filepath.Base(filepath.Dir(path))
		pageName := strings.TrimSuffix(filepath.Base(path), ".dreego")
		jobs = append(jobs, job{path, pkgName, pageName})
		return nil
	})
	if err != nil {
		return err
	}

	for _, j := range jobs {
		input, err := os.ReadFile(j.path)
		if err != nil {
			return fmt.Errorf("error reading %s: %w", j.path, err)
		}

		hash := hashContent(input)
		outPath := strings.TrimSuffix(j.path, ".dreego") + "_dreego.go"

		if existingHash, ok := readHashFromFile(outPath); ok && existingHash == hash {
			skipped++
			continue
		}

		tokens, err := lexer.Lex(string(input))
		if err != nil {
			return fmt.Errorf("error lexing %s: %w", j.path, err)
		}

		p := parser.NewParser(tokens)
		file, err := p.Parse()
		if err != nil {
			return fmt.Errorf("error parsing %s: %w", j.path, err)
		}

		handlerCode, err := codegen.GenerateHandler(file, layout, j.pkgName, j.pageName)
		if err != nil {
			return fmt.Errorf("error generating code for %s: %w", j.path, err)
		}

		out := "// hash:" + hash + "\n" + handlerCode
		if err := os.WriteFile(outPath, []byte(out), 0644); err != nil {
			return fmt.Errorf("error writing %s: %w", outPath, err)
		}

		generated++
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

func findLayout() *ast.File {
	candidates := []string{
		"dreego/layouts/default.dreego",
		"layout.dreego",
	}
	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err == nil {
			tokens, err := lexer.Lex(string(data))
			if err == nil {
				p := parser.NewParser(tokens)
				f, err := p.Parse()
				if err == nil {
					return f
				}
			}
		}
	}

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

func isLayout(path string) bool {
	return path == "dreego/layouts/default.dreego" || filepath.Base(path) == "layout.dreego"
}

func hashContent(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func readHashFromFile(path string) (string, bool) {
	f, err := os.Open(path)
	if err != nil {
		return "", false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return "", false
	}
	line := scanner.Text()
	if !strings.HasPrefix(line, "// hash:") {
		return "", false
	}
	return strings.TrimPrefix(line, "// hash:"), true
}
