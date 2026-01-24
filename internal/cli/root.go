package cli

import (
	"fmt"
	"os"

	"github.com/paradise-runner/cider/internal/convert"
	"github.com/spf13/cobra"
)

var version = "dev"

// SetVersion sets the version string (called from main)
func SetVersion(v string) {
	version = v
}

var rootCmd = &cobra.Command{
	Use:   "cider",
	Short: "Bidirectionally sync Markdown files with Apple Notes",
	Long:  `Cider allows you to push Markdown files to Apple Notes and pull changes back, using front-matter to track state.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Check for pandoc availability (skip for help and completion commands)
		if cmd.Name() != "help" && cmd.Name() != "completion" && cmd.Name() != "__complete" {
			if err := convert.CheckPandocAvailable(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return err
			}
		}
		return nil
	},
}

// Execute runs the root command
func Execute() error {
	// Update version before executing (in case it was set via SetVersion)
	rootCmd.Version = version
	return rootCmd.Execute()
}

func init() {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(`{{.Version}}` + "\n")
}
