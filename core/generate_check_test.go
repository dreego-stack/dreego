package core

import (
	"strings"
	"testing"
)

func TestPlanDiffEmptyOnEqual(t *testing.T) {
	plan := genPlan{
		files: map[string]string{
			"dreego/gen/routes.go": "package gen\n",
			"dreego/gen/dree.go":   "package gen\n",
		},
	}
	disk := map[string]string{
		"dreego/gen/routes.go": "package gen\n",
		"dreego/gen/dree.go":   "package gen\n",
	}
	diff := plan.diff(disk)
	if len(diff) != 0 {
		t.Fatalf("expected no diff, got: %v", diff)
	}
}

func TestPlanDiffMissingFile(t *testing.T) {
	plan := genPlan{
		files: map[string]string{
			"dreego/gen/routes.go":     "package gen\n",
			"dreego/gen/components.go": "package gen\n",
		},
	}
	disk := map[string]string{
		"dreego/gen/routes.go": "package gen\n",
	}
	diff := plan.diff(disk)
	if len(diff) != 1 {
		t.Fatalf("expected 1 diff, got: %v", diff)
	}
	if diff[0].path != "dreego/gen/components.go" || diff[0].kind != diffMissing {
		t.Fatalf("expected missing components.go, got: %v", diff)
	}
}

func TestPlanDiffExtraFile(t *testing.T) {
	plan := genPlan{
		files: map[string]string{
			"dreego/gen/routes.go": "package gen\n",
		},
	}
	disk := map[string]string{
		"dreego/gen/routes.go": "package gen\n",
		"dreego/gen/stale.go":  "package gen\n",
	}
	diff := plan.diff(disk)
	if len(diff) != 1 {
		t.Fatalf("expected 1 diff, got: %v", diff)
	}
	if diff[0].path != "dreego/gen/stale.go" || diff[0].kind != diffExtra {
		t.Fatalf("expected extra stale.go, got: %v", diff)
	}
}

func TestPlanDiffStaleContent(t *testing.T) {
	plan := genPlan{
		files: map[string]string{
			"dreego/gen/routes.go": "package gen // new\n",
		},
	}
	disk := map[string]string{
		"dreego/gen/routes.go": "package gen // old\n",
	}
	diff := plan.diff(disk)
	if len(diff) != 1 {
		t.Fatalf("expected 1 diff, got: %v", diff)
	}
	if diff[0].path != "dreego/gen/routes.go" || diff[0].kind != diffStale {
		t.Fatalf("expected stale routes.go, got: %v", diff)
	}
}

func TestPlanDiffReportFormat(t *testing.T) {
	plan := genPlan{
		files: map[string]string{
			"dreego/gen/routes.go": "new\n",
		},
	}
	disk := map[string]string{
		"dreego/gen/routes.go": "old\n",
		"dreego/gen/extra.go":  "x\n",
	}
	diff := plan.diff(disk)
	rep := diffReport(diff)
	if !strings.Contains(rep, "routes.go") {
		t.Fatalf("report missing routes.go: %s", rep)
	}
	if !strings.Contains(rep, "extra.go") {
		t.Fatalf("report missing extra.go: %s", rep)
	}
}
