package artifact

import (
	"strings"
	"testing"
)

func TestParseDeltaSpec(t *testing.T) {
	content := `## ADDED Requirements

### Requirement: User Auth
The system SHALL authenticate users via JWT.
#### Scenario: Valid login
- **WHEN** user submits valid credentials
- **THEN** system returns JWT token

## MODIFIED Requirements

### Requirement: Rate Limiting
The system MUST limit API calls to 100 per minute per user.

## REMOVED Requirements
- ` + "`" + `### Requirement: Legacy Auth` + "`" + `

## RENAMED Requirements
- FROM: ` + "`" + `### Requirement: Old Logging` + "`" + `
- TO: ` + "`" + `### Requirement: Structured Logging` + "`" + `
`

	plan, presence, err := ParseDeltaSpec(content)
	if err != nil {
		t.Fatalf("ParseDeltaSpec failed: %v", err)
	}

	if !presence.Added || !presence.Modified || !presence.Removed || !presence.Renamed {
		t.Error("expected all section presences to be true")
	}

	if len(plan.Added) != 1 || plan.Added[0].Name != "User Auth" {
		t.Errorf("expected 1 added requirement 'User Auth', got %v", plan.Added)
	}
	if len(plan.Modified) != 1 || plan.Modified[0].Name != "Rate Limiting" {
		t.Errorf("expected 1 modified requirement 'Rate Limiting', got %v", plan.Modified)
	}
	if len(plan.Removed) != 1 || plan.Removed[0] != "Legacy Auth" {
		t.Errorf("expected 1 removed requirement 'Legacy Auth', got %v", plan.Removed)
	}
	if len(plan.Renamed) != 1 || plan.Renamed[0].From != "Old Logging" || plan.Renamed[0].To != "Structured Logging" {
		t.Errorf("expected 1 rename from 'Old Logging' to 'Structured Logging', got %v", plan.Renamed)
	}
}

func TestParseDeltaSpecEmpty(t *testing.T) {
	content := `## ADDED Requirements

No requirements here.
`
	_, _, err := ParseDeltaSpec(content)
	if err == nil {
		t.Error("expected error for delta spec with no operations")
	}
}

func TestExtractRequirementsSection(t *testing.T) {
	content := `# My Spec

## Purpose
Some purpose text.

## Requirements

### Requirement: Auth
The system SHALL authenticate.

### Requirement: Logging
The system SHALL log events.

## Decisions

Some decisions.
`

	before, headerLine, _, blocks, after := ExtractRequirementsSection(content)

	if !strings.Contains(before, "## Purpose") {
		t.Error("before should contain Purpose section")
	}
	if headerLine != "## Requirements" {
		t.Errorf("expected header '## Requirements', got %q", headerLine)
	}
	if len(blocks) != 2 {
		t.Fatalf("expected 2 requirement blocks, got %d", len(blocks))
	}
	if blocks[0].Name != "Auth" {
		t.Errorf("expected first block 'Auth', got %q", blocks[0].Name)
	}
	if blocks[1].Name != "Logging" {
		t.Errorf("expected second block 'Logging', got %q", blocks[1].Name)
	}
	if !strings.Contains(after, "## Decisions") {
		t.Error("after should contain Decisions section")
	}
}

func TestExtractRequirementsSectionMissing(t *testing.T) {
	content := `# My Spec

## Purpose
No requirements yet.
`
	before, headerLine, _, blocks, _ := ExtractRequirementsSection(content)
	if headerLine != "" {
		t.Errorf("expected empty headerLine, got %q", headerLine)
	}
	if len(blocks) != 0 {
		t.Errorf("expected 0 blocks, got %d", len(blocks))
	}
	if before != content {
		t.Error("before should equal entire content when no Requirements section")
	}
}

func TestMergeDeltasAdded(t *testing.T) {
	base := `# My Spec

## Requirements

### Requirement: Auth
The system SHALL authenticate.
`

	plan := &DeltaPlan{
		Added: []RequirementBlock{
			{
				HeaderLine: "### Requirement: Logging",
				Name:       "Logging",
				Raw:        "### Requirement: Logging\nThe system SHALL log events.",
			},
		},
	}

	result, stats, err := MergeDeltas(base, plan)
	if err != nil {
		t.Fatalf("MergeDeltas failed: %v", err)
	}
	if stats.Added != 1 {
		t.Errorf("expected 1 added, got %d", stats.Added)
	}
	if !strings.Contains(result, "### Requirement: Logging") {
		t.Error("result should contain the added requirement")
	}
	if !strings.Contains(result, "### Requirement: Auth") {
		t.Error("result should still contain the original requirement")
	}
}

