package api

import (
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
