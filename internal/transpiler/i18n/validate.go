package i18n

import (
	"fmt"
	"sort"
	"strconv"
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
		if argument.TimeZoneArgument != "" {
			if _, exists := message.Arguments[argument.TimeZoneArgument]; !exists {
				return fmt.Errorf("argument %q: time zone argument %q is not defined", name, argument.TimeZoneArgument)
			}
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
		names, err := placeholders(*value.Text)
		if err != nil {
			return err
		}
		for _, name := range names {
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
	selectorFormat := arguments[selector.Argument].Format
	if value.Plural != nil && selectorFormat != "number" && selectorFormat != "integer" && selectorFormat != "percent" && selectorFormat != "currency" {
		return fmt.Errorf("plural argument %q must use a numeric format", selector.Argument)
	}
	if value.Select != nil && selectorFormat != "string" {
		return fmt.Errorf("select argument %q must use string format", selector.Argument)
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
		if value.Plural != nil && !validPluralCase(name) {
			return fmt.Errorf("unsupported plural case %q", name)
		}
		if err := validateValue(child, arguments); err != nil {
			return fmt.Errorf("case %q: %w", name, err)
		}
	}
	return nil
}

func validPluralCase(name string) bool {
	switch name {
	case "zero", "one", "two", "few", "many", "other":
		return true
	}
	if !strings.HasPrefix(name, "=") {
		return false
	}
	_, err := strconv.ParseFloat(strings.TrimPrefix(name, "="), 64)
	return err == nil
}

func inferArguments(value Value, arguments map[string]Argument) error {
	if value.Text != nil {
		names, err := placeholders(*value.Text)
		if err != nil {
			return err
		}
		for _, name := range names {
			if _, exists := arguments[name]; !exists {
				arguments[name] = Argument{Format: "string"}
			}
		}
		return nil
	}
	selector := value.Plural
	if selector == nil {
		selector = value.Select
	}
	if selector == nil {
		return nil
	}
	if _, exists := arguments[selector.Argument]; !exists {
		format := "string"
		if value.Plural != nil {
			format = "number"
		}
		arguments[selector.Argument] = Argument{Format: format}
	}
	for _, child := range selector.Cases {
		if err := inferArguments(child, arguments); err != nil {
			return err
		}
	}
	return nil
}

func inferArgumentFormats(value Value, arguments map[string]Argument) {
	if value.Text != nil {
		return
	}
	selector := value.Plural
	format := "number"
	if selector == nil {
		selector = value.Select
		format = "string"
	}
	if selector == nil {
		return
	}
	argument := arguments[selector.Argument]
	if argument.Format == "" {
		argument.Format = format
		arguments[selector.Argument] = argument
	}
	for _, child := range selector.Cases {
		inferArgumentFormats(child, arguments)
	}
}

func placeholders(text string) ([]string, error) {
	seen := map[string]struct{}{}
	for index := 0; index < len(text); index++ {
		if text[index] == '}' {
			if index+1 < len(text) && text[index+1] == '}' {
				index++
				continue
			}
			return nil, fmt.Errorf("unexpected closing placeholder brace")
		}
		if text[index] != '{' {
			continue
		}
		if index+1 < len(text) && text[index+1] == '{' {
			index++
			continue
		}
		end := strings.IndexByte(text[index+1:], '}')
		if end < 0 {
			return nil, fmt.Errorf("unclosed placeholder")
		}
		name := text[index+1 : index+1+end]
		if !validName(name, false) {
			return nil, fmt.Errorf("invalid placeholder %q", name)
		}
		seen[name] = struct{}{}
		index += end + 1
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

func validateCoverage(set Set) error {
	defaultCatalog, exists := set.Locales[set.DefaultLocale]
	if !exists {
		return fmt.Errorf("default locale %q has no catalog", set.DefaultLocale)
	}
	locales := make([]string, 0, len(set.Locales))
	for locale := range set.Locales {
		locales = append(locales, locale)
	}
	sort.Strings(locales)
	for _, locale := range locales {
		catalog := set.Locales[locale]
		keys := make([]string, 0, len(catalog.Messages))
		for key := range catalog.Messages {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			message := catalog.Messages[key]
			defaultMessage, exists := defaultCatalog.Messages[key]
			if !exists {
				return fmt.Errorf("message %q is missing from default locale %q (referenced by %q)", key, set.DefaultLocale, locale)
			}
			got := argumentNames(message.Arguments)
			want := argumentNames(defaultMessage.Arguments)
			if strings.Join(got, "\x00") != strings.Join(want, "\x00") {
				return fmt.Errorf("message %q in locale %q has arguments %v, want %v", key, locale, got, want)
			}
			for _, name := range want {
				gotFormat := message.Arguments[name]
				wantFormat := defaultMessage.Arguments[name]
				if gotFormat.Format != wantFormat.Format || gotFormat.CurrencyArgument != wantFormat.CurrencyArgument || gotFormat.TimeZoneArgument != wantFormat.TimeZoneArgument {
					return fmt.Errorf("message %q argument %q in locale %q has an incompatible format contract", key, name, locale)
				}
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
