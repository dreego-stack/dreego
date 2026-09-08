package output

import (
	"strings"
	"testing"
)

func TestGenClientEmitsNormalizedJavaScriptArtifact(t *testing.T) {
	t.Parallel()

	artifact := Artifact{Code: `if (a < b) { console.log("ready") }`}
	want := "\tb.WriteString(\"<script>\")\n" +
		"\tb.WriteString(`if (a < b) { console.log(\"ready\") }`)\n" +
		"\tb.WriteString(\"</script>\")\n"

	if got := GenClient(artifact); got != want {
		t.Fatalf("GenClient() = %q, want %q", got, want)
	}
}

func TestGenClientWritesRuntimeBeforeCode(t *testing.T) {
	got := GenClient(Artifact{RuntimePath: "/_dreego/lua.js", Code: "client();"})
	if !strings.Contains(got, `<script src="/_dreego/lua.js"></script>`) || strings.Index(got, "lua.js") > strings.Index(got, "client();") {
		t.Fatalf("generated client = %s", got)
	}
}

func TestGenClientEscapesHTMLScriptEndSequences(t *testing.T) {
	got := GenClient(Artifact{Code: `const lower = "</script>"; const upper = "</SCRIPT>";`})
	if strings.Count(strings.ToLower(got), `</script>`) != 1 {
		t.Fatalf("generated client contains an early script end: %s", got)
	}
	if !strings.Contains(got, `<\/script>`) {
		t.Fatalf("generated client does not preserve the JavaScript value safely: %s", got)
	}
}
