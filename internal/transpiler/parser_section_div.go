package transpiler

import (
	"fmt"
	"strings"
)

func (p *Parser) parseDivSection() (*TemplateSection, error) {
	tok := p.current()
	if err := checkAttrControlFlow(tok.Attr, tok.Pos); err != nil {
		return nil, err
	}
	p.advance()
	nodes, err := p.parseDivNodes()
	if err != nil {
		return nil, err
	}
	return &TemplateSection{Nodes: nodes}, nil
}

func (p *Parser) parseDivNodes() ([]TemplateNode, error) {
	var nodes []TemplateNode
	depth := 0

	for {
		tok := p.current()
		if tok.Type == TokenEOF {
			return nil, fmt.Errorf("unclosed <div>")
		}
		if tok.Type == TokenTagOpen && tok.Tag == "div" {
			if err := checkAttrControlFlow(tok.Attr, tok.Pos); err != nil {
				return nil, err
			}
			p.advance()
			depth++
			content := fmt.Sprintf("<%s", tok.Tag)
			if tok.Attr != "" {
				content += " " + tok.Attr
			}
			content += ">"
			nodes = append(nodes, TemplateNode{Type: NodeText, Content: content})
			continue
		}
		if tok.Type == TokenTagClose && tok.Tag == "div" {
			p.advance()
			depth--
			if depth < 0 {
				return nodes, nil
			}
			nodes = append(nodes, TemplateNode{Type: NodeText, Content: fmt.Sprintf("</%s>", tok.Tag)})
			continue
		}

		node, err := p.parseTemplateNode("div")
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
}

func (p *Parser) parseTemplateNode(parent string) (TemplateNode, error) {
	tok := p.current()

	switch tok.Type {
	case TokenText:
		p.advance()
		return TemplateNode{Type: NodeText, Content: tok.Value, Pos: tok.Pos}, nil
	case TokenExpression:
		p.advance()
		expr, filters := parseExpression(tok.Value)
		return TemplateNode{Type: NodeExpression, Content: expr, Filters: filters, Pos: tok.Pos}, nil
	case TokenIfOpen:
		cond := tok.Value
		openPos := tok.Pos
		p.advance()
		children, err := p.parseIfNodes(openPos)
		if err != nil {
			return TemplateNode{}, err
		}
		var elseChildren []TemplateNode
		for p.current().Type == TokenElse || p.current().Type == TokenElseIf {
			if p.current().Type == TokenElse {
				p.advance()
				elseChildren, err = p.parseElseNodes()
				if err != nil {
					return TemplateNode{}, err
				}
				break
			}
			elseCond := p.current().Value
			elsePos := p.current().Pos
			p.advance()
			elseIfChildren, err := p.parseIfNodes(elsePos)
			if err != nil {
				return TemplateNode{}, err
			}
			nested := TemplateNode{Type: NodeIf, Cond: elseCond, Children: elseIfChildren}
			if p.current().Type == TokenElse {
				p.advance()
				nested.ElseChildren, err = p.parseElseNodes()
				if err != nil {
					return TemplateNode{}, err
				}
			}
			elseChildren = append(elseChildren, nested)
			if p.current().Type == TokenIfClose {
				p.advance()
				break
			}
		}
		if p.current().Type == TokenIfClose {
			p.advance()
		}
		return TemplateNode{Type: NodeIf, Cond: cond, Children: children, ElseChildren: elseChildren}, nil
	case TokenEachOpen:
		items, item, err := parseEachClause(tok.Value)
		if err != nil {
			return TemplateNode{}, err
		}
		p.advance()
		children, err := p.parseEachNodes()
		if err != nil {
			return TemplateNode{}, err
		}
		var elseChildren []TemplateNode
		if p.current().Type == TokenEachElse {
			p.advance()
			elseChildren, err = p.parseEachElseNodes()
			if err != nil {
				return TemplateNode{}, err
			}
		}
		if p.current().Type == TokenEachClose {
			p.advance()
		}
		return TemplateNode{Type: NodeEach, Items: items, Item: item, Children: children, ElseChildren: elseChildren}, nil
	case TokenTagClose:
		p.advance()
		return TemplateNode{Type: NodeText, Content: fmt.Sprintf("</%s>", tok.Tag), Pos: tok.Pos}, nil
	case TokenSlot:
		p.advance()
		return TemplateNode{Type: NodeSlot, Content: tok.Value}, nil
	case TokenSlotOpen:
		p.advance()
		children, err := p.parseSlotNodes()
		if err != nil {
			return TemplateNode{}, err
		}
		return TemplateNode{Type: NodeSlot, Content: tok.Value, Children: children, Pos: tok.Pos}, nil
	case TokenTagOpen:
		if err := checkAttrControlFlow(tok.Attr, tok.Pos); err != nil {
			return TemplateNode{}, err
		}
		p.advance()
		content := fmt.Sprintf("<%s", tok.Tag)
		if tok.Attr != "" {
			content += " " + tok.Attr
		}
		content += ">"
		return TemplateNode{Type: NodeText, Content: content, Pos: tok.Pos}, nil
	case TokenComponentSelfClose:
		if err := checkAttrControlFlow(tok.Attr, tok.Pos); err != nil {
			return TemplateNode{}, err
		}
		p.advance()
		return TemplateNode{Type: NodeComponentCall, Tag: tok.Tag, Attrs: tok.Attr, SelfClose: true, Pos: tok.Pos}, nil
	case TokenComponentTagOpen:
		if err := checkAttrControlFlow(tok.Attr, tok.Pos); err != nil {
			return TemplateNode{}, err
		}
		p.advance()
		children, err := p.parseComponentNodes(tok.Tag)
		if err != nil {
			return TemplateNode{}, err
		}
		return TemplateNode{Type: NodeComponentCall, Tag: tok.Tag, Attrs: tok.Attr, Children: children, Pos: tok.Pos}, nil
	case TokenSlotClose:
		return TemplateNode{}, fmt.Errorf("unexpected {/slot} at position %d", tok.Pos)
	case TokenElse:
		return TemplateNode{}, fmt.Errorf("unexpected {#else} outside {#if} at position %d", tok.Pos)
	case TokenElseIf:
		return TemplateNode{}, fmt.Errorf("unexpected {#else if} outside {#if} at position %d", tok.Pos)
	case TokenVerbatim:
		p.advance()
		return TemplateNode{Type: NodeVerbatim, Content: tok.Value}, nil
	default:
		return TemplateNode{}, fmt.Errorf("unexpected token %s in template at position %d", tok.Type, tok.Pos)
	}
}

func (p *Parser) parseSlotNodes() ([]TemplateNode, error) {
	var nodes []TemplateNode
	for {
		tok := p.current()
		if tok.Type == TokenEOF {
			return nil, fmt.Errorf("unclosed {#slot}")
		}
		if tok.Type == TokenSlotClose {
			p.advance()
			return nodes, nil
		}
		node, err := p.parseTemplateNode("slot")
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
}

func (p *Parser) parseComponentNodes(tag string) ([]TemplateNode, error) {
	var nodes []TemplateNode

	for {
		tok := p.current()
		if tok.Type == TokenEOF {
			return nil, fmt.Errorf("unclosed <@%s>", tag)
		}
		if tok.Type == TokenComponentTagClose && tok.Tag == tag {
			p.advance()
			return nodes, nil
		}
		if tok.Type == TokenComponentTagClose {
			return nil, fmt.Errorf("unexpected </@%s>, expected </@%s>", tok.Tag, tag)
		}

		node, err := p.parseTemplateNode("component")
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
}

// checkAttrControlFlow rejects {#if}/{#each} control flow inside tag attribute
// values. The lexer scans a whole tag as one token, so control flow inside a
// quoted attribute would either stay literal (route path: cond never referenced)
// or be misparsed as an expression (component path: broken Go). Fail fast with
// a clear error instead of silently generating corrupt code.
func checkAttrControlFlow(attrs string, pos int) error {
	if attrs == "" {
		return nil
	}
	inQuote := false
	for i := 0; i < len(attrs); i++ {
		switch attrs[i] {
		case '"', '\'':
			inQuote = !inQuote
		case '{':
			if inQuote && strings.HasPrefix(attrs[i:], "{#if ") {
				return fmt.Errorf("{#if} inside attribute value at position %d is not supported; wrap the whole tag in {#if} instead", pos)
			}
			if inQuote && strings.HasPrefix(attrs[i:], "{#each ") {
				return fmt.Errorf("{#each} inside attribute value at position %d is not supported; wrap the whole tag in {#each} instead", pos)
			}
		}
	}
	return nil
}
