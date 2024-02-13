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
			dashdbs, err := api.DashDBs(server)
			if err != nil {
				return err
			}
			fmt.Println("DashDBs:")
			for _, item := range dashdbs {
				fmt.Printf("%v\n", item)
			}
			fmt.Println("DashboardVersions:")
			for _, item := range dashdbs {
				dashversions, err := api.DashboardVersions(server, item.UID)
				if err != nil {
					return err
				}
				fmt.Printf("DashboardVersions (%s):\n", item.UID)
				for _, subItem := range dashversions {
					fmt.Printf("%v\n", subItem)
				}
			}
			ds, err := api.DataSources(server)
			if err != nil {
				return err
			}
			fmt.Println("Data sources:")
			for _, item := range ds {
				fmt.Printf("%v\n", item)
			}
			return nil
		},
	}
	return cmd
}
