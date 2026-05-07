package artifact

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultArtifactGraph(t *testing.T) {
	g := DefaultArtifactGraph()
	if g.Name != "define" {
		t.Errorf("expected name 'define', got %q", g.Name)
	}
	if len(g.Artifacts) != 4 {
		t.Errorf("expected 4 artifacts, got %d", len(g.Artifacts))
	}
}

func TestBuildOrder(t *testing.T) {
	g := DefaultArtifactGraph()
	order, err := g.BuildOrder()
	if err != nil {
		t.Fatalf("BuildOrder failed: %v", err)
	}
	if len(order) != 4 {
		t.Fatalf("expected 4 artifacts in order, got %d", len(order))
	}

	// spec comes first (no dependencies)
	if order[0].ID != "spec" {
		t.Errorf("expected first artifact 'spec', got %q", order[0].ID)
	}
	// plan depends on spec
	if order[1].ID != "plan" {
		t.Errorf("expected second artifact 'plan', got %q", order[1].ID)
	}
	// task depends on plan
	if order[2].ID != "task" {
		t.Errorf("expected third artifact 'task', got %q", order[2].ID)
	}
	// review depends on task
	if order[3].ID != "review" {
		t.Errorf("expected fourth artifact 'review', got %q", order[3].ID)
	}
}

func TestBuildOrderCycle(t *testing.T) {
	g := &ArtifactGraph{
		Name: "cyclic",
		Artifacts: []Artifact{
			{ID: "a", Generates: "a.md", Requires: []string{"b"}},
			{ID: "b", Generates: "b.md", Requires: []string{"a"}},
		},
	}
	_, err := g.BuildOrder()
	if err == nil {
		t.Error("expected error for cyclic graph")
	}
}

func TestBuildOrderMissingDependency(t *testing.T) {
	g := &ArtifactGraph{
		Name: "missing-dep",
		Artifacts: []Artifact{
			{ID: "a", Generates: "a.md", Requires: []string{"nonexistent"}},
		},
	}
	_, err := g.BuildOrder()
	if err == nil {
		t.Error("expected error for missing dependency")
	}
}

func TestBuildOrderDiamond(t *testing.T) {
	g := &ArtifactGraph{
		Name: "diamond",
		Artifacts: []Artifact{
			{ID: "a", Generates: "a.md", Requires: []string{}},
			{ID: "b", Generates: "b.md", Requires: []string{"a"}},
			{ID: "c", Generates: "c.md", Requires: []string{"a"}},
			{ID: "d", Generates: "d.md", Requires: []string{"b", "c"}},
		},
	}
	order, err := g.BuildOrder()
	if err != nil {
		t.Fatalf("BuildOrder failed: %v", err)
	}
	if len(order) != 4 {
		t.Fatalf("expected 4 artifacts, got %d", len(order))
	}
	if order[0].ID != "a" {
		t.Errorf("expected 'a' first, got %q", order[0].ID)
	}
	if order[3].ID != "d" {
		t.Errorf("expected 'd' last, got %q", order[3].ID)
	}
}

func TestNextArtifacts(t *testing.T) {
	g := DefaultArtifactGraph()
	dir := t.TempDir()

	// Nothing exists yet — spec should be available (no requirements)
	next := g.NextArtifacts(dir)
	if len(next) != 1 || next[0].ID != "spec" {
		t.Errorf("expected [spec], got %v", artifactIDs(next))
	}

	// Create spec.md — plan should be available
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec"), 0644)
	next = g.NextArtifacts(dir)
	if len(next) != 1 || next[0].ID != "plan" {
		t.Errorf("expected [plan], got %v", artifactIDs(next))
	}

	// Create plan.md — task should be available
	os.WriteFile(filepath.Join(dir, "plan.md"), []byte("# Plan"), 0644)
	next = g.NextArtifacts(dir)
	if len(next) != 1 || next[0].ID != "task" {
		t.Errorf("expected [task], got %v", artifactIDs(next))
	}

	// Create task.md — review should be available
	os.WriteFile(filepath.Join(dir, "task.md"), []byte("# Tasks"), 0644)
	next = g.NextArtifacts(dir)
	if len(next) != 1 || next[0].ID != "review" {
		t.Errorf("expected [review], got %v", artifactIDs(next))
	}

	// Create review.md — nothing should be available
	os.WriteFile(filepath.Join(dir, "review.md"), []byte("# Review"), 0644)
	next = g.NextArtifacts(dir)
	if len(next) != 0 {
		t.Errorf("expected no next artifacts, got %v", artifactIDs(next))
	}
}

func TestBlockedArtifacts(t *testing.T) {
	g := DefaultArtifactGraph()
	dir := t.TempDir()

	// Nothing exists — plan, task, review are blocked
	blocked := g.BlockedArtifacts(dir)
	if len(blocked) != 3 {
		t.Fatalf("expected 3 blocked artifacts, got %d", len(blocked))
	}

	// plan is blocked by spec
	if blocked[0].Artifact.ID != "plan" || len(blocked[0].MissingRequires) != 1 || blocked[0].MissingRequires[0] != "spec" {
		t.Errorf("expected plan blocked by [spec], got %v", blocked[0])
	}
}

func TestCurrentPhase(t *testing.T) {
	g := DefaultArtifactGraph()
	dir := t.TempDir()

	// Nothing exists — DEFINE
	if phase := g.CurrentPhase(dir); phase != "DEFINE" {
		t.Errorf("expected DEFINE, got %q", phase)
	}

	// spec.md exists — PLAN
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec"), 0644)
	if phase := g.CurrentPhase(dir); phase != "PLAN" {
		t.Errorf("expected PLAN, got %q", phase)
	}

	os.WriteFile(filepath.Join(dir, "plan.md"), []byte("# Plan"), 0644)
	os.WriteFile(filepath.Join(dir, "task.md"), []byte("# Tasks"), 0644)
	if phase := g.CurrentPhase(dir); phase != "REVIEW" {
		t.Errorf("expected REVIEW, got %q", phase)
	}

	os.WriteFile(filepath.Join(dir, "review.md"), []byte("# Review"), 0644)
	if phase := g.CurrentPhase(dir); phase != "SHIP" {
		t.Errorf("expected SHIP, got %q", phase)
	}
}

func TestLoadArtifactGraph(t *testing.T) {
	dir := t.TempDir()
	jsonContent := `{
  "name": "test-workflow",
  "artifacts": [
    {"id": "a", "generates": "a.md", "requires": []},
    {"id": "b", "generates": "b.md", "requires": ["a"]}
  ]
}`
	path := filepath.Join(dir, "schema.json")
	os.WriteFile(path, []byte(jsonContent), 0644)

	g, err := LoadArtifactGraph(path)
	if err != nil {
		t.Fatalf("LoadArtifactGraph failed: %v", err)
	}
	if g.Name != "test-workflow" {
		t.Errorf("expected name 'test-workflow', got %q", g.Name)
	}
	if len(g.Artifacts) != 2 {
		t.Errorf("expected 2 artifacts, got %d", len(g.Artifacts))
	}
}

func artifactIDs(arts []Artifact) []string {
	ids := make([]string, len(arts))
	for i, a := range arts {
		ids[i] = a.ID
	}
	return ids
}
