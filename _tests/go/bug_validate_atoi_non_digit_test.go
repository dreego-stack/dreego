package tests

import (
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
)

func TestBugValidateAtoiNonDigit(t *testing.T) {
	type Form struct {
		Name string `validate:"min=abc"`
	}
	errs := dreego.ValidateForm(Form{Name: "x"})
	if errs == nil || errs["name"] == "" {
		t.Fatal("non-digit min rule silently accepted")
	}
}
