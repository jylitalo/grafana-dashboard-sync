package api

import (
	"encoding/json"
	"fmt"

	"github.com/jylitalo/grafana-dashboard-sync/config"
)

type DashDB struct {
	Id        int      `json:"id"`
	UID       string   `json:"uid"`
	Title     string   `json:"title"`
	URI       string   `json:"uri"`
	URL       string   `json:"url"`
	Slug      string   `json:"slug"`
	Type      string   `json:"type"`
	Tags      []string `json:"tags"`
	IsStarred bool     `json:"isStarred"`
	SortMeta  int      `json:"sortMeta"`
}

type DashboardVersion struct {
	Id          int    `json:"id"`
	DashboardId int    `json:"dashboardId"`
	UID         string `json:"uid"`
	Version     int    `json:"version"`
	Message     string `json:"message"`
	// ignore parentVersion, restoredFrom, created, createdBy
}

func DashDBs(target config.Grafana) ([]DashDB, error) {
	body, err := getBody(target, "/api/search?query=&type=dash-db")
	if err != nil {
		return nil, err
	}
	sources := []DashDB{}
	err = json.Unmarshal(body, &sources)
	return sources, err
}

func DashboardVersions(target config.Grafana, uid string) ([]DashboardVersion, error) {
	path := fmt.Sprintf("/api/dashboards/uid/%s/versions", uid)
	body, err := getBody(target, path)
	if err != nil {
		return nil, err
	}
	sources := []DashboardVersion{}
	err = json.Unmarshal(body, &sources)
	fmt.Printf("%v\n", string(body))
	return sources, err
}
