package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/engnhn/hostbook/core"
	"github.com/engnhn/hostbook/storage"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new host",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Host Name (e.g. my-server): ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		fmt.Print("Hostname (IP or Domain): ")
		hostname, _ := reader.ReadString('\n')
		hostname = strings.TrimSpace(hostname)

		fmt.Print("User: ")
		user, _ := reader.ReadString('\n')
		user = strings.TrimSpace(user)

		fmt.Print("Port [22]: ")
		port, _ := reader.ReadString('\n')
		port = strings.TrimSpace(port)
		if port == "" {
			port = "22"
		}

		fmt.Print("Identity File Path (optional): ")
		identityFile, _ := reader.ReadString('\n')
		identityFile = strings.TrimSpace(identityFile)

		fmt.Print("Tags (comma separated, e.g. prod,db): ")
		tagsStr, _ := reader.ReadString('\n')
		tagsStr = strings.TrimSpace(tagsStr)
		var tags []string
		if tagsStr != "" {
			parts := strings.Split(tagsStr, ",")
			for _, p := range parts {
				tags = append(tags, strings.TrimSpace(p))
			}
		}

		var password string
		var savePassword bool
		prompt := &survey.Confirm{
			Message: "Do you want to save a password for this host securely?",
		}
		survey.AskOne(prompt, &savePassword)

		if savePassword {
			pwdPrompt := &survey.Password{
				Message: "Password:",
			}
			survey.AskOne(pwdPrompt, &password)
		}

		host := core.Host{
			Name:         name,
			Hostname:     hostname,
			User:         user,
			Port:         port,
			IdentityFile: identityFile,
			Tags:         tags,
		}

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

		for _, h := range hosts {
			if h.Name == name {
				fmt.Printf("Host with name '%s' already exists.\n", name)
				return
			}
		}

		hosts = append(hosts, host)

		if err := s.SaveHosts(hosts); err != nil {
			fmt.Printf("Error saving host: %v\n", err)
			return
		}

		sshConfig := core.GenerateSSHConfig(hosts)
		if err := s.SaveSSHConfig(sshConfig); err != nil {
			fmt.Printf("Error saving SSH config: %v\n", err)
			return
		}

		if savePassword && password != "" {
			if err := core.SavePassword(name, password); err != nil {
				fmt.Printf("Warning: Failed to save password securely: %v\n", err)
			} else {
				fmt.Println("Password saved securely in system keyring.")
			}
		}

		fmt.Printf("Host '%s' added successfully!\n", name)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
