// Tool: render-graphs
// Extracts Graphviz dot blocks from SKILL.md files and renders them to SVG.
//
// Usage:
//   go run tools/render-graphs.go <skill-directory>           # Render each diagram separately
//   go run tools/render-graphs.go <skill-directory> --combine # Combine all into one diagram
//
// Requires: graphviz (dot) installed on system

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var dotBlockRe = regexp.MustCompile("(?s)```dot\n(.*?)```")
var digraphNameRe = regexp.MustCompile(`digraph\s+(\w+)`)
var digraphBodyRe = regexp.MustCompile("(?s)digraph\\s+\\w+\\s*\\{(.*?)\\}")
var rankdirRe = regexp.MustCompile(`(?m)^\s*rankdir\s*=\s*\w+\s*;?\s*$`)

type dotBlock struct {
	Name    string
	Content string
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: render-graphs <skill-directory> [--combine]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Options:")
		fmt.Fprintln(os.Stderr, "  --combine    Combine all diagrams into one SVG")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Example:")
		fmt.Fprintln(os.Stderr, "  go run tools/render-graphs.go skills/subagent")
		os.Exit(1)
	}

	combine := false
	var skillDirArg string
	for _, a := range args {
		if a == "--combine" {
			combine = true
		} else {
			skillDirArg = a
		}
	}

	skillDir, err := filepath.Abs(skillDirArg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}

	skillFile := filepath.Join(skillDir, "SKILL.md")
	skillName := strings.ReplaceAll(filepath.Base(skillDir), "-", "_")

	data, err := os.ReadFile(skillFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s not found\n", skillFile)
		os.Exit(1)
	}

	// Check if dot is available
	if _, err := exec.LookPath("dot"); err != nil {
		fmt.Fprintln(os.Stderr, "Error: graphviz (dot) not found. Install with:")
		fmt.Fprintln(os.Stderr, "  brew install graphviz    # macOS")
		fmt.Fprintln(os.Stderr, "  apt install graphviz     # Linux")
		fmt.Fprintln(os.Stderr, "  choco install graphviz   # Windows")
		os.Exit(1)
	}

	blocks := extractDotBlocks(string(data))
	if len(blocks) == 0 {
		fmt.Printf("No ```dot blocks found in %s\n", skillFile)
		os.Exit(0)
	}

	fmt.Printf("Found %d diagram(s) in %s/SKILL.md\n", len(blocks), filepath.Base(skillDir))

	outputDir := filepath.Join(skillDir, "diagrams")
	os.MkdirAll(outputDir, 0755)

	if combine {
		combined := combineGraphs(blocks, skillName)
		svg, err := renderToSVG(combined)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Failed to render combined diagram: %v\n", err)
		} else {
			outPath := filepath.Join(outputDir, skillName+"_combined.svg")
			os.WriteFile(outPath, []byte(svg), 0644)
			fmt.Printf("  Rendered: %s_combined.svg\n", skillName)

			dotPath := filepath.Join(outputDir, skillName+"_combined.dot")
			os.WriteFile(dotPath, []byte(combined), 0644)
			fmt.Printf("  Source: %s_combined.dot\n", skillName)
		}
	} else {
		for _, block := range blocks {
			svg, err := renderToSVG(block.Content)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  Failed: %s: %v\n", block.Name, err)
			} else {
				outPath := filepath.Join(outputDir, block.Name+".svg")
				os.WriteFile(outPath, []byte(svg), 0644)
				fmt.Printf("  Rendered: %s.svg\n", block.Name)
			}
		}
	}

	fmt.Printf("\nOutput: %s/\n", outputDir)
}

func extractDotBlocks(markdown string) []dotBlock {
	var blocks []dotBlock
	matches := dotBlockRe.FindAllStringSubmatch(markdown, -1)

	for i, m := range matches {
		content := strings.TrimSpace(m[1])
		nameMatch := digraphNameRe.FindStringSubmatch(content)
		name := fmt.Sprintf("graph_%d", i+1)
		if nameMatch != nil {
			name = nameMatch[1]
		}
		blocks = append(blocks, dotBlock{Name: name, Content: content})
	}

	return blocks
}

func extractGraphBody(dotContent string) string {
	m := digraphBodyRe.FindStringSubmatch(dotContent)
	if m == nil {
		return ""
	}
	body := m[1]
	body = rankdirRe.ReplaceAllString(body, "")
	return strings.TrimSpace(body)
}

func combineGraphs(blocks []dotBlock, skillName string) string {
	var bodies []string
	for i, block := range blocks {
		body := extractGraphBody(block.Content)
		indented := ""
		for _, line := range strings.Split(body, "\n") {
			indented += "    " + line + "\n"
		}
		bodies = append(bodies, fmt.Sprintf("  subgraph cluster_%d {\n    label=\"%s\";\n%s  }", i, block.Name, indented))
	}

	return fmt.Sprintf("digraph %s_combined {\n  rankdir=TB;\n  compound=true;\n  newrank=true;\n\n%s\n}", skillName, strings.Join(bodies, "\n\n"))
}

func renderToSVG(dotContent string) (string, error) {
	cmd := exec.Command("dot", "-Tsvg")
	cmd.Stdin = strings.NewReader(dotContent)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("dot failed: %v", err)
	}
	return string(output), nil
}
