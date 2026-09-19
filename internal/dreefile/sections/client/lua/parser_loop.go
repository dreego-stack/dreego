package lua

func (p *sourceParser) whileStatement() (statement, error) {
	condition, err := p.expression(0)
	if err != nil {
		return nil, err
	}
	if _, err := p.require(tokenDo, "expected do after while condition"); err != nil {
		return nil, err
	}
	p.loopDepth++
	body, err := p.statements(map[tokenKind]bool{tokenEnd: true, tokenEOF: true})
	p.loopDepth--
	if err != nil {
		return nil, err
	}
	if _, err := p.require(tokenEnd, "expected end to close while"); err != nil {
		return nil, err
	}
	return whileStatement{condition: condition, body: body}, nil
}

func (p *sourceParser) numericForStatement() (statement, error) {
	name, err := p.require(tokenIdentifier, "expected loop variable after for")
	if err != nil {
		return nil, err
	}
	if _, err := p.require(tokenAssign, "expected = after numeric for variable"); err != nil {
		return nil, err
	}
	initial, err := p.expression(0)
	if err != nil {
		return nil, err
	}
	if _, err := p.require(tokenComma, "expected , after numeric for initial value"); err != nil {
		return nil, err
	}
	limit, err := p.expression(0)
	if err != nil {
		return nil, err
	}
	step := expression(literalExpression{kind: tokenNumber, value: "1"})
	if p.match(tokenComma) {
		step, err = p.expression(0)
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.require(tokenDo, "expected do after numeric for values"); err != nil {
		return nil, err
	}
	p.loopDepth++
	body, err := p.statements(map[tokenKind]bool{tokenEnd: true, tokenEOF: true})
	p.loopDepth--
	if err != nil {
		return nil, err
	}
	if _, err := p.require(tokenEnd, "expected end to close numeric for"); err != nil {
		return nil, err
	}
	return numericForStatement{name: name.value, initial: initial, limit: limit, step: step, body: body}, nil
}
