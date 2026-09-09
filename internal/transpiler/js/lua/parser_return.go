package lua

func (p *sourceParser) returnStatement() (statement, error) {
	if p.functionDepth == 0 {
		return nil, p.errorf(p.previous(), "return is only valid inside a function")
	}
	switch p.current().kind {
	case tokenSeparator, tokenEnd, tokenElseIf, tokenElse, tokenEOF:
		return returnStatement{}, nil
	}
	value, err := p.expression(0)
	if err != nil {
		return nil, err
	}
	if p.current().kind == tokenComma {
		return nil, p.errorf(p.current(), "multiple return values are not supported")
	}
	return returnStatement{value: value}, nil
}
