package i18n

import (
	"encoding/json"
	"io"
	"sort"
)

type Extraction struct {
	FormatVersion int                `json:"formatVersion"`
	SourceLocale  string             `json:"sourceLocale"`
	Locales       []string           `json:"locales"`
	Messages      []ExtractedMessage `json:"messages"`
}

type ExtractedMessage struct {
	ID        string              `json:"id"`
	Value     Value               `json:"value"`
	Arguments map[string]Argument `json:"arguments,omitempty"`
}

func Extract(writer io.Writer, set Set) error {
	locales := make([]string, 0, len(set.Locales))
	for locale := range set.Locales {
		locales = append(locales, locale)
	}
	sort.Strings(locales)
	source := set.Locales[set.DefaultLocale]
	keys := make([]string, 0, len(source.Messages))
	for key := range source.Messages {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	document := Extraction{FormatVersion: 1, SourceLocale: set.DefaultLocale, Locales: locales}
	for _, key := range keys {
		message := source.Messages[key]
		document.Messages = append(document.Messages, ExtractedMessage{ID: key, Value: message.Value, Arguments: message.Arguments})
	}
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	return encoder.Encode(document)
}
