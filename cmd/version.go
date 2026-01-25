package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "1.26.1"
	BuildTime = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hostbook",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Hostbook v%s (Built: %s)\n", Version, BuildTime)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
