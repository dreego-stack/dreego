package core

import "testing"

func TestExprKind(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want exprKind
	}{
		{name: "double quoted string", expr: `"hello"`, want: exprKindStringLiteral},
		{name: "raw backtick string", expr: "`hello`", want: exprKindStringLiteral},
		{name: "rune", expr: `'x'`, want: exprKindStringLiteral},
		{name: "positive int", expr: `42`, want: exprKindIntLiteral},
		{name: "negative int", expr: `-42`, want: exprKindIntLiteral},
		{name: "hex", expr: `0x2a`, want: exprKindIntLiteral},
		{name: "binary", expr: `0b1010`, want: exprKindIntLiteral},
		{name: "octal", expr: `0o52`, want: exprKindIntLiteral},
		{name: "float", expr: `3.14`, want: exprKindOther},
		{name: "bool", expr: `true`, want: exprKindOther},
		{name: "non-literal expression", expr: `value + 1`, want: exprKindOther},
		{name: "empty", expr: ``, want: exprKindOther},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := classifyExpression(tt.expr)
			if got != tt.want {
				t.Errorf("classifyExpression(%q) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}
}
