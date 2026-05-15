package convert

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

// SkillGlobs maps skill names to file glob patterns for auto-attached Cursor rules.
var SkillGlobs = map[string][]string{
	"frontend":    {"src/**/*.tsx", "src/**/*.jsx", "src/**/*.css", "src/**/*.scss"},
	"tdd":         {"**/*_test.go", "**/*_test.ts", "**/*.test.*", "**/*.spec.*"},
	"security":    {"**/auth/**", "**/middleware/**"},
	"performance": {"**/api/**", "**/query*"},
	"migration":   {"**/migrations/**", "**/db/**"},
}

var frontmatterRe = regexp.MustCompile(`(?s)^---\s*\n(.*?)\n---\s*\n?(.*)`)

// ParseFrontmatter extracts YAML frontmatter from markdown content.
// Returns a map of key-value pairs and the body after the frontmatter.
func ParseFrontmatter(content string) (map[string]string, string) {
	matches := frontmatterRe.FindStringSubmatch(content)
	if matches == nil {
		return nil, content
	}
	fm := make(map[string]string)
	for _, line := range strings.Split(matches[1], "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		fm[key] = val
	}
	return fm, matches[2]
}

// SkillToMDC converts a SKILL.md to Cursor .mdc format.
func SkillToMDC(name, description, body string, globs []string) string {
	return buildMDC(description, false, globs, body)
}

// CommandToMDC converts a command .md to Cursor .mdc format.
func CommandToMDC(description, body string) string {
	return buildMDC(description, false, nil, body)
}

// AgentToMDC converts an agent .md to Cursor .mdc format.
func AgentToMDC(name, description, body string) string {
	return buildMDC(description, false, nil, body)
}

// AlwaysApplyRule creates an always-apply .mdc rule.
func AlwaysApplyRule(description, body string) string {
	return buildMDC(description, true, nil, body)
}

func buildMDC(description string, alwaysApply bool, globs []string, body string) string {
	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("description: %s\n", description))
	if len(globs) > 0 {
		sb.WriteString("globs:\n")
		for _, g := range globs {
			sb.WriteString(fmt.Sprintf("  - %q\n", g))
		}
	}
	sb.WriteString(fmt.Sprintf("alwaysApply: %v\n", alwaysApply))
	sb.WriteString("---\n")
	sb.WriteString(body)
	return sb.String()
}

// SkillNameFromFile derives the skill name from its directory path.
// e.g., "skills/code-review/SKILL.md" -> "code-review"
func SkillNameFromFile(path string) string {
	dir := filepath.Dir(path)
	return filepath.Base(dir)
}

// ConvertSkillFile reads a SKILL.md and converts it to .mdc format.
func ConvertSkillFile(content string, name string) string {
	fm, body := ParseFrontmatter(content)
	description := name
	if d, ok := fm["description"]; ok && d != "" {
		description = d
	}
	globs := SkillGlobs[name]
	return SkillToMDC(name, description, body, globs)
}

// ConvertCommandFile reads a command .md and converts it to .mdc format.
func ConvertCommandFile(content string, name string) string {
	fm, body := ParseFrontmatter(content)
	description := name
	if d, ok := fm["description"]; ok && d != "" {
		description = d
	}
	return CommandToMDC(description, body)
}

// ConvertAgentFile reads an agent .md and converts it to .mdc format.
func ConvertAgentFile(content string, name string) string {
	fm, body := ParseFrontmatter(content)
	description := name
	if d, ok := fm["description"]; ok && d != "" {
		description = d
	}
	agentName := name
	if n, ok := fm["name"]; ok && n != "" {
		agentName = n
	}
	return AgentToMDC(agentName, description, body)
}

// --- Codex Conversion ---

