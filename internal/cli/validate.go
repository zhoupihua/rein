package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhoupihua/rein/internal/output"
	"github.com/zhoupihua/rein/internal/project"
)

var validateCmd = &cobra.Command{
	Use:   "validate [feature]",
	Short: "Validate artifact completeness for a feature",
	Long: `Validate that all required artifacts exist for each workflow phase.

Required artifacts:
  DEFINE: spec.md
  PLAN:   plan.md, task.md
  REVIEW: review.md

A phase is COMPLETE only when ALL its required artifacts exist.
A feature is READY for its current phase only when all prior phases are complete.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	p, err := project.Resolve()
	if err != nil {
		output.PrintError(err, isJSON())
		return nil
	}

	if len(args) > 0 {
		vr := project.Validate(p, args[0])
		output.Print(vr, isJSON())
		return nil
	}

	if !project.HasArtifacts(p) {
		output.Print(map[string]string{
			"status":     "NONE",
			"suggestion": "No features found. Run /spec or /triage to start.",
		}, isJSON())
		return nil
	}

	var results []project.ValidateResult
	allReady := true
	for _, name := range p.Changes {
		vr := project.Validate(p, name)
		results = append(results, vr)
		if !vr.Ready {
			allReady = false
		}
	}

	if isJSON() {
		output.PrintJSON(results)
	} else {
		for _, vr := range results {
			fmt.Printf("Feature: %s\n", vr.Feature)
			fmt.Printf("  Current phase: %s\n", vr.Current)
			if vr.Ready {
				fmt.Printf("  Status: READY\n")
			} else {
				fmt.Printf("  Status: NOT READY — %s\n", project.FormatMissing(vr.Phases))
			}
			for _, pr := range vr.Phases {
				status := "✓"
				if !pr.Complete {
					status = "✗"
				}
				fmt.Printf("  %s %s: ", status, pr.Phase)
				for i, a := range pr.Artifacts {
					if i > 0 {
						fmt.Print(", ")
					}
					if a.Exists {
						fmt.Printf("%s ✓", a.File)
					} else {
						fmt.Printf("%s ✗", a.File)
					}
				}
				fmt.Println()
			}
			fmt.Println()
		}
		if allReady {
			fmt.Println("All features ready.")
		}
	}
	return nil
}
