package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall hostbook",
	Long:  `Removes the .hostbook data directory and tries to remove the hostbook binary.`,
	Run: func(cmd *cobra.Command, args []string) {
		confirm := false
		prompt := &survey.Confirm{
			Message: "Are you sure you want to uninstall hostbook? This will delete all your saved hosts.",
		}
		survey.AskOne(prompt, &confirm)

		if !confirm {
			fmt.Println("Uninstall cancelled.")
			return
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			return
		}

		configDir := filepath.Join(homeDir, ".hostbook")
		if err := os.RemoveAll(configDir); err != nil {
			fmt.Printf("Error removing data directory: %v\n", err)
		} else {
			fmt.Printf("Removed data directory: %s\n", configDir)
		}

		exePath, err := os.Executable()
		if err != nil {
			fmt.Printf("Could not determine executable path: %v\n", err)
			fmt.Println("Please remove the binary manually.")
			return
		}

		if strings.Contains(exePath, "go-build") {
			fmt.Println("Running in development mode (go run). Skipping binary removal.")
			return
		}

		if err := os.Remove(exePath); err != nil {
			fmt.Printf("Error removing binary: %v\n", err)
			fmt.Printf("Please remove it manually: %s\n", exePath)
		} else {
			fmt.Printf("Removed binary: %s\n", exePath)
		}

		fmt.Println("Hostbook uninstalled successfully.")
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
