package transpiler

import (
	"strings"
	"testing"
)

func TestGenTemplateNodeFilterRawUpper(t *testing.T) {
	n := TemplateNode{
		Type:    NodeExpression,
		Content: "name",
		Filters: []string{"raw", "upper"},
	}
	result, err := genTemplateNode(NewGenerator(), n, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "strings.ToUpper") {
		t.Errorf("upper filter must wrap in strings.ToUpper, got:\n%s", result)
	}
	if strings.Contains(result, "dreego.SafeText") {
		t.Errorf("raw filter must skip dreego.SafeText, got:\n%s", result)
	}
}
