package hook

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/zhoupihua/rein/internal/artifact"
)

func TaskProgress() {
	input := ReadToolInput()
	if input == "" {
		return
	}

	target := ExtractFilePath(input)
	if target == "" {
		return
	}

	// Skip task.md edits
	if strings.Contains(target, "docs/rein/changes/") && strings.HasSuffix(target, "task.md") {
		return
	}

	editedFile := filepath.Base(target)
	changesDir := filepath.Join(ProjectDir(), "docs", "rein", "changes")
	if _, err := os.Stat(changesDir); os.IsNotExist(err) {
		return
	}

	// Scan each feature directory for task.md that references the edited file
	entries, err := os.ReadDir(changesDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		taskPath := filepath.Join(changesDir, entry.Name(), "task.md")
		matchedID := matchTaskFileSingle(taskPath, editedFile)
		if matchedID != "" {
			tf, err := artifact.ParseTaskFile(taskPath)
			if err != nil {
				continue
			}
			id, ok := artifact.ParseTaskID(matchedID)
			if !ok {
				continue
			}
			task := tf.FindTask(id)
			if task != nil && !task.Done {
				tf.CheckTask(id)
				OutputAdditional("PostToolUse", "Auto-checked task "+matchedID+" (file match: "+editedFile+")")
			}
			return
		}

		// Phase 2: scan plan.md **Files:** fields
		planPath := filepath.Join(changesDir, entry.Name(), "plan.md")
		matchedID = matchPlanFileSingle(planPath, taskPath, editedFile)
		if matchedID != "" {
			tf, err := artifact.ParseTaskFile(taskPath)
			if err != nil {
				continue
			}
			id, ok := artifact.ParseTaskID(matchedID)
			if !ok {
				continue
			}
			task := tf.FindTask(id)
			if task != nil && !task.Done {
				tf.CheckTask(id)
				OutputAdditional("PostToolUse", "Auto-checked task "+matchedID+" (file match: "+editedFile+")")
			}
			return
		}
	}
}

func matchTaskFileSingle(taskPath, editedFile string) string {
	if _, err := os.Stat(taskPath); err != nil {
		return ""
	}

	codeExts := map[string]bool{
		"go": true, "ts": true, "tsx": true, "js": true, "jsx": true,
		"py": true, "rs": true, "java": true, "rb": true, "sql": true,
		"yaml": true, "yml": true, "json": true, "toml": true, "proto": true,
		"graphql": true, "css": true, "scss": true, "html": true, "sh": true,
		"ps1": true, "mod": true, "sum": true, "env": true, "conf": true,
		"xml": true, "dart": true, "swift": true, "kt": true, "c": true,
		"cpp": true, "h": true, "hpp": true, "php": true, "tf": true,
		"lock": true, "txt": true, "md": true,
	}

	tf, err := artifact.ParseTaskFile(taskPath)
	if err != nil {
		return ""
	}

	for _, task := range tf.AllTasks() {
		if task.Done {
			continue
		}
		refs := extractRefs(task.Description, codeExts)
		for _, ref := range refs {
			if filepath.Base(ref) == editedFile {
				return task.ID.String()
			}
		}
	}
	return ""
}

func matchPlanFileSingle(planPath, taskPath, editedFile string) string {
	if _, err := os.Stat(planPath); err != nil {
		return ""
	}

	codeExts := map[string]bool{
		"go": true, "ts": true, "tsx": true, "js": true, "jsx": true,
		"py": true, "rs": true, "java": true, "rb": true, "sql": true,
		"yaml": true, "yml": true, "json": true, "toml": true, "proto": true,
		"graphql": true, "css": true, "scss": true, "html": true, "sh": true,
		"ps1": true, "mod": true, "sum": true, "env": true, "conf": true,
		"xml": true, "dart": true, "swift": true, "kt": true, "c": true,
		"cpp": true, "h": true, "hpp": true, "php": true, "tf": true,
		"lock": true, "txt": true, "md": true,
	}

	plan, err := artifact.ParsePlanFile(planPath)
	if err != nil {
		return ""
	}

	for _, detail := range plan.TaskDetails {
		if detail.Files == "" {
			continue
		}
		refs := extractRefs(detail.Files, codeExts)
		for _, ref := range refs {
			if filepath.Base(ref) == editedFile {
				// Find matching unchecked task
				tf, err := artifact.ParseTaskFile(taskPath)
				if err != nil {
					continue
				}
				task := tf.FindTask(detail.ID)
				if task != nil && !task.Done {
					return detail.ID.String()
				}
			}
		}
	}
	return ""
}

func extractRefs(line string, codeExts map[string]bool) []string {
	var refs []string

	// Backtick-enclosed references
	inBacktick := false
	var current strings.Builder
	for _, ch := range line {
		if ch == '`' {
			if inBacktick {
				if current.Len() > 0 {
					refs = append(refs, current.String())
				}
				current.Reset()
				inBacktick = false
			} else {
				inBacktick = true
			}
		} else if inBacktick {
			current.WriteRune(ch)
		}
	}

	// Plain filenames with code extensions
	words := strings.Fields(line)
	for _, word := range words {
		word = strings.TrimRight(word, ".,;:")
		dot := strings.LastIndex(word, ".")
		if dot > 0 && dot < len(word)-1 {
			ext := word[dot+1:]
			if codeExts[ext] {
				refs = append(refs, word)
			}
		}
	}

	return refs
}
