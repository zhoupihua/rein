package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhoupihua/rein/internal/artifact"
	"github.com/zhoupihua/rein/internal/output"
	"github.com/zhoupihua/rein/internal/project"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage tasks",
}

var taskNextCmd = &cobra.Command{
	Use:   "next",
	Short: "Show next unchecked task",
	RunE:  runTaskNext,
}

var taskDoneCmd = &cobra.Command{
	Use:   "done <id>",
	Short: "Mark task as done (e.g., rein task done 1.1)",
	Args:  cobra.ExactArgs(1),
	RunE:  runTaskDone,
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks with status",
	RunE:  runTaskList,
}

func init() {
	rootCmd.AddCommand(taskCmd)
	taskCmd.AddCommand(taskNextCmd)
	taskCmd.AddCommand(taskDoneCmd)
	taskCmd.AddCommand(taskListCmd)
}

type TaskResult struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type TaskListResult struct {
	Feature string       `json:"feature"`
	Tasks   []TaskResult `json:"tasks"`
	Done    int          `json:"done"`
	Total   int          `json:"total"`
}

func resolveFeatureWithTask(p *project.Project) (string, error) {
	name := project.FindFeatureWithTask(p)
	if name == "" {
		return "", fmt.Errorf("no task.md found in any feature under %s", project.ChangesDir)
	}
	return name, nil
}

func runTaskNext(cmd *cobra.Command, args []string) error {
	p, err := project.Resolve()
	if err != nil {
		output.PrintError(err, isJSON())
		return nil
	}

	name, err := resolveFeatureWithTask(p)
	if err != nil {
		output.PrintError(err, isJSON())
		return nil
	}

	tf, err := artifact.ParseTaskFile(project.TaskFilePath(p.Dir, name))
	if err != nil {
		output.PrintError(err, isJSON())
		return nil
	}

	task := tf.FirstUnchecked()
	if task == nil {
		output.Print(map[string]string{
			"feature": name,
			"message": "all tasks complete",
		}, isJSON())
		return nil
	}

	output.Print(TaskResult{
		ID:          task.ID.String(),
		Description: task.Description,
		Done:        false,
	}, isJSON())
	return nil
}

func runTaskDone(cmd *cobra.Command, args []string) error {
	id, ok := artifact.ParseTaskID(args[0])
	if !ok {
		output.PrintError(fmt.Errorf("invalid task ID: %s (expected format like 1.1)", args[0]), isJSON())
		return nil
	}

	p, err := project.Resolve()
	if err != nil {
		output.PrintError(err, isJSON())
		return nil
	}

	name, err := resolveFeatureWithTask(p)
	if err != nil {
		output.PrintError(err, isJSON())
		return nil
	}

	tf, err := artifact.ParseTaskFile(project.TaskFilePath(p.Dir, name))
	if err != nil {
		output.PrintError(err, isJSON())
		return nil
	}

	if tf.FindTask(id) == nil {
		output.PrintError(fmt.Errorf("task %s not found", id), isJSON())
		return nil
	}

	if tf.CheckTask(id) {
		if isJSON() {
			output.PrintJSON(map[string]string{"status": "done", "task": id.String(), "feature": name})
		} else {
			fmt.Printf("Task %s marked as done (feature: %s)\n", id, name)
		}
	} else {
		output.PrintError(fmt.Errorf("task %s already done or not found", id), isJSON())
	}
	return nil
}

func runTaskList(cmd *cobra.Command, args []string) error {
	p, err := project.Resolve()
	if err != nil {
		output.PrintError(err, isJSON())
		return nil
	}

	name, err := resolveFeatureWithTask(p)
	if err != nil {
		output.PrintError(err, isJSON())
		return nil
	}

	tf, err := artifact.ParseTaskFile(project.TaskFilePath(p.Dir, name))
	if err != nil{
		output.PrintError(err, isJSON())
		return nil
	}

	var results []TaskResult
	done, total := tf.CountDone()
	for _, t := range tf.AllTasks() {
		results = append(results, TaskResult{
			ID:          t.ID.String(),
			Description: t.Description,
			Done:        t.Done,
		})
	}

	output.Print(TaskListResult{
		Feature: name,
		Tasks:   results,
		Done:    done,
		Total:   total,
	}, isJSON())
	return nil
}
