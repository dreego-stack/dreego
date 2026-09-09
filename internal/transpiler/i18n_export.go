package transpiler

import (
	"fmt"
	"io"
	"path/filepath"

	transpileri18n "github.com/dreego-stack/dreego/internal/transpiler/i18n"
)

func ExtractI18n(writer io.Writer) error {
	roots, err := findWebsiteRoots()
	if err != nil {
		return err
	}
	if len(roots) != 1 {
		return fmt.Errorf("i18n extraction requires exactly one website root, found %d", len(roots))
	}
	settings, err := loadSettings(roots[0])
	if err != nil {
		return err
	}
	if settings == nil || !settings.I18n.Enabled {
		return fmt.Errorf("i18n is not enabled in %s", configFileName)
	}
	set, err := transpileri18n.Load(filepath.Join(roots[0], "locales"), settings.I18n.Locales, settings.I18n.DefaultLocale)
	if err != nil {
		return err
	}
	return transpileri18n.Extract(writer, set)
}
