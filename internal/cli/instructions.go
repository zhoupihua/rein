package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zhoupihua/rein/internal/artifact"
	"github.com/zhoupihua/rein/internal/output"
	"github.com/zhoupihua/rein/internal/project"
)

var instructionsCmd = &cobra.Command{
	Use:   "instructions",
	Short: "Generate instructions for AI agents",
}

var instructionsApplyCmd = &cobra.Command{
	Use:   "apply [feature]",
	Short: "Get full context for current task (spec + plan + task)",
	RunE:  runInstructionsApply,
}

var instructionsSpecsCmd = &cobra.Command{
	Use:   "specs",
	Short: "Get instructions for writing specs",
	RunE:  runInstructionsSpecs,
}

var instructionsTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Get instructions for creating plan and tasks from spec",
	RunE:  runInstructionsTasks,
}

func init() {
	rootCmd.AddCommand(instructionsCmd)
	instructionsCmd.AddCommand(instructionsApplyCmd)
	instructionsCmd.AddCommand(instructionsSpecsCmd)
	instructionsCmd.AddCommand(instructionsTasksCmd)
}

type InstructionsResult struct {
	Phase        string               `json:"phase"`
	Feature      string               `json:"feature"`
	Ready        bool                 `json:"ready"`
	Missing      string               `json:"missing,omitempty"`
	CurrentTask  *artifact.Task       `json:"currentTask,omitempty"`
	TaskDetail   *artifact.TaskDetail `json:"taskDetail,omitempty"`
	SpecContext  string               `json:"specContext,omitempty"`
	PlanGoal     string               `json:"planGoal,omitempty"`
	Progress     string               `json:"progress"`
}

func resolveFeatureName(p *project.Project, args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return project.FirstFeature(p)
}

func runInstructionsApply(cmd *cobra.Command, args []string) error {
	p, err := project.Resolve()
	if err != nil {
		output.PrintError(err, isJSON())
		return nil
	}

	name := resolveFeatureName(p, args)
	if name == "" {
		output.Print(map[string]string{
			"phase":      "NONE",
			"suggestion": "No artifacts found. Start with /spec or /triage.",
		}, isJSON())
		return nil
	}

	vr := project.Validate(p, name)
	fa := project.ResolveFeature(p, name)

	result := InstructionsResult{
		Phase:   vr.Current,
		Feature: name,
		Ready:   vr.Ready,
	}

	if !vr.Ready {
		result.Missing = project.FormatMissing(vr.Phases)
		result.Progress = fmt.Sprintf("Phase %s incomplete", vr.Current)
		output.Print(result, isJSON())
		return nil
	}

	// Build phase: find current task from task.md
	if fa.Task {
		tf, err := artifact.ParseTaskFile(project.TaskFilePath(p.Dir, name))
		if err == nil {
			task := tf.FirstUnchecked()
			if task == nil {
				result.Phase = "SHIP"
				result.Progress = "all tasks complete"
				output.Print(result, isJSON())
				return nil
			}
			result.CurrentTask = task
			done, total := tf.CountDone()
			result.Progress = fmt.Sprintf("%d/%d tasks done", done, total)

			// Find matching plan detail
			if fa.Plan {
				plan, err := artifact.ParsePlanFile(filepath.Join(p.Dir, project.ChangesDir, name, "plan.md"))
				if err == nil {
					result.PlanGoal = plan.Goal
					detail := plan.FindTaskDetail(task.ID)
					if detail != nil {
						result.TaskDetail = detail
					}
				}
			}
		}
	}

	// Find matching spec context
	if fa.Spec {
		spec, err := artifact.ParseSpecFile(filepath.Join(p.Dir, project.ChangesDir, name, "spec.md"))
		if err == nil {
			var ctx strings.Builder
			ctx.WriteString("Goal: " + spec.Goals + "\n")
			ctx.WriteString("Requirements:\n")
			for _, r := range spec.Requirements {
				ctx.WriteString(fmt.Sprintf("- %s\n", r.Name))
				for _, s := range r.Scenarios {
					ctx.WriteString(fmt.Sprintf("  - %s: WHEN %s THEN %s\n", s.Name, s.When, s.Then))
				}
			}
			result.SpecContext = ctx.String()
		}
	}

	output.Print(result, isJSON())
	return nil
}

func runInstructionsSpecs(cmd *cobra.Command, args []string) error {
	p, err := project.Resolve()
	if err != nil {
		output.PrintError(err, isJSON())
		return nil
	}

	instruction := map[string]string{
		"phase": "DEFINE",
		"task":  "Write refine.md, spec.md, and design.md",
		"template": `# Feature Name — Spec

## Context
Describe the current state and why this change is needed.

## Goals
- What this feature achieves

## Non-Goals
- What this feature explicitly does NOT cover

### Requirement: <name>
#### Scenario: <name>
- **WHEN** <condition>
- **THEN** <expected result>

## Decisions
- Key technical decisions and their rationale

## Risks
- Potential issues and mitigations`,
	}

	name := resolveFeatureName(p, args)
	if name != "" {
		fa := project.ResolveFeature(p, name)
		if fa.Spec {
			instruction["note"] = fmt.Sprintf("Feature '%s' already has spec.md. Consider updating or creating a new one.", name)
		}
	}

	output.Print(instruction, isJSON())
	return nil
}

func runInstructionsTasks(cmd *cobra.Command, args []string) error {
	p, err := project.Resolve()
	if err != nil {
		output.PrintError(err, isJSON())
		return nil
	}

	instruction := map[string]string{
		"phase": "PLAN",
		"task":  "Break spec into plan.md + task.md",
	}

	name := resolveFeatureName(p, args)
	if name != "" {
		fa := project.ResolveFeature(p, name)
		if fa.Spec {
			spec, err := artifact.ParseSpecFile(filepath.Join(p.Dir, project.ChangesDir, name, "spec.md"))
			if err == nil {
				var reqs strings.Builder
				for _, r := range spec.Requirements {
					reqs.WriteString(fmt.Sprintf("- %s (%d scenarios)\n", r.Name, len(r.Scenarios)))
				}
				instruction["specRequirements"] = reqs.String()
			}
		}
	}

	instruction["planTemplate"] = `# Feature Name — Plan

**Goal:** <one-line goal>

### 1.1 <task title>
- **Acceptance:** <how to verify>
- **Verification:** <test method>
- **Dependencies:** <prerequisite task IDs>
- **Files:** <files to modify>
- **Scope:** <what's in/out>
- **Notes:** <implementation tips>`

	instruction["taskTemplate"] = "# Feature Name\n\n## 1. Define\n- [ ] 1.1 <task description> `" + "file.go" + "`\n- [ ] 1.2 <task description>\n\n## 2. Build\n- [ ] 2.1 <task description> `" + "file.go" + "`"

	output.Print(instruction, isJSON())
	return nil
}
