package core

import (
	"strings"
)

func (p *Parser) parseStyleSection() (*StyleSection, error) {
	content, err := p.readUntilClose("style")
	if err != nil {
		return nil, err
	}
	return &StyleSection{Code: strings.TrimSpace(content)}, nil
}
