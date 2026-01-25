package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hostbook",
	Short: "Hostbook manages your SSH connections",
	Long: `Hostbook is a CLI tool to manage your SSH connections locally.
It generates SSH configuration files and helps you connect to your servers easily.
No data is sent to external servers.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
}
