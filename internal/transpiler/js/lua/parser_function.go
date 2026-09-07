package lua

func (p *sourceParser) functionExpression() (expression, error) {
	if _, err := p.require(tokenLeftParen, "expected ( after function"); err != nil {
		return nil, err
	}
	var params []string
	if !p.match(tokenRightParen) {
		for {
			name, err := p.require(tokenIdentifier, "expected parameter name")
			if err != nil {
				return nil, err
			}
			params = append(params, name.value)
			if p.match(tokenRightParen) {
				break
			}
			if _, err := p.require(tokenComma, "expected , or ) after parameter"); err != nil {
				return nil, err
			}
		}
	}
	body, err := p.statements(map[tokenKind]bool{tokenEnd: true, tokenEOF: true})
	if err != nil {
		return nil, err
	}
	if _, err := p.require(tokenEnd, "expected end to close function"); err != nil {
		return nil, err
	}
	return functionExpression{params: params, body: body}, nil
}
