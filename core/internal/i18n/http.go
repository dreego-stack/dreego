package i18n

import (
	"net/http"

	"golang.org/x/text/language"
)

type Negotiator struct {
	config     Config
	localizer  Localizer
	locales    []string
	matcher    language.Matcher
	cookieName string
}

func NewNegotiator(config Config, localizer Localizer) (*Negotiator, error) {
	tags := make([]language.Tag, 0, len(config.Locales))
	locales := make([]string, 0, len(config.Locales))
	for _, catalog := range config.Locales {
		tag, err := language.Parse(catalog.Locale)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
		locales = append(locales, tag.String())
	}
	cookieName := config.CookieName
	if cookieName == "" {
		cookieName = "dreego_locale"
	}
	if len(config.Detection) == 0 {
		config.Detection = []string{"cookie", "browser", "custom", "default"}
	}
	return &Negotiator{config: config, localizer: localizer, locales: locales, matcher: language.NewMatcher(tags), cookieName: cookieName}, nil
}

func (n *Negotiator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		locale := n.Resolve(r)
		ctx := WithContext(r.Context(), n.localizer, locale)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (n *Negotiator) Resolve(r *http.Request) string {
	for _, detector := range n.config.Detection {
		switch detector {
		case "account":
			if locale := n.resolveCandidate(n.config.Account, r); locale != "" {
				return locale
			}
		case "cookie":
			if cookie, err := r.Cookie(n.cookieName); err == nil {
				if locale := n.match([]string{cookie.Value}); locale != "" {
					return locale
				}
			}
		case "browser":
			tags, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
			if err == nil {
				preferred := make([]string, len(tags))
				for index, tag := range tags {
					preferred[index] = tag.String()
				}
				if locale := n.match(preferred); locale != "" {
					return locale
				}
			}
		case "custom":
			for _, resolver := range n.config.Resolvers {
				if locale := n.resolveCandidate(resolver, r); locale != "" {
					return locale
				}
			}
		case "default":
			return n.config.DefaultLocale
		}
	}
	return n.config.DefaultLocale
}

func (n *Negotiator) resolveCandidate(resolver Resolver, request *http.Request) string {
	if resolver == nil {
		return ""
	}
	return n.match([]string{resolver(request)})
}

func (n *Negotiator) match(preferred []string) string {
	if len(preferred) == 0 || preferred[0] == "" {
		return ""
	}
	_, index, confidence := n.matcher.Match(parseTags(preferred)...)
	if confidence == language.No || index < 0 || index >= len(n.locales) {
		return ""
	}
	return n.locales[index]
}

func parseTags(values []string) []language.Tag {
	tags := make([]language.Tag, 0, len(values))
	for _, value := range values {
		if tag, err := language.Parse(value); err == nil {
			tags = append(tags, tag)
		}
	}
	return tags
}
