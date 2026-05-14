package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zhoupihua/rein/internal/convert"
)

var (
	convertIDE      string
	convertSource   string
	convertSourceDir string
	convertOutputDir string
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert artifacts to IDE-specific formats",
	Long:  "Convert SKILL.md, command, and agent files to IDE-specific formats (e.g., Cursor .mdc rules, Codex CODEX.md).",
	RunE:  runConvert,
}

func init() {
	convertCmd.Flags().StringVar(&convertIDE, "ide", "cursor", "target IDE (cursor, codex)")
	convertCmd.Flags().StringVar(&convertSource, "source", "", "single source file to convert")
	convertCmd.Flags().StringVar(&convertSourceDir, "source-dir", "", "source directory for batch conversion")
	convertCmd.Flags().StringVar(&convertOutputDir, "output-dir", "", "output directory for batch conversion (stdout if empty)")

	rootCmd.AddCommand(convertCmd)
}

func runConvert(cmd *cobra.Command, args []string) error {
	switch convertIDE {
	case "cursor":
		return runConvertCursor()
	case "codex":
		return runConvertCodex()
	default:
		return fmt.Errorf("unsupported IDE: %s (supported: cursor, codex)", convertIDE)
	}
}

func runConvertCursor() error {
	if convertSource != "" {
		return convertSingle(convertSource)
	}
	if convertSourceDir != "" {
		return convertBatch(convertSourceDir, convertOutputDir)
	}
	return fmt.Errorf("specify --source for single file or --source-dir for batch conversion")
}

func runConvertCodex() error {
	if convertSourceDir == "" {
		return fmt.Errorf("--source-dir is required for Codex conversion")
	}
	if convertOutputDir == "" {
		return fmt.Errorf("--output-dir is required for Codex conversion")
	}
	os.MkdirAll(convertOutputDir, 0o755)

	// Collect all artifact content
	skills, commands, agents := collectArtifacts(convertSourceDir)

	// Generate CODEX.md
	codexMd := convert.ConvertToCODEXMd(skills, commands, agents)
	codexMdPath := filepath.Join(convertOutputDir, "CODEX.md")
	if err := os.WriteFile(codexMdPath, []byte(codexMd), 0o644); err != nil {
		return fmt.Errorf("writing CODEX.md: %w", err)
	}

	// Generate .codex/config.toml with hooks
	codexDir := filepath.Join(convertOutputDir, ".codex")
	os.MkdirAll(codexDir, 0o755)
	reinBin := findReinBinary()
	configTOML := convert.CodexConfigTOML(reinBin)
	configPath := filepath.Join(codexDir, "config.toml")
	if err := os.WriteFile(configPath, []byte(configTOML), 0o644); err != nil {
		return fmt.Errorf("writing config.toml: %w", err)
	}

	count := len(skills) + len(commands) + len(agents)
	fmt.Fprintf(os.Stderr, "Converted %d artifacts to %s (CODEX.md + .codex/config.toml)\n", count, convertOutputDir)
	return nil
}

func collectArtifacts(sourceDir string) (skills, commands, agents map[string]string) {
	skills = make(map[string]string)
	commands = make(map[string]string)
	agents = make(map[string]string)

	// Skills
	skillsDir := filepath.Join(sourceDir, "skills")
	if entries, err := os.ReadDir(skillsDir); err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			skillFile := filepath.Join(skillsDir, e.Name(), "SKILL.md")
			if data, err := os.ReadFile(skillFile); err == nil {
				skills[e.Name()] = string(data)
			}
		}
	}

	// Commands
	commandsDir := filepath.Join(sourceDir, "commands")
	if entries, err := os.ReadDir(commandsDir); err == nil {
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			if data, err := os.ReadFile(filepath.Join(commandsDir, e.Name())); err == nil {
				name := strings.TrimSuffix(e.Name(), ".md")
				commands[name] = string(data)
			}
		}
	}

	// Agents
	agentsDir := filepath.Join(sourceDir, "agents")
	if entries, err := os.ReadDir(agentsDir); err == nil {
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			if data, err := os.ReadFile(filepath.Join(agentsDir, e.Name())); err == nil {
				name := strings.TrimSuffix(e.Name(), ".md")
				agents[name] = string(data)
			}
		}
	}

	return skills, commands, agents
}

