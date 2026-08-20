package validate

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

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
	if err := BindForm(r, f); err == nil {
		t.Fatal("expected error for non-pointer target")
	}
}

func TestSaveOldEmptyStruct(t *testing.T) {
	c := &testSetter{}
	type empty struct{}
	SaveOld(c, empty{})
	if len(c.data) != 0 {
		t.Error("expected no data for empty struct")
	}
}

func TestSaveErrorsOverwrite(t *testing.T) {
	c := &testSetter{}
	c.Set("error_name", "old")
	SaveErrors(c, map[string]string{"name": "new"})
	if c.data["error_name"] != "new" {
		t.Errorf("expected 'new', got '%v'", c.data["error_name"])
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
	c := &testSetter{}
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
	c := &testSetter{}
	SaveErrors(c, nil)
	if len(c.data) != 0 {
		t.Error("expected no data for nil errors")
	}
}
