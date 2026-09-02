package parser

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/html/md"
	"github.com/dreego-stack/dreego/internal/transpiler/ir"
	"github.com/dreego-stack/dreego/internal/transpiler/tokens"
)

func (p *Parser) parseBodySection() (*ir.BodySection, error) {
	tok := p.current()
	if err := checkAttrControlFlow(tok.Attr, tok.Pos); err != nil {
		return nil, err
	}
	p.advance()
	nodes, err := p.parseBodyNodes()
	if err != nil {
		return nil, err
	}
	language := sectionLanguage(tok.Attr)
	if language == "md" {
		if err := rejectMdTagInMdBody(nodes); err != nil {
			return nil, err
		}
		nodes, err = md.TransformNodes(nodes)
		if err != nil {
			return nil, err
		}
	} else {
		nodes, err = transformInlineMd(nodes)
		if err != nil {
			return nil, err
		}
	}
	method, explicit := parseBodyMethod(tok.Attr)
	return &ir.BodySection{Nodes: nodes, Method: method, MethodExplicit: explicit}, nil
}

func parseBodyMethod(attrs string) (string, bool) {
	for _, part := range strings.Fields(attrs) {
		if strings.HasPrefix(part, "method=") {
			return strings.ToUpper(strings.Trim(strings.TrimPrefix(part, "method="), "\"'")), true
		}
	}
	return "GET", false
}

