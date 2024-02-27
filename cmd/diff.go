package cmd

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/jylitalo/grafana-dashboard-sync/api"
	"github.com/jylitalo/grafana-dashboard-sync/config"
	"github.com/spf13/cobra"
)

type board struct {
	dashDB api.DashDB
	board  api.Dashboard
}

func panelTitles(panels []api.Panel) string {
	titles := []string{}
	for _, panel := range panels {
		titles = append(titles, panel.Title)
	}
	return strings.Join(titles, ",")
}

func dbToMap(target config.Grafana, dashDBs []api.DashDB) (map[string]board, error) {
	m := map[string]board{}
	for _, item := range dashDBs {
		dashboard, err := api.GetDashboard(target, item.UID)
		if err != nil {
			return m, err
		}
		m[item.Title] = board{dashDB: item, board: dashboard}
	}
	return m, nil
}

func diffDashboards(server1, server2 config.Grafana) error {
	dashdb1, err1 := api.GetDashDBs(server1)
	dashdb2, err2 := api.GetDashDBs(server2)
	if err := errors.Join(err1, err2); err != nil {
		return err
	}
	identical := true
	if len(dashdb1) != len(dashdb2) {
		slog.Info("Different number of dashboards", server1.Name, len(dashdb1), server2.Name, len(dashdb2))
		identical = false
	}
	dbMap1, err1 := dbToMap(server1, dashdb1)
	dbMap2, err2 := dbToMap(server2, dashdb2)
	if err := errors.Join(err1, err2); err != nil {
		return err
	}
	for key, value1 := range dbMap1 {
		value2, ok := dbMap2[key]
		if !ok {
			slog.Info("only in "+server1.Name, "dashboard", key)
			delete(dbMap1, key)
			identical = false
			continue
		}
		if len(value1.board.Dashboard.Panels) != len(value2.board.Dashboard.Panels) {
			slog.Info(
				"different number of panels", "dashboard", key,
				"server1_len", len(value1.board.Dashboard.Panels),
				"server2_len", len(value2.board.Dashboard.Panels),
				"server1_panels", panelTitles(value1.board.Dashboard.Panels),
				"server2_panels", panelTitles(value2.board.Dashboard.Panels),
			)
		}
		delete(dbMap1, key)
		delete(dbMap2, key)
	}
	for key := range dbMap2 {
		slog.Info("only in "+server2.Name, "dashboard", key)
		identical = false
	}
	if identical {
		slog.Info("dashboards are identical", server1.Name, server2.Name)
	}
	return nil
}

func dsToMap(ds []api.DataSource) map[string]api.DataSource {
	m := map[string]api.DataSource{}
	for _, item := range ds {
		m[item.Name] = item
	}
	return m
}

func diffDatasources(server1, server2 config.Grafana) error {
	ds1, err1 := api.GetDataSources(server1)
	ds2, err2 := api.GetDataSources(server2)
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
			err = diffDatasources(server1, server2)
			if err != nil {
				return err
			}
			err = diffDashboards(server1, server2)
			if err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
