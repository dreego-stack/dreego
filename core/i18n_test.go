package core

import (
	"context"
	"testing"

	corei18n "github.com/dreego-stack/dreego/core/internal/i18n"
)

func TestMessageUsesRequestLocalizer(t *testing.T) {
	localizer, err := NewLocalizer(I18nConfig{DefaultLocale: "en", Locales: []LocaleCatalog{{Locale: "en", Messages: map[string]LocalizedMessage{
		"hello": {Value: MessageText("Hello, {name}!")},
	}}}})
	if err != nil {
		t.Fatal(err)
	}
	ctx := corei18n.WithContext(context.Background(), localizer, "en")
	if got := Message(ctx, "hello", MessageArg{Name: "name", Value: "Ada"}); got != "Hello, Ada!" {
		t.Fatalf("Message = %q", got)
	}
	if got := Locale(ctx); got != "en" {
		t.Fatalf("Locale = %q", got)
	}
}

func TestMessageWithoutLocalizerReturnsKey(t *testing.T) {
	if got := Message(context.Background(), "home.title"); got != "home.title" {
		t.Fatalf("Message = %q", got)
	}
}
