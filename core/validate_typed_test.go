package core

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

// typed-forms.1: int binding. Currently RED — BindForm only supports reflect.String.
func TestBindFormIntField(t *testing.T) {
	type form struct {
		Count int
	}
	f := form{}
	r := &http.Request{Form: url.Values{"count": {"42"}}}
	if err := BindForm(r, &f); err != nil {
		t.Fatalf("BindForm returned error for int field: %v", err)
	}
	if f.Count != 42 {
		t.Errorf("expected 42, got %d", f.Count)
	}
}

// typed-forms.1: bool binding (checkbox "on"). Currently RED.
func TestBindFormBoolFieldOn(t *testing.T) {
	type form struct {
		Remember bool
	}
	f := form{}
	r := &http.Request{Form: url.Values{"remember": {"on"}}}
	if err := BindForm(r, &f); err != nil {
		t.Fatalf("BindForm returned error for bool field: %v", err)
	}
	if !f.Remember {
		t.Error("expected Remember=true for 'on' checkbox value")
	}
}

// typed-forms.1: bool binding — absent/empty value stays false. Currently RED.
func TestBindFormBoolFieldAbsent(t *testing.T) {
	type form struct {
		Remember bool
	}
	f := form{Remember: false}
	r := &http.Request{Form: url.Values{}}
	if err := BindForm(r, &f); err != nil {
		t.Fatalf("BindForm returned error: %v", err)
	}
	if f.Remember {
		t.Error("expected Remember=false when checkbox absent")
	}
}

// typed-forms.1: []string slice binding — multiple values. Currently RED.
func TestBindFormSliceField(t *testing.T) {
	type form struct {
		Tags []string
	}
	f := form{}
	r := &http.Request{Form: url.Values{"tags": {"go", "htmx", "alpine"}}}
	if err := BindForm(r, &f); err != nil {
		t.Fatalf("BindForm returned error for slice field: %v", err)
	}
	want := []string{"go", "htmx", "alpine"}
	if !reflect.DeepEqual(f.Tags, want) {
		t.Errorf("expected %v, got %v", want, f.Tags)
	}
}

// typed-forms.1: custom validator registration + application. Currently RED —
// no RegisterRule/registered-rule dispatch exists in applyRule.
func TestRegisterRuleCustom(t *testing.T) {
	app := New()
	app.RegisterRule("even", func(val string) string {
		n, err := atoi(val)
		if err != nil {
			return "must be a number"
		}
		if n%2 != 0 {
			return "must be even"
		}
		return ""
	})

	type form struct {
		Num string `validate:"even"`
	}

	bad := form{Num: "3"}
	errs := ValidateForm(bad)
	if errs == nil || errs["num"] == "" {
		t.Error("expected custom validator to reject odd value")
	}

	good := form{Num: "4"}
	errs = ValidateForm(good)
	if errs != nil {
		t.Errorf("expected no errors for even value, got %v", errs)
	}
}

// typed-forms.1: custom validator not clobbered by built-ins. Regression guard.
func TestRegisterRuleCustomDoesNotBreakBuiltins(t *testing.T) {
	type form struct {
		Name string `validate:"required"`
	}
	f := form{Name: ""}
	errs := ValidateForm(f)
	if errs == nil || errs["name"] == "" {
		t.Error("expected required error still to work alongside custom rules")
	}
}

// typed-forms.1: email rejects missing @ (already works). GREEN regression.
func TestEmailRejectsMissingAt(t *testing.T) {
	if applyRule("email", "user.example.com") == "" {
		t.Error("expected email error for address without @")
	}
}

// typed-forms.1: email rejects missing dot (already works). GREEN regression.
func TestEmailRejectsMissingDot(t *testing.T) {
	if applyRule("email", "user@example") == "" {
		t.Error("expected email error for address without dot")
	}
}

// typed-forms.1: email accepts valid address (already works). GREEN regression.
func TestEmailAcceptsValid(t *testing.T) {
	if applyRule("email", "user@example.com") != "" {
		t.Error("expected no error for valid email")
	}
}
