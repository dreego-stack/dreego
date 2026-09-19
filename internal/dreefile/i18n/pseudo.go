package i18n

import "strings"

func isPseudoLocale(locale string) bool {
	return locale == "en-XA" || locale == "ar-XB"
}

func pseudoCatalog(source Catalog, locale string) Catalog {
	messages := make(map[string]Message, len(source.Messages))
	for key, message := range source.Messages {
		message.Value = pseudoValue(message.Value, locale)
		messages[key] = message
	}
	return Catalog{Locale: locale, Messages: messages}
}

func pseudoValue(value Value, locale string) Value {
	if value.Text != nil {
		text := pseudoText(*value.Text, locale)
		return Value{Text: &text}
	}
	selector := value.Plural
	result := Value{}
	if selector == nil {
		selector = value.Select
		result.Select = &Selector{Argument: selector.Argument, Type: selector.Type, Cases: map[string]Value{}}
		selector = result.Select
	} else {
		result.Plural = &Selector{Argument: selector.Argument, Type: selector.Type, Cases: map[string]Value{}}
		selector = result.Plural
	}
	var sourceCases map[string]Value
	if value.Plural != nil {
		sourceCases = value.Plural.Cases
	} else {
		sourceCases = value.Select.Cases
	}
	for name, child := range sourceCases {
		selector.Cases[name] = pseudoValue(child, locale)
	}
	return result
}

func pseudoText(source, locale string) string {
	if locale == "ar-XB" {
		return "\u202e" + source + "\u202c"
	}
	replacer := strings.NewReplacer(
		"a", "á", "e", "é", "i", "í", "o", "ó", "u", "ú",
		"A", "Á", "E", "É", "I", "Í", "O", "Ó", "U", "Ú",
	)
	var output strings.Builder
	for source != "" {
		start := strings.IndexByte(source, '{')
		if start < 0 {
			output.WriteString(replacer.Replace(source))
			break
		}
		output.WriteString(replacer.Replace(source[:start]))
		end := strings.IndexByte(source[start:], '}')
		if end < 0 {
			output.WriteString(replacer.Replace(source[start:]))
			break
		}
		end += start
		output.WriteString(source[start : end+1])
		source = source[end+1:]
	}
	return "[!! " + output.String() + " ~~~~ !!]"
}
