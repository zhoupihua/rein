package convert

import (
	"strings"
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	content := `---
name: tdd
description: Test-driven development skill
---

# TDD Skill

Some body content.
`
	fm, body := ParseFrontmatter(content)
	if fm["name"] != "tdd" {
		t.Errorf("expected name=tdd, got %q", fm["name"])
	}
	if fm["description"] != "Test-driven development skill" {
		t.Errorf("unexpected description: %q", fm["description"])
	}
	if !strings.Contains(body, "# TDD Skill") {
		t.Errorf("body should contain heading, got: %q", body)
	}
}

func TestParseFrontmatter_NoFrontmatter(t *testing.T) {
	content := "# Just a heading\n\nSome text."
	fm, body := ParseFrontmatter(content)
	if fm != nil {
		t.Errorf("expected nil frontmatter, got %v", fm)
	}
	if body != content {
		t.Errorf("body should equal input without frontmatter")
	}
}

func TestSkillToMDC(t *testing.T) {
	result := SkillToMDC("tdd", "TDD skill description", "# TDD\n\nBody here.", nil)
	if !strings.Contains(result, "description: TDD skill description") {
		t.Errorf("missing description in output: %q", result)
	}
	if !strings.Contains(result, "alwaysApply: false") {
		t.Errorf("missing alwaysApply in output: %q", result)
	}
	if !strings.Contains(result, "# TDD") {
		t.Errorf("missing body content in output: %q", result)
	}
}

func TestSkillToMDC_WithGlobs(t *testing.T) {
	globs := []string{"**/*_test.go", "**/*.test.*"}
	result := SkillToMDC("tdd", "TDD description", "body", globs)
	if !strings.Contains(result, "globs:") {
		t.Errorf("missing globs section: %q", result)
	}
	if !strings.Contains(result, `"**/*_test.go"`) {
		t.Errorf("missing glob entry: %q", result)
	}
}

func TestCommandToMDC(t *testing.T) {
	result := CommandToMDC("Run the feature workflow", "# Feature\n\nSteps here.")
	if !strings.Contains(result, "description: Run the feature workflow") {
		t.Errorf("missing description: %q", result)
	}
	if !strings.Contains(result, "alwaysApply: false") {
		t.Errorf("should not always apply: %q", result)
	}
}

func TestAgentToMDC(t *testing.T) {
	result := AgentToMDC("code-reviewer", "Senior code reviewer", "# Code Reviewer\n\nFramework.")
	if !strings.Contains(result, "description: Senior code reviewer") {
		t.Errorf("missing description: %q", result)
	}
	if !strings.Contains(result, "alwaysApply: false") {
		t.Errorf("should not always apply: %q", result)
	}
}

func TestAlwaysApplyRule(t *testing.T) {
	result := AlwaysApplyRule("Project conventions", "Always mark tasks done.")
	if !strings.Contains(result, "alwaysApply: true") {
		t.Errorf("should always apply: %q", result)
	}
}

func TestConvertSkillFile(t *testing.T) {
	content := `---
name: code-review
description: Five-axis code review
---

# Code Review

Review content.
`
	result := ConvertSkillFile(content, "code-review")
	if !strings.Contains(result, "description: Five-axis code review") {
		t.Errorf("missing description from frontmatter: %q", result)
	}
	if !strings.Contains(result, "# Code Review") {
		t.Errorf("missing body: %q", result)
	}
}

func TestConvertCommandFile(t *testing.T) {
	content := `---
description: Run the feature workflow
---

# Feature Workflow

Step 1...
`
	result := ConvertCommandFile(content, "feature")
	if !strings.Contains(result, "description: Run the feature workflow") {
		t.Errorf("missing description: %q", result)
	}
}

func TestConvertAgentFile(t *testing.T) {
	content := `---
name: code-reviewer
description: Senior code reviewer persona
---

# Code Reviewer

Review framework.
`
	result := ConvertAgentFile(content, "code-reviewer")
	if !strings.Contains(result, "description: Senior code reviewer persona") {
		t.Errorf("missing description: %q", result)
	}
}

func TestSkillNameFromFile(t *testing.T) {
	tests := []struct {
		path, expected string
	}{
		{"skills/code-review/SKILL.md", "code-review"},
		{"skills/tdd/SKILL.md", "tdd"},
		{"skills/frontend/SKILL.md", "frontend"},
	}
	for _, tt := range tests {
		got := SkillNameFromFile(tt.path)
		if got != tt.expected {
			t.Errorf("SkillNameFromFile(%q) = %q, want %q", tt.path, got, tt.expected)
		}
	}
}

