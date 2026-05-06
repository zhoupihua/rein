package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhoupihua/rein/internal/output"
	"github.com/zhoupihua/rein/internal/project"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current workflow phase and artifact completeness",
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

type StatusResult struct {
	Feature  string               `json:"feature"`
	Phase    string               `json:"phase"`
	Ready    bool                 `json:"ready"`
	Phases   []project.PhaseResult `json:"phases"`
	Build    *project.BuildProgress `json:"build,omitempty"`
	Suggestion string             `json:"suggestion"`
}

func runStatus(cmd *cobra.Command, args []string) error {
	p, err := project.Resolve()
	if err != nil {
		output.PrintError(err, isJSON())
		return nil
	}

	if !project.HasArtifacts(p) {
		result := StatusResult{
			Phase:       "NONE",
			Ready:       false,
			Suggestion:  "Run /spec or /triage to start a workflow",
		}
		output.Print(result, isJSON())
		return nil
	}

	// If a feature name is provided as argument, use it; otherwise show all
	features := p.Changes
	if len(args) > 0 {
		features = args
	}

	var results []StatusResult
	for _, name := range features {
		fs := project.StatusOf(p, name)
		results = append(results, StatusResult{
			Feature:    fs.Name,
			Phase:      fs.Phase,
			Ready:      fs.Ready,
			Phases:     fs.Phases,
			Build:      fs.Build,
			Suggestion: suggestNext(fs.Phase, fs.Ready, fs.Phases),
		})
	}

	if len(results) == 1 {
		output.Print(results[0], isJSON())
	} else {
		output.Print(results, isJSON())
	}
	return nil
}

func suggestNext(phase string, ready bool, phases []project.PhaseResult) string {
	if !ready {
		missing := project.FormatMissing(phases)
		if missing != "" {
			return fmt.Sprintf("Phase %s incomplete: %s", phase, missing)
		}
	}

	switch phase {
	case "NONE":
		return "Run /spec or /triage to start a workflow"
	case "DEFINE":
		return "Run /refine then /spec-driven then /spec to produce all DEFINE artifacts"
	case "PLAN":
		return "Run /plan to create plan.md + task.md"
	case "BUILD":
		return "Run rein task next to see current task, or /do to execute"
	case "REVIEW":
		return "Run /code-review for 5-axis review"
	case "SHIP":
		return "Run /ship to finalize and commit"
	default:
		return ""
	}
}
