package md

import shared "github.com/dreego-stack/dreego/internal/md"

type Renderer = shared.Renderer

var (
	atxHeading     = shared.AtxHeading
	ulItem         = shared.ULItem
	olItem         = shared.OLItem
	htmlBlockStart = shared.HTMLBlockStart
	footnoteDefRe  = shared.FootnoteDefRe
)

func newRenderer(mode Mode) *Renderer {
	return shared.NewRenderer(mode)
}

func safeFenceLanguage(info string) string {
	return shared.SafeFenceLanguage(info)
}

func isTableSeparator(s string) bool {
	return shared.IsTableSeparator(s)
}

func isHR(s string) bool {
	return shared.IsHR(s)
}

func safeURL(raw string, isImage bool) string {
	return shared.SafeURL(raw, isImage)
}
