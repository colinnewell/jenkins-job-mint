package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version number that is baked in as the program is built.
//nolint:gochecknoglobals
var Version = "No version defined at build time"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Args:  standardValidation,
	Use:   "version",
	Short: "Display version number",
	Long:  `Provide information about the version of the tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
