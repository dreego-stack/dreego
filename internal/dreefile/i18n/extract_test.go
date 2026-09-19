package i18n

import (
	"bytes"
	"strings"
	"testing"
)

func TestExtractIsDeterministicAndMachineReadable(t *testing.T) {
	text := "Hallo, {name}!"
	set := Set{DefaultLocale: "de", Locales: map[string]Catalog{
		"en": {Locale: "en"},
		"de": {Locale: "de", Messages: map[string]Message{
			"home.greeting": {Value: Value{Text: &text}, Arguments: map[string]Argument{"name": {}}},
		}},
	}}
	var first, second bytes.Buffer
	if err := Extract(&first, set); err != nil {
		t.Fatal(err)
	}
	if err := Extract(&second, set); err != nil {
		t.Fatal(err)
	}
	if first.String() != second.String() {
		t.Fatal("extraction is not deterministic")
	}
	for _, want := range []string{`"formatVersion": 1`, `"sourceLocale": "de"`, `"id": "home.greeting"`, `"name"`} {
		if !strings.Contains(first.String(), want) {
			t.Errorf("extraction missing %q:\n%s", want, first.String())
		}
	}
}
