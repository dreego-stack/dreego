package transpiler

import (
	"fmt"
	"strings"
)

func (p *Parser) parseIfNodes(openPos int) ([]TemplateNode, error) {
	var nodes []TemplateNode

	for {
		tok := p.current()
		if tok.Type == TokenEOF {
			return nil, fmt.Errorf("unclosed {#if} at position %d", openPos)
		}
		if tok.Type == TokenIfClose || tok.Type == TokenElse || tok.Type == TokenElseIf {
			break
		}
		if tok.Type == TokenTagClose {
			node, err := p.parseTemplateNode("if")
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
			continue
		}

		node, err := p.parseTemplateNode("if")
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (p *Parser) parseElseNodes() ([]TemplateNode, error) {
	var nodes []TemplateNode

	for {
		tok := p.current()
		if tok.Type == TokenEOF {
			return nil, fmt.Errorf("unclosed {#else}")
		}
		if tok.Type == TokenIfClose {
			break
		}
		if tok.Type == TokenElse || tok.Type == TokenElseIf {
			return nil, fmt.Errorf("unexpected {#else} or {#else if} inside {#else}")
		}
		if tok.Type == TokenTagClose {
			node, err := p.parseTemplateNode("if")
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
			continue
		}

		node, err := p.parseTemplateNode("if")
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (p *Parser) parseEachNodes() ([]TemplateNode, error) {
	var nodes []TemplateNode

	for {
		tok := p.current()
		if tok.Type == TokenEOF {
			return nil, fmt.Errorf("unclosed {#each}")
		}
		if tok.Type == TokenEachClose || tok.Type == TokenEachElse {
			break
		}
		if tok.Type == TokenTagClose {
			node, err := p.parseTemplateNode("each")
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
			continue
		}

		node, err := p.parseTemplateNode("each")
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (p *Parser) parseEachElseNodes() ([]TemplateNode, error) {
	var nodes []TemplateNode

	for {
		tok := p.current()
		if tok.Type == TokenEOF {
			return nil, fmt.Errorf("unclosed {#each else}")
		}
		if tok.Type == TokenEachClose {
			break
		}

		node, err := p.parseTemplateNode("each")
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func parseEachClause(clause string) (items, item string, err error) {
	parts := strings.SplitN(clause, " as ", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("expected 'items as item', got %q", clause)
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
}

func parseExpression(raw string) (expr string, filters []string) {
	if !strings.Contains(raw, "|") {
		return raw, nil
	}
	var parts []string
	start := 0
	var quote byte
	escaped := false
	for i := 0; i < len(raw); i++ {
		c := raw[i]
		if quote != 0 {
			if quote == '`' {
				if c == '`' {
					quote = 0
				}
				continue
			}
			if escaped {
				escaped = false
				continue
			}
			if c == '\\' {
				escaped = true
				continue
			}
			if c == quote {
				quote = 0
			}
			continue
		}
		switch c {
		case '"', '\'', '`':
			quote = c
		case '|':
			if (i+1 >= len(raw) || raw[i+1] != '|') && (i == 0 || raw[i-1] != '|') {
				parts = append(parts, raw[start:i])
				start = i + 1
			}
		}
	}
	parts = append(parts, raw[start:])
	expr = strings.TrimSpace(parts[0])
	for _, f := range parts[1:] {
		filters = append(filters, strings.TrimSpace(f))
	}
	return
}
