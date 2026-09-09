package i18n

import (
	"fmt"
	"sort"
	"strings"
)

func validateMessage(message Message) error {
	if err := validateValue(message.Value, message.Arguments); err != nil {
		return err
	}
	for name, argument := range message.Arguments {
		if !validName(name, false) {
			return fmt.Errorf("invalid argument name %q", name)
		}
		switch argument.Format {
		case "", "string", "number", "integer", "percent", "currency", "date", "time", "datetime":
		default:
			return fmt.Errorf("argument %q has unsupported format %q", name, argument.Format)
		}
		if argument.Format == "currency" && argument.CurrencyArgument == "" {
			return fmt.Errorf("argument %q: currencyArgument is required", name)
		}
		if argument.CurrencyArgument != "" {
			if _, exists := message.Arguments[argument.CurrencyArgument]; !exists {
				return fmt.Errorf("argument %q: currency argument %q is not defined", name, argument.CurrencyArgument)
			}
		}
		if argument.CurrencyArgument != "" && !validName(argument.CurrencyArgument, false) {
			return fmt.Errorf("argument %q has invalid currencyArgument %q", name, argument.CurrencyArgument)
		}
		if argument.TimeZoneArgument != "" && !validName(argument.TimeZoneArgument, false) {
			return fmt.Errorf("argument %q has invalid timeZoneArgument %q", name, argument.TimeZoneArgument)
		}
	}
	return nil
}

func validateValue(value Value, arguments map[string]Argument) error {
	variants := 0
	if value.Text != nil {
		variants++
	}
	if value.Plural != nil {
		variants++
	}
	if value.Select != nil {
		variants++
	}
	if variants != 1 {
		return fmt.Errorf("message value must define exactly one of text, plural, or select")
	}
	if value.Text != nil {
		for _, name := range placeholders(*value.Text) {
			if _, exists := arguments[name]; !exists {
				return fmt.Errorf("placeholder %q has no argument definition", name)
			}
		}
		return nil
	}
	selector := value.Plural
	if selector == nil {
		selector = value.Select
	}
	if !validName(selector.Argument, false) {
		return fmt.Errorf("invalid selector argument %q", selector.Argument)
	}
	if _, exists := arguments[selector.Argument]; !exists {
		arguments[selector.Argument] = Argument{}
	}
	if _, exists := selector.Cases["other"]; !exists {
		return fmt.Errorf("selector requires an %q case", "other")
	}
	if value.Plural != nil && selector.Type != "" && selector.Type != "cardinal" && selector.Type != "ordinal" {
		return fmt.Errorf("unsupported plural type %q", selector.Type)
	}
	for name, child := range selector.Cases {
		if name == "" {
			return fmt.Errorf("selector case must not be empty")
		}
		if err := validateValue(child, arguments); err != nil {
			return fmt.Errorf("case %q: %w", name, err)
		}
	}
	return nil
}

func inferArguments(value Value, arguments map[string]Argument) {
	if value.Text != nil {
		for _, name := range placeholders(*value.Text) {
			arguments[name] = Argument{}
		}
		return
	}
	selector := value.Plural
	if selector == nil {
		selector = value.Select
	}
	if selector == nil {
		return
	}
	arguments[selector.Argument] = Argument{}
	for _, child := range selector.Cases {
		inferArguments(child, arguments)
	}
}

func placeholders(text string) []string {
	seen := map[string]struct{}{}
	for index := 0; index < len(text); index++ {
		if text[index] != '{' {
			continue
		}
		end := strings.IndexByte(text[index+1:], '}')
		if end < 0 {
			continue
		}
		name := text[index+1 : index+1+end]
		if validName(name, false) {
			seen[name] = struct{}{}
		}
		index += end + 1
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func validateCoverage(set Set) error {
	defaultCatalog, exists := set.Locales[set.DefaultLocale]
	if !exists {
		return fmt.Errorf("default locale %q has no catalog", set.DefaultLocale)
	}
	for locale, catalog := range set.Locales {
		for key, message := range catalog.Messages {
			defaultMessage, exists := defaultCatalog.Messages[key]
			if !exists {
				return fmt.Errorf("message %q is missing from default locale %q (referenced by %q)", key, set.DefaultLocale, locale)
			}
			got := argumentNames(message.Arguments)
			want := argumentNames(defaultMessage.Arguments)
			if strings.Join(got, "\x00") != strings.Join(want, "\x00") {
				return fmt.Errorf("message %q in locale %q has arguments %v, want %v", key, locale, got, want)
			}
		}
	}
	return nil
}

func argumentNames(arguments map[string]Argument) []string {
	names := make([]string, 0, len(arguments))
	for name := range arguments {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
