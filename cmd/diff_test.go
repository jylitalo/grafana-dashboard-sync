package cmd

import (
	"testing"

	"github.com/jylitalo/grafana-dashboard-sync/api"
)

// {
// 	"datasource": {
//		"type": "prometheus",
//		"uid": "f5976bf5-7c7a-4606-b2f5-311e2c9a02d9"
//	},
//	"editorMode": "code",
//	"expr": "count(kube_deployment_created) - count(kube_deployment_status_replicas_available{instance=~\"$cluster.*\"})",
//	"interval": "",
//	"legendFormat": "Failed",
//	"range": true,
//	"refId": "A"
// },
// {
//	"datasource": {
//		"type": "prometheus",
//		"uid": "f5976bf5-7c7a-4606-b2f5-311e2c9a02d9"
//	},
//	"expr": "count(kube_deployment_created)",
//	"interval": "",
//	"legendFormat": "Running",
//	"refId": "B"
//}

var targetFailed api.Target = api.Target{
	DataSource:   api.DashDataSource{},
	EditorMode:   "code",
	Expr:         "count(kube_deployment_created) - count(kube_deployment_status_replicas_available{instance=~\"$cluster.*\"})",
	LegendFormat: "Failed",
	Range:        true,
	RefId:        "A",
}

var targetRunning api.Target = api.Target{
	DataSource:   api.DashDataSource{},
	Expr:         "count(kube_deployment_created)",
	LegendFormat: "Running",
	RefId:        "B",
}

func TestDiffTargets(t *testing.T) {
	one := []api.Target{targetFailed, targetRunning}
	two := []api.Target{targetFailed, targetRunning}
	diff := diffTargets(one, two)
	if len(diff) != 0 {
		t.Errorf("identical lists returned %#v", diff)
	}
	two = []api.Target{targetFailed}
	diff = diffTargets(one, two)
	if len(diff) != 2 {
		t.Errorf("different lists returned wrong number of lines")
	}
	if diff[1][1] != one[1].RefId || diff[1][2] != "" {
		t.Errorf("returned wrong refIds")
	}
	one = []api.Target{targetFailed}
	two = []api.Target{targetFailed, targetRunning}
	diff = diffTargets(one, two)
	if len(diff) != 2 {
		t.Errorf("different lists returned wrong number of lines")
	}
	if diff[1][1] != "" || diff[1][2] != two[1].RefId {
		t.Errorf("returned wrong refIds")
	}
}
