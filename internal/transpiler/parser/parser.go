package parser

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
	"github.com/dreego-stack/dreego/internal/transpiler/tokens"
)

type Parser struct {
	tokens           []tokens.Token
	pos              int
	templateFromBody bool
}

func NewParser(tokens []tokens.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() (*ir.File, error) {
	file := &ir.File{}

	for p.pos < len(p.tokens) {
		tok := p.current()

		if tok.Type == tokens.TokenEOF {
			break
		}

		if tok.Type == tokens.TokenText && strings.TrimSpace(tok.Value) == "" {
			p.advance()
			continue
		}

		if tok.Type != tokens.TokenTagOpen {
			return nil, fmt.Errorf("expected root section, got %s at position %d", tok.Type, tok.Pos)
		}

		switch tok.Tag {
		case "server":
			language, err := parseSectionLanguage(tok, "go")
			if err != nil {
				return nil, err
			}
			section, err := p.parseServerSection()
			if err != nil {
				return nil, err
			}
			section.Method = "GET"
			section.Language = language
			section.ContentType = parseServerAttrs(tok.Attr)
			if m := parseServerMethod(tok.Attr); m != "" {
				section.Method = m
				section.MethodExplicit = true
			}
			file.Server = append(file.Server, *section)
		case "body":
			language, err := parseSectionLanguage(tok, "html")
			if err != nil {
				return nil, err
			}
			section, err := p.parseBodySection()
			if err != nil {
				return nil, err
			}
			section.Language = language
			for _, existing := range file.Bodies {
				if existing.Method == section.Method {
					return nil, fmt.Errorf("duplicate <body> section for method %s at position %d", section.Method, tok.Pos)
				}
			}
			if file.Body == nil && section.Method == "GET" {
				file.Body = section
			}
			file.Bodies = append(file.Bodies, *section)
			p.templateFromBody = true
		case "head":
			language, err := parseSectionLanguage(tok, "html")
			if err != nil {
				return nil, err
			}
			section, err := p.parseRawSection("head")
			if err != nil {
				return nil, err
			}
			if file.Head != nil {
				return nil, fmt.Errorf("duplicate <head> section at position %d", tok.Pos)
			}
			file.Head = &ir.HeadSection{Content: strings.TrimSpace(section), Language: language}
		case "client":
			language, err := parseSectionLanguage(tok, "js")
			if err != nil {
				return nil, err
			}
			section, err := p.parseRawSection("client")
			if err != nil {
				return nil, err
			}
			if file.Client != nil {
				return nil, fmt.Errorf("duplicate <client> section at position %d", tok.Pos)
			}
			file.Client = &ir.ClientSection{Code: strings.TrimSpace(section), Language: language}
		case "style":
			language, err := parseSectionLanguage(tok, "css")
			if err != nil {
				return nil, err
			}
			section, err := p.parseRawSection("style")
			if err != nil {
				return nil, err
			}
			if file.Style != nil {
				return nil, fmt.Errorf("duplicate <style> section at position %d", tok.Pos)
			}
			file.Style = &ir.StyleSection{Code: strings.TrimSpace(section), Language: language}
		case "go":
			return nil, fmt.Errorf("legacy root <go> at position %d: replace root <go> with <server>", tok.Pos)
		case "div":
			return nil, fmt.Errorf("legacy root <div> at position %d: replace root <div> with <body>", tok.Pos)
		case "script":
			return nil, fmt.Errorf("legacy root <script> at position %d: replace root <script> with <client>", tok.Pos)
		default:
			return nil, fmt.Errorf("expected root section, got <%s> at position %d", tok.Tag, tok.Pos)
		}
	}

	return file, nil
}

func parseSectionLanguage(tok tokens.Token, defaultLanguage string) (string, error) {
	language := sectionLanguage(tok.Attr)
	if language == "" {
		return defaultLanguage, nil
	}
	if language == defaultLanguage {
		return language, nil
	}
	return "", fmt.Errorf("unsupported language %q for <%s> at position %d; install a processor for this section and language", language, tok.Tag, tok.Pos)
}

func sectionLanguage(attrs string) string {
	for _, part := range strings.Fields(attrs) {
		if strings.HasPrefix(part, "lang=") {
			return strings.ToLower(strings.Trim(strings.TrimPrefix(part, "lang="), "\"'"))
		}
	}
	return ""
}

func (p *Parser) advance() {
	p.pos++
}

func (p *Parser) current() tokens.Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return tokens.Token{Type: tokens.TokenEOF}
}

func (p *Parser) peek() tokens.Token {
	if p.pos+1 < len(p.tokens) {
		return p.tokens[p.pos+1]
	}
	return tokens.Token{Type: tokens.TokenEOF}
}

func parseServerAttrs(attrs string) string {
	if attrs == "" {
		return ""
	}
	for _, part := range strings.Fields(attrs) {
		if strings.HasPrefix(part, "type=") {
			v := strings.TrimPrefix(part, "type=")
			return strings.Trim(v, "\"'")
		}
	}
	return ""
}

// parseServerMethod extracts an explicit method= attribute from a <server> section's
// attributes. It returns "" when no method attribute is present.
func parseServerMethod(attrs string) string {
	if attrs == "" {
		return ""
	}
	for _, part := range strings.Fields(attrs) {
		if strings.HasPrefix(part, "method=") {
			v := strings.Trim(strings.TrimPrefix(part, "method="), "\"'")
			return strings.ToUpper(v)
		}
	}
	return ""
}
