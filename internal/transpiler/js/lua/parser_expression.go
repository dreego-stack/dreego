package lua

import "strconv"

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
	case tokenMinus, tokenNot, tokenHash:
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
		var err error
		value, err = p.tableExpression()
		if err != nil {
			return nil, err
		}
	case tokenFunction:
		var err error
		value, err = p.functionExpression()
		if err != nil {
			return nil, err
		}
	default:
		return nil, p.errorf(current, "expected an expression")
	}
	return p.postfix(value)
}

func (p *sourceParser) postfix(value expression) (expression, error) {
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
		case p.match(tokenLeftBracket):
			index, err := p.expression(0)
			if err != nil {
				return nil, err
			}
			if _, err := p.require(tokenRightBracket, "expected ] after table index"); err != nil {
				return nil, err
			}
			value = indexExpression{object: value, index: index}
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
