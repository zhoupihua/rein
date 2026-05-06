package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	jsonMode bool
	ver      string
)

func SetVersion(v string) {
	ver = v
	rootCmd.Version = v
}

var rootCmd = &cobra.Command{
	Use:     "rein",
	Short:   "rein — AI coding workflow CLI",
	Version: ver,
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonMode, "json", false, "output as JSON")
}

func Execute() error {
	return rootCmd.Execute()
}

func isJSON() bool {
	return jsonMode
}

func exitError(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
