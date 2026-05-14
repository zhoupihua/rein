package hook

import (
	"path/filepath"
)

func Guard() {
	target := FilePath()
	if target == "" {
		return
	}

	// Resolve relative to project dir
	if !filepath.IsAbs(target) {
		target = filepath.Join(ProjectDir(), target)
	}

	if ManifestContains(target) {
		BlockExit("This file is protected by rein. Remove its entry from .rein-manifest to allow edits.")
	}
}
