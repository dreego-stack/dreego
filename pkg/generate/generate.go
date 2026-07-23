package generate

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"codeberg.org/dreego/dreego/pkg/codegen"
	"codeberg.org/dreego/dreego/pkg/lexer"
	"codeberg.org/dreego/dreego/pkg/parser"
)

func Run(routesDir string) error {
	if routesDir == "" {
		routesDir = "routes"
	}

	return filepath.WalkDir(routesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".dreego" {
			return nil
		}

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

		pageName := strings.TrimSuffix(filepath.Base(path), ".dreego")
		handlerCode, err := codegen.GenerateHandler(file, pageName)
		if err != nil {
			return fmt.Errorf("error generating code for %s: %w", path, err)
		}

		outPath := strings.TrimSuffix(path, ".dreego") + "_dreego.go"
		if err := os.WriteFile(outPath, []byte(handlerCode), 0644); err != nil {
			return fmt.Errorf("error writing %s: %w", outPath, err)
		}

		fmt.Printf("generated: %s -> %s\n", path, outPath)
		return nil
	})
}
