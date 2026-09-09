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
		"www/layouts/default.dreego": `<body><html lang="de"><head>{#head}</head><body>{#slot}</body></html></body>`,
		"www/routes/index.dreego":    `<body><p>[[ home.greeting name="Ada" ]]</p><p>[[ cart.items count=2 ]]</p></body>`,
	})
	headers := map[string]string{"Accept-Language": "en-GB,en;q=0.8"}
	code, body, _ := client.Request(t, http.MethodGet, "/", "", headers)
	if code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", code, body)
	}
	if !strings.Contains(body, "Hello, Ada!") || !strings.Contains(body, "2 items") {
		t.Fatalf("unexpected body: %s", body)
	}
	if !strings.Contains(body, `<html lang="en">`) {
		t.Fatalf("document locale metadata missing: %s", body)
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

func TestI18nGenerationRequiresEnabledConfiguration(t *testing.T) {
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/dreego.config.json":  `{}`,
		"www/routes/index.dreego": `<body>[[ home.title ]]</body>`,
	})
	output, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil || !strings.Contains(output, "enable i18n") {
		t.Fatalf("output = %q, error = %v", output, err)
	}
}

func TestI18nExtraction(t *testing.T) {
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/dreego.config.json":       `{"i18n":{"enabled":true,"defaultLocale":"de","locales":["de","en"]}}`,
		"www/locales/de/messages.json": `{"home.title":"Willkommen"}`,
		"www/locales/en/messages.json": `{"home.title":"Welcome"}`,
		"www/routes/index.dreego":      `<body>[[ home.title ]]</body>`,
	})
	output, err := dreegotest.RunCLI(t, dir, "i18n", "extract")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output, `"sourceLocale": "de"`) || !strings.Contains(output, `"id": "home.title"`) {
		t.Fatalf("unexpected extraction: %s", output)
	}
}

func TestI18nGeneratedArgumentsRemainTyped(t *testing.T) {
	dreegotest.MustBuildFail(t, map[string]string{
		"www/dreego.config.json":       `{"i18n":{"enabled":true,"defaultLocale":"en","locales":["en"]}}`,
		"www/locales/en/messages.json": `{"cart.items":{"plural":{"argument":"count","cases":{"one":"One","other":"Many"}}}}`,
		"www/routes/index.dreego":      `<body>[[ cart.items count="two" ]]</body>`,
	})
}

func TestI18nBuildsAcrossTemplateContextsAndMethods(t *testing.T) {
	config := `{"i18n":{"enabled":true,"defaultLocale":"en","locales":["en"]}}`
	catalog := `{
		"layout.footer":"Footer",
		"component.label":"Label",
		"page.title":"Title",
		"page.description":"Description",
		"page.heading":"Heading",
		"method.saved":"Saved"
	}`
	dreegotest.MustBuild(t, map[string]string{
		"www/dreego.config.json":       config,
		"www/locales/en/messages.json": catalog,
		"www/layouts/default.dreego":   `<body><html><head>{#head}</head><body>{#slot}<footer>[[ layout.footer ]]</footer></body></html></body>`,
		"www/components/Label.dreego":  `Component Label ()` + "\n" + `<body><span>[[ component.label ]]</span></body>`,
		"www/routes/+page.dreego":      `<head><title>[[ page.title ]]</title><meta name="description" content="[[ page.description ]]"/></head><body><h1>[[ page.heading ]]</h1><@Label/></body><body method="post">[[ method.saved ]]</body>`,
		"www/routes/docs/+page.dreego": `<body lang="md"># [[ page.heading ]]</body>`,
	})
}
