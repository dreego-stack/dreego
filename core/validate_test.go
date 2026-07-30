package core

import (
	"net/http"
	"net/url"
	"strings"
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
	if atoi("42") != 42 {
		t.Errorf("expected 42, got %d", atoi("42"))
	}
	if atoi("0") != 0 {
		t.Errorf("expected 0, got %d", atoi("0"))
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
	Count int
}

func TestBindFormNonStringFieldReturnsError(t *testing.T) {
	f := bindFormNonString{Count: 5}
	r := &http.Request{Form: url.Values{"count": {"99"}}}
	err := BindForm(r, &f)
	if err == nil {
		t.Error("expected error for non-string field (B2)")
	}
}

func TestSaveOld(t *testing.T) {
	c := &SSRContext{data: map[string]any{}}
	f := bindFormStruct{Name: "Alice", Email: "a@b.c"}
	SaveOld(c, f)
	if c.Get("old_name") != "Alice" {
		t.Errorf("expected 'Alice', got '%v'", c.Get("old_name"))
	}
	if c.Get("old_email") != "a@b.c" {
		t.Errorf("expected 'a@b.c', got '%v'", c.Get("old_email"))
	}
}

func TestSaveErrors(t *testing.T) {
	c := &SSRContext{data: map[string]any{}}
	errs := map[string]string{"name": "is required", "email": "must be valid"}
	SaveErrors(c, errs)
	if c.Get("error_name") != "is required" {
		t.Errorf("expected 'is required', got '%v'", c.Get("error_name"))
	}
	if c.Get("error_email") != "must be valid" {
		t.Errorf("expected 'must be valid', got '%v'", c.Get("error_email"))
	}
}

func TestFieldError(t *testing.T) {
	fe := FieldError{Field: "name", Message: "is required"}
	if fe.Field != "name" {
		t.Errorf("expected 'name', got '%s'", fe.Field)
	}
	if fe.Message != "is required" {
		t.Errorf("expected 'is required', got '%s'", fe.Message)
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
	if applyRule("min=abc", "x") != "" {
		t.Error("expected no error when min value is non-numeric")
	}
}

func TestApplyRuleMaxNonNumeric(t *testing.T) {
	if applyRule("max=abc", "x") == "" {
		t.Error("expected error when max value is non-numeric (atoi returns 0)")
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
	if atoi("") != 0 {
		t.Errorf("expected 0, got %d", atoi(""))
	}
}

func TestAtoiNonDigits(t *testing.T) {
	if atoi("abc") != 0 {
		t.Errorf("expected 0, got %d", atoi("abc"))
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

func TestSaveOldWithFormTag(t *testing.T) {
	c := &SSRContext{data: map[string]any{}}
	f := bindFormStruct{Name: "Bob"}
	SaveOld(c, f)
	if c.Get("old_name") != "Bob" {
		t.Errorf("expected 'Bob', got '%v'", c.Get("old_name"))
	}
}

func TestSaveErrorsEmpty(t *testing.T) {
	c := &SSRContext{data: map[string]any{}}
	SaveErrors(c, map[string]string{})
	if len(c.data) != 0 {
		t.Error("expected no data for empty errors")
	}
}

func TestValidateFormTagFallback(t *testing.T) {
	type form struct {
		UserName string `validate:"required"`
	}
	f := form{UserName: ""}
	errs := ValidateForm(f)
	if errs == nil || errs["username"] == "" {
		t.Error("expected required error with lowercase field name fallback")
	}
}

func TestBindFormTagFallback(t *testing.T) {
	type form struct {
		FullName string
	}
	f := form{}
	r := &http.Request{Form: url.Values{"fullname": {"Charlie"}}}
	BindForm(r, &f)
	if f.FullName != "Charlie" {
		t.Errorf("expected 'Charlie', got '%s'", f.FullName)
	}
}

func TestApplyRuleMinExact(t *testing.T) {
	if applyRule("min=3", "abc") != "" {
		t.Error("expected no error for exact min length")
	}
}

func TestApplyRuleMaxExact(t *testing.T) {
	if applyRule("max=3", "abc") != "" {
		t.Error("expected no error for exact max length")
	}
}

func TestApplyRuleMinZero(t *testing.T) {
	if applyRule("min=0", "") != "" {
		t.Error("expected no error for min=0 with empty value")
	}
}

func TestApplyRuleMaxZero(t *testing.T) {
	if applyRule("max=0", "x") == "" {
		t.Error("expected error for max=0 with non-empty value")
	}
}

func TestValidateFormMultipleTags(t *testing.T) {
	type form struct {
		Field string `validate:"required,min=3"`
	}
	f := form{Field: ""}
	errs := ValidateForm(f)
	if errs == nil || errs["field"] == "" {
		t.Error("expected required error")
	}
	f.Field = "ab"
	errs = ValidateForm(f)
	if errs == nil || errs["field"] == "" {
		t.Error("expected min error")
	}
}

func TestBindFormNoFields(t *testing.T) {
	type empty struct{}
	f := empty{}
	r := &http.Request{Form: url.Values{}}
	err := BindForm(r, &f)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateFormWhitespaceTag(t *testing.T) {
	type form struct {
		Field string `validate:" required , email "`
	}
	f := form{Field: ""}
	errs := ValidateForm(f)
	if errs == nil || errs["field"] == "" {
		t.Error("expected required error with whitespace in tag")
	}
}

func TestApplyRuleMinLarge(t *testing.T) {
	if applyRule("min=100", strings.Repeat("x", 50)) == "" {
		t.Error("expected error for value shorter than large min")
	}
}

func TestApplyRuleMaxLarge(t *testing.T) {
	if applyRule("max=100", strings.Repeat("x", 150)) == "" {
		t.Error("expected error for value longer than large max")
	}
}

func TestValidateFormPtrNil(t *testing.T) {
	var f *validateFormRequired
	func() {
		defer func() {
			if recover() == nil {
				t.Error("expected panic for nil pointer form")
			}
		}()
		ValidateForm(f)
	}()
}

func TestBindFormNonPtr(t *testing.T) {
	f := bindFormStruct{}
	r := &http.Request{Form: url.Values{"name": {"Alice"}}}
	func() {
		defer func() {
			if recover() == nil {
				t.Error("expected panic for non-pointer target")
			}
		}()
		BindForm(r, f)
	}()
}

func TestSaveOldEmptyStruct(t *testing.T) {
	c := &SSRContext{data: map[string]any{}}
	type empty struct{}
	SaveOld(c, empty{})
	if len(c.data) != 0 {
		t.Error("expected no data for empty struct")
	}
}

func TestSaveErrorsOverwrite(t *testing.T) {
	c := &SSRContext{data: map[string]any{}}
	c.Set("error_name", "old")
	SaveErrors(c, map[string]string{"name": "new"})
	if c.Get("error_name") != "new" {
		t.Errorf("expected 'new', got '%v'", c.Get("error_name"))
	}
}

func TestValidateFormEmptyTag(t *testing.T) {
	type form struct {
		Field string `validate:""`
	}
	f := form{Field: "x"}
	errs := ValidateForm(f)
	if errs != nil {
		t.Error("expected no errors for empty validate tag")
	}
}

func TestBindFormExtraFormValues(t *testing.T) {
	f := bindFormStruct{}
	r := &http.Request{Form: url.Values{"name": {"Alice"}, "extra": {"ignored"}}}
	BindForm(r, &f)
	if f.Name != "Alice" {
		t.Errorf("expected 'Alice', got '%s'", f.Name)
	}
}

func TestApplyRuleMinWithSpaces(t *testing.T) {
	if applyRule("min= 3", "ab") == "" {
		t.Error("expected error for min with space")
	}
}

func TestApplyRuleMaxWithSpaces(t *testing.T) {
	if applyRule("max= 3", "abcd") == "" {
		t.Error("expected error for max with space")
	}
}

func TestValidateFormMultipleErrorsSameField(t *testing.T) {
	type form struct {
		Field string `validate:"required,email"`
	}
	f := form{Field: ""}
	errs := ValidateForm(f)
	if errs == nil || errs["field"] == "" {
		t.Error("expected required error (first rule breaks)")
	}
}

func TestBindFormCaseInsensitive(t *testing.T) {
	type form struct {
		NAME string
	}
	f := form{}
	r := &http.Request{Form: url.Values{"name": {"Alice"}}}
	BindForm(r, &f)
	if f.NAME != "Alice" {
		t.Errorf("expected 'Alice', got '%s'", f.NAME)
	}
}

func TestSaveOldPtr(t *testing.T) {
	c := &SSRContext{data: map[string]any{}}
	f := &bindFormStruct{Name: "Alice"}
	func() {
		defer func() {
			if recover() == nil {
				t.Error("expected panic for pointer to SaveOld")
			}
		}()
		SaveOld(c, f)
	}()
}

func TestSaveErrorsNil(t *testing.T) {
	c := &SSRContext{data: map[string]any{}}
	SaveErrors(c, nil)
	if len(c.data) != 0 {
		t.Error("expected no data for nil errors")
	}
}