func TestConvertToCODEXMd(t *testing.T) {
	skills := map[string]string{
		"tdd": "---\nname: tdd\ndescription: TDD skill\n---\n\n# TDD\n\nWrite tests first.",
	}
	commands := map[string]string{
		"feature": "---\ndescription: Feature workflow\n---\n\n# Feature\n\nSix steps.",
	}
	agents := map[string]string{
		"code-reviewer": "---\nname: code-reviewer\ndescription: Reviewer\n---\n\n# Code Reviewer\n\nFive-axis review.",
	}

	result := ConvertToCODEXMd(skills, commands, agents)

	if !strings.Contains(result, "# rein Project Instructions") {
		t.Errorf("missing header: %q", result)
	}
	if !strings.Contains(result, "## Task Progress") {
		t.Errorf("missing task progress section: %q", result)
	}
	if !strings.Contains(result, "## Skills") {
		t.Errorf("missing skills section: %q", result)
	}
	if !strings.Contains(result, "### tdd") {
		t.Errorf("missing tdd skill: %q", result)
	}
	if !strings.Contains(result, "Write tests first") {
		t.Errorf("missing skill body: %q", result)
	}
	if !strings.Contains(result, "## Commands") {
		t.Errorf("missing commands section: %q", result)
	}
	if !strings.Contains(result, "### feature") {
		t.Errorf("missing feature command: %q", result)
	}
	if !strings.Contains(result, "## Agents") {
		t.Errorf("missing agents section: %q", result)
	}
	if !strings.Contains(result, "### code-reviewer") {
		t.Errorf("missing code-reviewer agent: %q", result)
	}
	// Frontmatter should be stripped from body
	if strings.Contains(result, "description:") {
		t.Errorf("frontmatter leaked into output: %q", result)
	}
}

func TestConvertToCODEXMd_Empty(t *testing.T) {
	result := ConvertToCODEXMd(nil, nil, nil)
	if !strings.Contains(result, "# rein Project Instructions") {
		t.Errorf("missing header for empty input")
	}
	if strings.Contains(result, "## Skills") {
		t.Errorf("should not have skills section when empty")
	}
}

func TestCodexConfigTOML(t *testing.T) {
	result := CodexConfigTOML("/usr/local/bin/rein")
	if !strings.Contains(result, "multi_agent = true") {
		t.Errorf("missing multi_agent feature: %q", result)
	}
	if !strings.Contains(result, "hook guard") {
		t.Errorf("missing guard hook: %q", result)
	}
	if !strings.Contains(result, "hook gate") {
		t.Errorf("missing gate hook: %q", result)
	}
	if !strings.Contains(result, "hook format") {
		t.Errorf("missing format hook: %q", result)
	}
	if !strings.Contains(result, "/usr/local/bin/rein") {
		t.Errorf("missing rein path: %q", result)
	}
}

func TestCodexConfigTOML_EscapesWindowsPath(t *testing.T) {
	result := CodexConfigTOML(`C:\Users\admin\.claude\bin\rein.exe`)

	if strings.Contains(result, `command = "C:\Users`) {
		t.Errorf("Windows path was not escaped for TOML: %q", result)
	}
	if !strings.Contains(result, `command = "C:\\Users\\admin\\.claude\\bin\\rein.exe hook guard"`) {
		t.Errorf("missing escaped guard command: %q", result)
	}
	if !strings.Contains(result, `command = "C:\\Users\\admin\\.claude\\bin\\rein.exe hook gate"`) {
		t.Errorf("missing escaped gate command: %q", result)
	}
	if !strings.Contains(result, `command = "C:\\Users\\admin\\.claude\\bin\\rein.exe hook format"`) {
		t.Errorf("missing escaped format command: %q", result)
	}
}

func TestCodexConfigTOML_EscapesQuotes(t *testing.T) {
	result := CodexConfigTOML(`C:\Program Files\rein "stable"\rein.exe`)

	if !strings.Contains(result, `C:\\Program Files\\rein \"stable\"\\rein.exe hook guard`) {
		t.Errorf("missing escaped quote in command: %q", result)
	}
}
