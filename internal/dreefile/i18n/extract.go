package i18n

import (
	"encoding/json"
	"io"
	"maps"
	"slices"
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
	locales := slices.Sorted(maps.Keys(set.Locales))
	source := set.Locales[set.DefaultLocale]
	keys := slices.Sorted(maps.Keys(source.Messages))
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
