package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure shell PATH for hostbook",
	Long:  `Adds the Go binary directory to your shell's PATH configuration file (.bashrc or .zshrc).`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			return
		}

		goBinPath := filepath.Join(homeDir, "go", "bin")
		pathExport := fmt.Sprintf("\nexport PATH=$PATH:%s", goBinPath)

		shell := os.Getenv("SHELL")
		var configFile string

		if strings.Contains(shell, "zsh") {
			configFile = filepath.Join(homeDir, ".zshrc")
		} else if strings.Contains(shell, "bash") {
			configFile = filepath.Join(homeDir, ".bashrc")
		} else {
			fmt.Println("Could not detect supported shell (bash/zsh). Please add manually:")
			fmt.Printf("export PATH=$PATH:%s\n", goBinPath)
			return
		}

		content, err := os.ReadFile(configFile)
		if err == nil {
			if strings.Contains(string(content), goBinPath) {
				fmt.Printf("PATH configuration already exists in %s\n", configFile)
				return
			}
		}

		f, err := os.OpenFile(configFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening config file: %v\n", err)
			return
		}
		defer f.Close()

		if _, err := f.WriteString(pathExport); err != nil {
			fmt.Printf("Error writing to config file: %v\n", err)
			return
		}

		fmt.Printf("Successfully added configuration to %s\n", configFile)
		fmt.Println("Please restart your terminal or run this command to apply changes:")
		fmt.Printf("source %s\n", configFile)

		// Check for sshpass
		if _, err := exec.LookPath("sshpass"); err != nil {
			fmt.Println("\n[Suggestion] Install 'sshpass' to enable auto-connect feature:")
			fmt.Println("  sudo apt install sshpass  (Debian/Ubuntu)")
			fmt.Println("  brew install sshpass      (macOS)")
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
