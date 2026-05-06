package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zhoupihua/rein/internal/project"
)

func ArtifactValidate() {
	input := ReadToolInput()
	if input == "" {
		return
	}

	target := ExtractFilePath(input)
	if target == "" {
		return
	}

	// Only trigger for files under docs/rein/changes/
	if !strings.Contains(target, "docs/rein/changes/") {
		return
	}

	// Extract feature name from path: docs/rein/changes/<name>/...
	parts := strings.Split(target, "/")
	featureName := ""
	for i, p := range parts {
		if p == "changes" && i+1 < len(parts) {
			featureName = parts[i+1]
			break
		}
	}
	if featureName == "" {
		return
	}

	// Resolve project and validate
	p, err := project.Resolve()
	if err != nil {
		return
	}

	vr := project.Validate(p, featureName)
	fa := project.ResolveFeature(p, featureName)

	// Determine the actual current phase considering BUILD progress
	actualPhase := vr.Current
	buildDone, buildTotal := 0, 0
	if fa.Task {
		taskPath := filepath.Join(p.Dir, project.ChangesDir, featureName, "task.md")
		if content, err := os.ReadFile(taskPath); err == nil {
			for _, line := range strings.Split(string(content), "\n") {
				t := strings.TrimSpace(line)
				if strings.HasPrefix(t, "- [x]") {
					buildDone++
					buildTotal++
				} else if strings.HasPrefix(t, "- [ ]") {
					buildTotal++
				}
			}
		}
		if buildTotal > 0 && buildDone < buildTotal {
			actualPhase = "BUILD"
		} else if buildTotal > 0 && buildDone == buildTotal {
			// All tasks done — phase is REVIEW or SHIP
			if !fa.Review {
				actualPhase = "REVIEW"
			} else {
				actualPhase = "SHIP"
			}
		}
	}

	// Check which artifact phases are incomplete (excluding BUILD/SHIP which have no file artifacts)
	var incompletePhases []string
	for _, pr := range vr.Phases {
		if !pr.Complete {
			incompletePhases = append(incompletePhases, fmt.Sprintf("%s: missing %s", pr.Phase, strings.Join(pr.Missing, ", ")))
		}
	}

	if len(incompletePhases) == 0 {
		// All artifact-based phases complete — congratulate and suggest next
		switch actualPhase {
		case "BUILD":
			if buildTotal > 0 {
				OutputAdditional("PostToolUse",
					fmt.Sprintf("PLAN 阶段完成 ✓ 下一步: 运行 /do 开始实现任务 (%d/%d done)", buildDone, buildTotal))
			} else {
				OutputAdditional("PostToolUse",
					"PLAN 阶段完成 ✓ 下一步: 运行 /do 开始实现任务")
			}
		case "REVIEW":
			OutputAdditional("PostToolUse",
				"BUILD 完成 ✓ 下一步: 运行 /code-review 进行五轴审查")
		case "SHIP":
			OutputAdditional("PostToolUse",
				"REVIEW 阶段完成 ✓ 下一步: 运行 /ship 提交并归档")
		}
		return
	}

	// Report the first incomplete artifact phase
	// Only warn about phases that should have been completed by now
	switch actualPhase {
	case "DEFINE", "PLAN":
		// Only warn about the current phase's missing artifacts
		for _, pr := range vr.Phases {
			if pr.Phase == actualPhase && !pr.Complete {
				OutputAdditional("PostToolUse",
					fmt.Sprintf("⚠️ %s 阶段不完整，缺少: %s", pr.Phase, strings.Join(pr.Missing, ", ")))
				return
			}
		}
	case "BUILD", "REVIEW", "SHIP":
		// Prior phases should all be complete — report first incomplete
		for _, pr := range vr.Phases {
			if !pr.Complete {
				OutputAdditional("PostToolUse",
					fmt.Sprintf("⚠️ %s 阶段不完整，缺少: %s", pr.Phase, strings.Join(pr.Missing, ", ")))
				return
			}
		}
	}
}
