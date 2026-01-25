package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/engnhn/hostbook/storage"
	"github.com/spf13/cobra"
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

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tHOSTNAME\tUSER\tPORT\tTAGS")
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
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", h.Name, h.Hostname, h.User, h.Port, tags)
		}
		w.Flush()
	},
}

func init() {
	listCmd.Flags().String("tag", "", "Filter by tag")
	rootCmd.AddCommand(listCmd)
}
