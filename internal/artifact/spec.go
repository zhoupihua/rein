package artifact

import (
	"bufio"
	"regexp"
	"strings"
)

var (
	sectionRe     = regexp.MustCompile(`^##\s+(.+)$`)
	requirementRe = regexp.MustCompile(`^###\s+Requirement:\s*(.+)$`)
	scenarioRe    = regexp.MustCompile(`^####\s+Scenario:\s*(.+)$`)
	whenRe        = regexp.MustCompile(`^-\s+\*\*WHEN\*\*\s*(.+)$`)
	thenRe        = regexp.MustCompile(`^-\s+\*\*THEN\*\*\s*(.+)$`)
	testRe        = regexp.MustCompile(`^-\s+\*\*TEST\*\*\s+(.+)$`)
	decisionRe    = regexp.MustCompile(`^-\s+\*\*Decision:\*\*\s*(.+?)\s*—\s*\*\*Rationale:\*\*\s*(.+)$`)
)

type Spec struct {
	Sections     []string
	Requirements []Requirement
	Decisions    []Decision
	Risks        string
}

type Requirement struct {
	Name      string
	Scenarios []Scenario
}

type Scenario struct {
	Name string
	When string
	Then string
	Test string
}

type Decision struct {
	Choice    string
	Rationale string
}

func ParseSpecFile(path string) (*Spec, error) {
	content, err := ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseSpecContent(content), nil
}

func ParseSpecContent(content string) *Spec {
	spec := &Spec{}
	var currentReq *Requirement
	var currentScenario *Scenario
	var currentSection string
	var sectionBuf *strings.Builder

	flushReq := func() {
		if currentScenario != nil && currentReq != nil {
			currentReq.Scenarios = append(currentReq.Scenarios, *currentScenario)
			currentScenario = nil
		}
		if currentReq != nil {
			spec.Requirements = append(spec.Requirements, *currentReq)
			currentReq = nil
		}
	}

	flushSection := func() {
		if sectionBuf == nil {
			return
		}
		text := strings.TrimSpace(sectionBuf.String())
		switch currentSection {
		case "Risks", "Risks / Trade-offs":
			spec.Risks = text
		}
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		if m := sectionRe.FindStringSubmatch(line); m != nil {
			flushReq()
			flushSection()
			currentSection = m[1]
			sectionBuf = &strings.Builder{}
			spec.Sections = append(spec.Sections, m[1])
			continue
		}

		// Parse structured entries
		if m := requirementRe.FindStringSubmatch(line); m != nil {
			flushReq()
			currentReq = &Requirement{Name: m[1]}
			continue
		}

		if m := scenarioRe.FindStringSubmatch(line); m != nil {
			if currentScenario != nil && currentReq != nil {
				currentReq.Scenarios = append(currentReq.Scenarios, *currentScenario)
			}
			currentScenario = &Scenario{Name: m[1]}
			continue
		}

		if currentScenario != nil {
			if m := whenRe.FindStringSubmatch(line); m != nil {
				currentScenario.When = m[1]
			} else if m := thenRe.FindStringSubmatch(line); m != nil {
				currentScenario.Then = m[1]
			} else if m := testRe.FindStringSubmatch(line); m != nil {
				currentScenario.Test = strings.Trim(m[1], "`")
			}
			continue
		}

		if m := decisionRe.FindStringSubmatch(line); m != nil {
			spec.Decisions = append(spec.Decisions, Decision{
				Choice:    m[1],
				Rationale: m[2],
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

	flushReq()
	flushSection()

	return spec
}

func (s *Spec) SectionNames() []string {
	return s.Sections
}

func (s *Spec) HasSection(name string) bool {
	for _, sec := range s.Sections {
		if sec == name {
			return true
		}
	}
	return false
}
