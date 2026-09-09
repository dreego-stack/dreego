package transpiler

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

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
	Enabled       bool                `json:"enabled"`
	DefaultLocale string              `json:"defaultLocale"`
	Locales       []string            `json:"locales"`
	URLStrategy   string              `json:"urlStrategy"`
	Detection     []string            `json:"detection"`
	Domains       map[string]string   `json:"domains"`
	Fallbacks     map[string][]string `json:"fallbacks"`
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
	canonicalDomains := make(map[string]string, len(i.Domains))
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
		domain := i.Domains[locale]
		if domain == "" {
			domain = i.Domains[canonical]
		}
		canonicalDomains[canonical] = strings.TrimSpace(domain)
	}
	i.Domains = canonicalDomains
	canonicalFallbacks := make(map[string][]string, len(i.Fallbacks))
	for source, targets := range i.Fallbacks {
		sourceTag, err := language.Parse(source)
		if err != nil {
			return fmt.Errorf("invalid fallback locale %q: %w", source, err)
		}
		canonicalSource := sourceTag.String()
		if _, exists := locales[canonicalSource]; !exists {
			return fmt.Errorf("fallback source locale %q is not listed in locales", canonicalSource)
		}
		seen := map[string]struct{}{}
		for _, target := range targets {
			targetTag, err := language.Parse(target)
			if err != nil {
				return fmt.Errorf("invalid fallback locale %q: %w", target, err)
			}
			canonicalTarget := targetTag.String()
			if _, exists := locales[canonicalTarget]; !exists {
				return fmt.Errorf("fallback locale %q is not listed in locales", canonicalTarget)
			}
			if _, exists := seen[canonicalTarget]; exists {
				return fmt.Errorf("duplicate fallback locale %q for %q", canonicalTarget, canonicalSource)
			}
			seen[canonicalTarget] = struct{}{}
			canonicalFallbacks[canonicalSource] = append(canonicalFallbacks[canonicalSource], canonicalTarget)
		}
	}
	i.Fallbacks = canonicalFallbacks
	if err := validateFallbackCycles(i.Fallbacks); err != nil {
		return err
	}
	defaultTag, err := language.Parse(i.DefaultLocale)
	if err != nil {
		return fmt.Errorf("invalid defaultLocale %q: %w", i.DefaultLocale, err)
	}
	i.DefaultLocale = defaultTag.String()
	if _, exists := locales[i.DefaultLocale]; !exists {
		return fmt.Errorf("defaultLocale %q is not listed in locales", i.DefaultLocale)
	}
	if i.URLStrategy == "domain" {
		seenDomains := map[string]string{}
		for _, locale := range i.Locales {
			domain := strings.ToLower(i.Domains[locale])
			if domain == "" {
				return fmt.Errorf("domain is required for locale %q", locale)
			}
			if !validLocaleDomain(domain) {
				return fmt.Errorf("invalid domain %q for locale %q", domain, locale)
			}
			if existing := seenDomains[domain]; existing != "" {
				return fmt.Errorf("domain %q is used by locales %q and %q", domain, existing, locale)
			}
			seenDomains[domain] = locale
			i.Domains[locale] = domain
		}
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

func validLocaleDomain(domain string) bool {
	if domain == "" || strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return false
	}
	for _, label := range strings.Split(domain, ".") {
		if label == "" || label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		for index := 0; index < len(label); index++ {
			char := label[index]
			if (char < 'a' || char > 'z') && (char < '0' || char > '9') && char != '-' {
				return false
			}
		}
	}
	return true
}

func validateFallbackCycles(fallbacks map[string][]string) error {
	visiting := map[string]bool{}
	visited := map[string]bool{}
	var visit func(string) error
	visit = func(locale string) error {
		if visiting[locale] {
			return fmt.Errorf("fallback cycle includes locale %q", locale)
		}
		if visited[locale] {
			return nil
		}
		visiting[locale] = true
		for _, target := range fallbacks[locale] {
			if err := visit(target); err != nil {
				return err
			}
		}
		visiting[locale] = false
		visited[locale] = true
		return nil
	}
	for locale := range fallbacks {
		if err := visit(locale); err != nil {
			return err
		}
	}
	return nil
}
