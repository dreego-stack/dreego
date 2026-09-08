package output

import (
	"fmt"
	"regexp"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

var scriptEnd = regexp.MustCompile(`(?i)</script`)

type Artifact struct {
	Code        string
	RuntimePath string
}

func GenClient(artifact Artifact) string {
	return GenClientTo(artifact, "b", "\t")
}

func GenClientTo(artifact Artifact, builder, indent string) string {
	var runtime string
	if artifact.RuntimePath != "" {
		runtime = fmt.Sprintf("%s%s.WriteString(%s)\n", indent, builder, ir.GoLiteral(`<script src="`+artifact.RuntimePath+`"></script>`))
	}
	code := scriptEnd.ReplaceAllString(artifact.Code, `<\/script`)
	return runtime + fmt.Sprintf("%s%s.WriteString(\"<script>\")\n%s%s.WriteString(%s)\n%s%s.WriteString(\"</script>\")\n", indent, builder, indent, builder, ir.GoLiteral(code), indent, builder)
}
