package artifact

import (
	"path/filepath"
	"testing"
)

func TestParseTaskWithSubTasks(t *testing.T) {
	path := filepath.Join("testdata", "sample-task.md")
	tf, err := ParseTaskFile(path)
	if err != nil {
		t.Fatalf("ParseTaskFile: %v", err)
	}

	// Task 3.1 has 3 sub-tasks (RED, GREEN, REFACTOR)
	task := tf.FindTask(TaskID{Phase: 3, Seq: 1})
	if task == nil {
		t.Fatal("task 3.1 not found")
	}
	if len(task.SubTasks) != 3 {
		t.Fatalf("expected 3 sub-tasks, got %d", len(task.SubTasks))
	}
	if task.SubTasks[0].Description != "RED: Test user model fields" {
		t.Errorf("sub-task 0 = %q", task.SubTasks[0].Description)
	}
	if task.SubTasks[0].Done {
		t.Error("sub-task 0 should be unchecked")
	}
	if !task.SubTasks[1].Done {
		t.Error("sub-task 1 should be checked (GREEN)")
	}
	if task.SubTasks[1].Description != "GREEN: Implement user struct" {
		t.Errorf("sub-task 1 = %q", task.SubTasks[1].Description)
	}
	if task.SubTasks[2].Description != "REFACTOR: Extract validation helpers" {
		t.Errorf("sub-task 2 = %q", task.SubTasks[2].Description)
	}
}

func TestParseTaskWithoutSubTasks(t *testing.T) {
	path := filepath.Join("testdata", "sample-task.md")
	tf, err := ParseTaskFile(path)
	if err != nil {
		t.Fatalf("ParseTaskFile: %v", err)
	}

	// Task 3.2 has no sub-tasks
	task := tf.FindTask(TaskID{Phase: 3, Seq: 2})
	if task == nil {
		t.Fatal("task 3.2 not found")
	}
	if len(task.SubTasks) != 0 {
		t.Errorf("expected 0 sub-tasks, got %d", len(task.SubTasks))
	}
}

func TestFirstUncheckedWithSubTasks(t *testing.T) {
	content := `# Test Tasks

## 1. Build

- [x] 1.1 Done task
  - [x] RED: test
  - [x] GREEN: impl
- [ ] 1.2 Incomplete task
  - [x] RED: test
  - [ ] GREEN: impl
- [ ] 1.3 Next task
`
	tf := ParseTaskContent("test.md", content)
	task := tf.FirstUnchecked()
	if task == nil {
		t.Fatal("expected unchecked task")
	}
	if task.ID != (TaskID{Phase: 1, Seq: 2}) {
		t.Errorf("FirstUnchecked ID = %v, want 1.2", task.ID)
	}
}

func TestCountDoneWithSubTasks(t *testing.T) {
	content := `# Test Tasks

## 1. Build

- [x] 1.1 Done task
  - [x] RED: test
  - [x] GREEN: impl
- [ ] 1.2 Partial task
  - [x] RED: test
  - [ ] GREEN: impl
- [ ] 1.3 No sub-tasks
`
	tf := ParseTaskContent("test.md", content)
	done, total := tf.CountDone()
	// 1.1: 2 sub-tasks done → 2/2
	// 1.2: 1 sub-task done → 1/2
	// 1.3: no sub-tasks, not done → 0/1
	// Total: 3 done, 5 total
	if done != 3 || total != 5 {
		t.Errorf("CountDone = %d/%d, want 3/5", done, total)
	}
}

