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
)

type Plan struct {
	Goal         string
	TaskDetails  []TaskDetail
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

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		if m := taskDetailRe.FindStringSubmatch(line); m != nil {
			if current != nil {
				plan.TaskDetails = append(plan.TaskDetails, *current)
			}
			id, _ := ParseTaskID(m[1])
			current = &TaskDetail{ID: id, Title: m[2]}
			continue
		}

		if current == nil {
			if m := goalRe.FindStringSubmatch(line); m != nil {
				plan.Goal = m[1]
			}
			continue
		}

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
		}
	}

	if current != nil {
		plan.TaskDetails = append(plan.TaskDetails, *current)
	}

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
