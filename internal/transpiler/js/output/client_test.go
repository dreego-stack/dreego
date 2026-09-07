package output

import "testing"

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
