package js

import (
	"fmt"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

// GenClient emits the <client> section as a passthrough <script> block. The
// client code is JavaScript emitted as-is (identity mapping); the surrounding
// <script> tags are the only generated wrapper.
func GenClient(code string) string {
	return fmt.Sprintf("\tb.WriteString(\"<script>\")\n\tb.WriteString(%s)\n\tb.WriteString(\"</script>\")\n", ir.GoLiteral(code))
}
