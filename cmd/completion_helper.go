package cmd

import (
	"strings"

	"github.com/engnhn/hostbook/storage"
	"github.com/spf13/cobra"
)

func getHostCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	s, err := storage.NewStorage()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	hosts, err := s.LoadHosts()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string
	for _, h := range hosts {
		if strings.HasPrefix(h.Name, toComplete) {
			names = append(names, h.Name)
		}
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
