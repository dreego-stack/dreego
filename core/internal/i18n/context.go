package i18n

import "context"

type contextState struct {
	localizer Localizer
	locale    string
}

type contextKey struct{}

func WithContext(ctx context.Context, localizer Localizer, locale string) context.Context {
	return context.WithValue(ctx, contextKey{}, contextState{localizer: localizer, locale: locale})
}

func FromContext(ctx context.Context) (Localizer, string, bool) {
	state, ok := ctx.Value(contextKey{}).(contextState)
	return state.localizer, state.locale, ok && state.localizer != nil
}

func Locale(ctx context.Context) string {
	state, _ := ctx.Value(contextKey{}).(contextState)
	return state.locale
}
