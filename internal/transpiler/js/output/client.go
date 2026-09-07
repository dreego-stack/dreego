package output

import (
	"fmt"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

type Artifact struct {
	Code string
}

func GenClient(artifact Artifact) string {
	return fmt.Sprintf("\tb.WriteString(\"<script>\")\n\tb.WriteString(%s)\n\tb.WriteString(\"</script>\")\n", ir.GoLiteral(artifact.Code))
}
