package artifact

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

var (
	phaseHeadingRe = regexp.MustCompile(`^##\s+(\d+)\.\s+(.+)$`)
	taskUncheckedRe = regexp.MustCompile(`^(\s*)- \[ \] (\d+\.\d+)\s+(.+)$`)
	taskCheckedRe   = regexp.MustCompile(`^(\s*)- \[x\] (\d+\.\d+)\s+(.+)$`)
	subTaskRe       = regexp.MustCompile(`^\s{2,}- \[ \] (.+)$`)
)

type TaskFile struct {
	Path   string
	Phases []Phase
}

type Phase struct {
	Number int
	Name   string
	Tasks  []Task
}

type Task struct {
	ID          TaskID
	Description string
	Done        bool
	SubTasks    []string
	Lines       []string
}

func ParseTaskFile(path string) (*TaskFile, error) {
	content, err := ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseTaskContent(path, content), nil
}

func ParseTaskContent(path, content string) *TaskFile {
	tf := &TaskFile{Path: path}
	var currentPhase *Phase
	var currentTask *Task

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		if m := phaseHeadingRe.FindStringSubmatch(line); m != nil {
			if currentTask != nil && currentPhase != nil {
				currentPhase.Tasks = append(currentPhase.Tasks, *currentTask)
				currentTask = nil
			}
			var num int
			fmt.Sscanf(m[1], "%d", &num)
			currentPhase = &Phase{Number: num, Name: m[2]}
			tf.Phases = append(tf.Phases, *currentPhase)
			currentPhase = &tf.Phases[len(tf.Phases)-1]
			continue
		}

		if m := taskUncheckedRe.FindStringSubmatch(line); m != nil {
			if currentTask != nil && currentPhase != nil {
				currentPhase.Tasks = append(currentPhase.Tasks, *currentTask)
			}
			id, _ := ParseTaskID(m[2])
			currentTask = &Task{
				ID:          id,
				Description: m[3],
				Done:        false,
				Lines:       []string{line},
			}
			continue
		}

		if m := taskCheckedRe.FindStringSubmatch(line); m != nil {
			if currentTask != nil && currentPhase != nil {
				currentPhase.Tasks = append(currentPhase.Tasks, *currentTask)
			}
			id, _ := ParseTaskID(m[2])
			currentTask = &Task{
				ID:          id,
				Description: m[3],
				Done:        true,
				Lines:       []string{line},
			}
			continue
		}

		if currentTask != nil {
			if m := subTaskRe.FindStringSubmatch(line); m != nil {
				currentTask.SubTasks = append(currentTask.SubTasks, m[1])
			}
			currentTask.Lines = append(currentTask.Lines, line)
		}
	}

	if currentTask != nil && currentPhase != nil {
		currentPhase.Tasks = append(currentPhase.Tasks, *currentTask)
	}

	return tf
}

func (tf *TaskFile) AllTasks() []Task {
	var tasks []Task
	for _, p := range tf.Phases {
		tasks = append(tasks, p.Tasks...)
	}
	return tasks
}

func (tf *TaskFile) FirstUnchecked() *Task {
	for i := range tf.Phases {
		for j := range tf.Phases[i].Tasks {
			if !tf.Phases[i].Tasks[j].Done {
				return &tf.Phases[i].Tasks[j]
			}
		}
	}
	return nil
}

func (tf *TaskFile) CountDone() (done, total int) {
	for _, p := range tf.Phases {
		for _, t := range p.Tasks {
			total++
			if t.Done {
				done++
			}
		}
	}
	return
}

func (tf *TaskFile) TaskIDs() []TaskID {
	var ids []TaskID
	for _, p := range tf.Phases {
		for _, t := range p.Tasks {
			ids = append(ids, t.ID)
		}
	}
	return ids
}

func (tf *TaskFile) FindTask(id TaskID) *Task {
	for i := range tf.Phases {
		for j := range tf.Phases[i].Tasks {
			if tf.Phases[i].Tasks[j].ID == id {
				return &tf.Phases[i].Tasks[j]
			}
		}
	}
	return nil
}

func (tf *TaskFile) CheckTask(id TaskID) bool {
	found := false
	content, err := ReadFile(tf.Path)
	if err != nil {
		return false
	}

	var lines []string
	for _, line := range strings.Split(content, "\n") {
		if !found {
			if m := taskUncheckedRe.FindStringSubmatch(line); m != nil {
				parsedID, _ := ParseTaskID(m[2])
				if parsedID == id {
					lines = append(lines, m[1]+"- [x] "+id.String()+" "+m[3])
					found = true
					continue
				}
			}
		}
		lines = append(lines, line)
	}

	if found {
		WriteFile(tf.Path, strings.Join(lines, "\n"))
	}
	return found
}
