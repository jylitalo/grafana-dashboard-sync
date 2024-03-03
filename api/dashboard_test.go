package api

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

var flatPanelList = []Panel{
	{Id: 1},
	{Id: 2},
	{Id: 3},
	{Id: 4},
	{Id: 5},
	{Id: 6},
	{Id: 7},
	{Id: 8},
	{Id: 9},
}

//go:embed test-data/instances_closed.json
var instancesClosed string

//go:embed test-data/instances_open.json
var instancesOpen string

var treePanelList = []Panel{
	{Id: 1, Panels: []Panel{{Id: 2}}},
	{Id: 3, Panels: []Panel{{Id: 4}, {Id: 5}}},
	{Id: 6, Panels: []Panel{{Id: 7, Panels: []Panel{{Id: 8}}}, {Id: 9}}},
}

// TestFlattenPanel should return panels 0..10
func TestFlattenPanel(t *testing.T) {
	values := [][]Panel{flatPanelList, treePanelList}
	for idx, value := range values {
		t.Run(fmt.Sprintf("TestFlattenPanel.%d", idx), func(t *testing.T) {
			p := Panel{Id: 0, Panels: value}
			for idx, item := range p.Flatten() {
				if idx != item.Id {
					t.Errorf("Mismatch on %v\n", item)
				}
			}
		})
	}
}

// TestFlattenDashboard should return panels 0..10
func TestFlattenDashboard(t *testing.T) {
	values := [][]Panel{flatPanelList, treePanelList}
	for idx, value := range values {
		t.Run(fmt.Sprintf("TestFlattenDashboard.%d", idx), func(t *testing.T) {
			p := DashboardJSON{}
			p.Dashboard.Panels = append([]Panel{{Id: 0}}, value...)
			for idx, item := range p.Flatten() {
				if idx != item.Id {
					t.Errorf("Mismatch on %v\n", item)
				}
			}
		})
	}
}

func TestRealFlattenDashboard(t *testing.T) {
	closedRows, errClosed := parseDashboardJSON([]byte(instancesClosed))
	openRows, errOpen := parseDashboardJSON([]byte(instancesOpen))
	if err := errors.Join(errOpen, errClosed); err != nil {
		t.Errorf("errs found errOpen=%v, errClosed=%v", errOpen, errClosed)
	}
	closedPanels := closedRows.Flatten()
	openPanels := openRows.Flatten()
	if len(closedPanels) != len(openPanels) {
		t.Errorf("panel list length don't match (%d vs. %d)", len(closedPanels), len(openPanels))
	}
}

func TestParseDashboardJSON(t *testing.T) {
	items := []string{instancesClosed, instancesOpen}
	for idx, item := range items {
		t.Run(fmt.Sprintf("TestParseDashboardJSON.%d", idx), func(t *testing.T) {
			inStruct, err := parseDashboardJSON([]byte(item))
			if err != nil {
				t.Errorf("parseDashboardJSON failed due to %v", err)
			}
			inJSON, err := json.Marshal(inStruct)
			if err != nil {
				t.Errorf("Marshal failed due to %v", err)
			}
			if string(inJSON) != item {
				t.Errorf("Not match ... %s", inJSON)
			}
		})
	}
}
