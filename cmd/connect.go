package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"github.com/engnhn/hostbook/core"
	"github.com/engnhn/hostbook/storage"
	"golang.org/x/term"

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

		password, _ := core.GetPassword(name)

		c := exec.Command(sshPath, "-F", sshConfigPath, name)

		if password == "" {

			env := os.Environ()
			argsSSH := []string{"ssh", "-F", sshConfigPath, name}
			if err := syscall.Exec(sshPath, argsSSH, env); err != nil {
				fmt.Printf("Error executing ssh: %v\n", err)
			}
			return
		}

		fmt.Println("Auto-connecting with stored password...")

		ptmx, err := pty.Start(c)
		if err != nil {
			fmt.Printf("Error starting PTY: %v\n", err)
			return
		}
		defer func() { _ = ptmx.Close() }()

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGWINCH)
		go func() {
			for range ch {
				if err := pty.InheritSize(os.Stdin, ptmx); err != nil {

				}
			}
		}()
		ch <- syscall.SIGWINCH

		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {

			fmt.Printf("Warning: failed to set raw mode: %v\n", err)
		} else {
			defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()
		}

		go func() {
			_, _ = io.Copy(ptmx, os.Stdin)
		}()

		buf := make([]byte, 1024)
		passwordSent := false

		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}

				if pathErr, ok := err.(*os.PathError); ok && pathErr.Err == syscall.EIO {
					break
				}
				break
			}

			os.Stdout.Write(buf[:n])

			if !passwordSent {
				output := string(buf[:n])

				if containsPasswordPrompt(output) {
					_, _ = ptmx.Write([]byte(password + "\n"))
					passwordSent = true
				}
			}
		}
	},
}

func containsPasswordPrompt(s string) bool {

	lower := ""
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			lower += string(r + 32)
		} else {
			lower += string(r)
		}
	}

	return (len(lower) > 5 && (contains(lower, "password:") || contains(lower, "passphrase")))
}

func contains(s, substr string) bool {

	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
