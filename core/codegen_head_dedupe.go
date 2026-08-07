package core

import "strings"

// dedupeHeadMerge removes layout head tags the route head overrides.
// The route wins: when the route head defines a <title> or a
// <meta name="description">, the layout's corresponding tag is dropped
// from the layout head prefix so the merged output carries exactly one.
// Detection is static (tag structure only), so dynamic tag content like
// <title>{pageTitle}</title> still triggers dedupe.
func dedupeHeadMerge(layoutPrefix, routeHead string) string {
	if strings.Contains(routeHead, "<title") {
		layoutPrefix = stripTitleTag(layoutPrefix)
	}
	if hasMetaDescription(routeHead) {
		layoutPrefix = stripMetaDescriptionTag(layoutPrefix)
	}
	return layoutPrefix
}

func hasMetaDescription(s string) bool {
	return strings.Contains(s, `name="description"`) || strings.Contains(s, `name='description'`)
}

func stripTitleTag(s string) string {
	for {
		open := strings.Index(s, "<title")
		if open < 0 {
			return s
		}
		closeIdx := strings.Index(s[open:], "</title>")
		if closeIdx < 0 {
			return s
		}
		end := open + closeIdx + len("</title>")
		s = s[:open] + s[end:]
	}
}

func stripMetaDescriptionTag(s string) string {
	offset := 0
	for {
		open := strings.Index(s[offset:], "<meta")
		if open < 0 {
			return s
		}
		open += offset
		end := strings.IndexByte(s[open:], '>')
		if end < 0 {
			return s
		}
		tag := s[open : open+end+1]
		if hasMetaDescription(tag) {
			s = s[:open] + s[open+end+1:]
			offset = open
			continue
		}
		offset = open + end + 1
	}
}
