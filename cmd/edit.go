package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/engnhn/hostbook/core"
	"github.com/engnhn/hostbook/storage"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:               "edit [name]",
	Short:             "Edit a host",
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
				Message: "Select a host to edit:",
				Options: options,
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				fmt.Printf("Selection failed: %v\n", err)
				return
			}
		}

		var targetHost *core.Host
		targetIndex := -1
		for i, h := range hosts {
			if h.Name == name {
				targetHost = &hosts[i]
				targetIndex = i
				break
			}
		}

		if targetHost == nil {
			fmt.Printf("Host '%s' not found.\n", name)
			return
		}

		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("Host Name [%s]: ", targetHost.Name)
		newName, _ := reader.ReadString('\n')
		newName = strings.TrimSpace(newName)
		if newName != "" {
			targetHost.Name = newName
		}

		fmt.Printf("Hostname [%s]: ", targetHost.Hostname)
		newHostname, _ := reader.ReadString('\n')
		newHostname = strings.TrimSpace(newHostname)
		if newHostname != "" {
			targetHost.Hostname = newHostname
		}

		fmt.Printf("User [%s]: ", targetHost.User)
		newUser, _ := reader.ReadString('\n')
		newUser = strings.TrimSpace(newUser)
		if newUser != "" {
			targetHost.User = newUser
		}

		fmt.Printf("Port [%s]: ", targetHost.Port)
		newPort, _ := reader.ReadString('\n')
		newPort = strings.TrimSpace(newPort)
		if newPort != "" {
			targetHost.Port = newPort
		}

		fmt.Printf("Identity File Path [%s]: ", targetHost.IdentityFile)
		newIdentityFile, _ := reader.ReadString('\n')
		newIdentityFile = strings.TrimSpace(newIdentityFile)
		if newIdentityFile != "" {
			targetHost.IdentityFile = newIdentityFile
		}

		var modifyPassword bool
		prompt := &survey.Confirm{
			Message: "Do you want to modify/set the password?",
		}
		survey.AskOne(prompt, &modifyPassword)

		var newPassword string
		if modifyPassword {
			pwdPrompt := &survey.Password{
				Message: "New Password (leave empty to delete):",
			}
			survey.AskOne(pwdPrompt, &newPassword)
		}

		hosts[targetIndex] = *targetHost

		if err := s.SaveHosts(hosts); err != nil {
			fmt.Printf("Error saving hosts: %v\n", err)
			return
		}

		sshConfig := core.GenerateSSHConfig(hosts)
		if err := s.SaveSSHConfig(sshConfig); err != nil {
			fmt.Printf("Error saving SSH config: %v\n", err)
			return
		}

		if modifyPassword {
			if newPassword == "" {

				core.DeletePassword(targetHost.Name)
				fmt.Println("Password removed from keyring.")
			} else {
				if err := core.SavePassword(targetHost.Name, newPassword); err != nil {
					fmt.Printf("Warning: Failed to update password: %v\n", err)
				} else {
					fmt.Println("Password updated securely.")
				}
			}
		}

		fmt.Printf("Host '%s' updated successfully.\n", targetHost.Name)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
