package js

import "github.com/dreego-stack/dreego/internal/transpiler/js/js"

// GenClient emits the <client> section as a passthrough <script> block. The
// client code is JavaScript emitted as-is (identity mapping); the surrounding
// <script> tags are the only generated wrapper.
func GenClient(code string) string {
	return js.GenClient(code)
}
