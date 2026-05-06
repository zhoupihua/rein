package project

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	ArtifactDir = "docs/rein"
	ChangesDir  = "docs/rein/changes"
	ArchiveDir  = "docs/rein/archive"
)

// PhaseArtifact defines which artifacts each phase must produce.
var PhaseArtifact = map[string][]string{
	"DEFINE": {"spec.md"},
	"PLAN":   {"plan.md", "task.md"},
	"REVIEW": {"review.md"},
}

// OptionalArtifact lists artifacts that are checked when present but not required.
var OptionalArtifact = map[string][]string{
	"DEFINE": {"refine.md", "design.md"},
}

// PhaseOrder is the canonical order of phases.
var PhaseOrder = []string{"DEFINE", "PLAN", "BUILD", "REVIEW", "SHIP"}

type Project struct {
	Dir     string
	Changes []string // feature names with artifacts in changes/
}

func Resolve() (*Project, error) {
	dir := os.Getenv("CLAUDE_PROJECT_DIR")
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	p := &Project{Dir: dir}
	p.Changes = findFeatures(dir, ChangesDir)
	return p, nil
}

func HasArtifacts(p *Project) bool {
	return len(p.Changes) > 0
}

// FeatureArtifacts tracks which artifact files exist for a feature.
type FeatureArtifacts struct {
	Name    string
	Refine  bool
	Spec    bool
	Design  bool
	Plan    bool
	Task    bool
	Review  bool
}

func ResolveFeature(p *Project, name string) FeatureArtifacts {
	fa := FeatureArtifacts{Name: name}
	base := filepath.Join(p.Dir, ChangesDir, name)
	if _, err := os.Stat(base); err != nil {
		return fa
	}
	fa.Refine = fileExists(filepath.Join(base, "refine.md"))
	fa.Spec = fileExists(filepath.Join(base, "spec.md"))
	fa.Design = fileExists(filepath.Join(base, "design.md"))
	fa.Plan = fileExists(filepath.Join(base, "plan.md"))
	fa.Task = fileExists(filepath.Join(base, "task.md"))
	fa.Review = fileExists(filepath.Join(base, "review.md"))
	return fa
}

// ArtifactStatus reports whether a single artifact exists.
type ArtifactStatus struct {
	File    string `json:"file"`
	Exists  bool   `json:"exists"`
	Dir     string `json:"dir"`
}

// PhaseResult reports completeness of a single phase.
type PhaseResult struct {
	Phase       string            `json:"phase"`
	Complete    bool              `json:"complete"`
	Artifacts   []ArtifactStatus  `json:"artifacts"`
	Optional    []ArtifactStatus  `json:"optional,omitempty"`
	Missing     []string          `json:"missing,omitempty"`
}

// ValidateResult reports completeness of all phases for a feature.
type ValidateResult struct {
	Feature string        `json:"feature"`
	Phases  []PhaseResult `json:"phases"`
	Current string        `json:"currentPhase"`
	Ready   bool          `json:"ready"` // all phases up to current are complete
}

// Validate checks all phase artifacts for a feature.
func Validate(p *Project, name string) ValidateResult {
	fa := ResolveFeature(p, name)
	result := ValidateResult{Feature: name}

	artifactMap := map[string]map[string]bool{
		"refine.md": {"refine": fa.Refine},
		"spec.md":   {"spec": fa.Spec},
		"design.md": {"design": fa.Design},
		"plan.md":   {"plan": fa.Plan},
		"task.md":   {"task": fa.Task},
		"review.md": {"review": fa.Review},
	}

	lookupExists := func(file string) bool {
		if m, ok := artifactMap[file]; ok {
			for _, v := range m {
				return v
			}
		}
		return false
	}

	currentPhase := DeterminePhase(fa)
	result.Current = currentPhase

	for _, phase := range PhaseOrder {
		required, ok := PhaseArtifact[phase]
		if !ok {
			continue // BUILD and SHIP have no file artifacts to check
		}

		pr := PhaseResult{Phase: phase}
		for _, f := range required {
			exists := lookupExists(f)
			pr.Artifacts = append(pr.Artifacts, ArtifactStatus{
				File:   f,
				Exists: exists,
				Dir:    filepath.Join(ChangesDir, name),
			})
			if !exists {
				pr.Missing = append(pr.Missing, f)
			}
		}
		// Add optional artifacts for visibility without affecting completeness
		if optional, ok := OptionalArtifact[phase]; ok {
			for _, f := range optional {
				exists := lookupExists(f)
				pr.Optional = append(pr.Optional, ArtifactStatus{
					File:   f,
					Exists: exists,
					Dir:    filepath.Join(ChangesDir, name),
				})
			}
		}
		pr.Complete = len(pr.Missing) == 0
		result.Phases = append(result.Phases, pr)
	}

	// "ready" means all phases before and including current have no missing artifacts
	result.Ready = isReadyForPhase(result, currentPhase)

	return result
}

