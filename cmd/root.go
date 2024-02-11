package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context) error {
	rootCmd := &cobra.Command{
		Use:   "grafana-dashboard-sync [dashboard-name]",
		Short: "Sync dashboard with two grafana instances",
	}
	rootCmd.AddCommand(diffCmd(), listCmd())
	return rootCmd.ExecuteContext(ctx)
}
