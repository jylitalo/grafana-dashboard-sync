package cmd

import (
	"fmt"
	"log/slog"

	"github.com/jylitalo/grafana-dashboard-sync/api"
	"github.com/jylitalo/grafana-dashboard-sync/config"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [server]",
		Short: "list configuration",
		Long:  "Fetch configuration from server and show it on screen",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			cfg, err := config.Get(ctx)
			if err != nil {
				return err
			}
			server, ok := cfg[args[0]]
			if !ok {
				slog.Error("Server not found from config", "server", server)
			}
			ds, err := api.DataSources(server)
			if err != nil {
				return err
			}
			for _, item := range ds {
				fmt.Printf("%v\n", item)
			}
			return nil
		},
	}
	return cmd
}
