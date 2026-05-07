package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSkillFilesExist verifies all SKILL.md files are present and non-empty.
func TestSkillFilesExist(t *testing.T) {
	skillsDir := findSkillsDir(t)
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		t.Fatalf("cannot read skills dir: %v", err)
	}

	found := 0
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		skillFile := filepath.Join(skillsDir, e.Name(), "SKILL.md")
		data, err := os.ReadFile(skillFile)
		if err != nil {
			t.Errorf("skill %q: missing SKILL.md", e.Name())
			continue
		}
		if len(data) == 0 {
			t.Errorf("skill %q: SKILL.md is empty", e.Name())
		}
		found++
	}

	if found == 0 {
		t.Error("no skills found")
	}
	t.Logf("found %d skills with SKILL.md", found)
}

// TestSkillFrontmatter verifies SKILL.md files have YAML frontmatter.
func TestSkillFrontmatter(t *testing.T) {
	skillsDir := findSkillsDir(t)
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		t.Fatalf("cannot read skills dir: %v", err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(skillsDir, e.Name(), "SKILL.md"))
		if err != nil {
			continue
		}
		content := string(data)
		if !strings.HasPrefix(content, "---") {
			t.Errorf("skill %q: SKILL.md missing frontmatter delimiter", e.Name())
		}
	}
}

// TestSkillSupportingFiles verifies referenced supporting files exist.
func TestSkillSupportingFiles(t *testing.T) {
	skillsDir := findSkillsDir(t)

	// Check known supporting files that skills reference
	supportingFiles := []string{
		"debugging/root-cause-tracing.md",
		"debugging/defense-in-depth.md",
		"debugging/condition-based-waiting.md",
		"debugging/testing-anti-patterns.md",
		"define/frameworks.md",
		"define/refinement-criteria.md",
		"define/examples.md",
		"define/visual-thinking.md",
		"define/spec-reviewer-prompt.md",
		"define/frame-template.html",
		"define/helper.js",
		"writing-skills/anthropic-best-practices.md",
		"writing-skills/persuasion-principles.md",
		"writing-skills/testing-skills-with-subagents.md",
		"writing-skills/graphviz-conventions.dot",
		"subagent/implementer-prompt.md",
		"subagent/spec-reviewer-prompt.md",
		"subagent/code-quality-reviewer-prompt.md",
		"planning/plan-reviewer-prompt.md",
	}

	for _, rel := range supportingFiles {
		path := filepath.Join(skillsDir, rel)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing supporting file: %s", rel)
		}
	}
}

func findSkillsDir(t *testing.T) string {
	t.Helper()
	candidates := []string{
		"../../skills",
		"../../../skills",
	}

	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			abs, _ := filepath.Abs(c)
			return abs
		}
	}
	t.Fatal("cannot find skills directory")
	return ""
}
