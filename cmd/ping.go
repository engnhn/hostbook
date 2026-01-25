package cmd

import (
	"fmt"
	"net"
	"os"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/engnhn/hostbook/core"
	"github.com/engnhn/hostbook/storage"
	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:               "ping",
	Short:             "Check health of all hosts",
	ValidArgsFunction: getHostCompletion,
	Run: func(cmd *cobra.Command, args []string) {
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

		fmt.Printf("Pinging %d hosts...\n\n", len(hosts))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tHOSTNAME\tSTATUS\tLATENCY")

		var wg sync.WaitGroup
		results := make(chan string, len(hosts))

		for _, h := range hosts {
			wg.Add(1)
			go func(host core.Host) {
				defer wg.Done()
				status, latency := checkHost(host)
				results <- fmt.Sprintf("%s\t%s\t%s\t%s", host.Name, host.Hostname, status, latency)
			}(h)
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		var rows []string
		for res := range results {
			rows = append(rows, res)
		}

		for _, row := range rows {
			fmt.Fprintln(w, row)
		}
		w.Flush()
	},
}

func checkHost(h core.Host) (string, string) {
	timeout := 2 * time.Second
	target := net.JoinHostPort(h.Hostname, h.Port)
	if h.Port == "" {
		target = net.JoinHostPort(h.Hostname, "22")
	}

	start := time.Now()
	conn, err := net.DialTimeout("tcp", target, timeout)
	duration := time.Since(start)

	if err != nil {
		return "OFFLINE", "-"
	}
	defer conn.Close()

	return "ONLINE", duration.Round(time.Millisecond).String()
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
