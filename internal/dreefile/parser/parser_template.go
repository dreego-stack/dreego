package parser

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/dreefile/ir"
	"github.com/dreego-stack/dreego/internal/dreefile/tokens"
)

func (p *Parser) parseIfNodes(openPos int) ([]ir.TemplateNode, error) {
	var nodes []ir.TemplateNode

	for {
		tok := p.current()
		if tok.Type == tokens.TokenEOF {
			return nil, fmt.Errorf("unclosed {#if} at position %d", openPos)
		}
		if tok.Type == tokens.TokenIfClose || tok.Type == tokens.TokenElse || tok.Type == tokens.TokenElseIf {
			break
		}
		if tok.Type == tokens.TokenTagClose {
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

func (p *Parser) parseElseNodes() ([]ir.TemplateNode, error) {
	var nodes []ir.TemplateNode

	for {
		tok := p.current()
		if tok.Type == tokens.TokenEOF {
			return nil, fmt.Errorf("unclosed {#else}")
		}
		if tok.Type == tokens.TokenIfClose {
			break
		}
		if tok.Type == tokens.TokenElse || tok.Type == tokens.TokenElseIf {
			return nil, fmt.Errorf("unexpected {#else} or {#else if} inside {#else}")
		}
		if tok.Type == tokens.TokenTagClose {
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

func (p *Parser) parseEachNodes() ([]ir.TemplateNode, error) {
	var nodes []ir.TemplateNode

	for {
		tok := p.current()
		if tok.Type == tokens.TokenEOF {
			return nil, fmt.Errorf("unclosed {#each}")
		}
		if tok.Type == tokens.TokenEachClose || tok.Type == tokens.TokenEachElse {
			break
		}
		if tok.Type == tokens.TokenTagClose {
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

func (p *Parser) parseEachElseNodes() ([]ir.TemplateNode, error) {
	var nodes []ir.TemplateNode

	for {
		tok := p.current()
		if tok.Type == tokens.TokenEOF {
			return nil, fmt.Errorf("unclosed {#each else}")
		}
		if tok.Type == tokens.TokenEachClose {
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
