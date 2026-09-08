package lua

import (
	"fmt"
	"strconv"
)

type sourceParser struct {
	tokens        []token
	index         int
	functionDepth int
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
		if _, ok := value.(returnStatement); ok {
			p.separators()
			if !stop[p.current().kind] {
				return nil, p.errorf(p.current(), "return must be the final statement in its block")
			}
			return statements, nil
		}
	}
}

func (p *sourceParser) statement() (statement, error) {
	if p.match(tokenReturn) {
		return p.returnStatement()
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

func (p *sourceParser) expression(minimum int) (expression, error) {
	left, err := p.prefix()
	if err != nil {
		return nil, err
	}
	for {
		precedence := binaryPrecedence(p.current().kind)
		if precedence < minimum {
			break
		}
		operator := p.advance().kind
		next := precedence + 1
		if operator == tokenConcat {
			next = precedence
		}
		right, err := p.expression(next)
		if err != nil {
			return nil, err
		}
		left = binaryExpression{left: left, op: operator, right: right}
	}
	return left, nil
}

func (p *sourceParser) prefix() (expression, error) {
	current := p.advance()
	var value expression
	switch current.kind {
	case tokenNumber:
		if _, err := strconv.ParseFloat(current.value, 64); err != nil {
			return nil, p.errorf(current, "invalid number %q", current.value)
		}
		value = literalExpression{kind: current.kind, value: current.value}
	case tokenString, tokenTrue, tokenFalse, tokenNil:
		value = literalExpression{kind: current.kind, value: current.value}
	case tokenIdentifier:
		value = nameExpression{name: current.value}
	case tokenMinus, tokenNot:
		right, err := p.expression(7)
		if err != nil {
			return nil, err
		}
		value = unaryExpression{op: current.kind, right: right}
	case tokenLeftParen:
		var err error
		value, err = p.expression(0)
		if err != nil {
			return nil, err
		}
		if _, err := p.require(tokenRightParen, "expected )"); err != nil {
			return nil, err
		}
	case tokenLeftBrace:
		return nil, p.errorf(current, "tables are not supported in the Lua MVP")
	case tokenFunction:
		var err error
		value, err = p.functionExpression()
		if err != nil {
			return nil, err
		}
	default:
		return nil, p.errorf(current, "expected an expression")
	}
	for {
		switch {
		case p.match(tokenDot):
			field, err := p.require(tokenIdentifier, "expected field name after .")
			if err != nil {
				return nil, err
			}
			value = memberExpression{object: value, field: field.value}
		case p.match(tokenLeftParen):
			args, err := p.arguments()
			if err != nil {
				return nil, err
			}
			value = callExpression{callee: value, args: args}
		case p.match(tokenColon):
			field, err := p.require(tokenIdentifier, "expected method name after :")
			if err != nil {
				return nil, err
			}
			if _, err := p.require(tokenLeftParen, "expected ( after method name"); err != nil {
				return nil, err
			}
			args, err := p.arguments()
			if err != nil {
				return nil, err
			}
			value = callExpression{callee: memberExpression{object: value, field: field.value}, args: args}
		default:
			return value, nil
		}
	}
}

func (p *sourceParser) arguments() ([]expression, error) {
	var result []expression
	if p.match(tokenRightParen) {
		return result, nil
	}
	for {
		value, err := p.expression(0)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
		if p.match(tokenRightParen) {
			return result, nil
		}
		if _, err := p.require(tokenComma, "expected , or ) in arguments"); err != nil {
			return nil, err
		}
	}
}

func binaryPrecedence(kind tokenKind) int {
	switch kind {
	case tokenOr:
		return 1
	case tokenAnd:
		return 2
	case tokenEqual, tokenNotEqual, tokenLess, tokenLessEqual, tokenGreater, tokenGreaterEqual:
		return 3
	case tokenConcat:
		return 4
	case tokenPlus, tokenMinus:
		return 5
	case tokenStar, tokenSlash, tokenPercent:
		return 6
	default:
		return -1
	}
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
