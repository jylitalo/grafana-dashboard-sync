package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/jylitalo/grafana-dashboard-sync/config"
)

type DataSource struct {
	Name string `json:"name"`
	Type string `json:"type"`
	UID  string `json:"uid"`
}

func DataSources(target config.Grafana) ([]DataSource, error) {
	bearer := "Bearer " + target.Bearer
	url := target.URL + "/api/datasources"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	sources := []DataSource{}
	err = json.Unmarshal(body, &sources)
	return sources, err
}
