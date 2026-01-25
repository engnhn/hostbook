package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/engnhn/hostbook/storage"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export hosts to JSON or YAML",
	Long:  `Export all your hosts to a standard format (JSON or YAML) for backup or migration.`,
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")

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

		var output []byte
		var marchalErr error

		if format == "yaml" {
			output, marchalErr = yaml.Marshal(hosts)
		} else {
			output, marchalErr = json.MarshalIndent(hosts, "", "  ")
		}

		if marchalErr != nil {
			fmt.Printf("Error exporting data: %v\n", marchalErr)
			return
		}

		fmt.Println(string(output))
	},
}

func init() {
	exportCmd.Flags().StringP("format", "f", "json", "Output format (json or yaml)")
	rootCmd.AddCommand(exportCmd)
}
