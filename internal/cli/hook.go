package cli

import (
	"github.com/spf13/cobra"
	"github.com/zhoupihua/rein/internal/hook"
)

var hookCmd = &cobra.Command{
	Use:   "hook <name>",
	Short: "Run a hook handler (called by Claude Code events)",
	Args:  cobra.ExactArgs(1),
	RunE:  runHook,
}

func init() {
	rootCmd.AddCommand(hookCmd)
}

func runHook(cmd *cobra.Command, args []string) error {
	name := args[0]
	switch name {
	case "session-start":
		hook.SessionStart()
	case "guard":
		hook.Guard()
	case "guard-bash":
		hook.GuardBash()
	case "gate":
		hook.Gate()
	case "format":
		hook.Format()
	case "checkbox-guard":
		hook.CheckboxGuard()
	case "task-progress":
		hook.TaskProgress()
	case "leak-guard":
		hook.LeakGuard()
	case "inject":
		hook.Inject()
	default:
		return nil
	}
	return nil
}
