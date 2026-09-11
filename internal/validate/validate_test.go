package validate

import (
	"net/http"
	"net/url"
	"testing"
)

func TestApplyRuleRequired(t *testing.T) {
	if applyRule("required", "") == "" {
		t.Error("expected error for empty value")
	}
	if applyRule("required", "x") != "" {
		t.Error("expected no error for non-empty value")
	}
}

func TestApplyRuleEmail(t *testing.T) {
	if applyRule("email", "bad") == "" {
		t.Error("expected error for invalid email")
	}
	if applyRule("email", "a@b.c") != "" {
		t.Error("expected no error for valid email")
	}
}

func TestApplyRuleMin(t *testing.T) {
	if applyRule("min=3", "ab") == "" {
		t.Error("expected error for too-short value")
	}
	if applyRule("min=3", "abc") != "" {
		t.Error("expected no error for long-enough value")
	}
}

func TestApplyRuleMax(t *testing.T) {
	if applyRule("max=3", "abcd") == "" {
		t.Error("expected error for too-long value")
	}
	if applyRule("max=3", "abc") != "" {
		t.Error("expected no error for short-enough value")
	}
}

func TestApplyRuleUnknown(t *testing.T) {
	if applyRule("bogus", "x") != "" {
		t.Error("expected no error for unknown rule")
	}
}

func TestAtoi(t *testing.T) {
	if n, err := Atoi("42"); n != 42 || err != nil {
		t.Errorf("expected 42, got %d, err %v", n, err)
	}
	if n, err := Atoi("0"); n != 0 || err != nil {
		t.Errorf("expected 0, got %d, err %v", n, err)
	}
}

type validateFormRequired struct {
	Name string `validate:"required"`
}

type validateFormEmail struct {
	Mail string `validate:"email"`
}

type validateFormMulti struct {
	Name  string `validate:"required"`
	Email string `validate:"email"`
}

type validateFormNoTag struct {
	Name string
}

func TestValidateFormRequired(t *testing.T) {
	f := validateFormRequired{Name: ""}
	errs := ValidateForm(f)
	if errs == nil || errs["name"] == "" {
		t.Error("expected required error")
	}
	f.Name = "x"
	errs = ValidateForm(f)
	if errs != nil {
		t.Error("expected no errors")
	}
}

func TestValidateFormEmail(t *testing.T) {
	f := validateFormEmail{Mail: "bad"}
	errs := ValidateForm(f)
	if errs == nil || errs["mail"] == "" {
		t.Error("expected email error")
	}
	f.Mail = "a@b.c"
	errs = ValidateForm(f)
	if errs != nil {
		t.Error("expected no errors")
	}
}

func TestValidateFormMultipleRules(t *testing.T) {
	f := validateFormMulti{Name: "", Email: "bad"}
	errs := ValidateForm(f)
	if errs == nil {
		t.Fatal("expected errors")
	}
	if errs["name"] == "" {
		t.Error("expected required error on name")
	}
	if errs["email"] == "" {
		t.Error("expected email error on email")
	}
}

func TestValidateFormNoTag(t *testing.T) {
	f := validateFormNoTag{Name: ""}
	errs := ValidateForm(f)
	if errs != nil {
		t.Error("expected no errors for field without validate tag")
	}
}

type bindFormStruct struct {
	Name  string `form:"name"`
	Email string
}

func TestBindFormStringField(t *testing.T) {
	f := bindFormStruct{}
	r := &http.Request{Form: url.Values{"name": {"Alice"}}}
	BindForm(r, &f)
	if f.Name != "Alice" {
		t.Errorf("expected 'Alice', got '%s'", f.Name)
	}
}

func TestBindFormFallbackToLowercase(t *testing.T) {
	f := bindFormStruct{}
	r := &http.Request{Form: url.Values{"email": {"a@b.c"}}}
	BindForm(r, &f)
	if f.Email != "a@b.c" {
		t.Errorf("expected 'a@b.c', got '%s'", f.Email)
	}
}

type bindFormNonString struct {
	Count map[string]string
}

func TestBindFormNonStringFieldReturnsError(t *testing.T) {
	f := bindFormNonString{Count: map[string]string{"x": "y"}}
	r := &http.Request{Form: url.Values{"count": {"99"}}}
	err := BindForm(r, &f)
	if err == nil {
		t.Error("expected error for non-string field (B2)")
	}
}

