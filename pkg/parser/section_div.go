package parser

import (
	"fmt"

	"codeberg.org/dreego/dreego/pkg/ast"
	"codeberg.org/dreego/dreego/pkg/lexer"
)

func (p *Parser) parseDivSection() (*ast.TemplateSection, error) {
	p.advance()
	nodes, err := p.parseDivNodes()
	if err != nil {
		return nil, err
	}
	return &ast.TemplateSection{Nodes: nodes}, nil
}

func (p *Parser) parseDivNodes() ([]ast.TemplateNode, error) {
	var nodes []ast.TemplateNode

	for {
		tok := p.current()
		if tok.Type == lexer.TokenEOF {
			return nil, fmt.Errorf("unclosed <div>")
		}
		if tok.Type == lexer.TokenTagClose {
			p.advance()
			return nodes, nil
		}

		node, err := p.parseTemplateNode("div")
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
}

func (p *Parser) parseTemplateNode(parent string) (ast.TemplateNode, error) {
	tok := p.current()

	switch tok.Type {
	case lexer.TokenText:
		p.advance()
		return ast.TemplateNode{Type: ast.NodeText, Content: tok.Value}, nil
	case lexer.TokenExpression:
		p.advance()
		return ast.TemplateNode{Type: ast.NodeExpression, Content: tok.Value}, nil
	case lexer.TokenIfOpen:
		p.advance()
		children, err := p.parseIfNodes()
		if err != nil {
			return ast.TemplateNode{}, err
		}
		return ast.TemplateNode{Type: ast.NodeIf, Cond: tok.Value, Children: children}, nil
	case lexer.TokenEachOpen:
		p.advance()
		items, item, err := parseEachClause(tok.Value)
		if err != nil {
			return ast.TemplateNode{}, err
		}
		children, err := p.parseEachNodes()
		if err != nil {
			return ast.TemplateNode{}, err
		}
		return ast.TemplateNode{Type: ast.NodeEach, Items: items, Item: item, Children: children}, nil
	case lexer.TokenTagOpen:
		return ast.TemplateNode{}, fmt.Errorf("unexpected <%s> inside <%s> at position %d", tok.Tag, parent, tok.Pos)
	default:
		return ast.TemplateNode{}, fmt.Errorf("unexpected token %s in template at position %d", tok.Type, tok.Pos)
	}
}
