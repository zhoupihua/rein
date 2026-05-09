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
	Phase       string               `json:"phase"`
	Feature     string               `json:"feature"`
	Ready       bool                 `json:"ready"`
	Missing     string               `json:"missing,omitempty"`
	CurrentTask *artifact.Task       `json:"currentTask,omitempty"`
	TaskDetail  *artifact.TaskDetail `json:"taskDetail,omitempty"`
	SpecContext string               `json:"specContext,omitempty"`
	PlanGoal    string               `json:"planGoal,omitempty"`
	Progress    string               `json:"progress"`
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

	// Find matching proposal and spec context
	var ctx strings.Builder
	if fa.Proposal {
		prop, err := artifact.ParseProposalFile(filepath.Join(p.Dir, project.ChangesDir, name, "proposal.md"))
		if err == nil {
			if prop.Goals != "" {
				ctx.WriteString("Goal: " + prop.Goals + "\n")
			}
			if prop.Why != "" {
				ctx.WriteString("Why: " + prop.Why + "\n")
			}
			if prop.NonGoals != "" {
				ctx.WriteString("Non-Goals: " + prop.NonGoals + "\n")
			}
		}
	}
	if fa.Spec {
		spec, err := artifact.ParseSpecFile(filepath.Join(p.Dir, project.ChangesDir, name, "spec.md"))
		if err == nil {
			ctx.WriteString("Requirements:\n")
			for _, r := range spec.Requirements {
				ctx.WriteString(fmt.Sprintf("- %s\n", r.Name))
				for _, s := range r.Scenarios {
					line := fmt.Sprintf("  - %s: WHEN %s THEN %s", s.Name, s.When, s.Then)
					if s.Test != "" {
						line += " TEST " + s.Test
					}
					ctx.WriteString(line + "\n")
				}
			}
			if len(spec.Decisions) > 0 {
				ctx.WriteString("Decisions:\n")
				for _, d := range spec.Decisions {
					ctx.WriteString(fmt.Sprintf("- %s (rationale: %s)\n", d.Choice, d.Rationale))
				}
			}
			if spec.Risks != "" {
				ctx.WriteString("Risks: " + spec.Risks + "\n")
			}
		}
	}
	result.SpecContext = ctx.String()

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
		"task":  "编写 spec.md — 需求、决策和风险",
		"template": `### Requirement: <名称>
#### Scenario: <场景名称>
- **WHEN** <条件>
- **THEN** <期望结果>
- **TEST** ` + "`" + `<测试函数名>` + "`" + ` (可选)

## Decisions
- **Decision:** <关键技术选择> — **Rationale:** <选择理由>

## Risks / Trade-offs
- 潜在问题及应对措施`,
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
		"task":  "将规格拆分为 plan.md + task.md",
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

**Goal:** <一句话目标>

## Architecture Overview
<描述高层架构>

## Dependency Graph
<ASCII树形图展示任务依赖>

## Vertical Slice Strategy
<如何将工作切分为垂直功能切片>

## Risk/Mitigation Table
| Risk | Mitigation |
|------|------------|
| <风险> | <应对策略> |

## Parallelization
| Task | Classification | Notes |
|------|---------------|-------|
| <任务> | safe/sequential/needs-coordination | <备注> |

## Self-Audit Checklist
- [ ] 所有任务都有验收条件
- [ ] 没有占位符值
- [ ] 依赖按顺序满足
- [ ] 每个任务完成后系统可工作

## Handoff
准备执行。所有任务都包含具体的文件路径和函数名。

## Task Details
### 1.1 <任务标题>
- **Acceptance:** <如何验证>
- **Verification:** <测试方法>
- **Dependencies:** <前置任务ID>
- **Files:** <需修改的文件>
- **Scope:** <范围说明>
- **Notes:** <实现提示>
- **Approach:** <实现策略>
- **Edge Cases:** <需处理的边界情况>
- **Rollback:** <回滚方式>`

	instruction["taskTemplate"] = "# Feature Name\n\n## 1. Define\n- [ ] 1.1 <任务描述> `" + "file.go" + "`\n  - [ ] RED: <测试描述>\n  - [ ] GREEN: <实现描述>\n  - [ ] REFACTOR: <重构描述>\n- [ ] 1.2 <任务描述>\n\n## 2. Build\n- [ ] 2.1 <任务描述> `" + "file.go" + "`"

	output.Print(instruction, isJSON())
	return nil
}
