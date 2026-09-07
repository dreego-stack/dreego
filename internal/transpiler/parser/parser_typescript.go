package parser

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
	"github.com/dreego-stack/dreego/internal/transpiler/tokens"
)

func (p *Parser) parseTypeScriptNode(open tokens.Token) (ir.TemplateNode, error) {
	p.advance()
	contentPos := p.current().Pos
	var code strings.Builder
	for {
		tok := p.current()
		switch {
		case tok.Type == tokens.TokenEOF:
			return ir.TemplateNode{}, fmt.Errorf("unclosed <script lang=\"ts\"> at position %d", open.Pos)
		case tok.Type == tokens.TokenTagClose && tok.Tag == "script":
			p.advance()
			raw := code.String()
			trimmed := strings.TrimSpace(raw)
			contentPos += strings.Index(raw, trimmed)
			return ir.TemplateNode{Type: ir.NodeClientScript, Content: trimmed, Language: "ts", Pos: contentPos}, nil
		case tok.Type == tokens.TokenText:
			code.WriteString(tok.Value)
			p.advance()
		default:
			return ir.TemplateNode{}, fmt.Errorf("unexpected token %s in TypeScript at position %d", tok.Type, tok.Pos)
		}
	}
}
