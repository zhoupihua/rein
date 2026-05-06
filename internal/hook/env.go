package hook

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func ReadToolInput() string {
	input := os.Getenv("CLAUDE_TOOL_INPUT")
	if input == "" {
		path := os.Getenv("CLAUDE_TOOL_INPUT_FILE_PATH")
		if path != "" {
			data, err := os.ReadFile(path)
			if err == nil {
				input = string(data)
			}
		}
	}
	return input
}

func ExtractFilePath(input string) string {
	var obj map[string]any
	if err := json.Unmarshal([]byte(input), &obj); err != nil {
		return ""
	}
	fp, _ := obj["file_path"].(string)
	fp = strings.ReplaceAll(fp, "\\\\", "\\")
	fp = strings.ReplaceAll(fp, "\\", "/")
	return fp
}

func ExtractBashCommand(input string) string {
	var obj map[string]any
	if err := json.Unmarshal([]byte(input), &obj); err != nil {
		return ""
	}
	cmd, _ := obj["command"].(string)
	return cmd
}

func ConfigDir() string {
	dir := os.Getenv("CLAUDE_CONFIG_DIR")
	if dir != "" {
		return dir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude")
}

func ProjectDir() string {
	dir := os.Getenv("CLAUDE_PROJECT_DIR")
	if dir != "" {
		return dir
	}
	dir, _ = os.Getwd()
	return dir
}

func OutputAdditional(event, context string) {
	msg, _ := json.Marshal(map[string]any{
		"hookSpecificOutput": map[string]string{
			"hookEventName":    event,
			"additionalContext": context,
		},
	})
	os.Stdout.Write(msg)
}

func OutputBlock(reason string) {
	msg, _ := json.Marshal(map[string]string{
		"decision": "block",
		"reason":   reason,
	})
	os.Stdout.Write(msg)
}

func ManifestContains(filePath string) bool {
	manifestPath := filepath.Join(ConfigDir(), ".rein-manifest")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return false
	}
	rel := strings.TrimPrefix(filePath, ConfigDir()+"/")
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.TrimSuffix(line, "/") == rel || strings.HasPrefix(rel, line) {
			return true
		}
	}
	return false
}
