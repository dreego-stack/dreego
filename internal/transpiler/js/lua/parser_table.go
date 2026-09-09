package lua

import "strconv"

func (p *sourceParser) tableExpression() (expression, error) {
	var fields []tableField
	nextIndex := 1
	if p.match(tokenRightBrace) {
		return tableExpression{}, nil
	}
	for {
		field, implicit, err := p.tableField(nextIndex)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
		if implicit {
			nextIndex++
		}
		if p.match(tokenRightBrace) {
			return tableExpression{fields: fields}, nil
		}
		if !p.match(tokenComma) && !p.match(tokenSeparator) {
			return nil, p.errorf(p.current(), "expected , or ; or } in table")
		}
		if p.match(tokenRightBrace) {
			return tableExpression{fields: fields}, nil
		}
	}
}

func (p *sourceParser) tableField(nextIndex int) (tableField, bool, error) {
	if p.match(tokenLeftBracket) {
		key, err := p.expression(0)
		if err != nil {
			return tableField{}, false, err
		}
		if _, err := p.require(tokenRightBracket, "expected ] after table key"); err != nil {
			return tableField{}, false, err
		}
		if _, err := p.require(tokenAssign, "expected = after table key"); err != nil {
			return tableField{}, false, err
		}
		value, err := p.expression(0)
		return tableField{key: key, value: value}, false, err
	}
	if p.current().kind == tokenIdentifier && p.peekKind(1) == tokenAssign {
		name := p.advance().value
		p.advance()
		value, err := p.expression(0)
		return tableField{key: literalExpression{kind: tokenString, value: name}, value: value}, false, err
	}
	value, err := p.expression(0)
	return tableField{key: literalExpression{kind: tokenNumber, value: strconv.Itoa(nextIndex)}, value: value}, true, err
}

func (p *sourceParser) peekKind(offset int) tokenKind {
	if p.index+offset >= len(p.tokens) {
		return tokenEOF
	}
	return p.tokens[p.index+offset].kind
}
