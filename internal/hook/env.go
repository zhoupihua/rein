package hook

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// IDE returns the current IDE: "cursor" if CURSOR_SESSION is set,
// "codex" if CODEX_SESSION is set, else "claude".
func IDE() string {
	if os.Getenv("CURSOR_SESSION") != "" {
		return "cursor"
	}
	if os.Getenv("CODEX_SESSION") != "" {
		return "codex"
	}
	return "claude"
}

// FilePath returns the target file path from the appropriate source:
// Cursor provides FILE_PATH env var; Claude Code provides JSON tool input.
func FilePath() string {
	if IDE() == "cursor" {
		return os.Getenv("FILE_PATH")
	}
	return ExtractFilePath(ReadToolInput())
}

// BlockExit blocks an action. For Claude Code it outputs JSON; for Cursor/Codex
// it prints to stderr and exits with code 1.
func BlockExit(reason string) {
	if IDE() != "claude" {
		fmt.Fprintln(os.Stderr, reason)
		os.Exit(1)
	}
	OutputBlock(reason)
}

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
	if dir = os.Getenv("CURSOR_PROJECT_DIR"); dir != "" {
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

// manifestContainsIn checks if filePath is listed in the manifest at manifestPath.
func manifestContainsIn(filePath, manifestPath string) bool {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return false
	}
	baseDir := filepath.Dir(manifestPath)
	rel := strings.TrimPrefix(filePath, baseDir+string(filepath.Separator))
	rel = strings.ReplaceAll(rel, "\\", "/")
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

func ManifestContains(filePath string) bool {
	// Check .claude/.rein-manifest
	if manifestContainsIn(filePath, filepath.Join(ConfigDir(), ".rein-manifest")) {
		return true
	}
	pd := ProjectDir()
	// Check .cursor/.rein-manifest (project-level Cursor install)
	if manifestContainsIn(filePath, filepath.Join(pd, ".cursor", ".rein-manifest")) {
		return true
	}
	// Check .codex/.rein-manifest (project-level Codex install)
	if manifestContainsIn(filePath, filepath.Join(pd, ".codex", ".rein-manifest")) {
		return true
	}
	return false
}
