package tests

import (
	"net/http"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestI18nTemplateRendering(t *testing.T) {
	t.Parallel()
	client := dreegotest.Serve(t, map[string]string{
		"www/dreego.config.json": `{
			"i18n": {
				"enabled": true,
				"defaultLocale": "de",
				"locales": ["de", "en"],
				"detection": ["cookie", "browser", "default"]
			}
		}`,
		"www/locales/de/messages.json": `{
			"home.greeting": "Hallo, {name}!",
			"cart.items": {"plural":{"argument":"count","cases":{"one":"Ein Artikel","other":"{count} Artikel"}}}
		}`,
		"www/locales/en/messages.json": `{
			"home.greeting": "Hello, {name}!",
			"cart.items": {"plural":{"argument":"count","cases":{"one":"One item","other":"{count} items"}}}
		}`,
		"www/routes/index.dreego": `<body><p>[[ home.greeting name="Ada" ]]</p><p>[[ cart.items count=2 ]]</p></body>`,
	})
	headers := map[string]string{"Accept-Language": "en-GB,en;q=0.8"}
	code, body, _ := client.Request(t, http.MethodGet, "/", "", headers)
	if code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", code, body)
	}
	if !strings.Contains(body, "Hello, Ada!") || !strings.Contains(body, "2 items") {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestI18nGenerationRejectsMessageArgumentMismatch(t *testing.T) {
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/dreego.config.json":       `{"i18n":{"enabled":true,"defaultLocale":"en","locales":["en"]}}`,
		"www/locales/en/messages.json": `{"home.greeting":"Hello, {name}!"}`,
		"www/routes/index.dreego":      `<body>[[ home.greeting ]]</body>`,
	})
	output, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil || !strings.Contains(output, `message "home.greeting" uses arguments [], want [name]`) {
		t.Fatalf("output = %q, error = %v", output, err)
	}
}
