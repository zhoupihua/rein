package hook

import (
	"os/exec"
	"strings"
)

func Format() {
	target := FilePath()
	if target == "" {
		return
	}

	// Only format web files
	ext := strings.ToLower(target[strings.LastIndex(target, ".")+1:])
	webExts := map[string]bool{"js": true, "jsx": true, "ts": true, "tsx": true, "css": true, "scss": true, "html": true, "json": true, "md": true}
	if !webExts[ext] {
		return
	}

	// Try prettier
	cmd := exec.Command("npx", "prettier", "--write", target)
	cmd.Run() // Ignore errors (prettier might not be installed)
}
