package ir

import "testing"

func TestExprKind(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want ExprKind
	}{
		{name: "double quoted string", expr: `"hello"`, want: ExprKindStringLiteral},
		{name: "raw backtick string", expr: "`hello`", want: ExprKindStringLiteral},
		{name: "rune", expr: `'x'`, want: ExprKindStringLiteral},
		{name: "positive int", expr: `42`, want: ExprKindIntLiteral},
		{name: "negative int", expr: `-42`, want: ExprKindIntLiteral},
		{name: "hex", expr: `0x2a`, want: ExprKindIntLiteral},
		{name: "binary", expr: `0b1010`, want: ExprKindIntLiteral},
		{name: "octal", expr: `0o52`, want: ExprKindIntLiteral},
		{name: "float", expr: `3.14`, want: ExprKindOther},
		{name: "bool", expr: `true`, want: ExprKindOther},
		{name: "non-literal expression", expr: `value + 1`, want: ExprKindOther},
		{name: "empty", expr: ``, want: ExprKindOther},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyExpression(tt.expr)
			if got != tt.want {
				t.Errorf("ClassifyExpression(%q) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}
}
