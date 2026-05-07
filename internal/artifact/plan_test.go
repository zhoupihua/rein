package artifact

import (
	"path/filepath"
	"testing"
)

func TestParsePlanWithSections(t *testing.T) {
	path := filepath.Join("testdata", "sample-plan.md")
	plan, err := ParsePlanFile(path)
	if err != nil {
		t.Fatalf("ParsePlanFile: %v", err)
	}

	if plan.Goal != "Implement user authentication with JWT tokens" {
		t.Errorf("Goal = %q", plan.Goal)
	}
	if plan.Architecture == "" {
		t.Error("Architecture should be populated")
	}
	if plan.DependencyGraph == "" {
		t.Error("DependencyGraph should be populated")
	}
	if plan.SliceStrategy == "" {
		t.Error("SliceStrategy should be populated")
	}
	if plan.RiskTable == "" {
		t.Error("RiskTable should be populated")
	}
	if plan.Parallelization == "" {
		t.Error("Parallelization should be populated")
	}
	if plan.SelfAudit == "" {
		t.Error("SelfAudit should be populated")
	}
	if plan.Handoff == "" {
		t.Error("Handoff should be populated")
	}
}

func TestParsePlanWithNewTaskDetailFields(t *testing.T) {
	path := filepath.Join("testdata", "sample-plan.md")
	plan, err := ParsePlanFile(path)
	if err != nil {
		t.Fatalf("ParsePlanFile: %v", err)
	}

	detail := plan.FindTaskDetail(TaskID{Phase: 1, Seq: 1})
	if detail == nil {
		t.Fatal("task 1.1 not found")
	}
	if detail.Approach != "Interview stakeholders and review OWASP" {
		t.Errorf("Approach = %q", detail.Approach)
	}
	if detail.EdgeCases != "Multi-tenant auth requirements" {
		t.Errorf("EdgeCases = %q", detail.EdgeCases)
	}
	if detail.Rollback != "Delete requirements doc" {
		t.Errorf("Rollback = %q", detail.Rollback)
	}
}

func TestParsePlanBackwardCompat(t *testing.T) {
	content := `# Simple Plan

**Goal:** Fix the bug

### 1.1 Fix null pointer

- **Acceptance:** No more NPE
- **Verification:** Test passes
- **Dependencies:** None
- **Files:** ` + "`main.go`" + `
- **Scope:** Single function
- **Notes:** Check nil before access
`
	plan := ParsePlanContent(content)
	if plan.Goal != "Fix the bug" {
		t.Errorf("Goal = %q", plan.Goal)
	}
	if len(plan.TaskDetails) != 1 {
		t.Fatalf("expected 1 task detail, got %d", len(plan.TaskDetails))
	}
	td := plan.TaskDetails[0]
	if td.Title != "Fix null pointer" {
		t.Errorf("Title = %q", td.Title)
	}
	if td.Approach != "" {
		t.Errorf("Approach should be empty for old format, got %q", td.Approach)
	}
	// New section fields should be empty
	if plan.Architecture != "" || plan.DependencyGraph != "" {
		t.Error("new section fields should be empty for old format plan")
	}
}

func TestParsePlanSectionAlternatives(t *testing.T) {
	content := `# Plan

**Goal:** Test

## Architecture

Some architecture text.

## Risks and Mitigations

Some risk text.
`
	plan := ParsePlanContent(content)
	if plan.Architecture == "" {
		t.Error("Architecture should match 'Architecture' heading")
	}
	if plan.RiskTable == "" {
		t.Error("RiskTable should match 'Risks and Mitigations' heading")
	}
}
