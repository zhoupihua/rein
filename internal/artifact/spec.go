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
)

type Spec struct {
	Sections     []string
	Requirements []Requirement
	Goals        string
	NonGoals     string
	Context      string
}

type Requirement struct {
	Name      string
	Scenarios []Scenario
}

type Scenario struct {
	Name string
	When string
	Then string
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

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		if m := sectionRe.FindStringSubmatch(line); m != nil {
			if currentScenario != nil && currentReq != nil {
				currentReq.Scenarios = append(currentReq.Scenarios, *currentScenario)
				currentScenario = nil
			}
			if currentReq != nil {
				spec.Requirements = append(spec.Requirements, *currentReq)
				currentReq = nil
			}
			spec.Sections = append(spec.Sections, m[1])
			switch m[1] {
			case "Context":
				spec.Context = ""
			case "Goals":
				spec.Goals = ""
			case "Non-Goals":
				spec.NonGoals = ""
			}
			continue
		}

		if m := requirementRe.FindStringSubmatch(line); m != nil {
			if currentScenario != nil && currentReq != nil {
				currentReq.Scenarios = append(currentReq.Scenarios, *currentScenario)
				currentScenario = nil
			}
			if currentReq != nil {
				spec.Requirements = append(spec.Requirements, *currentReq)
			}
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
			}
		}
	}

	if currentScenario != nil && currentReq != nil {
		currentReq.Scenarios = append(currentReq.Scenarios, *currentScenario)
	}
	if currentReq != nil {
		spec.Requirements = append(spec.Requirements, *currentReq)
	}

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
