package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize rein in current project (create docs/rein/ directories)",
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	projectDir, err := os.Getwd()
	if err != nil {
		return err
	}
	dirs := []string{"changes", "archive"}
	for _, sub := range dirs {
		path := filepath.Join(projectDir, "docs", "rein", sub)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("create %s: %w", path, err)
		}
	}
	fmt.Println("docs/rein/{changes,archive} created")
	fmt.Println("Run 'rein status' to check current phase")
	return nil
}
