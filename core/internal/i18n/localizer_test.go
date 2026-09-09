package i18n

import (
	"context"
	"strings"
	"testing"
	"time"

	"golang.org/x/text/language"
)

func testLocalizer(t *testing.T) *CatalogLocalizer {
	t.Helper()
	localizer, err := NewCatalogLocalizer(Config{
		DefaultLocale: "de",
		Locales: []LocaleCatalog{
			{Locale: "de", Messages: map[string]Message{
				"hello":      {Value: Text("Hallo, {name}!")},
				"items":      {Value: Plural("count", "cardinal", map[string]Value{"one": Text("Ein Artikel"), "other": Text("{count} Artikel")})},
				"position":   {Value: Plural("position", "ordinal", map[string]Value{"one": Text("{position}. Platz"), "other": Text("{position}. Platz")})},
				"salutation": {Value: Select("form", map[string]Value{"formal": Text("Guten Tag"), "other": Text("Hallo")})},
			}},
			{Locale: "en", Messages: map[string]Message{
				"hello": {Value: Text("Hello, {name}!")},
				"items": {Value: Plural("count", "cardinal", map[string]Value{"one": Text("One item"), "other": Text("{count} items")})},
			}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	return localizer
}

func TestCatalogLocalizerRendersMessages(t *testing.T) {
	localizer := testLocalizer(t)
	cases := []struct {
		key  string
		args []Argument
		want string
	}{
		{"hello", []Argument{{Name: "name", Value: "Ada"}}, "Hallo, Ada!"},
		{"items", []Argument{{Name: "count", Value: 1}}, "Ein Artikel"},
		{"items", []Argument{{Name: "count", Value: 3}}, "3 Artikel"},
		{"salutation", []Argument{{Name: "form", Value: "formal"}}, "Guten Tag"},
	}
	for _, tc := range cases {
		got, err := localizer.Localize(context.Background(), "de", tc.key, tc.args)
		if err != nil || got != tc.want {
			t.Errorf("Localize(%q) = %q, %v; want %q", tc.key, got, err, tc.want)
		}
	}
}

func TestCatalogLocalizerSupportsEscapedBraces(t *testing.T) {
	localizer, err := NewCatalogLocalizer(Config{DefaultLocale: "en", Locales: []LocaleCatalog{{Locale: "en", Messages: map[string]Message{
		"braces": {Value: Text("Use {{name}} for {name}")},
	}}}})
	if err != nil {
		t.Fatal(err)
	}
	got, err := localizer.Localize(context.Background(), "en", "braces", []Argument{{Name: "name", Value: "Ada"}})
	if err != nil || got != "Use {name} for Ada" {
		t.Fatalf("Localize = %q, %v", got, err)
	}
}

func TestCatalogLocalizerSupportsDecimalAndExactPluralCases(t *testing.T) {
	localizer, err := NewCatalogLocalizer(Config{DefaultLocale: "en", Locales: []LocaleCatalog{{Locale: "en", Messages: map[string]Message{
		"items": {Value: Plural("count", "cardinal", map[string]Value{"=0": Text("Empty"), "one": Text("One"), "other": Text("{count}")})},
	}}}})
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		value any
		want  string
	}{{0, "Empty"}, {1.5, "1.5"}} {
		got, err := localizer.Localize(context.Background(), "en", "items", []Argument{{Name: "count", Value: tc.value}})
		if err != nil || got != tc.want {
			t.Errorf("value %v: got %q, %v; want %q", tc.value, got, err, tc.want)
		}
	}
}

func TestCatalogLocalizerFallsBackToDefault(t *testing.T) {
	localizer := testLocalizer(t)
	got, err := localizer.Localize(context.Background(), "en", "salutation", []Argument{{Name: "form", Value: "other"}})
	if err != nil || got != "Hallo" {
		t.Fatalf("Localize = %q, %v", got, err)
	}
}

func TestCatalogLocalizerUsesExplicitFallbackChain(t *testing.T) {
	localizer, err := NewCatalogLocalizer(Config{
		DefaultLocale: "de",
		Fallbacks:     map[string][]string{"de-CH": {"fr"}},
		Locales: []LocaleCatalog{
			{Locale: "de", Messages: map[string]Message{"legal": {Value: Text("Deutsch")}}},
			{Locale: "fr", Messages: map[string]Message{"legal": {Value: Text("Français")}}},
			{Locale: "de-CH", Messages: map[string]Message{}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	got, err := localizer.Localize(context.Background(), "de-CH", "legal", nil)
	if err != nil || got != "Français" {
		t.Fatalf("Localize = %q, %v", got, err)
	}
}

func TestCatalogLocalizerFormatsWithoutCurrencyConversion(t *testing.T) {
	localizer, err := NewCatalogLocalizer(Config{DefaultLocale: "de", Locales: []LocaleCatalog{{Locale: "de", Messages: map[string]Message{
		"total": {
			Value:     Text("Gesamt: {amount}"),
			Arguments: map[string]ArgumentFormat{"amount": {Format: "currency", CurrencyArgument: "currency"}},
		},
	}}}})
	if err != nil {
		t.Fatal(err)
	}
	got, err := localizer.Localize(context.Background(), "de", "total", []Argument{{Name: "amount", Value: 1234.5}, {Name: "currency", Value: "USD"}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "1.234,50") || !strings.Contains(got, "$") {
		t.Fatalf("localized currency = %q", got)
	}
}

func TestMatchSupportedLocale(t *testing.T) {
	matcher := language.NewMatcher([]language.Tag{language.German, language.English})
	if got := MatchSupported(matcher, []string{"fr", "en-GB"}); got != "en" {
		t.Fatalf("locale = %q", got)
	}
}

func TestCatalogLocalizerFormatsDateInExplicitTimeZone(t *testing.T) {
	localizer, err := NewCatalogLocalizer(Config{DefaultLocale: "de", Locales: []LocaleCatalog{{Locale: "de", Messages: map[string]Message{
		"date": {Value: Text("Datum: {instant}"), Arguments: map[string]ArgumentFormat{
			"instant": {Format: "datetime", TimeZoneArgument: "zone", Layout: "02.01.2006 15:04 MST"},
			"zone":    {},
		}},
	}}}})
	if err != nil {
		t.Fatal(err)
	}
	instant := time.Date(2026, 1, 2, 12, 0, 0, 0, time.UTC)
	got, err := localizer.Localize(context.Background(), "de", "date", []Argument{{Name: "instant", Value: instant}, {Name: "zone", Value: "Europe/Berlin"}})
	if err != nil || got != "Datum: 02.01.2026 13:00 CET" {
		t.Fatalf("Localize = %q, %v", got, err)
	}
}
