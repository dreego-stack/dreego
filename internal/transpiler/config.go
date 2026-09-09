package transpiler

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"golang.org/x/text/language"
)

var ErrInvalidI18n = errors.New("invalid i18n configuration")

type Redirect struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Status int    `json:"status"`
}

type Rewrite struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Settings struct {
	Logging   Logging    `json:"logging"`
	Redirects []Redirect `json:"redirects"`
	Rewrites  []Rewrite  `json:"rewrites"`
	I18n      I18n       `json:"i18n"`
}

type Logging struct {
	Enabled bool `json:"enabled"`
}

type I18n struct {
	Enabled       bool     `json:"enabled"`
	DefaultLocale string   `json:"defaultLocale"`
	Locales       []string `json:"locales"`
	URLStrategy   string   `json:"urlStrategy"`
	Detection     []string `json:"detection"`
}

func LoadConfig(path string) (*Settings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	if err := s.I18n.validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidI18n, err)
	}
	return &s, nil
}

func (i *I18n) validate() error {
	if !i.Enabled {
		return nil
	}
	if i.DefaultLocale == "" {
		return fmt.Errorf("defaultLocale is required")
	}
	if i.URLStrategy == "" {
		i.URLStrategy = "none"
	}
	if i.URLStrategy != "none" && i.URLStrategy != "prefix" && i.URLStrategy != "domain" {
		return fmt.Errorf("unsupported urlStrategy %q", i.URLStrategy)
	}
	locales := make(map[string]struct{}, len(i.Locales))
	for index, locale := range i.Locales {
		tag, err := language.Parse(locale)
		if err != nil {
			return fmt.Errorf("invalid locale %q: %w", locale, err)
		}
		canonical := tag.String()
		if _, exists := locales[canonical]; exists {
			return fmt.Errorf("duplicate locale %q", canonical)
		}
		i.Locales[index] = canonical
		locales[canonical] = struct{}{}
	}
	defaultTag, err := language.Parse(i.DefaultLocale)
	if err != nil {
		return fmt.Errorf("invalid defaultLocale %q: %w", i.DefaultLocale, err)
	}
	i.DefaultLocale = defaultTag.String()
	if _, exists := locales[i.DefaultLocale]; !exists {
		return fmt.Errorf("defaultLocale %q is not listed in locales", i.DefaultLocale)
	}
	if len(i.Detection) == 0 {
		i.Detection = []string{"cookie", "browser", "custom", "default"}
	}
	detectors := make(map[string]struct{}, len(i.Detection))
	for index, detector := range i.Detection {
		switch detector {
		case "account", "cookie", "browser", "custom", "default":
		default:
			return fmt.Errorf("unsupported locale detector %q", detector)
		}
		if _, exists := detectors[detector]; exists {
			return fmt.Errorf("duplicate locale detector %q", detector)
		}
		if detector == "default" && index != len(i.Detection)-1 {
			return fmt.Errorf("locale detector %q must be last", detector)
		}
		detectors[detector] = struct{}{}
	}
	return nil
}
