package hook

import (
	"os"
	"regexp"
)

var secretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`AKIA[0-9A-Z]{16}`),           // AWS Access Key
	regexp.MustCompile(`sk-[a-zA-Z0-9]{32,}`),         // OpenAI API Key
	regexp.MustCompile(`ghp_[a-zA-Z0-9]{36}`),         // GitHub PAT
	regexp.MustCompile(`-----BEGIN (?:RSA |EC )?PRIVATE KEY-----`), // SSH/Private Key
	regexp.MustCompile(`sk_live_[a-zA-Z0-9]{24,}`),    // Stripe Live Key
}

func LeakGuard() {
	// Non-Claude IDEs have no PostRead/PostBash hook equivalent
	if IDE() != "claude" {
		return
	}

	// Read tool result from environment or file
	result := os.Getenv("CLAUDE_TOOL_RESULT")
	if result == "" {
		resultPath := os.Getenv("CLAUDE_TOOL_RESULT_FILE_PATH")
		if resultPath != "" {
			data, err := os.ReadFile(resultPath)
			if err == nil {
				result = string(data)
			}
		}
	}

	if result == "" {
		return
	}

	for _, pattern := range secretPatterns {
		if pattern.MatchString(result) {
			OutputBlock("Potential secret detected in output. Review before proceeding.")
			return
		}
	}
}
