package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func SessionStart() {
	configDir := ConfigDir()
	skillPath := filepath.Join(configDir, "skills", "using-rein", "SKILL.md")

	var context strings.Builder

	// Inject using-rein skill
	data, err := os.ReadFile(skillPath)
	if err == nil {
		context.Write(data)
		context.WriteString("\n\n")
	}

	// Scan for active tasks
	projectDir := ProjectDir()
	tasksDir := filepath.Join(projectDir, "docs", "rein", "tasks")
	entries, err := os.ReadDir(tasksDir)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), "task.md") {
				continue
			}
			content, err := os.ReadFile(filepath.Join(tasksDir, entry.Name()))
			if err != nil {
				continue
			}
			lines := strings.Split(string(content), "\n")
			var unchecked []string
			for _, line := range lines {
				if strings.Contains(line, "- [ ]") {
					unchecked = append(unchecked, strings.TrimSpace(line))
				}
			}
			if len(unchecked) > 0 {
				context.WriteString(fmt.Sprintf("Active tasks in %s:\n", entry.Name()))
				for _, t := range unchecked {
					context.WriteString(t + "\n")
				}
				context.WriteString("\n")
			}
		}
	}

	if context.Len() > 0 {
		OutputAdditional("SessionStart", context.String())
	}
}
