package cmd

import (
	"fmt"

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
				return fmt.Errorf("server (%s) not found from config", args[0])
			}
			dashdbs, err := api.GetDashDBs(server)
			if err != nil {
				return err
			}
			fmt.Println("Dashboards:")
			for _, item := range dashdbs {
				fmt.Printf("DashDB (%s): %v\n", item.Title, item)
				dboard, err := api.GetDashboard(server, item.UID)
				if err != nil {
					return err
				}
				fmt.Printf("Dashboard (%s): %#v\n", item.Title, dboard)
			}
			ds, err := api.GetDataSources(server)
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
