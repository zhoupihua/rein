package hook

import (
	"os"
	"os/exec"
	"strings"
)

func Gate() {
	if IDE() != "claude" {
		// Cursor/Codex PreCommit: run tests unconditionally
		runTests()
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

	runTests()
}

func runTests() {
	testCmd := exec.Command("npm", "test")
	testCmd.Dir = ProjectDir()
	testCmd.Stdout = os.Stderr
	testCmd.Stderr = os.Stderr

	if err := testCmd.Run(); err != nil {
		BlockExit("Tests failed. Fix tests before deploying.")
	}
}
