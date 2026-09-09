package lua

import (
	"fmt"
)

type sourceParser struct {
	tokens        []token
	index         int
	functionDepth int
	loopDepth     int
}

func parse(tokens []token) (program, error) {
	parser := sourceParser{tokens: tokens}
	statements, err := parser.statements(map[tokenKind]bool{tokenEOF: true})
	return program{statements: statements}, err
}

func (p *sourceParser) statements(stop map[tokenKind]bool) ([]statement, error) {
	var statements []statement
	for {
		p.separators()
		if stop[p.current().kind] {
			return statements, nil
		}
		value, err := p.statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, value)
		switch value.(type) {
		case returnStatement:
			p.separators()
			if !stop[p.current().kind] {
				return nil, p.errorf(p.current(), "return must be the final statement in its block")
			}
			return statements, nil
		case breakStatement:
			p.separators()
			if !stop[p.current().kind] {
				return nil, p.errorf(p.current(), "break must be the final statement in its block")
			}
			return statements, nil
		}
	}
}

func (p *sourceParser) statement() (statement, error) {
	if p.match(tokenReturn) {
		return p.returnStatement()
	}
	if p.match(tokenBreak) {
		if p.loopDepth == 0 {
			return nil, p.errorf(p.previous(), "break is only valid inside a loop")
		}
		return breakStatement{}, nil
	}
	if p.match(tokenLocal) {
		if p.match(tokenFunction) {
			name, err := p.require(tokenIdentifier, "expected a name after local function")
			if err != nil {
				return nil, err
			}
			value, err := p.functionExpression()
			return localStatement{name: name.value, value: value, recursive: true}, err
		}
		name, err := p.require(tokenIdentifier, "expected a name after local")
		if err != nil {
			return nil, err
		}
		if _, err := p.require(tokenAssign, "expected = in local declaration"); err != nil {
			return nil, err
		}
		value, err := p.expression(0)
		return localStatement{name: name.value, value: value}, err
	}
	if p.match(tokenIf) {
		return p.ifStatement()
	}
	if p.match(tokenWhile) {
		return p.whileStatement()
	}
	if p.match(tokenFor) {
		return p.numericForStatement()
	}
	value, err := p.expression(0)
	if err != nil {
		return nil, err
	}
	if p.match(tokenAssign) {
		if !assignable(value) {
			return nil, p.errorf(p.previous(), "invalid assignment target")
		}
		right, err := p.expression(0)
		return assignStatement{target: value, value: right}, err
	}
	if _, ok := value.(callExpression); !ok {
		return nil, p.errorf(p.previous(), "only function calls may be used as expression statements")
	}
	return expressionStatement{value: value}, nil
}

func (p *sourceParser) ifStatement() (statement, error) {
	condition, err := p.expression(0)
	if err != nil {
		return nil, err
	}
	if _, err := p.require(tokenThen, "expected then after if condition"); err != nil {
		return nil, err
	}
	body, err := p.statements(map[tokenKind]bool{tokenElseIf: true, tokenElse: true, tokenEnd: true, tokenEOF: true})
	if err != nil {
		return nil, err
	}
	result := ifStatement{branches: []ifBranch{{condition: condition, body: body}}}
	for p.match(tokenElseIf) {
		condition, err = p.expression(0)
		if err != nil {
			return nil, err
		}
		if _, err := p.require(tokenThen, "expected then after elseif condition"); err != nil {
			return nil, err
		}
		body, err = p.statements(map[tokenKind]bool{tokenElseIf: true, tokenElse: true, tokenEnd: true, tokenEOF: true})
		if err != nil {
			return nil, err
		}
		result.branches = append(result.branches, ifBranch{condition: condition, body: body})
	}
	if p.match(tokenElse) {
		result.otherwise, err = p.statements(map[tokenKind]bool{tokenEnd: true, tokenEOF: true})
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.require(tokenEnd, "expected end to close if"); err != nil {
		return nil, err
	}
	return result, nil
}

func (p *sourceParser) separators() {
	for p.match(tokenSeparator) {
	}
}

func (p *sourceParser) require(kind tokenKind, message string) (token, error) {
	if p.current().kind != kind {
		return token{}, p.errorf(p.current(), message)
	}
	return p.advance(), nil
}

func (p *sourceParser) match(kind tokenKind) bool {
	if p.current().kind != kind {
		return false
	}
	p.index++
	return true
}

func (p *sourceParser) advance() token {
	value := p.current()
	if value.kind != tokenEOF {
		p.index++
	}
	return value
}

func (p *sourceParser) current() token  { return p.tokens[p.index] }
func (p *sourceParser) previous() token { return p.tokens[p.index-1] }

func (p *sourceParser) errorf(at token, format string, values ...any) error {
	return fmt.Errorf("Lua %d:%d: %s", at.line, at.column, fmt.Sprintf(format, values...))
}