// DeterminePhase returns the first incomplete phase based on missing artifacts.
// Optional artifacts (refine.md, design.md) do not block phase progression.
func DeterminePhase(fa FeatureArtifacts) string {
	switch {
	case !fa.Spec:
		return "DEFINE"
	case !fa.Plan:
		return "PLAN"
	case !fa.Task:
		return "PLAN"
	case !fa.Review:
		return "REVIEW"
	default:
		return "SHIP"
	}
}

// BuildProgress returns task completion info if task.md exists.
type BuildProgress struct {
	Done  int `json:"done"`
	Total int `json:"total"`
}

// FeatureStatus combines validation + build progress for display.
type FeatureStatus struct {
	Name    string        `json:"name"`
	Phase   string        `json:"phase"`
	Ready   bool          `json:"ready"`
	Phases  []PhaseResult `json:"phases"`
	Build   *BuildProgress `json:"build,omitempty"`
}

// StatusOf returns full status for a feature including build progress.
func StatusOf(p *Project, name string) FeatureStatus {
	vr := Validate(p, name)
	fs := FeatureStatus{
		Name:   name,
		Phase:  vr.Current,
		Ready:  vr.Ready,
		Phases: vr.Phases,
	}

	// Add build progress if task.md exists
	taskPath := filepath.Join(p.Dir, ChangesDir, name, "task.md")
	if fileExists(taskPath) {
		fs.Build = &BuildProgress{}
		// Parse task counts from file
		content, err := os.ReadFile(taskPath)
		if err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if strings.HasPrefix(trimmed, "- [x]") {
					fs.Build.Done++
					fs.Build.Total++
				} else if strings.HasPrefix(trimmed, "- [ ]") {
					fs.Build.Total++
				}
			}
		}
		// Update phase: if all tasks done but no review, still in REVIEW
		if fs.Build.Done == fs.Build.Total && fs.Build.Total > 0 && vr.Current == "REVIEW" {
			// Already correct
		} else if fs.Build.Done < fs.Build.Total && fs.Build.Total > 0 {
			fs.Phase = "BUILD"
		} else if fs.Build.Done == fs.Build.Total && fs.Build.Total > 0 && !vr.Ready {
			fs.Phase = "REVIEW"
		}
	}

	return fs
}

// FeaturePath returns the full path to a feature directory.
func FeaturePath(projectDir, name string) string {
	return filepath.Join(projectDir, ChangesDir, name)
}

// ArchivePath returns the full path to an archive directory.
func ArchivePath(projectDir, name string) string {
	return filepath.Join(projectDir, ArchiveDir, name)
}

// TaskFilePath returns the path to a feature's task.md.
func TaskFilePath(projectDir, name string) string {
	return filepath.Join(projectDir, ChangesDir, name, "task.md")
}

// ContainsReinPath checks if a path references rein-managed files.
func ContainsReinPath(path string) bool {
	return strings.Contains(path, "docs/rein/changes/") ||
		strings.Contains(path, "docs/rein/archive/")
}

// IsTaskFile checks if a path is a task.md under changes/.
func IsTaskFile(path string) bool {
	return strings.Contains(path, "docs/rein/changes/") && strings.HasSuffix(path, "task.md")
}

// FirstFeature returns the first feature name, or empty string.
func FirstFeature(p *Project) string {
	if len(p.Changes) == 0 {
		return ""
	}
	return p.Changes[0]
}

// FindFeatureWithTask returns the first feature name that has a task.md.
func FindFeatureWithTask(p *Project) string {
	for _, name := range p.Changes {
		if fileExists(filepath.Join(p.Dir, ChangesDir, name, "task.md")) {
			return name
		}
	}
	return ""
}

// isReadyForPhase checks that all phases before and including the target have no missing artifacts.
func isReadyForPhase(vr ValidateResult, target string) bool {
	for _, pr := range vr.Phases {
		if !pr.Complete {
			return false
		}
		if pr.Phase == target {
			return true
		}
	}
	// BUILD/SHIP have no artifact checks — ready if all checked phases complete
	return true
}

func findFeatures(root, subdir string) []string {
	dir := filepath.Join(root, subdir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var features []string
	for _, e := range entries {
		if e.IsDir() {
			features = append(features, e.Name())
		}
	}
	sort.Strings(features)
	return features
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// FormatMissing returns a human-readable string listing missing artifacts.
func FormatMissing(phases []PhaseResult) string {
	var parts []string
	for _, pr := range phases {
		if len(pr.Missing) > 0 {
			parts = append(parts, fmt.Sprintf("%s: missing %s", pr.Phase, strings.Join(pr.Missing, ", ")))
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, "; ")
}
