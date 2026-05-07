package artifact

import (
	"bufio"
	"regexp"
	"strings"
)

var (
	taskDetailRe    = regexp.MustCompile(`^###\s+(\d+\.\d+)\s+(.+)$`)
	goalRe          = regexp.MustCompile(`^\*\*Goal:\*\*\s*(.+)$`)
	acceptanceRe    = regexp.MustCompile(`^-\s+\*\*Acceptance:\*\*\s*(.+)$`)
	verificationRe  = regexp.MustCompile(`^-\s+\*\*Verification:\*\*\s*(.+)$`)
	dependenciesRe  = regexp.MustCompile(`^-\s+\*\*Dependencies:\*\*\s*(.+)$`)
	filesRe         = regexp.MustCompile(`^-\s+\*\*Files:\*\*\s*(.+)$`)
	scopeRe         = regexp.MustCompile(`^-\s+\*\*Scope:\*\*\s*(.+)$`)
	notesRe         = regexp.MustCompile(`^-\s+\*\*Notes:\*\*\s*(.+)$`)
	approachRe      = regexp.MustCompile(`^-\s+\*\*Approach:\*\*\s*(.+)$`)
	edgeCasesRe     = regexp.MustCompile(`^-\s+\*\*Edge Cases:\*\*\s*(.+)$`)
	rollbackRe      = regexp.MustCompile(`^-\s+\*\*Rollback:\*\*\s*(.+)$`)
	planSectionRe   = regexp.MustCompile(`^##\s+(.+)$`)
)

type Plan struct {
	Goal            string
	Architecture    string
	DependencyGraph string
	SliceStrategy   string
	RiskTable       string
	Parallelization string
	SelfAudit       string
	Handoff         string
	TaskDetails     []TaskDetail
}

type TaskDetail struct {
	ID           TaskID
	Title        string
	Acceptance   string
	Verification string
	Dependencies string
	Files        string
	Scope        string
	Notes        string
	Approach     string
	EdgeCases    string
	Rollback     string
}

func ParsePlanFile(path string) (*Plan, error) {
	content, err := ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParsePlanContent(content), nil
}

func ParsePlanContent(content string) *Plan {
	plan := &Plan{}
	var current *TaskDetail
	var currentSection string
	var sectionBuf *strings.Builder

	flushSection := func() {
		if sectionBuf == nil {
			return
		}
		text := strings.TrimSpace(sectionBuf.String())
		switch currentSection {
		case "Architecture Overview", "Architecture":
			plan.Architecture = text
		case "Dependency Graph":
			plan.DependencyGraph = text
		case "Vertical Slice Strategy":
			plan.SliceStrategy = text
		case "Risk/Mitigation Table", "Risks and Mitigations":
			plan.RiskTable = text
		case "Parallelization":
			plan.Parallelization = text
		case "Self-Audit Checklist":
			plan.SelfAudit = text
		case "Handoff":
			plan.Handoff = text
		}
	}

	flushTask := func() {
		if current != nil {
			plan.TaskDetails = append(plan.TaskDetails, *current)
			current = nil
		}
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		// Section headings
		if m := planSectionRe.FindStringSubmatch(line); m != nil {
			flushTask()
			flushSection()
			currentSection = m[1]
			sectionBuf = &strings.Builder{}

			switch currentSection {
			case "Task Details":
				// entering task details section
			}
			continue
		}

		// Task detail headings (only within Task Details section or when no section tracking)
		if m := taskDetailRe.FindStringSubmatch(line); m != nil {
			if current != nil {
				plan.TaskDetails = append(plan.TaskDetails, *current)
			}
			id, _ := ParseTaskID(m[1])
			current = &TaskDetail{ID: id, Title: m[2]}
			continue
		}

		// Within a task detail
		if current != nil {
			if m := acceptanceRe.FindStringSubmatch(line); m != nil {
				current.Acceptance = m[1]
			} else if m := verificationRe.FindStringSubmatch(line); m != nil {
				current.Verification = m[1]
			} else if m := dependenciesRe.FindStringSubmatch(line); m != nil {
				current.Dependencies = m[1]
			} else if m := filesRe.FindStringSubmatch(line); m != nil {
				current.Files = m[1]
			} else if m := scopeRe.FindStringSubmatch(line); m != nil {
				current.Scope = m[1]
			} else if m := notesRe.FindStringSubmatch(line); m != nil {
				current.Notes = m[1]
			} else if m := approachRe.FindStringSubmatch(line); m != nil {
				current.Approach = m[1]
			} else if m := edgeCasesRe.FindStringSubmatch(line); m != nil {
				current.EdgeCases = m[1]
			} else if m := rollbackRe.FindStringSubmatch(line); m != nil {
				current.Rollback = m[1]
			}
			continue
		}

		// Goal (before any section or task detail)
		if m := goalRe.FindStringSubmatch(line); m != nil {
			plan.Goal = m[1]
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

	flushTask()
	flushSection()

	return plan
}

func (p *Plan) FindTaskDetail(id TaskID) *TaskDetail {
	for i := range p.TaskDetails {
		if p.TaskDetails[i].ID == id {
			return &p.TaskDetails[i]
		}
	}
	return nil
}

func (p *Plan) TaskIDs() []TaskID {
	var ids []TaskID
	for _, td := range p.TaskDetails {
		ids = append(ids, td.ID)
	}
	return ids
}