func TestSaveOld(t *testing.T) {
	c := &testSetter{}
	f := bindFormStruct{Name: "Alice", Email: "a@b.c"}
	SaveOld(c, f)
	if c.data["old_name"] != "Alice" {
		t.Errorf("expected 'Alice', got '%v'", c.data["old_name"])
	}
	if c.data["old_email"] != "a@b.c" {
		t.Errorf("expected 'a@b.c', got '%v'", c.data["old_email"])
	}
}

func TestSaveErrors(t *testing.T) {
	c := &testSetter{}
	errs := map[string]string{"name": "is required", "email": "must be valid"}
	SaveErrors(c, errs)
	if c.data["error_name"] != "is required" {
		t.Errorf("expected 'is required', got '%v'", c.data["error_name"])
	}
	if c.data["error_email"] != "must be valid" {
		t.Errorf("expected 'must be valid', got '%v'", c.data["error_email"])
	}
}

func TestValidateFormPtr(t *testing.T) {
	f := &validateFormRequired{Name: ""}
	errs := ValidateForm(f)
	if errs == nil || errs["name"] == "" {
		t.Error("expected required error for pointer form")
	}
}

func TestBindFormEmptyValue(t *testing.T) {
	f := bindFormStruct{Name: "original"}
	r := &http.Request{Form: url.Values{"name": {""}}}
	BindForm(r, &f)
	if f.Name != "original" {
		t.Errorf("expected 'original', got '%s'", f.Name)
	}
}

func TestApplyRuleMinNonNumeric(t *testing.T) {
	if applyRule("min=abc", "x") == "" {
		t.Error("expected error when min value is non-numeric")
	}
}

func TestApplyRuleMaxNonNumeric(t *testing.T) {
	if applyRule("max=abc", "x") == "" {
		t.Error("expected error when max value is non-numeric")
	}
}

func TestValidateFormFirstRuleBreaks(t *testing.T) {
	f := validateFormMulti{Name: "", Email: "bad"}
	errs := ValidateForm(f)
	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d", len(errs))
	}
}

func TestBindFormParseError(t *testing.T) {
	f := bindFormStruct{}
	r := &http.Request{Method: "POST", Header: http.Header{"Content-Type": {"application/x-www-form-urlencoded"}}, Body: http.NoBody}
	r.Body = nil
	err := BindForm(r, &f)
	if err == nil {
		t.Error("expected error for nil body")
	}
}

func TestValidateFormEmptyStruct(t *testing.T) {
	type empty struct{}
	f := empty{}
	errs := ValidateForm(f)
	if errs != nil {
		t.Error("expected nil for empty struct")
	}
}

func TestAtoiEmpty(t *testing.T) {
	if n, err := Atoi(""); n != 0 || err == nil {
		t.Errorf("expected error for empty string, got %d, err %v", n, err)
	}
}

func TestAtoiNonDigits(t *testing.T) {
	if n, err := Atoi("abc"); n != 0 || err == nil {
		t.Errorf("expected error for non-digits, got %d, err %v", n, err)
	}
}

func TestApplyRuleRequiredWhitespace(t *testing.T) {
	if applyRule("required", "   ") == "" {
		t.Error("expected error for whitespace-only value")
	}
}

func TestApplyRuleEmailNoAt(t *testing.T) {
	if applyRule("email", "a.b") == "" {
		t.Error("expected error for email without @")
	}
}

func TestApplyRuleEmailNoDot(t *testing.T) {
	if applyRule("email", "a@b") == "" {
		t.Error("expected error for email without dot")
	}
}

func TestValidateFormMultipleRulesFirstPasses(t *testing.T) {
	type form struct {
		Field string `validate:"required,email"`
	}
	f := form{Field: "x"}
	errs := ValidateForm(f)
	if errs == nil || errs["field"] == "" {
		t.Error("expected email error when required passes but email fails")
	}
}

func TestBindFormMultipleFields(t *testing.T) {
	f := bindFormStruct{}
	r := &http.Request{Form: url.Values{"name": {"Alice"}, "email": {"a@b.c"}}}
	BindForm(r, &f)
	if f.Name != "Alice" {
		t.Errorf("expected 'Alice', got '%s'", f.Name)
	}
	if f.Email != "a@b.c" {
		t.Errorf("expected 'a@b.c', got '%s'", f.Email)
	}
}
