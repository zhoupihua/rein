package hook

import (
	"os"
	"path/filepath"
	"strings"
)

func Guard() {
	input := ReadToolInput()
	if input == "" {
		return
	}

	target := ExtractFilePath(input)
	if target == "" {
		return
	}

	// Resolve relative to project dir
	if !filepath.IsAbs(target) {
		target = filepath.Join(ProjectDir(), target)
	}

	// Check manifest
	manifestPath := filepath.Join(ConfigDir(), ".rein-manifest")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return // No manifest, allow
	}

	rel := strings.TrimPrefix(target, ConfigDir()+string(filepath.Separator))
	rel = strings.ReplaceAll(rel, "\\", "/")

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Directory entry (ends with /)
		if strings.HasSuffix(line, "/") {
			prefix := strings.TrimSuffix(line, "/")
			if strings.HasPrefix(rel, prefix+"/") || rel == prefix {
				OutputBlock("This file is protected by rein. Remove its entry from .rein-manifest to allow edits.")
				return
			}
		} else {
			if rel == line {
				OutputBlock("This file is protected by rein. Remove its entry from .rein-manifest to allow edits.")
				return
			}
		}
	}
}
