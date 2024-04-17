package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/jylitalo/grafana-dashboard-sync/api"
	"github.com/jylitalo/grafana-dashboard-sync/config"
)

type board struct {
	db   api.Dashboard
	json api.DashboardJSON
}

type panel struct {
	panel api.Panel
	index int
}

func varsToMap(vars []api.Variable) map[string]api.Variable {
	m := map[string]api.Variable{}
	for _, item := range vars {
		m[item.Name] = item
	}
	return m
}

func diffVars(one, two []api.Variable) [][]string {
	if len(one) != len(two) {
		return [][]string{
			{"variables mismatch", fmt.Sprintf("#%d", len(one)), fmt.Sprintf("#%d", len(two))},
		}
	}
	oneVars := varsToMap(one)
	twoVars := varsToMap(two)
	diff := [][]string{}
	oneMissing := []string{}
	twoMissing := []string{}
	for idx := range oneVars {
		if _, ok := twoVars[idx]; !ok {
			oneMissing = append(oneMissing, idx)
			delete(oneVars, idx)
			continue
		}
		if oneVars[idx].Definition != twoVars[idx].Definition {
			diff = append(diff, []string{idx + " definitions don't match", oneVars[idx].Definition, twoVars[idx].Definition})
		}
		if oneVars[idx].Regex != twoVars[idx].Regex {
			diff = append(diff, []string{idx + " regex don't match", oneVars[idx].Regex, twoVars[idx].Regex})
		}
		delete(oneVars, idx)
		delete(twoVars, idx)
	}
	for idx := range twoVars {
		twoMissing = append(twoMissing, idx)
	}
	if len(oneMissing) > 0 || len(twoMissing) > 0 {
		diff = append(diff, []string{"unique variables", strings.Join(oneMissing, ","), strings.Join(twoMissing, ",")})
	}
	return diff
}

func panelToMap(panels []api.Panel) map[string]panel {
	m := map[string]panel{}
	for idx, item := range panels {
		m[strings.TrimSpace(item.Title)] = panel{index: idx, panel: item}
	}
	return m
}

func dbToMap(dashboards []api.Dashboard) (map[string]board, error) {
	m := map[string]board{}
	for _, item := range dashboards {
		dashboard, err := item.GetJSON()
		if err != nil {
			return m, err
		}
		m[item.Title] = board{db: item, json: dashboard}
	}
	return m, nil
}

func truncLine(cell string) string {
	if len(cell) > 60 {
		return cell[:55] + "..."
	}
	return cell
}

func uniqTargetRefIds(targets []api.Target, start, end int) []string {
	ret := []string{}
	if len(targets) < end {
		return ret
	}
	for idx := start; idx < end; idx++ {
		ret = append(ret, targets[idx].RefId)
	}
	return ret
}

func diffTargets(one, two []api.Target) [][]string {
	oneLen := len(one)
	twoLen := len(two)
	diff := [][]string{}
	if len(one) != len(two) {
		diff = append(diff, []string{
			"targets mismatch", fmt.Sprintf("%d targets", oneLen), fmt.Sprintf("%d targets", twoLen),
		})
	}
	minItems := min(oneLen, twoLen)
	for idx := 0; idx < minItems; idx++ {
		if one[idx].RefId != two[idx].RefId {
			diff = append(diff, []string{"refId mismatch", one[idx].RefId, two[idx].RefId})
		}
		if one[idx].Expr != two[idx].Expr {
			diff = append(diff, []string{
				fmt.Sprintf("refId: %s expr mismatch", one[idx].RefId),
				truncLine(one[idx].Expr), truncLine(two[idx].Expr)})
		}
	}
	maxItems := max(oneLen, twoLen)
	onePlus := uniqTargetRefIds(one, minItems, maxItems)
	twoPlus := uniqTargetRefIds(two, minItems, maxItems)
	if len(onePlus) > 0 || len(twoPlus) > 0 {
		diff = append(diff, []string{"unique refIds", strings.Join(onePlus, ","), strings.Join(twoPlus, ",")})
	}
	return diff
}

