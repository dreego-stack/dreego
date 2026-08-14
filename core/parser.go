package core

import (
	"fmt"
	"strings"
)

type Parser struct {
	tokens          []Token
	pos             int
	templateFromDiv bool
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() (*File, error) {
	file := &File{}

	for p.pos < len(p.tokens) {
		tok := p.current()

		if tok.Type == TokenEOF {
			break
		}

		if tok.Type == TokenText && strings.TrimSpace(tok.Value) == "" {
			p.advance()
			continue
		}

		if tok.Type != TokenTagOpen {
			return nil, fmt.Errorf("expected root section, got %s at position %d", tok.Type, tok.Pos)
		}

		switch tok.Tag {
		case "go":
			section, err := p.parseGoSection()
			if err != nil {
				return nil, err
			}
			section.Method = "GET"
			section.ContentType = parseGoAttrs(tok.Attr)
			file.Go = append(file.Go, *section)
		case "div":
			section, err := p.parseDivSection()
			if err != nil {
				return nil, err
			}
			if p.templateFromDiv {
				return nil, fmt.Errorf("duplicate <div> section at position %d", tok.Pos)
			}
			if file.Template == nil {
				file.Template = section
				p.templateFromDiv = true
			} else {
				file.Template.Nodes = append(file.Template.Nodes, section.Nodes...)
			}
		case "head":
			section, err := p.parseNonDivSection("head")
			if err != nil {
				return nil, err
			}
			if file.Head != nil {
				return nil, fmt.Errorf("duplicate <head> section at position %d", tok.Pos)
			}
			file.Head = &HeadSection{Content: strings.TrimSpace(section)}
		case "script":
			section, err := p.parseNonDivSection("script")
			if err != nil {
				return nil, err
			}
			if file.Script != nil {
				return nil, fmt.Errorf("duplicate <script> section at position %d", tok.Pos)
			}
			file.Script = &ScriptSection{Code: strings.TrimSpace(section)}
		case "style":
			section, err := p.parseNonDivSection("style")
			if err != nil {
				return nil, err
			}
			if file.Style != nil {
				return nil, fmt.Errorf("duplicate <style> section at position %d", tok.Pos)
			}
			file.Style = &StyleSection{Code: strings.TrimSpace(section)}
		default:
			return nil, fmt.Errorf("expected root section, got <%s> at position %d", tok.Tag, tok.Pos)
		}
	}

	return file, nil
}

func (p *Parser) advance() {
	p.pos++
}

func (p *Parser) current() Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return Token{Type: TokenEOF}
}

func (p *Parser) peek() Token {
	if p.pos+1 < len(p.tokens) {
		return p.tokens[p.pos+1]
	}
	return Token{Type: TokenEOF}
}

func isSectionTag(tag string) bool {
	switch tag {
	case "go", "div", "head", "script", "style":
		return true
	}
	return false
}

func parseGoAttrs(attrs string) string {
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
