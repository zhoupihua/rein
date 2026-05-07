package artifact

import (
	"fmt"
	"regexp"
	"strings"
)

// Delta operation types.
type DeltaOperation string

const (
	DeltaAdded    DeltaOperation = "ADDED"
	DeltaModified DeltaOperation = "MODIFIED"
	DeltaRemoved  DeltaOperation = "REMOVED"
	DeltaRenamed  DeltaOperation = "RENAMED"
)

// RequirementBlock represents a parsed requirement within a markdown spec.
type RequirementBlock struct {
	HeaderLine string // e.g. "### Requirement: Something"
	Name       string // e.g. "Something"
	Raw        string // full block including header and body
}

// RenamePair represents a FROM/TO rename operation.
type RenamePair struct {
	From string
	To   string
}

// DeltaPlan represents a parsed delta spec with operations to apply.
type DeltaPlan struct {
	Added    []RequirementBlock
	Modified []RequirementBlock
	Removed  []string    // requirement names
	Renamed  []RenamePair
}

// SectionPresence tracks which delta sections were found in the document.
type SectionPresence struct {
	Added    bool
	Modified bool
	Removed  bool
	Renamed  bool
}

var (
	reqSectionRe   = regexp.MustCompile(`(?i)^##\s+Requirements\s*$`)
	reqHeaderRe    = regexp.MustCompile(`(?i)^###\s*Requirement:\s*(.+?)\s*$`)
	deltaSectionRe = regexp.MustCompile(`(?i)^##\s+(ADDED|MODIFIED|REMOVED|RENAMED)\s+Requirements\s*$`)
	removedItemRe  = regexp.MustCompile(`(?i)^\s*-\s*` + "`" + `?###\s*Requirement:\s*(.+?)` + "`" + `?\s*$`)
	renameFromRe   = regexp.MustCompile(`(?i)^\s*-?\s*FROM:\s*` + "`" + `?###\s*Requirement:\s*(.+?)` + "`" + `?\s*$`)
	renameToRe     = regexp.MustCompile(`(?i)^\s*-?\s*TO:\s*` + "`" + `?###\s*Requirement:\s*(.+?)` + "`" + `?\s*$`)
)

// ParseDeltaSpec parses a delta spec document into a DeltaPlan.
func ParseDeltaSpec(content string) (*DeltaPlan, *SectionPresence, error) {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	plan := &DeltaPlan{}
	presence := &SectionPresence{}

	// Split into top-level sections
	sections := splitTopLevelSections(content)

	for title, body := range sections {
		upper := strings.ToUpper(strings.TrimSpace(title))
		switch {
		case strings.HasPrefix(upper, "ADDED"):
			presence.Added = true
			blocks := parseRequirementBlocks(body)
			plan.Added = blocks
		case strings.HasPrefix(upper, "MODIFIED"):
			presence.Modified = true
			blocks := parseRequirementBlocks(body)
			plan.Modified = blocks
		case strings.HasPrefix(upper, "REMOVED"):
			presence.Removed = true
			names := parseRemovedNames(body)
			plan.Removed = names
		case strings.HasPrefix(upper, "RENAMED"):
			presence.Renamed = true
			pairs := parseRenamedPairs(body)
			plan.Renamed = pairs
		}
	}

	if len(plan.Added)+len(plan.Modified)+len(plan.Removed)+len(plan.Renamed) == 0 {
		return nil, presence, fmt.Errorf("delta spec contains no operations")
	}

	return plan, presence, nil
}

