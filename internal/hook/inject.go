package hook

import (
	"os"
	"path/filepath"
)

func Inject() {
	// Non-Claude IDEs have no UserPromptExpansion hook equivalent
	if IDE() != "claude" {
		return
	}

	// Inject review checklist for /code-review
	configDir := ConfigDir()
	checklistPath := filepath.Join(configDir, "checklists", "review.md")

	data, err := os.ReadFile(checklistPath)
	if err != nil {
		return
	}

	OutputAdditional("UserPromptExpansion", string(data))
}
