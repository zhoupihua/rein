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
	EpicFile    = "epic.md"
)

// PhaseArtifact defines which artifacts each phase must produce.
var PhaseArtifact = map[string][]string{
	"DEFINE": {"spec.md"},
	"PLAN":   {"plan.md", "task.md"},
	"REVIEW": {"review.md"},
}

// PhaseOrder is the canonical order of phases.
var PhaseOrder = []string{"DEFINE", "PLAN", "BUILD", "REVIEW", "SHIP"}

type Project struct {
	Dir     string
	Changes []string // standalone feature names
	Epics   []Epic   // epics with their increments
}

// Epic represents an epic directory containing multiple increments.
type Epic struct {
	Name       string
	Increments []string // increment names (subdirectories)
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
	p.Changes, p.Epics = findFeaturesAndEpics(dir, ChangesDir)
	return p, nil
}

func HasArtifacts(p *Project) bool {
	return len(p.Changes) > 0 || len(p.Epics) > 0
}

// IsEpic checks whether a name refers to an epic in the project.
func IsEpic(p *Project, name string) bool {
	for _, e := range p.Epics {
		if e.Name == name {
			return true
		}
	}
	return false
}

// FindEpic returns the epic by name, or nil.
func FindEpic(p *Project, name string) *Epic {
	for i := range p.Epics {
		if p.Epics[i].Name == name {
			return &p.Epics[i]
		}
	}
	return nil
}

// FeatureArtifacts tracks which artifact files exist for a feature.
type FeatureArtifacts struct {
	Name   string
	Spec   bool
	Plan   bool
	Task   bool
	Review bool
}

func ResolveFeature(p *Project, name string) FeatureArtifacts {
	return resolveFeatureFromDir(filepath.Join(p.Dir, ChangesDir, name), name)
}

// ResolveIncrement resolves artifacts for an increment within an epic.
func ResolveIncrement(p *Project, epicName, incrementName string) FeatureArtifacts {
	return resolveFeatureFromDir(
		filepath.Join(p.Dir, ChangesDir, epicName, incrementName),
		epicName+"/"+incrementName,
	)
}

func resolveFeatureFromDir(base, name string) FeatureArtifacts {
	fa := FeatureArtifacts{Name: name}
	if _, err := os.Stat(base); err != nil {
		return fa
	}
	fa.Spec = fileExists(filepath.Join(base, "spec.md"))
	fa.Plan = fileExists(filepath.Join(base, "plan.md"))
	fa.Task = fileExists(filepath.Join(base, "task.md"))
	fa.Review = fileExists(filepath.Join(base, "review.md"))
	return fa
}

// ArtifactStatus reports whether a single artifact exists.
type ArtifactStatus struct {
	File   string `json:"file"`
	Exists bool   `json:"exists"`
	Dir    string `json:"dir"`
}

// PhaseResult reports completeness of a single phase.
type PhaseResult struct {
	Phase     string           `json:"phase"`
	Complete  bool             `json:"complete"`
	Artifacts []ArtifactStatus `json:"artifacts"`
	Missing   []string         `json:"missing,omitempty"`
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
	return validateFromArtifacts(fa, filepath.Join(ChangesDir, name))
}

// ValidateIncrement checks all phase artifacts for an increment within an epic.
func ValidateIncrement(p *Project, epicName, incrementName string) ValidateResult {
	fa := ResolveIncrement(p, epicName, incrementName)
	return validateFromArtifacts(fa, filepath.Join(ChangesDir, epicName, incrementName))
}

func validateFromArtifacts(fa FeatureArtifacts, artifactDir string) ValidateResult {
	result := ValidateResult{Feature: fa.Name}

	artifactMap := map[string]bool{
		"spec.md":   fa.Spec,
		"plan.md":   fa.Plan,
		"task.md":   fa.Task,
		"review.md": fa.Review,
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
			exists := artifactMap[f]
			pr.Artifacts = append(pr.Artifacts, ArtifactStatus{
				File:   f,
				Exists: exists,
				Dir:    artifactDir,
			})
			if !exists {
				pr.Missing = append(pr.Missing, f)
			}
		}
		pr.Complete = len(pr.Missing) == 0
		result.Phases = append(result.Phases, pr)
	}

	result.Ready = isReadyForPhase(result, currentPhase)
	return result
}

// DeterminePhase returns the first incomplete phase based on missing artifacts.
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
	Name   string         `json:"name"`
	Phase  string         `json:"phase"`
	Ready  bool           `json:"ready"`
	Phases []PhaseResult  `json:"phases"`
	Build  *BuildProgress `json:"build,omitempty"`
}

// StatusOf returns full status for a feature including build progress.
func StatusOf(p *Project, name string) FeatureStatus {
	vr := Validate(p, name)
	return statusFromValidateResult(vr, filepath.Join(p.Dir, ChangesDir, name, "task.md"))
}

// StatusOfIncrement returns full status for an increment within an epic.
func StatusOfIncrement(p *Project, epicName, incrementName string) FeatureStatus {
	vr := ValidateIncrement(p, epicName, incrementName)
	return statusFromValidateResult(vr, filepath.Join(p.Dir, ChangesDir, epicName, incrementName, "task.md"))
}

