package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/engnhn/hostbook/storage"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:               "connect [name]",
	Short:             "Connect to a host",
	ValidArgsFunction: getHostCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		var name string
		s, err := storage.NewStorage()
		if err != nil {
			fmt.Printf("Error initializing storage: %v\n", err)
			return
		}

		hosts, err := s.LoadHosts()
		if err != nil {
			fmt.Printf("Error loading hosts: %v\n", err)
			return
		}

		if len(hosts) == 0 {
			fmt.Println("No hosts found. Add one with 'hostbook add'.")
			return
		}

		if len(args) > 0 {
			name = args[0]
		} else {

			options := make([]string, len(hosts))
			for i, h := range hosts {
				options[i] = h.Name
			}

			prompt := &survey.Select{
				Message: "Select a host to connect:",
				Options: options,
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				fmt.Printf("Selection failed: %v\n", err)
				return
			}
		}

		found := false
		for _, h := range hosts {
			if h.Name == name {
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("Host '%s' not found.\n", name)
			return
		}

		sshConfigPath := s.GetSSHConfigPath()

		sshPath, err := exec.LookPath("ssh")
		if err != nil {
			fmt.Println("Error: ssh executable not found in PATH.")
			return
		}

		argsSSH := []string{"ssh", "-F", sshConfigPath, name}

		env := os.Environ()
		if err := syscall.Exec(sshPath, argsSSH, env); err != nil {
			fmt.Printf("Error executing ssh: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
