package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zhoupihua/rein/internal/artifact"
	"github.com/zhoupihua/rein/internal/project"
)

// TestDefaultArtifactGraphPhaseDetection verifies the default graph detects phases correctly.
func TestDefaultArtifactGraphPhaseDetection(t *testing.T) {
	g := artifact.DefaultArtifactGraph()
	dir := t.TempDir()

	if phase := g.CurrentPhase(dir); phase != "DEFINE" {
		t.Errorf("expected DEFINE with no artifacts, got %q", phase)
	}

	// Only proposal.md — still DEFINE (spec is the gate, not proposal)
	os.WriteFile(filepath.Join(dir, "proposal.md"), []byte("# Proposal"), 0644)
	if phase := g.CurrentPhase(dir); phase != "DEFINE" {
		t.Errorf("expected DEFINE with proposal only, got %q", phase)
	}

	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec"), 0644)
	if phase := g.CurrentPhase(dir); phase != "PLAN" {
		t.Errorf("expected PLAN with spec (proposal optional), got %q", phase)
	}

	os.WriteFile(filepath.Join(dir, "plan.md"), []byte("# Plan"), 0644)
	os.WriteFile(filepath.Join(dir, "task.md"), []byte("# Tasks"), 0644)
	if phase := g.CurrentPhase(dir); phase != "REVIEW" {
		t.Errorf("expected REVIEW with spec+plan+task, got %q", phase)
	}

	os.WriteFile(filepath.Join(dir, "review.md"), []byte("# Review"), 0644)
	if phase := g.CurrentPhase(dir); phase != "SHIP" {
		t.Errorf("expected SHIP with all artifacts, got %q", phase)
	}
}

// TestSchemaJSONLoading verifies the default schema.json can be loaded.
func TestSchemaJSONLoading(t *testing.T) {
	schemaPath := filepath.Join("..", "..", "docs", "rein", "schema.json")
	abs, err := filepath.Abs(schemaPath)
	if err != nil {
		t.Fatalf("cannot resolve schema path: %v", err)
	}

	g, err := artifact.LoadArtifactGraph(abs)
	if err != nil {
		t.Fatalf("LoadArtifactGraph failed: %v", err)
	}

	if g.Name != "define" {
		t.Errorf("expected name 'define', got %q", g.Name)
	}

	if len(g.Artifacts) != 4 {
		t.Fatalf("expected 4 artifacts, got %d", len(g.Artifacts))
	}

	// Verify build order
	order, err := g.BuildOrder()
	if err != nil {
		t.Fatalf("BuildOrder failed: %v", err)
	}
	if order[0].ID != "spec" {
		t.Errorf("expected spec first, got %q", order[0].ID)
	}
	if order[3].ID != "review" {
		t.Errorf("expected review last, got %q", order[3].ID)
	}
}

// TestDeltaSpecParsingAndMerge verifies delta operations work end-to-end.
func TestDeltaSpecParsingAndMerge(t *testing.T) {
	deltaContent := `## ADDED Requirements

### Requirement: User Auth
The system SHALL authenticate users via JWT.

## MODIFIED Requirements

### Requirement: Rate Limiting
The system MUST limit API calls to 100 per minute.

## REMOVED Requirements
- ` + "`" + `### Requirement: Legacy Auth` + "`" + `
`

	plan, presence, err := artifact.ParseDeltaSpec(deltaContent)
	if err != nil {
		t.Fatalf("ParseDeltaSpec failed: %v", err)
	}

	if !presence.Added || !presence.Modified || !presence.Removed {
		t.Error("expected Added, Modified, Removed to be true")
	}

	if len(plan.Added) != 1 || plan.Added[0].Name != "User Auth" {
		t.Errorf("unexpected Added: %v", plan.Added)
	}
	if len(plan.Modified) != 1 || plan.Modified[0].Name != "Rate Limiting" {
		t.Errorf("unexpected Modified: %v", plan.Modified)
	}
	if len(plan.Removed) != 1 || plan.Removed[0] != "Legacy Auth" {
		t.Errorf("unexpected Removed: %v", plan.Removed)
	}

	// Test merge
	base := `# My Spec

## Requirements

### Requirement: Rate Limiting
The system MUST limit API calls to 50 per minute.

### Requirement: Legacy Auth
The system SHALL use legacy authentication.
`

	result, stats, err := artifact.MergeDeltas(base, plan)
	if err != nil {
		t.Fatalf("MergeDeltas failed: %v", err)
	}

	if stats.Added != 1 || stats.Modified != 1 || stats.Removed != 1 {
		t.Errorf("unexpected stats: %+v", stats)
	}

	if result == "" {
		t.Error("expected non-empty result")
	}
}

// TestProjectResolve verifies project resolution works.
func TestProjectResolve(t *testing.T) {
	// Set CLAUDE_PROJECT_DIR to a temp dir
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "docs", "rein", "changes"), 0755)
	os.Setenv("CLAUDE_PROJECT_DIR", dir)
	defer os.Unsetenv("CLAUDE_PROJECT_DIR")

	p, err := project.Resolve()
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if p.Dir != dir {
		t.Errorf("expected Dir=%q, got %q", dir, p.Dir)
	}
	if p.Graph == nil {
		t.Error("expected non-nil Graph")
	}
}

// TestBuildOrderMatchesDefaultGraph verifies schema.json produces same order as default.
func TestBuildOrderMatchesDefaultGraph(t *testing.T) {
	schemaPath := filepath.Join("..", "..", "docs", "rein", "schema.json")
	abs, _ := filepath.Abs(schemaPath)

	loaded, err := artifact.LoadArtifactGraph(abs)
	if err != nil {
		t.Skip("schema.json not found, skipping")
	}

	default_ := artifact.DefaultArtifactGraph()

	loadedOrder, _ := loaded.BuildOrder()
	defaultOrder, _ := default_.BuildOrder()

	if len(loadedOrder) != len(defaultOrder) {
		t.Fatalf("order length mismatch: loaded=%d, default=%d", len(loadedOrder), len(defaultOrder))
	}

	for i := range loadedOrder {
		if loadedOrder[i].ID != defaultOrder[i].ID {
			t.Errorf("position %d: loaded=%q, default=%q", i, loadedOrder[i].ID, defaultOrder[i].ID)
		}
	}
}