// CodexProjectHeader is the header content for the rein project section in CODEX.md.
const CodexProjectHeader = `## Task Progress

When working on a feature with ` + "`docs/rein/changes/<name>/task.md`" + `, after completing
any task or sub-task, you MUST immediately mark it as done:

  rein task done <id>          # e.g., rein task done 1.2
  rein task done <subtask-id>  # e.g., rein task done 1.2.0

Do NOT skip this step. Marking progress is mandatory, not optional.

## rein Workflow

This project uses rein for structured development. Key commands:
- ` + "`rein status`" + ` — Check current workflow phase
- ` + "`rein task next`" + ` — Show next unchecked task
- ` + "`rein task done <id>`" + ` — Mark task complete
- ` + "`rein validate <feature>`" + ` — Validate artifact completeness`

// ConvertToCODEXMd merges all artifacts into a single CODEX.md for Codex CLI.
// skills, commands, agents are maps of name→raw file content (including frontmatter).
func ConvertToCODEXMd(skills, commands, agents map[string]string) string {
	var sb strings.Builder
	sb.WriteString("# rein Project Instructions\n\n")
	sb.WriteString(CodexProjectHeader)
	sb.WriteString("\n\n")

	if len(skills) > 0 {
		sb.WriteString("## Skills\n\n")
		for _, name := range sortedKeys(skills) {
			content := skills[name]
			_, body := ParseFrontmatter(content)
			sb.WriteString(fmt.Sprintf("### %s\n\n", name))
			sb.WriteString(strings.TrimSpace(body))
			sb.WriteString("\n\n")
		}
	}

	if len(commands) > 0 {
		sb.WriteString("## Commands\n\n")
		sb.WriteString("Reference these commands by name when asking Codex to perform a workflow.\n\n")
		for _, name := range sortedKeys(commands) {
			content := commands[name]
			_, body := ParseFrontmatter(content)
			sb.WriteString(fmt.Sprintf("### %s\n\n", name))
			sb.WriteString(strings.TrimSpace(body))
			sb.WriteString("\n\n")
		}
	}

	if len(agents) > 0 {
		sb.WriteString("## Agents\n\n")
		sb.WriteString("Reference these agents by name when asking Codex to adopt a perspective.\n\n")
		for _, name := range sortedKeys(agents) {
			content := agents[name]
			_, body := ParseFrontmatter(content)
			sb.WriteString(fmt.Sprintf("### %s\n\n", name))
			sb.WriteString(strings.TrimSpace(body))
			sb.WriteString("\n\n")
		}
	}

	return sb.String()
}

// CodexConfigTOML generates a .codex/config.toml with rein hooks.
func CodexConfigTOML(reinBinPath string) string {
	guardCommand := tomlBasicString(reinBinPath + " hook guard")
	gateCommand := tomlBasicString(reinBinPath + " hook gate")
	formatCommand := tomlBasicString(reinBinPath + " hook format")

	return fmt.Sprintf(`# rein Codex configuration
[features]
multi_agent = true

[[hooks.pre_command]]
command = %s
description = "Block edits to rein-managed files"

[[hooks.pre_command]]
command = %s
description = "Run tests before deploy commands"

[[hooks.post_command]]
command = %s
description = "Auto-format web files with prettier"
`, guardCommand, gateCommand, formatCommand)
}

func tomlBasicString(s string) string {
	var sb strings.Builder
	sb.Grow(len(s) + 2)
	sb.WriteByte('"')
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError && size == 1 {
			sb.WriteString(`\uFFFD`)
			s = s[size:]
			continue
		}

		switch r {
		case '\b':
			sb.WriteString(`\b`)
		case '\t':
			sb.WriteString(`\t`)
		case '\n':
			sb.WriteString(`\n`)
		case '\f':
			sb.WriteString(`\f`)
		case '\r':
			sb.WriteString(`\r`)
		case '"':
			sb.WriteString(`\"`)
		case '\\':
			sb.WriteString(`\\`)
		default:
			if r < 0x20 {
				sb.WriteString(fmt.Sprintf(`\u%04X`, r))
			} else {
				sb.WriteRune(r)
			}
		}
		s = s[size:]
	}
	sb.WriteByte('"')
	return sb.String()
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