func (p *Parser) parseBodyNodes() ([]ir.TemplateNode, error) {
	var nodes []ir.TemplateNode
	depth := 0

	for {
		tok := p.current()
		if tok.Type == tokens.TokenEOF {
			return nil, fmt.Errorf("unclosed <body>")
		}
		if tok.Type == tokens.TokenTagOpen && tok.Tag == "body" {
			if err := checkAttrControlFlow(tok.Attr, tok.Pos); err != nil {
				return nil, err
			}
			p.advance()
			depth++
			content := "<body"
			if tok.Attr != "" {
				content += " " + tok.Attr
			}
			content += ">"
			nodes = append(nodes, ir.TemplateNode{Type: ir.NodeText, Content: content, Pos: tok.Pos})
			continue
		}
		if tok.Type == tokens.TokenTagClose && tok.Tag == "body" {
			p.advance()
			depth--
			if depth < 0 {
				return nodes, nil
			}
			nodes = append(nodes, ir.TemplateNode{Type: ir.NodeText, Content: "</body>", Pos: tok.Pos})
			continue
		}

		node, err := p.parseTemplateNode("body")
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
}

func (p *Parser) parseTemplateNode(parent string) (ir.TemplateNode, error) {
	tok := p.current()

	switch tok.Type {
	case tokens.TokenText:
		p.advance()
		return ir.TemplateNode{Type: ir.NodeText, Content: tok.Value, Pos: tok.Pos}, nil
	case tokens.TokenExpression:
		p.advance()
		expr, filters := parseExpression(tok.Value)
		return ir.TemplateNode{Type: ir.NodeExpression, Content: expr, Filters: filters, Pos: tok.Pos}, nil
	case tokens.TokenIfOpen:
		cond := tok.Value
		openPos := tok.Pos
		p.advance()
		children, err := p.parseIfNodes(openPos)
		if err != nil {
			return ir.TemplateNode{}, err
		}
		var elseChildren []ir.TemplateNode
		for p.current().Type == tokens.TokenElse || p.current().Type == tokens.TokenElseIf {
			if p.current().Type == tokens.TokenElse {
				p.advance()
				elseChildren, err = p.parseElseNodes()
				if err != nil {
					return ir.TemplateNode{}, err
				}
				break
			}
			elseCond := p.current().Value
			elsePos := p.current().Pos
			p.advance()
			elseIfChildren, err := p.parseIfNodes(elsePos)
			if err != nil {
				return ir.TemplateNode{}, err
			}
			nested := ir.TemplateNode{Type: ir.NodeIf, Cond: elseCond, Children: elseIfChildren}
			if p.current().Type == tokens.TokenElse {
				p.advance()
				nested.ElseChildren, err = p.parseElseNodes()
				if err != nil {
					return ir.TemplateNode{}, err
				}
			}
			elseChildren = append(elseChildren, nested)
			if p.current().Type == tokens.TokenIfClose {
				p.advance()
				break
			}
		}
		if p.current().Type == tokens.TokenIfClose {
			p.advance()
		}
		return ir.TemplateNode{Type: ir.NodeIf, Cond: cond, Children: children, ElseChildren: elseChildren}, nil
	case tokens.TokenEachOpen:
		items, item, err := parseEachClause(tok.Value)
		if err != nil {
			return ir.TemplateNode{}, err
		}
		p.advance()
		children, err := p.parseEachNodes()
		if err != nil {
			return ir.TemplateNode{}, err
		}
		var elseChildren []ir.TemplateNode
		if p.current().Type == tokens.TokenEachElse {
			p.advance()
			elseChildren, err = p.parseEachElseNodes()
			if err != nil {
				return ir.TemplateNode{}, err
			}
		}
		if p.current().Type == tokens.TokenEachClose {
			p.advance()
		}
		return ir.TemplateNode{Type: ir.NodeEach, Items: items, Item: item, Children: children, ElseChildren: elseChildren}, nil
	case tokens.TokenTagClose:
		p.advance()
		return ir.TemplateNode{Type: ir.NodeText, Content: fmt.Sprintf("</%s>", tok.Tag), Pos: tok.Pos}, nil
	case tokens.TokenSlot:
		p.advance()
		return ir.TemplateNode{Type: ir.NodeSlot, Content: tok.Value}, nil
	case tokens.TokenSlotOpen:
		p.advance()
		children, err := p.parseSlotNodes()
		if err != nil {
			return ir.TemplateNode{}, err
		}
		return ir.TemplateNode{Type: ir.NodeSlot, Content: tok.Value, Children: children, Pos: tok.Pos}, nil
	case tokens.TokenTagOpen:
		if err := checkAttrControlFlow(tok.Attr, tok.Pos); err != nil {
			return ir.TemplateNode{}, err
		}
		p.advance()
		content := fmt.Sprintf("<%s", tok.Tag)
		if tok.Attr != "" {
			content += " " + tok.Attr
		}
		content += ">"
		return ir.TemplateNode{Type: ir.NodeText, Content: content, Pos: tok.Pos}, nil
	case tokens.TokenComponentSelfClose:
		if err := checkAttrControlFlow(tok.Attr, tok.Pos); err != nil {
			return ir.TemplateNode{}, err
		}
		p.advance()
		return ir.TemplateNode{Type: ir.NodeComponentCall, Tag: tok.Tag, Attrs: tok.Attr, SelfClose: true, Pos: tok.Pos}, nil
	case tokens.TokenComponentTagOpen:
		if err := checkAttrControlFlow(tok.Attr, tok.Pos); err != nil {
			return ir.TemplateNode{}, err
		}
		p.advance()
		children, err := p.parseComponentNodes(tok.Tag)
		if err != nil {
			return ir.TemplateNode{}, err
		}
		return ir.TemplateNode{Type: ir.NodeComponentCall, Tag: tok.Tag, Attrs: tok.Attr, Children: children, Pos: tok.Pos}, nil
	case tokens.TokenSlotClose:
		return ir.TemplateNode{}, fmt.Errorf("unexpected {/slot} at position %d", tok.Pos)
	case tokens.TokenElse:
		return ir.TemplateNode{}, fmt.Errorf("unexpected {#else} outside {#if} at position %d", tok.Pos)
	case tokens.TokenElseIf:
		return ir.TemplateNode{}, fmt.Errorf("unexpected {#else if} outside {#if} at position %d", tok.Pos)
	case tokens.TokenVerbatim:
		p.advance()
		return ir.TemplateNode{Type: ir.NodeVerbatim, Content: tok.Value}, nil
	default:
		return ir.TemplateNode{}, fmt.Errorf("unexpected token %s in template at position %d", tok.Type, tok.Pos)
	}
}

func (p *Parser) parseSlotNodes() ([]ir.TemplateNode, error) {
	var nodes []ir.TemplateNode
	for {
		tok := p.current()
		if tok.Type == tokens.TokenEOF {
			return nil, fmt.Errorf("unclosed {#slot}")
		}
		if tok.Type == tokens.TokenSlotClose {
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

func (p *Parser) parseComponentNodes(tag string) ([]ir.TemplateNode, error) {
	var nodes []ir.TemplateNode

	for {
		tok := p.current()
		if tok.Type == tokens.TokenEOF {
			return nil, fmt.Errorf("unclosed <@%s>", tag)
		}
		if tok.Type == tokens.TokenComponentTagClose && tok.Tag == tag {
			p.advance()
			return nodes, nil
		}
		if tok.Type == tokens.TokenComponentTagClose {
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
