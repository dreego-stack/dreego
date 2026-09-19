package parser

import "testing"

func TestDoctypeAtFileStartHintsBodyPlacement(t *testing.T) {
	parseExpectError(t,
		"<!DOCTYPE html>\n<body><p>hi</p></body>",
		"<body> section")
}

func TestHtmlAtFileStartHintsBodyPlacement(t *testing.T) {
	parseExpectError(t,
		"<html lang=\"en\">\n<body><p>hi</p></body>",
		"<body> section")
}
