package artifact

import (
	"bufio"
	"regexp"
	"strings"
)

var (
	epicSectionRe    = regexp.MustCompile(`^##\s+(.+)$`)
	sharedDecisionRe = regexp.MustCompile(`^-\s+\*\*Decision:\*\*\s*(.+?)\s*—\s*\*\*Rationale:\*\*\s*(.+)$`)
	incrementOrderRe = regexp.MustCompile(`^\d+\.\s+` + "`" + `(.+?)` + "`" + `\s*—\s*(.+)$`)
	dependencyRe     = regexp.MustCompile(`^-\s+` + "`" + `(.+?)` + "`" + `\s+depends\s+on\s+` + "`" + `(.+?)` + "`")
)

type Epic struct {
	ProblemStatement string
	Goals            string
	NonGoals         string
	SharedContext    string
	SharedDecisions  []SharedDecision
	IncrementOrder   []IncrementRef
	Dependencies     []Dependency
}

type SharedDecision struct {
	Decision  string
	Rationale string
}

type IncrementRef struct {
	Name        string
	Description string
}

type Dependency struct {
	From string
	On   string
}

func ParseEpicFile(path string) (*Epic, error) {
	content, err := ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseEpicContent(content), nil
}

func ParseEpicContent(content string) *Epic {
	epic := &Epic{}
	var currentSection string
	var sectionBuf *strings.Builder

	flushSection := func() {
		if sectionBuf == nil {
			return
		}
		text := strings.TrimSpace(sectionBuf.String())
		switch currentSection {
		case "Problem Statement":
			epic.ProblemStatement = text
		case "Goals", "Goals / Non-Goals":
			epic.Goals = text
		case "Non-Goals":
			epic.NonGoals = text
		case "Shared Context":
			epic.SharedContext = text
		}
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		if m := epicSectionRe.FindStringSubmatch(line); m != nil {
			flushSection()
			currentSection = m[1]
			sectionBuf = &strings.Builder{}
			continue
		}

		// Parse structured entries regardless of section
		if m := sharedDecisionRe.FindStringSubmatch(line); m != nil {
			epic.SharedDecisions = append(epic.SharedDecisions, SharedDecision{
				Decision:  m[1],
				Rationale: m[2],
			})
			continue
		}

		if m := incrementOrderRe.FindStringSubmatch(line); m != nil {
			epic.IncrementOrder = append(epic.IncrementOrder, IncrementRef{
				Name:        m[1],
				Description: m[2],
			})
			continue
		}

		if m := dependencyRe.FindStringSubmatch(line); m != nil {
			epic.Dependencies = append(epic.Dependencies, Dependency{
				From: m[1],
				On:   m[2],
			})
			continue
		}

		// Accumulate section text
		if sectionBuf != nil {
			if sectionBuf.Len() > 0 {
				sectionBuf.WriteString("\n")
			}
			sectionBuf.WriteString(line)
		}
	}

	flushSection()
	return epic
}

// IncrementNames returns the ordered list of increment names.
func (e *Epic) IncrementNames() []string {
	names := make([]string, len(e.IncrementOrder))
	for i, inc := range e.IncrementOrder {
		names[i] = inc.Name
	}
	return names
}

// DependenciesOf returns the names of increments that the given increment depends on.
func (e *Epic) DependenciesOf(incrementName string) []string {
	var deps []string
	for _, d := range e.Dependencies {
		if d.From == incrementName {
			deps = append(deps, d.On)
		}
	}
	return deps
}