func statusFromValidateResult(vr ValidateResult, taskPath string) FeatureStatus {
	fs := FeatureStatus{
		Name:   vr.Feature,
		Phase:  vr.Current,
		Ready:  vr.Ready,
		Phases: vr.Phases,
	}

	if fileExists(taskPath) {
		fs.Build = &BuildProgress{}
		content, err := os.ReadFile(taskPath)
		if err == nil {
			for _, line := range strings.Split(string(content), "\n") {
				trimmed := strings.TrimSpace(line)
				if strings.HasPrefix(trimmed, "- [x]") {
					fs.Build.Done++
					fs.Build.Total++
				} else if strings.HasPrefix(trimmed, "- [ ]") {
					fs.Build.Total++
				}
			}
		}
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

// EpicStatusResult reports the overall status of an epic.
type EpicStatusResult struct {
	Name       string          `json:"name"`
	Increments []FeatureStatus `json:"increments"`
	Overall    string          `json:"overall"` // "not started", "in progress", "all shipped"
}

// EpicStatus returns status for all increments in an epic.
func EpicStatus(p *Project, epicName string) EpicStatusResult {
	epic := FindEpic(p, epicName)
	if epic == nil {
		return EpicStatusResult{Name: epicName, Overall: "not found"}
	}

	result := EpicStatusResult{Name: epicName}
	allShipped := true
	anyStarted := false
	for _, inc := range epic.Increments {
		fs := StatusOfIncrement(p, epicName, inc)
		result.Increments = append(result.Increments, fs)
		if fs.Phase != "SHIP" {
			allShipped = false
		}
		if fs.Phase != "DEFINE" || fs.Ready {
			anyStarted = true
		}
	}

	switch {
	case allShipped && len(result.Increments) > 0:
		result.Overall = "all shipped"
	case anyStarted:
		result.Overall = "in progress"
	default:
		result.Overall = "not started"
	}

	return result
}

// FeaturePath returns the full path to a feature directory.
func FeaturePath(projectDir, name string) string {
	return filepath.Join(projectDir, ChangesDir, name)
}

// IncrementPath returns the full path to an increment directory within an epic.
func IncrementPath(projectDir, epicName, incrementName string) string {
	return filepath.Join(projectDir, ChangesDir, epicName, incrementName)
}

// ArchivePath returns the full path to an archive directory.
func ArchivePath(projectDir, name string) string {
	return filepath.Join(projectDir, ArchiveDir, name)
}

// TaskFilePath returns the path to a feature's task.md.
func TaskFilePath(projectDir, name string) string {
	return filepath.Join(projectDir, ChangesDir, name, "task.md")
}

// IncrementTaskFilePath returns the path to an increment's task.md.
func IncrementTaskFilePath(projectDir, epicName, incrementName string) string {
	return filepath.Join(projectDir, ChangesDir, epicName, incrementName, "task.md")
}

// EpicFilePath returns the path to an epic's epic.md.
func EpicFilePath(projectDir, epicName string) string {
	return filepath.Join(projectDir, ChangesDir, epicName, EpicFile)
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
	if len(p.Changes) > 0 {
		return p.Changes[0]
	}
	// Fall back to first increment of first epic
	for _, e := range p.Epics {
		if len(e.Increments) > 0 {
			return e.Name + "/" + e.Increments[0]
		}
	}
	return ""
}

// FindFeatureWithTask returns the first feature name that has a task.md.
func FindFeatureWithTask(p *Project) string {
	for _, name := range p.Changes {
		if fileExists(filepath.Join(p.Dir, ChangesDir, name, "task.md")) {
			return name
		}
	}
	// Also check epic increments
	for _, e := range p.Epics {
		for _, inc := range e.Increments {
			if fileExists(filepath.Join(p.Dir, ChangesDir, e.Name, inc, "task.md")) {
				return e.Name + "/" + inc
			}
		}
	}
	return ""
}

// ParseFeatureRef splits a feature reference like "epic/increment" into parts.
// Returns (epicName, incrementName, true) for epic increments,
// or (name, "", false) for standalone features.
func ParseFeatureRef(ref string) (epicName, incrementName string, isIncrement bool) {
	parts := strings.SplitN(ref, "/", 2)
	if len(parts) == 2 {
		return parts[0], parts[1], true
	}
	return ref, "", false
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

func findFeaturesAndEpics(root, subdir string) ([]string, []Epic) {
	dir := filepath.Join(root, subdir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil
	}

	var features []string
	var epics []Epic

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		entryPath := filepath.Join(dir, e.Name())

		// Check if this is an epic (has epic.md)
		if fileExists(filepath.Join(entryPath, EpicFile)) {
			epic := Epic{Name: e.Name()}
			// Scan for increment subdirectories
			subEntries, err := os.ReadDir(entryPath)
			if err == nil {
				for _, se := range subEntries {
					if !se.IsDir() {
						continue
					}
					// An increment is a subdirectory that has at least one standard artifact
					subPath := filepath.Join(entryPath, se.Name())
					if hasFeatureArtifacts(subPath) {
						epic.Increments = append(epic.Increments, se.Name())
					}
				}
			}
			sort.Strings(epic.Increments)
			epics = append(epics, epic)
		} else if hasFeatureArtifacts(entryPath) {
			features = append(features, e.Name())
		}
	}

	sort.Strings(features)
	sort.Slice(epics, func(i, j int) bool {
		return epics[i].Name < epics[j].Name
	})

	return features, epics
}

// hasFeatureArtifacts checks if a directory contains at least one standard feature artifact.
func hasFeatureArtifacts(dir string) bool {
	for _, f := range []string{"spec.md", "plan.md", "task.md", "review.md"} {
		if fileExists(filepath.Join(dir, f)) {
			return true
		}
	}
	return false
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
