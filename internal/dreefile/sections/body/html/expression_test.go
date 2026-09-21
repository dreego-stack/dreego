package html

import (
	"strings"
	"testing"
)

func TestConditionCodeMatrix(t *testing.T) {
	cases := []struct {
		name string
		cond string
		want string
	}{
		{"bool variable", "show", "dreego.Truthy(show)"},
		{"negation", "!show", "dreego.Truthy(!show)"},
		{"comparison", "score >= 80", "dreego.Truthy(score >= 80)"},
		{"string method", `c.Errors("email")`, `dreego.Truthy(c.Errors("email"))`},
		{"slice length", "len(items) > 0", "dreego.Truthy(len(items) > 0)"},
		{"init statement", `v := c.Get("x"); v != ""`, `v := c.Get("x"); dreego.Truthy(v != "")`},
		{"semicolon in string", `label == "a;b"`, `dreego.Truthy(label == "a;b")`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := conditionCode(c.cond); got != c.want {
				t.Errorf("conditionCode(%q) = %q, want %q", c.cond, got, c.want)
			}
		})
	}
}

func TestExpressionCodeMatrix(t *testing.T) {
	cases := []struct {
		name    string
		expr    string
		filters []string
		want    string
		wantRaw bool
		wantErr string
	}{
		{"no filter", "v", nil, `fmt.Sprintf("%v", v)`, false, ""},
		{"raw", "v", []string{"raw"}, `fmt.Sprintf("%v", v)`, true, ""},
		{"upper", "v", []string{"upper"}, `strings.ToUpper(fmt.Sprintf("%v", v))`, false, ""},
		{"raw upper", "v", []string{"raw", "upper"}, `strings.ToUpper(fmt.Sprintf("%v", v))`, true, ""},
		{"unknown", "v", []string{"nope"}, "", false, "unknown filter 'nope' at position 3"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			code, raw, err := expressionCode(c.expr, c.filters, 3)
			if c.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("expected error %q, got %v", c.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if code != c.want || raw != c.wantRaw {
				t.Errorf("expressionCode = (%q, %v), want (%q, %v)", code, raw, c.want, c.wantRaw)
			}
		})
	}
}
