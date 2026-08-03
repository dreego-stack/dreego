package core

func opt[T any](opts *Options, fn func(*Options) T) T {
	if opts == nil {
		var zero T
		return zero
	}
	return fn(opts)
}

func optStr(opts *Options, fn func(*Options) string, def string) string {
	if opts == nil {
		return def
	}
	v := fn(opts)
	if v == "" {
		return def
	}
	return v
}
