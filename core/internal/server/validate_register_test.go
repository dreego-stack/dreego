package server

import (
	"testing"

	"github.com/dreego-stack/dreego/core/internal/validate"
)

// typed-forms.1: custom validator registration + application. Currently RED —
// no RegisterRule/registered-rule dispatch exists in applyRule.
func TestRegisterRuleCustom(t *testing.T) {
	app1 := New()
	if err := app1.RegisterRule("even", func(val string) string {
		n, err := validate.Atoi(val)
		if err != nil {
			return "must be a number"
		}
		if n%2 != 0 {
			return "must be even"
		}
		return ""
	}); err != nil {
		t.Fatal(err)
	}

	type form struct {
		Num string `validate:"even"`
	}

	bad := form{Num: "3"}
	errs := app1.ValidateForm(bad)
	if errs == nil || errs["num"] == "" {
		t.Error("expected custom validator to reject odd value")
	}

	good := form{Num: "4"}
	errs = app1.ValidateForm(good)
	if errs != nil {
		t.Errorf("expected no errors for even value, got %v", errs)
	}

	app2 := New()
	if errs := app2.ValidateForm(bad); errs != nil {
		t.Errorf("custom rule from app1 leaked into app2: %v", errs)
	}
}

// typed-forms.1: custom validator not clobbered by built-ins. Regression guard.
func TestRegisterRuleCustomDoesNotBreakBuiltins(t *testing.T) {
	type form struct {
		Name string `validate:"required"`
	}
	f := form{Name: ""}
	errs := validate.ValidateForm(f)
	if errs == nil || errs["name"] == "" {
		t.Error("expected required error still to work alongside custom rules")
	}
}
