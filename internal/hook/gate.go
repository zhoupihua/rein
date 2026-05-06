package hook

import (
	"os"
	"os/exec"
	"strings"
)

func Gate() {
	input := ReadToolInput()
	if input == "" {
		return
	}

	cmd := ExtractBashCommand(input)
	if cmd == "" {
		return
	}

	// Detect deploy/push/publish commands
	deployKeywords := []string{"deploy", "push", "publish", "release"}
	cmdLower := strings.ToLower(cmd)

	shouldGate := false
	for _, kw := range deployKeywords {
		if strings.Contains(cmdLower, kw) {
			shouldGate = true
			break
		}
	}

	if !shouldGate {
		return
	}

	// Try running tests
	testCmd := exec.Command("npm", "test")
	testCmd.Dir = ProjectDir()
	testCmd.Stdout = os.Stderr
	testCmd.Stderr = os.Stderr

	if err := testCmd.Run(); err != nil {
		OutputBlock("Tests failed. Fix tests before deploying.")
	}
}
