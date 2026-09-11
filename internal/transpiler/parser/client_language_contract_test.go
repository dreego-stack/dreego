package parser

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
	"github.com/dreego-stack/dreego/internal/transpiler/lexer"
)

func TestClientLanguageContractRejectsUnsupportedLanguages(t *testing.T) {
	t.Parallel()
	for _, language := range []string{"python", "typescript", "coffee"} {
		t.Run(language, func(t *testing.T) {
			t.Parallel()
			source := `<body></body><client lang="` + language + `"></client>`
			tokens, err := lexer.Lex(source)
			if err != nil {
				t.Fatalf("Lex: %v", err)
			}
			_, err = NewParser(tokens).Parse()
			if err == nil || !strings.Contains(err.Error(), "unsupported language") {
				t.Fatalf("Parse() error = %v, want unsupported language", err)
			}
		})
	}
}

func TestClientLanguageContractNormalizesSupportedNames(t *testing.T) {
	t.Parallel()
	for _, language := range []string{"JS", "TS", "Lua"} {
		t.Run(language, func(t *testing.T) {
			t.Parallel()
			source := `<body></body><client lang="` + language + `"></client>`
			tokens, err := lexer.Lex(source)
			if err != nil {
				t.Fatalf("Lex: %v", err)
			}
			file, err := NewParser(tokens).Parse()
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if file.Client.Language != strings.ToLower(language) {
				t.Fatalf("language = %q", file.Client.Language)
			}
		})
	}
}

func TestInlineUnknownScriptLanguagesRemainOrdinaryHTML(t *testing.T) {
	t.Parallel()
	for _, language := range []string{"python", "typescript", "coffee"} {
		source := `<body><script lang="` + language + `">code</script></body>`
		tokens, err := lexer.Lex(source)
		if err != nil {
			t.Fatalf("Lex(%q): %v", language, err)
		}
		file, err := NewParser(tokens).Parse()
		if err != nil {
			t.Fatalf("Parse(%q): %v", language, err)
		}
		for _, node := range file.Body.Nodes {
			if node.Type == ir.NodeClientScript {
				t.Fatalf("inline %q script became a client processor node", language)
			}
		}
	}
}
