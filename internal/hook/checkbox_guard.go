package hook

import (
	"strings"
)

func CheckboxGuard() {
	input := ReadToolInput()
	if input == "" {
		return
	}

	target := ExtractFilePath(input)
	if target == "" {
		return
	}

	// Only check task.md files
	if !strings.Contains(target, "docs/rein/changes/") || !strings.HasSuffix(target, "task.md") {
		return
	}

	// Check if the edit included a checkbox toggle [ ] → [x]
	rawInput := ReadToolInput()

	hasOldCheck := strings.Contains(rawInput, "- [ ]")
	hasNewCheck := strings.Contains(rawInput, "- [x]")

	if hasOldCheck && !hasNewCheck {
		OutputAdditional("PostToolUse", "Warning: Edited task.md without toggling any checkboxes. Did you complete a task?")
	}
}
