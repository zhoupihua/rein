package artifact

import (
	"path/filepath"
	"testing"
)

func TestParseProposalAllSections(t *testing.T) {
	path := filepath.Join("testdata", "sample-proposal.md")
	prop, err := ParseProposalFile(path)
	if err != nil {
		t.Fatalf("ParseProposalFile: %v", err)
	}

	if prop.Why == "" {
		t.Error("Why is empty")
	}
	if prop.WhatChanges == "" {
		t.Error("WhatChanges is empty")
	}
	if prop.Goals == "" {
		t.Error("Goals is empty")
	}
	if prop.NonGoals == "" {
		t.Error("NonGoals is empty")
	}
	if prop.Assumptions == "" {
		t.Error("Assumptions is empty")
	}
	if prop.OpenQuestions == "" {
		t.Error("OpenQuestions is empty")
	}

	expectedSections := []string{"Why", "What Changes", "Goals", "Non-Goals", "Key Assumptions", "Open Questions"}
	for _, s := range expectedSections {
		if !prop.HasSection(s) {
			t.Errorf("missing section %q", s)
		}
	}
}

func TestParseProposalPartialSections(t *testing.T) {
	content := `# Test Proposal

## Why

Some motivation here.

## Goals

- Goal one
`
	prop := ParseProposalContent(content)
	if prop.Why == "" {
		t.Error("Why should be populated")
	}
	if prop.Goals == "" {
		t.Error("Goals should be populated")
	}
	if prop.WhatChanges != "" {
		t.Error("WhatChanges should be empty for partial proposal")
	}
	if len(prop.Sections) != 2 {
		t.Errorf("expected 2 sections, got %d", len(prop.Sections))
	}
}

func TestParseProposalEmpty(t *testing.T) {
	prop := ParseProposalContent("")
	if len(prop.Sections) != 0 {
		t.Errorf("expected 0 sections, got %d", len(prop.Sections))
	}
}

func TestParseProposalSectionContent(t *testing.T) {
	content := `# Test

## Why

Line one.
Line two.

## Goals

- Item A
- Item B
`
	prop := ParseProposalContent(content)
	if prop.Why != "Line one.\nLine two." {
		t.Errorf("Why = %q, want multi-line content", prop.Why)
	}
	if !contains(prop.Goals, "Item A") {
		t.Errorf("Goals = %q, want to contain 'Item A'", prop.Goals)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