func findReinBinary() string {
	home, _ := os.UserHomeDir()
	candidates := []string{
		filepath.Join(home, ".rein", "bin", "rein"),
		filepath.Join(home, ".rein", "bin", "rein.exe"),
		filepath.Join(home, ".claude", "bin", "rein"),
		filepath.Join(home, ".claude", "bin", "rein.exe"),
		"rein", // rely on PATH
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return "rein"
}

func convertSingle(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}
	content := string(data)
	base := filepath.Base(path)

	var result string
	switch {
	case base == "SKILL.md":
		name := convert.SkillNameFromFile(path)
		result = convert.ConvertSkillFile(content, name)
	case filepath.Ext(base) == ".md":
		name := strings.TrimSuffix(base, ".md")
		// Heuristic: if file is in agents/ dir, treat as agent; else command
		if strings.Contains(filepath.Dir(path), "agents") {
			result = convert.ConvertAgentFile(content, name)
		} else {
			result = convert.ConvertCommandFile(content, name)
		}
	default:
		return fmt.Errorf("unsupported file type: %s", base)
	}

	if convertOutputDir != "" {
		outName := deriveOutputName(path)
		outPath := filepath.Join(convertOutputDir, outName)
		os.MkdirAll(convertOutputDir, 0o755)
		if err := os.WriteFile(outPath, []byte(result), 0o644); err != nil {
			return fmt.Errorf("writing %s: %w", outPath, err)
		}
	} else {
		fmt.Print(result)
	}
	return nil
}

func convertBatch(sourceDir, outputDir string) error {
	if outputDir == "" {
		return fmt.Errorf("--output-dir is required for batch conversion")
	}
	os.MkdirAll(outputDir, 0o755)

	var count int

	// Skills: skills/<name>/SKILL.md
	skillsDir := filepath.Join(sourceDir, "skills")
	if entries, err := os.ReadDir(skillsDir); err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			skillFile := filepath.Join(skillsDir, e.Name(), "SKILL.md")
			data, err := os.ReadFile(skillFile)
			if err != nil {
				continue
			}
			result := convert.ConvertSkillFile(string(data), e.Name())
			outPath := filepath.Join(outputDir, e.Name()+".mdc")
			if err := os.WriteFile(outPath, []byte(result), 0o644); err != nil {
				return fmt.Errorf("writing %s: %w", outPath, err)
			}
			count++
		}
	}

	// Commands: commands/<name>.md
	commandsDir := filepath.Join(sourceDir, "commands")
	if entries, err := os.ReadDir(commandsDir); err == nil {
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			cmdFile := filepath.Join(commandsDir, e.Name())
			data, err := os.ReadFile(cmdFile)
			if err != nil {
				continue
			}
			name := strings.TrimSuffix(e.Name(), ".md")
			result := convert.ConvertCommandFile(string(data), name)
			outPath := filepath.Join(outputDir, name+".mdc")
			if err := os.WriteFile(outPath, []byte(result), 0o644); err != nil {
				return fmt.Errorf("writing %s: %w", outPath, err)
			}
			count++
		}
	}

	// Agents: agents/<name>.md
	agentsDir := filepath.Join(sourceDir, "agents")
	if entries, err := os.ReadDir(agentsDir); err == nil {
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			agentFile := filepath.Join(agentsDir, e.Name())
			data, err := os.ReadFile(agentFile)
			if err != nil {
				continue
			}
			name := strings.TrimSuffix(e.Name(), ".md")
			result := convert.ConvertAgentFile(string(data), name)
			outPath := filepath.Join(outputDir, name+".mdc")
			if err := os.WriteFile(outPath, []byte(result), 0o644); err != nil {
				return fmt.Errorf("writing %s: %w", outPath, err)
			}
			count++
		}
	}

	fmt.Fprintf(os.Stderr, "Converted %d files to %s\n", count, outputDir)
	return nil
}

func deriveOutputName(path string) string {
	base := filepath.Base(path)
	dir := filepath.Dir(path)
	switch {
	case base == "SKILL.md":
		return filepath.Base(dir) + ".mdc"
	case filepath.Ext(base) == ".md":
		return strings.TrimSuffix(base, ".md") + ".mdc"
	default:
		return base + ".mdc"
	}
}
