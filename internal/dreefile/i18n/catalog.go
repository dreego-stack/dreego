package i18n

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Set struct {
	DefaultLocale string
	Locales       map[string]Catalog
}

type Catalog struct {
	Locale   string
	Messages map[string]Message
}

type Message struct {
	Value     Value
	Arguments map[string]Argument
}

type Value struct {
	Text   *string   `json:"text,omitempty"`
	Plural *Selector `json:"plural,omitempty"`
	Select *Selector `json:"select,omitempty"`
}

type Selector struct {
	Argument string           `json:"argument"`
	Type     string           `json:"type,omitempty"`
	Cases    map[string]Value `json:"cases"`
}

type Argument struct {
	Format           string `json:"format,omitempty"`
	CurrencyArgument string `json:"currencyArgument,omitempty"`
	TimeZoneArgument string `json:"timeZoneArgument,omitempty"`
	Layout           string `json:"layout,omitempty"`
}

type messageJSON struct {
	Text      *string             `json:"text,omitempty"`
	Plural    *Selector           `json:"plural,omitempty"`
	Select    *Selector           `json:"select,omitempty"`
	Arguments map[string]Argument `json:"arguments,omitempty"`
}

func (v *Value) UnmarshalJSON(data []byte) error {
	var shorthand string
	if err := json.Unmarshal(data, &shorthand); err == nil {
		v.Text = &shorthand
		return nil
	}
	type rawValue Value
	var decoded rawValue
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&decoded); err != nil {
		return err
	}
	*v = Value(decoded)
	return nil
}

func Load(root string, locales []string, defaultLocale string) (Set, error) {
	set := Set{DefaultLocale: defaultLocale, Locales: make(map[string]Catalog, len(locales))}
	for _, locale := range locales {
		if isPseudoLocale(locale) {
			continue
		}
		catalog, err := loadCatalog(filepath.Join(root, locale), locale)
		if err != nil {
			return Set{}, err
		}
		set.Locales[locale] = catalog
	}
	defaultCatalog, exists := set.Locales[defaultLocale]
	if !exists {
		return Set{}, fmt.Errorf("default locale %q must have source catalogs", defaultLocale)
	}
	for _, locale := range locales {
		if isPseudoLocale(locale) {
			set.Locales[locale] = pseudoCatalog(defaultCatalog, locale)
		}
	}
	if err := validateCoverage(set); err != nil {
		return Set{}, err
	}
	return set, nil
}

func loadCatalog(dir, locale string) (Catalog, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return Catalog{}, fmt.Errorf("locale %q: %w", locale, err)
	}
	catalog := Catalog{Locale: locale, Messages: map[string]Message{}}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		messages, err := loadFile(path)
		if err != nil {
			return Catalog{}, err
		}
		for key, message := range messages {
			if _, exists := catalog.Messages[key]; exists {
				return Catalog{}, fmt.Errorf("%s: duplicate message %q", path, key)
			}
			catalog.Messages[key] = message
		}
	}
	return catalog, nil
}

func loadFile(path string) (map[string]Message, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := validateJSON(data); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	var raw map[string]json.RawMessage
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&raw); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	keys := slices.Sorted(maps.Keys(raw))
	messages := make(map[string]Message, len(raw))
	for _, key := range keys {
		if !validName(key, true) {
			return nil, fmt.Errorf("%s: invalid message key %q", path, key)
		}
		message, err := decodeMessage(raw[key])
		if err != nil {
			return nil, fmt.Errorf("%s: message %q: %w", path, key, err)
		}
		messages[key] = message
	}
	return messages, nil
}

func decodeMessage(raw json.RawMessage) (Message, error) {
	var shorthand string
	if err := json.Unmarshal(raw, &shorthand); err == nil {
		value := Value{Text: &shorthand}
		message := Message{Value: value, Arguments: map[string]Argument{}}
		if err := inferArguments(value, message.Arguments); err != nil {
			return Message{}, err
		}
		return message, nil
	}
	var source messageJSON
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&source); err != nil {
		return Message{}, err
	}
	message := Message{
		Value:     Value{Text: source.Text, Plural: source.Plural, Select: source.Select},
		Arguments: source.Arguments,
	}
	if message.Arguments == nil {
		message.Arguments = map[string]Argument{}
		if err := inferArguments(message.Value, message.Arguments); err != nil {
			return Message{}, err
		}
	}
	inferArgumentFormats(message.Value, message.Arguments)
	for name, argument := range message.Arguments {
		if argument.Format == "" {
			argument.Format = "string"
			message.Arguments[name] = argument
		}
	}
	if err := validateMessage(message); err != nil {
		return Message{}, err
	}
	return message, nil
}

func validName(value string, dots bool) bool {
	if value == "" {
		return false
	}
	for part := range strings.SplitSeq(value, ".") {
		if part == "" || !asciiLetter(part[0]) {
			return false
		}
		for index := 1; index < len(part); index++ {
			char := part[index]
			if !asciiLetter(char) && (char < '0' || char > '9') && char != '-' && char != '_' {
				return false
			}
		}
		if !dots && strings.Contains(value, ".") {
			return false
		}
	}
	return true
}

func asciiLetter(char byte) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z'
}