func diffPanels(one, two []api.Panel) [][]string {
	diff := [][]string{}
	onePanels := panelToMap(one)
	twoPanels := panelToMap(two)
	uniqOne := []string{}
	uniqTwo := []string{}
	for key := range onePanels {
		panel1 := onePanels[key]
		panel2, ok := twoPanels[key]
		if !ok {
			uniqOne = append(uniqOne, key)
			delete(onePanels, key)
			continue
		}
		if panel1.index != panel2.index {
			diff = append(diff, []string{
				"Panel: " + panel1.panel.Title + "\nIndex mismatch",
				fmt.Sprintf("#%d", panel1.index),
				fmt.Sprintf("#%d", panel2.index),
			})
		}
		for _, row := range diffTargets(panel1.panel.Targets, panel2.panel.Targets) {
			diff = append(diff, []string{"Panel: " + panel1.panel.Title + "\n" + row[0], row[1], row[2]})
		}
		delete(onePanels, key)
		delete(twoPanels, key)
	}
	for key := range twoPanels {
		uniqTwo = append(uniqTwo, key)
		delete(twoPanels, key)
	}
	if len(uniqOne) > 0 || len(uniqTwo) > 0 {
		diff = append(diff, []string{"Unique panels", strings.Join(uniqOne, "\n"), strings.Join(uniqTwo, "\n")})
	}
	return diff
}

func diffDashboards(server1, server2 config.Grafana) error {
	dashdb1, err1 := api.GetDashboards(server1)
	dashdb2, err2 := api.GetDashboards(server2)
	uniqOne := []string{}
	uniqTwo := []string{}
	if err := errors.Join(err1, err2); err != nil {
		return err
	}
	identical := true
	if len(dashdb1) != len(dashdb2) {
		slog.Warn("Different number of dashboards", server1.Name, len(dashdb1), server2.Name, len(dashdb2))
		identical = false
	}
	dbMap1, err1 := dbToMap(dashdb1)
	dbMap2, err2 := dbToMap(dashdb2)
	if err := errors.Join(err1, err2); err != nil {
		return err
	}
	diff := [][]string{}
	for key, value1 := range dbMap1 {
		value2, ok := dbMap2[key]
		if !ok {
			uniqOne = append(uniqOne, key)
			delete(dbMap1, key)
			identical = false
			continue
		}
		oneDiff := diffVars(value1.json.Dashboard.Templating.List, value2.json.Dashboard.Templating.List)
		for _, item := range oneDiff {
			diff = append(diff, []string{value1.db.Title + "\n" + item[0], item[1], item[2]})
		}
		oneDiff = diffPanels(value1.json.Flatten(), value2.json.Flatten())
		for _, item := range oneDiff {
			diff = append(diff, []string{value1.db.Title + "\n" + item[0], item[1], item[2]})
		}
		delete(dbMap1, key)
		delete(dbMap2, key)
	}
	for key := range dbMap2 {
		uniqTwo = append(uniqTwo, key)
		identical = false
	}
	if identical && len(diff) == 0 {
		slog.Info("dashboards are identical", server1.Name, server2.Name)
		return nil
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"", server1.Name, server2.Name})
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.Append([]string{"Unique Dashboards", strings.Join(uniqOne, "\n"), strings.Join(uniqTwo, "\n")})
	table.AppendBulk(diff)
	table.Render()
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
		slog.Warn("Different number of data sources", server1.Name, len(ds1), server2.Name, len(ds2))
		identical = false
	}
	dsMap1 := dsToMap(ds1)
	dsMap2 := dsToMap(ds2)
	for key, value1 := range dsMap1 {
		value2, ok := dsMap2[key]
		if !ok {
			slog.Warn("only in "+server1.Name, "datasource", key)
			delete(dsMap1, key)
			identical = false
			continue
		}
		if value1.Type != value2.Type {
			slog.Warn("type mismatch", server1.Name, value1.Type, server2.Name, value2.Type)
			identical = false
		}
		delete(dsMap1, key)
		delete(dsMap2, key)
	}
	for key := range dsMap2 {
		slog.Warn("only in "+server2.Name, "datasource", key)
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