func TestMergeDeltasRemoved(t *testing.T) {
	base := `# My Spec

## Requirements

### Requirement: Auth
The system SHALL authenticate.

### Requirement: Logging
The system SHALL log events.
`

	plan := &DeltaPlan{
		Removed: []string{"Logging"},
	}

	result, stats, err := MergeDeltas(base, plan)
	if err != nil {
		t.Fatalf("MergeDeltas failed: %v", err)
	}
	if stats.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", stats.Removed)
	}
	if strings.Contains(result, "### Requirement: Logging") {
		t.Error("result should NOT contain the removed requirement")
	}
	if !strings.Contains(result, "### Requirement: Auth") {
		t.Error("result should still contain the remaining requirement")
	}
}

func TestMergeDeltasRenamed(t *testing.T) {
	base := `# My Spec

## Requirements

### Requirement: Old Logging
The system SHALL log events.
`

	plan := &DeltaPlan{
		Renamed: []RenamePair{
			{From: "Old Logging", To: "Structured Logging"},
		},
	}

	result, stats, err := MergeDeltas(base, plan)
	if err != nil {
		t.Fatalf("MergeDeltas failed: %v", err)
	}
	if stats.Renamed != 1 {
		t.Errorf("expected 1 renamed, got %d", stats.Renamed)
	}
	if strings.Contains(result, "### Requirement: Old Logging") {
		t.Error("result should NOT contain old requirement name")
	}
	if !strings.Contains(result, "### Requirement: Structured Logging") {
		t.Error("result should contain new requirement name")
	}
}

func TestMergeDeltasModified(t *testing.T) {
	base := `# My Spec

## Requirements

### Requirement: Auth
The system SHALL authenticate via password.
`

	plan := &DeltaPlan{
		Modified: []RequirementBlock{
			{
				HeaderLine: "### Requirement: Auth",
				Name:       "Auth",
				Raw:        "### Requirement: Auth\nThe system SHALL authenticate via JWT.",
			},
		},
	}

	result, stats, err := MergeDeltas(base, plan)
	if err != nil {
		t.Fatalf("MergeDeltas failed: %v", err)
	}
	if stats.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", stats.Modified)
	}
	if !strings.Contains(result, "via JWT") {
		t.Error("result should contain the modified content")
	}
	if strings.Contains(result, "via password") {
		t.Error("result should NOT contain the old content")
	}
}

func TestMergeDeltasConflictModifiedAndRemoved(t *testing.T) {
	plan := &DeltaPlan{
		Modified: []RequirementBlock{
			{Name: "Auth", Raw: "### Requirement: Auth\nNew text."},
		},
		Removed: []string{"Auth"},
	}
	_, _, err := MergeDeltas("", plan)
	if err == nil {
		t.Error("expected error for MODIFIED + REMOVED conflict")
	}
}

func TestMergeDeltasConflictModifiedAndAdded(t *testing.T) {
	plan := &DeltaPlan{
		Modified: []RequirementBlock{
			{Name: "Auth", Raw: "### Requirement: Auth\nNew text."},
		},
		Added: []RequirementBlock{
			{Name: "Auth", Raw: "### Requirement: Auth\nNew requirement."},
		},
	}
	_, _, err := MergeDeltas("", plan)
	if err == nil {
		t.Error("expected error for MODIFIED + ADDED conflict")
	}
}

func TestMergeDeltasFullPipeline(t *testing.T) {
	base := `# My Spec

## Requirements

### Requirement: Auth
The system SHALL authenticate via password.

### Requirement: Logging
The system SHALL log events.

### Requirement: Legacy
The system SHALL use legacy format.
`

	plan := &DeltaPlan{
		Renamed: []RenamePair{
			{From: "Logging", To: "Structured Logging"},
		},
		Removed: []string{"Legacy"},
		Modified: []RequirementBlock{
			{
				HeaderLine: "### Requirement: Auth",
				Name:       "Auth",
				Raw:        "### Requirement: Auth\nThe system SHALL authenticate via JWT.",
			},
		},
		Added: []RequirementBlock{
			{
				HeaderLine: "### Requirement: Rate Limiting",
				Name:       "Rate Limiting",
				Raw:        "### Requirement: Rate Limiting\nThe system MUST limit API calls.",
			},
		},
	}

	result, stats, err := MergeDeltas(base, plan)
	if err != nil {
		t.Fatalf("MergeDeltas failed: %v", err)
	}

	if stats.Renamed != 1 || stats.Removed != 1 || stats.Modified != 1 || stats.Added != 1 {
		t.Errorf("expected renamed=1, removed=1, modified=1, added=1, got %+v", stats)
	}

	if !strings.Contains(result, "### Requirement: Auth") {
		t.Error("should contain Auth")
	}
	if !strings.Contains(result, "via JWT") {
		t.Error("Auth should be modified to JWT")
	}
	if strings.Contains(result, "### Requirement: Logging") {
		t.Error("should not contain old Logging name")
	}
	if !strings.Contains(result, "### Requirement: Structured Logging") {
		t.Error("should contain renamed Structured Logging")
	}
	if strings.Contains(result, "### Requirement: Legacy") {
		t.Error("should not contain removed Legacy")
	}
	if !strings.Contains(result, "### Requirement: Rate Limiting") {
		t.Error("should contain added Rate Limiting")
	}
}
