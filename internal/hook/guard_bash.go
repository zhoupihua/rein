package hook

import (
	"strings"
)

func GuardBash() {
	// Non-Claude IDEs have no PreBash hook equivalent
	if IDE() != "claude" {
		return
	}

	input := ReadToolInput()
	if input == "" {
		return
	}

	cmd := ExtractBashCommand(input)
	if cmd == "" {
		return
	}

	// Check for destructive commands targeting rein-managed files
	destructiveCmds := []string{"rm ", "rm -rf ", "rm -r ", "del ", "rmdir "}
	reinPaths := []string{"docs/rein/", ".rein/commands/", ".rein/skills/", ".rein/agents/", ".rein/hooks/", ".rein/checklists/", ".rein/.rein-manifest", ".rein/bin/"}

	for _, dc := range destructiveCmds {
		if strings.Contains(cmd, dc) {
			for _, rp := range reinPaths {
				if strings.Contains(cmd, rp) {
					OutputBlock("Destructive command targeting rein-managed files is not allowed.")
					return
				}
			}
		}
	}
}
