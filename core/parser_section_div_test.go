package core

import "testing"

// Tests for error paths in parseTemplateNode / parseDivNodes /
// parseComponentNodes / parseSlotNodes (core/parser_section_div.go).
// Each test lexes the input and drives the full parse chain via
// Lex + NewParser(tokens).Parse(), asserting the returned error.

func TestParseTemplateUnexpectedSlotClose(t *testing.T) {
	parseExpectError(t,
		"<div>{/slot}</div>",
		"unexpected {/slot}")
}

func TestParseTemplateElseOutsideIf(t *testing.T) {
	parseExpectError(t,
		"<div>{#else}</div>",
		"unexpected {#else} outside {#if}")
}

func TestParseTemplateUnexpectedToken(t *testing.T) {
	parseExpectError(t,
		"<div>{/if}</div>",
		"unexpected token IfClose in template")
}

func TestParseComponentMismatchedClose(t *testing.T) {
	parseExpectError(t,
		"<div><@a></@b></div>",
		"unexpected </@b>, expected </@a>")
}

func TestParseUnclosedDiv(t *testing.T) {
	parseExpectError(t,
		"<div>",
		"unclosed <div>")
}
