package cmd

import (
	"fmt"
	"net"
	"os"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/engnhn/hostbook/storage"
	"github.com/spf13/cobra"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all hosts",
	Run: func(cmd *cobra.Command, args []string) {
		filterTag, _ := cmd.Flags().GetString("tag")

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

		type hostStatus struct {
			name   string
			online bool
		}
		statusChan := make(chan hostStatus, len(hosts))
		var wg sync.WaitGroup

		for _, h := range hosts {
			wg.Add(1)
			go func(hostName, hostname, port string) {
				defer wg.Done()
				target := fmt.Sprintf("%s:%s", hostname, port)
				if port == "" {
					target = fmt.Sprintf("%s:22", hostname)
				}
				conn, err := net.DialTimeout("tcp", target, 500*time.Millisecond)
				online := false
				if err == nil {
					conn.Close()
					online = true
				}
				statusChan <- hostStatus{name: hostName, online: online}
			}(h.Name, h.Hostname, h.Port)
		}

		go func() {
			wg.Wait()
			close(statusChan)
		}()

		statusMap := make(map[string]bool)
		for status := range statusChan {
			statusMap[status.name] = status.online
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "STATUS\tNAME\tHOSTNAME\tUSER\tPORT\tTAGS")
		for _, h := range hosts {
			if filterTag != "" {
				found := false
				for _, t := range h.Tags {
					if t == filterTag {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			tags := ""
			if len(h.Tags) > 0 {
				tags = fmt.Sprintf("%v", h.Tags)
			}

			statusDot := fmt.Sprintf("%s●%s", colorRed, colorReset)
			if statusMap[h.Name] {
				statusDot = fmt.Sprintf("%s●%s", colorGreen, colorReset)
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", statusDot, h.Name, h.Hostname, h.User, h.Port, tags)
		}
		w.Flush()
	},
}

func init() {
	listCmd.Flags().String("tag", "", "Filter by tag")
	rootCmd.AddCommand(listCmd)
}
