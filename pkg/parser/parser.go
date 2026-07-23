package parser

import (
	"fmt"
	"strings"

	"codeberg.org/dreego/dreego/pkg/ast"
	"codeberg.org/dreego/dreego/pkg/lexer"
)

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() (*ast.File, error) {
	file := &ast.File{}

	for p.pos < len(p.tokens) {
		tok := p.current()

		if tok.Type == lexer.TokenEOF {
			break
		}

		if tok.Type == lexer.TokenText && strings.TrimSpace(tok.Value) == "" {
			p.advance()
			continue
		}

		if tok.Type != lexer.TokenTagOpen {
			return nil, fmt.Errorf("expected section tag, got %s at position %d", tok.Type, tok.Pos)
		}

		switch tok.Tag {
		case "go":
			section, err := p.parseGoSection()
			if err != nil {
				return nil, err
			}
			if file.Go != nil {
				return nil, fmt.Errorf("duplicate <go> section at position %d", tok.Pos)
			}
			file.Go = section
		case "div":
			section, err := p.parseDivSection()
			if err != nil {
				return nil, err
			}
			if file.Template != nil {
				return nil, fmt.Errorf("duplicate <div> section at position %d", tok.Pos)
			}
			file.Template = section
		case "head":
			section, err := p.parseNonDivSection("head")
			if err != nil {
				return nil, err
			}
			if file.Head != nil {
				return nil, fmt.Errorf("duplicate <head> section at position %d", tok.Pos)
			}
			file.Head = &ast.HeadSection{Content: strings.TrimSpace(section)}
		case "script":
			section, err := p.parseNonDivSection("script")
			if err != nil {
				return nil, err
			}
			if file.Script != nil {
				return nil, fmt.Errorf("duplicate <script> section at position %d", tok.Pos)
			}
			file.Script = &ast.ScriptSection{Code: strings.TrimSpace(section)}
		case "style":
			section, err := p.parseNonDivSection("style")
			if err != nil {
				return nil, err
			}
			if file.Style != nil {
				return nil, fmt.Errorf("duplicate <style> section at position %d", tok.Pos)
			}
			file.Style = &ast.StyleSection{Code: strings.TrimSpace(section)}
		default:
			return nil, fmt.Errorf("unknown section <%s> at position %d", tok.Tag, tok.Pos)
		}
	}

	return file, nil
}

func (p *Parser) advance() {
	p.pos++
}

func (p *Parser) current() lexer.Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return lexer.Token{Type: lexer.TokenEOF}
}

func (p *Parser) peek() lexer.Token {
	if p.pos+1 < len(p.tokens) {
		return p.tokens[p.pos+1]
	}
	return lexer.Token{Type: lexer.TokenEOF}
}
