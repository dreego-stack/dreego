package i18n

import (
	"context"
	"net/http"
)

type Argument struct {
	Name  string
	Value any
}

type ArgumentFormat struct {
	Format           string
	CurrencyArgument string
	TimeZoneArgument string
}

type Value struct {
	Text     *string
	Selector *Selector
}

type Selector struct {
	Argument string
	Kind     string
	Cases    map[string]Value
}

type Message struct {
	Value     Value
	Arguments map[string]ArgumentFormat
}

type LocaleCatalog struct {
	Locale   string
	Messages map[string]Message
}

type Config struct {
	DefaultLocale string
	Locales       []LocaleCatalog
	Detection     []string
	CookieName    string
	Account       Resolver
	Resolvers     []Resolver
}

type Resolver func(*http.Request) string

type Localizer interface {
	Localize(ctx context.Context, locale, key string, arguments []Argument) (string, error)
}

func Text(text string) Value {
	return Value{Text: &text}
}

func Plural(argument, kind string, cases map[string]Value) Value {
	return Value{Selector: &Selector{Argument: argument, Kind: kind, Cases: cases}}
}

func Select(argument string, cases map[string]Value) Value {
	return Value{Selector: &Selector{Argument: argument, Kind: "select", Cases: cases}}
}
