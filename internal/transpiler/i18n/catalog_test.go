package i18n

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeCatalog(t *testing.T, root, locale, name, content string) {
	t.Helper()
	dir := filepath.Join(root, "locales", locale)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestLoadCatalogs(t *testing.T) {
	root := t.TempDir()
	writeCatalog(t, root, "de", "common.json", `{
		"home.title": "Willkommen, {name}!",
		"cart.items": {"plural":{"argument":"count","type":"cardinal","cases":{"one":"Ein Artikel","other":"{count} Artikel"}}},
		"account.salutation": {"select":{"argument":"form","cases":{"formal":"Guten Tag, {name}","other":"Hallo, {name}"}}},
		"invoice.total": {"text":"Gesamt: {amount}","arguments":{"amount":{"format":"currency","currencyArgument":"currency"},"currency":{}}}
	}`)
	writeCatalog(t, root, "en", "common.json", `{
		"home.title": "Welcome, {name}!",
		"cart.items": {"plural":{"argument":"count","type":"cardinal","cases":{"one":"One item","other":"{count} items"}}},
		"account.salutation": {"select":{"argument":"form","cases":{"formal":"Good day, {name}","other":"Hello, {name}"}}},
		"invoice.total": {"text":"Total: {amount}","arguments":{"amount":{"format":"currency","currencyArgument":"currency"},"currency":{}}}
	}`)

	set, err := Load(filepath.Join(root, "locales"), []string{"de", "en"}, "de")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(set.Locales["de"].Messages) != 4 {
		t.Fatalf("messages = %d", len(set.Locales["de"].Messages))
	}
	message := set.Locales["de"].Messages["invoice.total"]
	if message.Arguments["amount"].Format != "currency" {
		t.Fatalf("unexpected format: %+v", message.Arguments)
	}
}

func TestLoadCatalogsRejectsInvalidData(t *testing.T) {
	cases := []struct {
		name    string
		content string
		want    string
	}{
		{"unknown field", `{"home.title":{"text":"Hi","html":true}}`, "unknown field"},
		{"duplicate JSON key", `{"home.title":"Hi","home.title":"Hello"}`, `duplicate JSON key "home.title"`},
		{"invalid key", `{"home title":"Hi"}`, "invalid message key"},
		{"missing other", `{"items":{"plural":{"argument":"count","cases":{"one":"One"}}}}`, `requires an "other" case`},
		{"invalid plural case", `{"items":{"plural":{"argument":"count","cases":{"others":"Items","other":"Items"}}}}`, `unsupported plural case "others"`},
		{"mixed variants", `{"items":{"text":"Items","plural":{"argument":"count","cases":{"other":"Items"}}}}`, "exactly one of text, plural, or select"},
		{"unknown format", `{"price":{"text":"{amount}","arguments":{"amount":{"format":"money"}}}}`, `unsupported format "money"`},
		{"missing currency argument", `{"price":{"text":"{amount}","arguments":{"amount":{"format":"currency"}}}}`, "currencyArgument is required"},
		{"undeclared currency argument", `{"price":{"text":"{amount}","arguments":{"amount":{"format":"currency","currencyArgument":"currency"}}}}`, `currency argument "currency" is not defined`},
		{"undeclared placeholder", `{"hello":{"text":"Hello, {name}","arguments":{}}}`, `placeholder "name" has no argument definition`},
		{"unclosed placeholder", `{"hello":"Hello, {name"}`, "unclosed placeholder"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			writeCatalog(t, root, "en", "messages.json", tc.content)
			_, err := Load(filepath.Join(root, "locales"), []string{"en"}, "en")
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want %q", err, tc.want)
			}
		})
	}
}

func TestLoadCatalogsRequiresMatchingArgumentContracts(t *testing.T) {
	root := t.TempDir()
	writeCatalog(t, root, "de", "messages.json", `{"hello":{"text":"Hallo {name}","arguments":{"name":{}}}}`)
	writeCatalog(t, root, "en", "messages.json", `{"hello":{"text":"Hello {person}","arguments":{"person":{}}}}`)
	_, err := Load(filepath.Join(root, "locales"), []string{"de", "en"}, "de")
	if err == nil || !strings.Contains(err.Error(), `message "hello" in locale "en" has arguments [person], want [name]`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadCatalogsRequiresMatchingFormatContracts(t *testing.T) {
	root := t.TempDir()
	writeCatalog(t, root, "de", "messages.json", `{"value":{"text":"{value}","arguments":{"value":{"format":"number"}}}}`)
	writeCatalog(t, root, "en", "messages.json", `{"value":{"text":"{value}","arguments":{"value":{"format":"date"}}}}`)
	_, err := Load(filepath.Join(root, "locales"), []string{"de", "en"}, "de")
	if err == nil || !strings.Contains(err.Error(), "incompatible format contract") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadCatalogsGeneratesPseudoLocales(t *testing.T) {
	root := t.TempDir()
	writeCatalog(t, root, "de", "messages.json", `{"hello":"Hallo, {name}!"}`)
	set, err := Load(filepath.Join(root, "locales"), []string{"de", "en-XA", "ar-XB"}, "de")
	if err != nil {
		t.Fatal(err)
	}
	accented := *set.Locales["en-XA"].Messages["hello"].Value.Text
	if !strings.Contains(accented, "á") || !strings.Contains(accented, "{name}") {
		t.Fatalf("accented pseudo message = %q", accented)
	}
	bidi := *set.Locales["ar-XB"].Messages["hello"].Value.Text
	if !strings.HasPrefix(bidi, "\u202e") || !strings.HasSuffix(bidi, "\u202c") {
		t.Fatalf("bidi pseudo message = %q", bidi)
	}
}

func TestLoadCatalogsRequiresDefaultCoverage(t *testing.T) {
	root := t.TempDir()
	writeCatalog(t, root, "de", "messages.json", `{"home.title":"Willkommen","only.de":"Nur Deutsch"}`)
	writeCatalog(t, root, "en", "messages.json", `{"home.title":"Welcome","only.en":"English only"}`)
	_, err := Load(filepath.Join(root, "locales"), []string{"de", "en"}, "de")
	if err == nil || !strings.Contains(err.Error(), `message "only.en" is missing from default locale "de"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}
