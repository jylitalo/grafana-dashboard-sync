package api

import (
	"encoding/json"

	"github.com/jylitalo/grafana-dashboard-sync/config"
)

type DataSource struct {
	Name string `json:"name"`
	Type string `json:"type"`
	UID  string `json:"uid"`
}

func DataSources(target config.Grafana) ([]DataSource, error) {
	body, err := getBody(target, "/api/datasources")
	if err != nil {
		return nil, err
	}
	sources := []DataSource{}
	err = json.Unmarshal(body, &sources)
	return sources, err
}
