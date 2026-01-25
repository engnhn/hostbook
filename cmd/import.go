package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/engnhn/hostbook/core"
	"github.com/engnhn/hostbook/storage"
	"github.com/kevinburke/ssh_config"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import hosts from your existing SSH config",
	Long:  `Reads ~/.ssh/config and imports hosts into Hostbook.`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			return
		}

		sshConfigPath := filepath.Join(homeDir, ".ssh", "config")
		f, err := os.Open(sshConfigPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("No SSH config file found at ~/.ssh/config")
				return
			}
			fmt.Printf("Error opening SSH config: %v\n", err)
			return
		}
		defer f.Close()

		cfg, err := ssh_config.Decode(f)
		if err != nil {
			fmt.Printf("Error parsing SSH config: %v\n", err)
			return
		}

		s, err := storage.NewStorage()
		if err != nil {
			fmt.Printf("Error initializing storage: %v\n", err)
			return
		}

		existingHosts, err := s.LoadHosts()
		if err != nil {
			fmt.Printf("Error loading hosts: %v\n", err)
			return
		}

		importedCount := 0
		skippedCount := 0

		for _, host := range cfg.Hosts {
			name := host.Patterns[0].String()
			if name == "*" || name == "" {
				continue
			}

			hostname := ""
			user := ""
			port := "22"
			identityFile := ""

			for _, node := range host.Nodes {
				line := strings.TrimSpace(node.String())
				parts := strings.Fields(line)
				if len(parts) < 2 {
					continue
				}
				key := strings.ToLower(parts[0])
				value := parts[1]

				switch key {
				case "hostname":
					hostname = value
				case "user":
					user = value
				case "port":
					port = value
				case "identityfile":
					identityFile = value
				}
			}

			if hostname == "" {
				continue
			}

			newHost := core.Host{
				Name:         name,
				Hostname:     hostname,
				User:         user,
				Port:         port,
				IdentityFile: identityFile,
			}

			exists := false
			for _, h := range existingHosts {
				if h.Name == name {
					exists = true
					break
				}
			}

			if exists {
				fmt.Printf("Host '%s' already exists. Skipping.\n", name)
				skippedCount++
				continue
			}

			fmt.Printf("Importing %s (%s)...\n", name, hostname)
			existingHosts = append(existingHosts, newHost)
			importedCount++
		}

		if importedCount > 0 {
			if err := s.SaveHosts(existingHosts); err != nil {
				fmt.Printf("Error saving hosts: %v\n", err)
				return
			}

			sshConfig := core.GenerateSSHConfig(existingHosts)
			if err := s.SaveSSHConfig(sshConfig); err != nil {
				fmt.Printf("Error saving SSH config: %v\n", err)
				return
			}
		}

		fmt.Printf("\nImport complete. Imported: %d, Skipped: %d\n", importedCount, skippedCount)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
