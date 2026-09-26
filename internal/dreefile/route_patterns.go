package dreefile

import "strings"

func buildPageName(rel string) string {
	parts := []string{}
	for seg := range strings.SplitSeq(rel, "/") {
		if seg == "" {
			continue
		}
		parts = append(parts, cleanSegment(seg))
	}
	if len(parts) == 0 {
		return "index"
	}
	return strings.Join(parts, "_")
}

func buildPattern(rel string) string {
	if rel == "" || rel == "." {
		return "/{$}"
	}
	segments := []string{}
	for seg := range strings.SplitSeq(rel, "/") {
		if seg == "" {
			continue
		}
		s := patternSegment(seg)
		if s != "" {
			segments = append(segments, s)
		}
	}
	if len(segments) == 0 {
		return "/{$}"
	}
	return "/" + strings.Join(segments, "/")
}

func errorCatchPattern(dirPattern string) string {
	if before, ok := strings.CutSuffix(dirPattern, "/{$}"); ok {
		return before + "/{p...}"
	}
	return dirPattern + "/{p...}"
}

func cleanSegment(seg string) string {
	for {
		if strings.HasPrefix(seg, "[[") && strings.HasSuffix(seg, "]]") {
			return ""
		}
		if strings.HasPrefix(seg, "[") && strings.HasSuffix(seg, "]") {
			seg = seg[1 : len(seg)-1]
			continue
		}
		if strings.HasPrefix(seg, "_") && strings.HasSuffix(seg, "_") {
			seg = seg[1 : len(seg)-1]
			continue
		}
		if strings.HasPrefix(seg, "(") && strings.HasSuffix(seg, ")") {
			seg = seg[1 : len(seg)-1]
			continue
		}
		return seg
	}
}

func patternSegment(seg string) string {
	if strings.HasPrefix(seg, "(") && strings.HasSuffix(seg, ")") {
		return ""
	}
	if strings.HasPrefix(seg, "[[") && strings.HasSuffix(seg, "]]") {
		return ""
	}
	wrapped := false
	for {
		if strings.HasPrefix(seg, "[") && strings.HasSuffix(seg, "]") {
			seg = seg[1 : len(seg)-1]
			wrapped = true
			continue
		}
		if strings.HasPrefix(seg, "_") && strings.HasSuffix(seg, "_") {
			seg = seg[1 : len(seg)-1]
			wrapped = true
			continue
		}
		break
	}
	if !wrapped {
		return seg
	}
	if seg == "" {
		return ""
	}
	if after, ok := strings.CutPrefix(seg, "..."); ok {
		return "{" + after + "...}"
	}
	return "{" + seg + "}"
}

func doubleBracketSegment(rel string) string {
	for seg := range strings.SplitSeq(rel, "/") {
		if strings.HasPrefix(seg, "[[") && strings.HasSuffix(seg, "]]") {
			return seg
		}
	}
	return ""
}
