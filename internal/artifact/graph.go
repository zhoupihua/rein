package artifact

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// Artifact defines a single artifact type in the dependency graph.
type Artifact struct {
	ID          string   `json:"id"`
	Generates   string   `json:"generates"`
	Requires    []string `json:"requires"`
	Description string   `json:"description,omitempty"`
}

// ArtifactGraph is a directed acyclic graph of artifact dependencies.
type ArtifactGraph struct {
	Name      string     `json:"name"`
	Artifacts []Artifact `json:"artifacts"`
}

// LoadArtifactGraph reads a schema.json file and returns the graph.
func LoadArtifactGraph(path string) (*ArtifactGraph, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read artifact graph: %w", err)
	}
	var g ArtifactGraph
	if err := json.Unmarshal(data, &g); err != nil {
		return nil, fmt.Errorf("parse artifact graph: %w", err)
	}
	return &g, nil
}

// DefaultArtifactGraph returns the built-in artifact graph used when no schema.json exists.
func DefaultArtifactGraph() *ArtifactGraph {
	return &ArtifactGraph{
		Name: "define",
		Artifacts: []Artifact{
			{ID: "spec", Generates: "spec.md", Requires: []string{}, Description: "Feature Specification with Requirements, Decisions, and Risks"},
			{ID: "plan", Generates: "plan.md", Requires: []string{"spec"}, Description: "Implementation Plan"},
			{ID: "task", Generates: "task.md", Requires: []string{"plan"}, Description: "Task Checklist"},
			{ID: "review", Generates: "review.md", Requires: []string{"task"}, Description: "Code Review Report"},
		},
	}
}

// BuildOrder returns artifacts in topological order (dependencies first).
// Uses Kahn's algorithm. Returns error if the graph contains a cycle.
func (g *ArtifactGraph) BuildOrder() ([]Artifact, error) {
	// Build adjacency list and in-degree map
	inDegree := make(map[string]int)
	dependents := make(map[string][]string) // dep -> list of artifacts that depend on it
	byID := make(map[string]Artifact)

	for _, a := range g.Artifacts {
		byID[a.ID] = a
		if _, ok := inDegree[a.ID]; !ok {
			inDegree[a.ID] = 0
		}
		for _, dep := range a.Requires {
			inDegree[a.ID]++
			dependents[dep] = append(dependents[dep], a.ID)
		}
	}

	// Validate all referenced dependencies exist
	for _, a := range g.Artifacts {
		for _, dep := range a.Requires {
			if _, ok := byID[dep]; !ok {
				return nil, fmt.Errorf("artifact %q requires %q which does not exist", a.ID, dep)
			}
		}
	}

	// Start with nodes that have no dependencies
	var queue []string
	for _, a := range g.Artifacts {
		if inDegree[a.ID] == 0 {
			queue = append(queue, a.ID)
		}
	}
	sort.Strings(queue) // deterministic ordering

	var order []Artifact
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		order = append(order, byID[id])

		deps := dependents[id]
		sort.Strings(deps)
		for _, dep := range deps {
			inDegree[dep]--
			if inDegree[dep] == 0 {
				queue = append(queue, dep)
			}
		}
		sort.Strings(queue)
	}

	if len(order) != len(g.Artifacts) {
		return nil, fmt.Errorf("artifact graph contains a cycle")
	}

	return order, nil
}

// NextArtifacts returns artifacts whose requirements are all satisfied
// (all required artifact files exist in the given directory).
func (g *ArtifactGraph) NextArtifacts(featureDir string) []Artifact {
	existing := g.existingArtifacts(featureDir)
	var next []Artifact
	for _, a := range g.Artifacts {
		if existing[a.ID] {
			continue
		}
		allSatisfied := true
		for _, dep := range a.Requires {
			if !existing[dep] {
				allSatisfied = false
				break
			}
		}
		if allSatisfied {
			next = append(next, a)
		}
	}
	return next
}

// BlockedArtifacts returns artifacts that cannot be created yet,
// along with which dependencies are missing.
type BlockedArtifact struct {
	Artifact        Artifact
	MissingRequires []string
}

func (g *ArtifactGraph) BlockedArtifacts(featureDir string) []BlockedArtifact {
	existing := g.existingArtifacts(featureDir)
	var blocked []BlockedArtifact
	for _, a := range g.Artifacts {
		if existing[a.ID] {
			continue
		}
		// Check if it's blocked (at least one requirement missing)
		var missing []string
		for _, dep := range a.Requires {
			if !existing[dep] {
				missing = append(missing, dep)
			}
		}
		if len(missing) > 0 {
			blocked = append(blocked, BlockedArtifact{
				Artifact:        a,
				MissingRequires: missing,
			})
		}
	}
	return blocked
}

// CurrentPhase returns the current phase name based on which artifacts exist.
func (g *ArtifactGraph) CurrentPhase(featureDir string) string {
	existing := g.existingArtifacts(featureDir)
	order, err := g.BuildOrder()
	if err != nil {
		return "DEFINE"
	}
	// Find the first artifact in build order that doesn't exist yet
	for _, a := range order {
		if !existing[a.ID] {
			return phaseForArtifact(a.ID)
		}
	}
	return "SHIP"
}

// existingArtifacts returns a map of artifact IDs that have their file present.
func (g *ArtifactGraph) existingArtifacts(featureDir string) map[string]bool {
	result := make(map[string]bool)
	for _, a := range g.Artifacts {
		path := fmt.Sprintf("%s/%s", featureDir, a.Generates)
		if _, err := os.Stat(path); err == nil {
			result[a.ID] = true
		}
	}
	return result
}

// phaseForArtifact maps artifact IDs to phase names.
func phaseForArtifact(id string) string {
	switch id {
	case "spec":
		return "DEFINE"
	case "plan", "task":
		return "PLAN"
	case "review":
		return "REVIEW"
	default:
		return "DEFINE"
	}
}
