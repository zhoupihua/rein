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
	subTaskCheckedRe = regexp.MustCompile(`^\s{2,}- \[x\] (.+)$`)
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
	SubTasks    []SubTask
	Lines       []string
}

type SubTask struct {
	Index       int
	Description string
	Done        bool
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
			if m := subTaskCheckedRe.FindStringSubmatch(line); m != nil {
				currentTask.SubTasks = append(currentTask.SubTasks, SubTask{
					Index:       len(currentTask.SubTasks),
					Description: m[1],
					Done:        true,
				})
			} else if m := subTaskRe.FindStringSubmatch(line); m != nil {
				currentTask.SubTasks = append(currentTask.SubTasks, SubTask{
					Index:       len(currentTask.SubTasks),
					Description: m[1],
					Done:        false,
				})
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
			if !tf.Phases[i].Tasks[j].IsComplete() {
				return &tf.Phases[i].Tasks[j]
			}
		}
	}
	return nil
}

func (tf *TaskFile) CountDone() (done, total int) {
	for _, p := range tf.Phases {
		for _, t := range p.Tasks {
			if len(t.SubTasks) > 0 {
				for _, st := range t.SubTasks {
					total++
					if st.Done {
						done++
					}
				}
			} else {
				total++
				if t.Done {
					done++
				}
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

func (t *Task) IsComplete() bool {
	if !t.Done {
		return false
	}
	for _, st := range t.SubTasks {
		if !st.Done {
			return false
		}
	}
	return true
}

func (t *Task) FirstUncheckedSubTask() int {
	for i, st := range t.SubTasks {
		if !st.Done {
			return i
		}
	}
	return -1
}

func (tf *TaskFile) CheckSubTask(parentID TaskID, index int) bool {
	content, err := ReadFile(tf.Path)
	if err != nil {
		return false
	}

	var lines []string
	inTarget := false
	subIdx := 0
	found := false

	for _, line := range strings.Split(content, "\n") {
		if found {
			lines = append(lines, line)
			continue
		}

		if m := taskUncheckedRe.FindStringSubmatch(line); m != nil {
			parsedID, _ := ParseTaskID(m[2])
			if parsedID == parentID {
				inTarget = true
			} else {
				inTarget = false
			}
			subIdx = 0
			lines = append(lines, line)
			continue
		}
		if m := taskCheckedRe.FindStringSubmatch(line); m != nil {
			parsedID, _ := ParseTaskID(m[2])
			if parsedID == parentID {
				inTarget = true
			} else {
				inTarget = false
			}
			subIdx = 0
			lines = append(lines, line)
			continue
		}

		if inTarget && !found {
			if m := subTaskRe.FindStringSubmatch(line); m != nil {
				if subIdx == index {
					lines = append(lines, "  - [x] "+m[1])
					found = true
					continue
				}
				subIdx++
			} else if m := subTaskCheckedRe.FindStringSubmatch(line); m != nil {
				subIdx++
				lines = append(lines, line)
				continue
			} else if !strings.HasPrefix(line, "  ") {
				inTarget = false
			}
		}

		lines = append(lines, line)
	}

	if found {
		WriteFile(tf.Path, strings.Join(lines, "\n"))
		task := tf.FindTask(parentID)
		if task != nil {
			task.SubTasks[index].Done = true
			allDone := true
			for _, st := range task.SubTasks {
				if !st.Done {
					allDone = false
					break
				}
			}
			if allDone {
				tf.CheckTask(parentID)
			}
		}
	}
	return found
}
