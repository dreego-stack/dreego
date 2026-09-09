package core

import (
	"context"
	"log/slog"

	corei18n "github.com/dreego-stack/dreego/core/internal/i18n"
)

type MessageArg = corei18n.Argument
type MessageArgumentFormat = corei18n.ArgumentFormat
type MessageValue = corei18n.Value
type MessageSelector = corei18n.Selector
type LocalizedMessage = corei18n.Message
type LocaleCatalog = corei18n.LocaleCatalog
type I18nConfig = corei18n.Config
type Localizer = corei18n.Localizer

var MessageText = corei18n.Text
var MessagePlural = corei18n.Plural
var MessageSelect = corei18n.Select

func NewLocalizer(config I18nConfig) (*corei18n.CatalogLocalizer, error) {
	return corei18n.NewCatalogLocalizer(config)
}

func Localize(ctx context.Context, key string, arguments ...MessageArg) (string, error) {
	localizer, locale, ok := corei18n.FromContext(ctx)
	if !ok {
		return key, nil
	}
	return localizer.Localize(ctx, locale, key, arguments)
}

func Message(ctx context.Context, key string, arguments ...MessageArg) string {
	value, err := Localize(ctx, key, arguments...)
	if err != nil {
		slog.Error("dreego: message localization failed", "key", key, "error", err)
		return key
	}
	return value
}

func Locale(ctx context.Context) string {
	return corei18n.Locale(ctx)
}