func TestCheckSubTask(t *testing.T) {
	content := `# Test Tasks

## 1. Build

- [ ] 1.1 Task with sub-tasks
  - [ ] RED: test
  - [ ] GREEN: impl
  - [ ] REFACTOR: cleanup
- [ ] 1.2 Other task
`
	tmpPath := filepath.Join(t.TempDir(), "task.md")
	if err := WriteFile(tmpPath, content); err != nil {
		t.Fatal(err)
	}

	tf, err := ParseTaskFile(tmpPath)
	if err != nil {
		t.Fatal(err)
	}

	ok := tf.CheckSubTask(TaskID{Phase: 1, Seq: 1}, 1)
	if !ok {
		t.Fatal("CheckSubTask returned false")
	}

	tf2, err := ParseTaskFile(tmpPath)
	if err != nil {
		t.Fatal(err)
	}
	task := tf2.FindTask(TaskID{Phase: 1, Seq: 1})
	if task == nil {
		t.Fatal("task not found after CheckSubTask")
	}
	if !task.SubTasks[1].Done {
		t.Error("sub-task 1 should be checked after CheckSubTask")
	}
	if task.SubTasks[0].Done {
		t.Error("sub-task 0 should still be unchecked")
	}
}

func TestCheckSubTaskAutoChecksParent(t *testing.T) {
	content := `# Test Tasks

## 1. Build

- [ ] 1.1 Task with sub-tasks
  - [x] RED: test
  - [ ] GREEN: impl
`
	tmpPath := filepath.Join(t.TempDir(), "task.md")
	if err := WriteFile(tmpPath, content); err != nil {
		t.Fatal(err)
	}

	tf, err := ParseTaskFile(tmpPath)
	if err != nil {
		t.Fatal(err)
	}

	ok := tf.CheckSubTask(TaskID{Phase: 1, Seq: 1}, 1)
	if !ok {
		t.Fatal("CheckSubTask returned false")
	}

	tf2, err := ParseTaskFile(tmpPath)
	if err != nil {
		t.Fatal(err)
	}
	task := tf2.FindTask(TaskID{Phase: 1, Seq: 1})
	if task == nil {
		t.Fatal("task not found")
	}
	if !task.Done {
		t.Error("parent task should be auto-checked when all sub-tasks done")
	}
}

func TestSubTaskIDParsing(t *testing.T) {
	tests := []struct {
		input   string
		want    SubTaskID
		wantOK  bool
	}{
		{"1.2.0", SubTaskID{Parent: TaskID{Phase: 1, Seq: 2}, Index: 0}, true},
		{"3.1.2", SubTaskID{Parent: TaskID{Phase: 3, Seq: 1}, Index: 2}, true},
		{"1.2", SubTaskID{}, false},   // only 2 segments
		{"1.2.0.3", SubTaskID{}, false}, // 4 segments
		{"a.b.c", SubTaskID{}, false},   // non-numeric
	}
	for _, tt := range tests {
		got, ok := ParseSubTaskID(tt.input)
		if ok != tt.wantOK {
			t.Errorf("ParseSubTaskID(%q) ok = %v, want %v", tt.input, ok, tt.wantOK)
			continue
		}
		if ok && got != tt.want {
			t.Errorf("ParseSubTaskID(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestSubTaskIDString(t *testing.T) {
	id := SubTaskID{Parent: TaskID{Phase: 1, Seq: 2}, Index: 0}
	if id.String() != "1.2.0" {
		t.Errorf("SubTaskID.String() = %q, want %q", id.String(), "1.2.0")
	}
}

func TestIsComplete(t *testing.T) {
	tests := []struct {
		name string
		task Task
		want bool
	}{
		{
			name: "done with all sub-tasks done",
			task: Task{Done: true, SubTasks: []SubTask{
				{Done: true},
				{Done: true},
			}},
			want: true,
		},
		{
			name: "done but sub-task unchecked",
			task: Task{Done: true, SubTasks: []SubTask{
				{Done: true},
				{Done: false},
			}},
			want: false,
		},
		{
			name: "not done no sub-tasks",
			task: Task{Done: false},
			want: false,
		},
		{
			name: "done no sub-tasks",
			task: Task{Done: true},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.task.IsComplete(); got != tt.want {
				t.Errorf("IsComplete() = %v, want %v", got, tt.want)
			}
		})
	}
}
