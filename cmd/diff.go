package cmd

import (
	"errors"
	"log/slog"

	"github.com/jylitalo/grafana-dashboard-sync/api"
	"github.com/jylitalo/grafana-dashboard-sync/config"
	"github.com/spf13/cobra"
)

func dsToMap(ds []api.DataSource) map[string]api.DataSource {
	m := map[string]api.DataSource{}
	for _, item := range ds {
		m[item.Name] = item
	}
	return m
}

func diffCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff [server1 server2]",
		Short: "diff two grafanas configuration",
		Long:  "Fetch configuration from two servers and create diff",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			cfg, err := config.Get(ctx)
			if err != nil {
				return err
			}
			server1, ok := cfg[args[0]]
			if !ok {
				slog.Error("Server not found from config", "server1", server1)
			}
			server2, ok := cfg[args[1]]
			if !ok {
				slog.Error("Server not found from config", "server2", server2)
			}
			ds1, err1 := api.DataSources(server1)
			ds2, err2 := api.DataSources(server2)
			if err := errors.Join(err1, err2); err != nil {
				return err
			}
			identical := true
			if len(ds1) != len(ds2) {
				slog.Info("Different number of data sources", server1.Name, len(ds1), server2.Name, len(ds2))
				identical = false
			}
			dsMap1 := dsToMap(ds1)
			dsMap2 := dsToMap(ds2)
			for key, value1 := range dsMap1 {
				value2, ok := dsMap2[key]
				if !ok {
					slog.Info("only in "+server1.Name, "datasource", key)
					delete(dsMap1, key)
					identical = false
					continue
				}
				if value1.Type != value2.Type {
					slog.Info("type mismatch", server1.Name, value1.Type, server2.Name, value2.Type)
					identical = false
				}
				delete(dsMap1, key)
				delete(dsMap2, key)
			}
			for key := range dsMap2 {
				slog.Info("only in "+server2.Name, "datasource", key)
				identical = false
			}
			if identical {
				slog.Info("data sources are identical", server1.Name, server2.Name)
			}
			return nil
		},
	}
	return cmd
}
