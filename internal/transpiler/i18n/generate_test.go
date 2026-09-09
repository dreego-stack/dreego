package i18n

import (
	"strings"
	"testing"
)

func TestValidateUses(t *testing.T) {
	text := "Hello, {name}"
	set := Set{DefaultLocale: "en", Locales: map[string]Catalog{"en": {Locale: "en", Messages: map[string]Message{
		"hello": {Value: Value{Text: &text}, Arguments: map[string]Argument{"name": {}}},
	}}}}
	if err := ValidateUses(set, []Use{{Key: "hello", Arguments: []string{"name"}}}); err != nil {
		t.Fatal(err)
	}
	if err := ValidateUses(set, []Use{{Key: "missing"}}); err == nil {
		t.Fatal("expected missing key error")
	}
	if err := ValidateUses(set, []Use{{Key: "hello"}}); err == nil {
		t.Fatal("expected argument mismatch")
	}
}

func TestGoConfig(t *testing.T) {
	text := "Hello"
	set := Set{DefaultLocale: "en", Locales: map[string]Catalog{"en": {Locale: "en", Messages: map[string]Message{
		"hello": {Value: Value{Text: &text}},
	}}}}
	got := GoConfig(set, []string{"browser", "default"}, "none", nil, nil)
	for _, want := range []string{`DefaultLocale: "en"`, `Detection: []string{"browser", "default"}`, `"hello": {Value: dreego.MessageText("Hello")}`} {
		if !strings.Contains(got, want) {
			t.Errorf("config missing %q:\n%s", want, got)
		}
	}
}
