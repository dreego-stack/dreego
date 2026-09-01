package transpiler

import (
	"strings"
	"testing"
)

// genTypedBlocks must emit a json branch and an xml branch when the file has
// typed go blocks.
func TestGenTypedBlocksJsonAndXml(t *testing.T) {
	file := &File{
		Server: []ServerSection{
			{Code: "c.W.Write([]byte(\"{}\"))", ContentType: "json"},
			{Code: "c.W.Write([]byte(\"<a/>\"))", ContentType: "xml"},
		},
	}
	out, err := genTypedBlocks(file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{
		`c.Wants("application/json")`,
		`"application/json; charset=utf-8"`,
		`c.Wants("application/xml")`,
		`"application/xml; charset=utf-8"`,
		`return "", nil`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("genTypedBlocks missing %q, got:\n%s", want, out)
		}
	}
}
