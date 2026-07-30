package core

import (
	"strings"
)

func (p *Parser) parseScriptSection() (*ScriptSection, error) {
	content, err := p.readUntilClose("script")
	if err != nil {
		return nil, err
	}
	return &ScriptSection{Code: strings.TrimSpace(content)}, nil
}
