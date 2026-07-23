package generate

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"codeberg.org/dreego/dreego/pkg/codegen"
	"codeberg.org/dreego/dreego/pkg/lexer"
	"codeberg.org/dreego/dreego/pkg/parser"
)

func Run() error {
	start := time.Now()
	var found, generated int

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

		found++
		outPath := strings.TrimSuffix(path, ".dreego") + "_dreego.go"

		srcInfo, _ := os.Stat(path)
		outInfo, _ := os.Stat(outPath)
		if outInfo != nil && !srcInfo.ModTime().After(outInfo.ModTime()) {
			return nil
		}

		pkgName := filepath.Base(filepath.Dir(path))
		pageName := strings.TrimSuffix(filepath.Base(path), ".dreego")

		input, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading %s: %w", path, err)
		}

		tokens, err := lexer.Lex(string(input))
		if err != nil {
			return fmt.Errorf("error lexing %s: %w", path, err)
		}

		p := parser.NewParser(tokens)
		file, err := p.Parse()
		if err != nil {
			return fmt.Errorf("error parsing %s: %w", path, err)
		}

		handlerCode, err := codegen.GenerateHandler(file, pkgName, pageName)
		if err != nil {
			return fmt.Errorf("error generating code for %s: %w", path, err)
		}

		if err := os.WriteFile(outPath, []byte(handlerCode), 0644); err != nil {
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
	fmt.Printf("Generated %d _dreego.go files (in %s)\n", generated, formatDuration(elapsed))
	return nil
}

func formatDuration(d time.Duration) string {
	switch {
	case d < time.Microsecond:
		return fmt.Sprintf("%dns", d.Nanoseconds())
	case d < time.Millisecond:
		return fmt.Sprintf("%dµs", d.Microseconds())
	default:
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
}
