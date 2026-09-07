package output

import (
	"fmt"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

type Artifact struct {
	Code string
}

func GenClient(artifact Artifact) string {
	return GenClientTo(artifact, "b", "\t")
}

func GenClientTo(artifact Artifact, builder, indent string) string {
	return fmt.Sprintf("%s%s.WriteString(\"<script>\")\n%s%s.WriteString(%s)\n%s%s.WriteString(\"</script>\")\n", indent, builder, indent, builder, ir.GoLiteral(artifact.Code), indent, builder)
}
