package i18n

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/text/currency"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type CatalogLocalizer struct {
	defaultLocale string
	catalogs      map[string]LocaleCatalog
	tags          map[string]language.Tag
}

func NewCatalogLocalizer(config Config) (*CatalogLocalizer, error) {
	localizer := &CatalogLocalizer{
		defaultLocale: config.DefaultLocale,
		catalogs:      make(map[string]LocaleCatalog, len(config.Locales)),
		tags:          make(map[string]language.Tag, len(config.Locales)),
	}
	for _, catalog := range config.Locales {
		tag, err := language.Parse(catalog.Locale)
		if err != nil {
			return nil, fmt.Errorf("locale %q: %w", catalog.Locale, err)
		}
		locale := tag.String()
		localizer.catalogs[locale] = catalog
		localizer.tags[locale] = tag
	}
	defaultTag, err := language.Parse(config.DefaultLocale)
	if err != nil {
		return nil, fmt.Errorf("default locale: %w", err)
	}
	localizer.defaultLocale = defaultTag.String()
	if _, exists := localizer.catalogs[localizer.defaultLocale]; !exists {
		return nil, fmt.Errorf("default locale %q has no catalog", localizer.defaultLocale)
	}
	return localizer, nil
}

func (l *CatalogLocalizer) Localize(_ context.Context, locale, key string, arguments []Argument) (string, error) {
	catalog, exists := l.catalogs[locale]
	if !exists {
		catalog = l.catalogs[l.defaultLocale]
		locale = l.defaultLocale
	}
	entry, exists := catalog.Messages[key]
	if !exists && locale != l.defaultLocale {
		catalog = l.catalogs[l.defaultLocale]
		entry, exists = catalog.Messages[key]
		locale = l.defaultLocale
	}
	if !exists {
		return "", fmt.Errorf("message %q is not defined", key)
	}
	values := make(map[string]any, len(arguments))
	for _, argument := range arguments {
		values[argument.Name] = argument.Value
	}
	tag := l.tags[locale]
	return renderValue(entry.Value, entry.Arguments, values, tag)
}

func renderValue(value Value, formats map[string]ArgumentFormat, arguments map[string]any, tag language.Tag) (string, error) {
	if value.Text != nil {
		return interpolate(*value.Text, formats, arguments, tag)
	}
	if value.Selector == nil {
		return "", fmt.Errorf("message value is empty")
	}
	raw, exists := arguments[value.Selector.Argument]
	if !exists {
		return "", fmt.Errorf("missing argument %q", value.Selector.Argument)
	}
	caseName, err := selectCase(value.Selector, raw, tag)
	if err != nil {
		return "", err
	}
	selected, exists := value.Selector.Cases[caseName]
	if !exists {
		selected = value.Selector.Cases["other"]
	}
	return renderValue(selected, formats, arguments, tag)
}

func selectCase(selector *Selector, value any, tag language.Tag) (string, error) {
	if selector.Kind == "select" {
		return fmt.Sprint(value), nil
	}
	number, ok := integer(value)
	if !ok {
		return "", fmt.Errorf("plural argument %q must be an integer", selector.Argument)
	}
	rules := plural.Cardinal
	if selector.Kind == "ordinal" {
		rules = plural.Ordinal
	}
	form := rules.MatchPlural(tag, number, 0, 0, 0, 0)
	switch form {
	case plural.Zero:
		return "zero", nil
	case plural.One:
		return "one", nil
	case plural.Two:
		return "two", nil
	case plural.Few:
		return "few", nil
	case plural.Many:
		return "many", nil
	default:
		return "other", nil
	}
}

func interpolate(text string, formats map[string]ArgumentFormat, arguments map[string]any, tag language.Tag) (string, error) {
	var output strings.Builder
	for {
		start := strings.IndexByte(text, '{')
		if start < 0 {
			output.WriteString(text)
			return output.String(), nil
		}
		output.WriteString(text[:start])
		end := strings.IndexByte(text[start+1:], '}')
		if end < 0 {
			return "", fmt.Errorf("unclosed placeholder")
		}
		name := text[start+1 : start+1+end]
		value, exists := arguments[name]
		if !exists {
			return "", fmt.Errorf("missing argument %q", name)
		}
		formatted, err := formatValue(value, formats[name], arguments, tag)
		if err != nil {
			return "", fmt.Errorf("argument %q: %w", name, err)
		}
		output.WriteString(formatted)
		text = text[start+end+2:]
	}
}

func formatValue(value any, format ArgumentFormat, arguments map[string]any, tag language.Tag) (string, error) {
	printer := message.NewPrinter(tag)
	switch format.Format {
	case "number", "integer", "":
		return printer.Sprint(value), nil
	case "percent":
		number, ok := decimal(value)
		if !ok {
			return "", fmt.Errorf("percent must be numeric")
		}
		return printer.Sprintf("%.2f%%", number*100), nil
	case "currency":
		code, ok := arguments[format.CurrencyArgument].(string)
		if !ok {
			return "", fmt.Errorf("currency argument %q must be an ISO 4217 string", format.CurrencyArgument)
		}
		unit, err := currency.ParseISO(code)
		if err != nil {
			return "", err
		}
		return printer.Sprint(currency.Symbol(unit.Amount(value))), nil
	default:
		return fmt.Sprint(value), nil
	}
}

func integer(value any) (int, bool) {
	switch number := value.(type) {
	case int:
		return number, true
	case int64:
		return int(number), true
	case int32:
		return int(number), true
	case uint:
		return int(number), true
	case uint64:
		return int(number), true
	default:
		return 0, false
	}
}

func decimal(value any) (float64, bool) {
	switch number := value.(type) {
	case float64:
		return number, true
	case float32:
		return float64(number), true
	default:
		integer, ok := integer(value)
		return float64(integer), ok
	}
}

func MatchSupported(matcher language.Matcher, preferred []string) string {
	tag, _ := language.MatchStrings(matcher, preferred...)
	base, _ := tag.Base()
	return base.String()
}
