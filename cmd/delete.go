package cmd

import (
	"fmt"

	"github.com/engnhn/hostbook/core"
	"github.com/engnhn/hostbook/storage"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:               "delete [name]",
	Short:             "Delete a host",
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
			fmt.Println("No hosts found.")
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
				Message: "Select a host to delete:",
				Options: options,
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				fmt.Printf("Selection failed: %v\n", err)
				return
			}
		}

		newHosts := []core.Host{}
		deleted := false
		for _, h := range hosts {
			if h.Name == name {
				deleted = true
				continue
			}
			newHosts = append(newHosts, h)
		}

		if !deleted {
			fmt.Printf("Host '%s' not found.\n", name)
			return
		}

		if err := s.SaveHosts(newHosts); err != nil {
			fmt.Printf("Error saving hosts: %v\n", err)
			return
		}

		sshConfig := core.GenerateSSHConfig(newHosts)
		if err := s.SaveSSHConfig(sshConfig); err != nil {
			fmt.Printf("Error saving SSH config: %v\n", err)
			return
		}

		fmt.Printf("Host '%s' deleted successfully.\n", name)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
