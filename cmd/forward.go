package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/engnhn/hostbook/storage"

	"github.com/spf13/cobra"
)

var forwardCmd = &cobra.Command{
	Use:               "forward <host> <local_port:remote_port>",
	Short:             "Start port forwarding (SSH Tunnel)",
	ValidArgsFunction: getHostCompletion,
	Args:              cobra.ExactArgs(2),
	Example:           "  hostbook forward my-server 8080:80",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		ports := args[1]

		if !strings.Contains(ports, ":") {
			fmt.Println("Invalid port format. Use local:remote (e.g. 8080:80)")
			return
		}

		parts := strings.Split(ports, ":")
		localPort := parts[0]
		remotePort := parts[1]

		forwardRule := fmt.Sprintf("%s:127.0.0.1:%s", localPort, remotePort)

		s, err := storage.NewStorage()
		if err != nil {
			fmt.Printf("Error initializing storage: %v\n", err)
			return
		}

		sshConfigPath := s.GetSSHConfigPath()

		fmt.Printf("Forwarding 127.0.0.1:%s -> %s:127.0.0.1:%s...\n", localPort, name, remotePort)
		fmt.Println("Press Ctrl+C to stop.")

		sshCmd := exec.Command("ssh", "-F", sshConfigPath, "-L", forwardRule, "-N", name)
		sshCmd.Stdout = os.Stdout
		sshCmd.Stderr = os.Stderr

		if err := sshCmd.Run(); err != nil {
			fmt.Printf("SSH command exited with error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(forwardCmd)
}
