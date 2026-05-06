package hook

import (
	"os"
	"path/filepath"
)

func Inject() {
	// Inject review checklist for /code-review
	configDir := ConfigDir()
	checklistPath := filepath.Join(configDir, "checklists", "review.md")

	data, err := os.ReadFile(checklistPath)
	if err != nil {
		return
	}

	OutputAdditional("UserPromptExpansion", string(data))
}
