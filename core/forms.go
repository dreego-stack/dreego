package core

import (
	"regexp"
)

var gActionRE = regexp.MustCompile(`g-action="([^"]+)"`)

func scanFormActions(nodes []TemplateNode) []string {
	var actions []string
	seen := map[string]bool{}
	for _, n := range nodes {
		extractFromNode(n, &actions, seen)
	}
	return actions
}

func extractFromNode(n TemplateNode, actions *[]string, seen map[string]bool) {
	if n.Type == NodeText {
		matches := gActionRE.FindAllStringSubmatch(n.Content, -1)
		for _, m := range matches {
			name := m[1]
			if !seen[name] {
				seen[name] = true
				*actions = append(*actions, name)
			}
		}
	}
	for _, child := range n.Children {
		extractFromNode(child, actions, seen)
	}
	for _, child := range n.ElseChildren {
		extractFromNode(child, actions, seen)
	}
}

func actionFuncName(action string) string {
	return action
}

func findFormStruct(goSections []GoSection, action string) string {
	combined := ""
	for _, g := range goSections {
		combined += g.Code + "\n"
	}
	re := regexp.MustCompile(`func\s+` + regexp.QuoteMeta(action) + `\s*\(\s*\w+\s+\S+,\s*\w+\s+(\w+)\s*\)`)
	matches := re.FindStringSubmatch(combined)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

func findFormHandler(goSections []GoSection, action string) bool {
	combined := ""
	for _, g := range goSections {
		combined += g.Code + "\n"
	}
	handlerRE := regexp.MustCompile(`func\s+` + regexp.QuoteMeta(action) + `\s*\(`)
	return handlerRE.MatchString(combined)
}

func hasValidateTag(goSections []GoSection, structName string) bool {
	return hasTagInStruct(goSections, structName, "validate")
}

func hasFormTag(goSections []GoSection, structName string) bool {
	return hasTagInStruct(goSections, structName, "form")
}

func hasTagInStruct(goSections []GoSection, structName, tagName string) bool {
	combined := ""
	for _, g := range goSections {
		combined += g.Code + "\n"
	}
	startRE := regexp.MustCompile(`type\s+` + regexp.QuoteMeta(structName) + `\s+struct\s*\{`)
	loc := startRE.FindStringIndex(combined)
	if loc == nil {
		return false
	}
	idx := loc[1]
	depth := 1
	for i := idx; i < len(combined) && depth > 0; i++ {
		if combined[i] == '{' {
			depth++
			continue
		}
		if combined[i] == '}' {
			depth--
			if depth == 0 {
				combined = combined[idx:i]
				break
			}
		}
	}
	tagRE := regexp.MustCompile(tagName + `:"[^"]*"`)
	return tagRE.MatchString(combined)
}
