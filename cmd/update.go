package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update hostbook to the latest version",
	Long:  `Updates hostbook to the latest version by running "go install github.com/engnhn/hostbook@latest".`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Checking for updates...")

		goBin, err := exec.LookPath("go")
		if err != nil {
			fmt.Println("Error: 'go' is not installed or not in PATH. Cannot update.")
			return
		}

		fmt.Println("Running: go install github.com/engnhn/hostbook@latest")
		installCmd := exec.Command(goBin, "install", "github.com/engnhn/hostbook@latest")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr

		if err := installCmd.Run(); err != nil {
			fmt.Printf("Update failed: %v\n", err)
			return
		}

		fmt.Println("Hostbook updated successfully!")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
