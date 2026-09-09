package i18n

import (
	"fmt"
	"sort"
	"strings"
)

func ValidateUses(set Set, uses map[string][]string) error {
	defaultCatalog := set.Locales[set.DefaultLocale]
	keys := make([]string, 0, len(uses))
	for key := range uses {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		message, exists := defaultCatalog.Messages[key]
		if !exists {
			return fmt.Errorf("message %q is not defined in default locale %q", key, set.DefaultLocale)
		}
		got := append([]string(nil), uses[key]...)
		sort.Strings(got)
		want := argumentNames(message.Arguments)
		if strings.Join(got, "\x00") != strings.Join(want, "\x00") {
			return fmt.Errorf("message %q uses arguments %v, want %v", key, got, want)
		}
	}
	return nil
}

func GoConfig(set Set, detection []string) string {
	var output strings.Builder
	output.WriteString("dreego.I18nConfig{DefaultLocale: ")
	output.WriteString(fmt.Sprintf("%q, Detection: %#v, Locales: []dreego.LocaleCatalog{", set.DefaultLocale, detection))
	locales := make([]string, 0, len(set.Locales))
	for locale := range set.Locales {
		locales = append(locales, locale)
	}
	sort.Strings(locales)
	for _, locale := range locales {
		catalog := set.Locales[locale]
		output.WriteString(fmt.Sprintf("{Locale: %q, Messages: map[string]dreego.LocalizedMessage{", locale))
		keys := make([]string, 0, len(catalog.Messages))
		for key := range catalog.Messages {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			output.WriteString(fmt.Sprintf("%q: %s,", key, goMessage(catalog.Messages[key])))
		}
		output.WriteString("}},")
	}
	output.WriteString("}}")
	return output.String()
}

func goMessage(message Message) string {
	var output strings.Builder
	output.WriteString("{Value: ")
	output.WriteString(goValue(message.Value))
	if len(message.Arguments) > 0 {
		output.WriteString(", Arguments: map[string]dreego.MessageArgumentFormat{")
		for _, name := range argumentNames(message.Arguments) {
			argument := message.Arguments[name]
			output.WriteString(fmt.Sprintf("%q: {Format: %q, CurrencyArgument: %q, TimeZoneArgument: %q},", name, argument.Format, argument.CurrencyArgument, argument.TimeZoneArgument))
		}
		output.WriteString("}")
	}
	output.WriteString("}")
	return output.String()
}

func goValue(value Value) string {
	if value.Text != nil {
		return fmt.Sprintf("dreego.MessageText(%q)", *value.Text)
	}
	selector := value.Plural
	function := "MessagePlural"
	kind := "cardinal"
	if selector == nil {
		selector = value.Select
		function = "MessageSelect"
		kind = ""
	} else if selector.Type != "" {
		kind = selector.Type
	}
	var output strings.Builder
	output.WriteString("dreego.")
	output.WriteString(function)
	output.WriteString(fmt.Sprintf("(%q", selector.Argument))
	if function == "MessagePlural" {
		output.WriteString(fmt.Sprintf(", %q", kind))
	}
	output.WriteString(", map[string]dreego.MessageValue{")
	cases := make([]string, 0, len(selector.Cases))
	for name := range selector.Cases {
		cases = append(cases, name)
	}
	sort.Strings(cases)
	for _, name := range cases {
		output.WriteString(fmt.Sprintf("%q: %s,", name, goValue(selector.Cases[name])))
	}
	output.WriteString("})")
	return output.String()
}