// ExtractRequirementsSection parses the ## Requirements section from a spec.
func ExtractRequirementsSection(content string) (before, headerLine, preamble string, blocks []RequirementBlock, after string) {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	lines := strings.Split(content, "\n")

	// Find the ## Requirements heading
	reqIdx := -1
	for i, line := range lines {
		if reqSectionRe.MatchString(line) {
			reqIdx = i
			break
		}
	}

	if reqIdx == -1 {
		before = content
		return
	}

	before = strings.Join(lines[:reqIdx], "\n")
	headerLine = lines[reqIdx]

	// Find section end (next ## at same or higher level)
	endIdx := len(lines)
	for i := reqIdx + 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "## ") {
			endIdx = i
			break
		}
	}

	sectionLines := lines[reqIdx+1 : endIdx]
	after = strings.Join(lines[endIdx:], "\n")

	// Separate preamble from requirement blocks
	preambleEnd := 0
	for i, line := range sectionLines {
		if reqHeaderRe.MatchString(line) {
			break
		}
		preambleEnd = i + 1
	}

	preamble = strings.Join(sectionLines[:preambleEnd], "\n")
	blockLines := sectionLines[preambleEnd:]

	// Parse requirement blocks
	var currentBlock *RequirementBlock
	for _, line := range blockLines {
		if m := reqHeaderRe.FindStringSubmatch(line); m != nil {
			if currentBlock != nil {
				currentBlock.Raw = strings.TrimSpace(currentBlock.Raw)
				blocks = append(blocks, *currentBlock)
			}
			currentBlock = &RequirementBlock{
				HeaderLine: line,
				Name:       strings.TrimSpace(m[1]),
				Raw:        line,
			}
		} else if currentBlock != nil {
			currentBlock.Raw += "\n" + line
		}
	}
	if currentBlock != nil {
		currentBlock.Raw = strings.TrimSpace(currentBlock.Raw)
		blocks = append(blocks, *currentBlock)
	}

	return
}

// splitTopLevelSections splits content into top-level sections keyed by title.
func splitTopLevelSections(content string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(content, "\n")

	var currentTitle string
	var currentBody []string

	for _, line := range lines {
		if m := deltaSectionRe.FindStringSubmatch(line); m != nil {
			if currentTitle != "" {
				result[currentTitle] = strings.TrimSpace(strings.Join(currentBody, "\n"))
			}
			currentTitle = m[1]
			currentBody = nil
		} else if currentTitle != "" {
			currentBody = append(currentBody, line)
		}
	}

	if currentTitle != "" {
		result[currentTitle] = strings.TrimSpace(strings.Join(currentBody, "\n"))
	}

	return result
}

// parseRequirementBlocks parses requirement blocks from section body text.
func parseRequirementBlocks(body string) []RequirementBlock {
	var blocks []RequirementBlock
	lines := strings.Split(body, "\n")
	var current *RequirementBlock

	for _, line := range lines {
		if m := reqHeaderRe.FindStringSubmatch(line); m != nil {
			if current != nil {
				current.Raw = strings.TrimSpace(current.Raw)
				blocks = append(blocks, *current)
			}
			current = &RequirementBlock{
				HeaderLine: line,
				Name:       strings.TrimSpace(m[1]),
				Raw:        line,
			}
		} else if current != nil {
			current.Raw += "\n" + line
		}
	}

	if current != nil {
		current.Raw = strings.TrimSpace(current.Raw)
		blocks = append(blocks, *current)
	}

	return blocks
}

// parseRemovedNames extracts requirement names from a REMOVED section.
func parseRemovedNames(body string) []string {
	var names []string
	seen := make(map[string]bool)

	lines := strings.Split(body, "\n")
	for _, line := range lines {
		// Try bullet list format first
		if m := removedItemRe.FindStringSubmatch(line); m != nil {
			name := strings.TrimSpace(m[1])
			if !seen[name] {
				names = append(names, name)
				seen[name] = true
			}
			continue
		}
		// Try requirement header format
		if m := reqHeaderRe.FindStringSubmatch(line); m != nil {
			name := strings.TrimSpace(m[1])
			if !seen[name] {
				names = append(names, name)
				seen[name] = true
			}
		}
	}

	return names
}

// parseRenamedPairs extracts FROM/TO rename pairs from a RENAMED section.
func parseRenamedPairs(body string) []RenamePair {
	var pairs []RenamePair
	var lastFrom string

	lines := strings.Split(body, "\n")
	for _, line := range lines {
		if m := renameFromRe.FindStringSubmatch(line); m != nil {
			lastFrom = strings.TrimSpace(m[1])
		} else if m := renameToRe.FindStringSubmatch(line); m != nil && lastFrom != "" {
			pairs = append(pairs, RenamePair{
				From: lastFrom,
				To:   strings.TrimSpace(m[1]),
			})
			lastFrom = ""
		}
	}

	return pairs
}
