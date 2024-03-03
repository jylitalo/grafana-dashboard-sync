package api

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/jylitalo/grafana-dashboard-sync/config"
)

// [
//
//	{
//		"id":3,"uid":"c0be4e42-43fc-4f37-8e5f-f7d70f58284e","title":"App Debug",
//		"uri":"db/app-debug","url":"/d/c0be4e42-43fc-4f37-8e5f-f7d70f58284e/app-debug",
//		"slug":"","type":"dash-db","tags":[],"isStarred":false,"sortMeta":0
//	},
//	...
//
// ]
type Dashboard struct {
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
	grafana   config.Grafana
}

//	{
//		"meta":{...},
//		"dashboard":{
//			"annotations":{
//		 		"list":[{
//					"builtIn":1,"datasource":{"type":"grafana","uid":"-- Grafana --"},
//					"enable":true,"hide":true,"iconColor":"rgba(0, 211, 255, 1)",
//					"name":"Annotations \u0026 Alerts",
//					"type":"dashboard"
//				}]
//			},
//			"editable":true,"fiscalYearStartMonth":0,"graphTooltip":0,"id":3,"links":[],
//			"liveNow":false,
//			"panels":[
//				{
//					"datasource":{
//						"type":"prometheus", "uid":"f5976bf5-7c7a-4606-b2f5-311e2c9a02d9"
//					},
//					"fieldConfig":{
//						"defaults":{
//							"mappings":[],"thresholds":{
//								"mode":"percentage","steps":[
//									{"color":"green","value":null},
//									{"color":"orange","value":70},
//									{"color":"red","value":85}
//								]
//							},
//							"unit":"percentunit","unitScale":true
//						},
//						"overrides":[]
//					},
//					"gridPos":{"h":6,"w":4,"x":0,"y":0},
//					"id":3,
//					"options":{
//						"minVizHeight":75,"minVizWidth":75,"orientation":"auto",
//						"reduceOptions":{
//							"calcs":["lastNotNull"],"fields":"","values":false
//						},
//						"showThresholdLabels":false,"showThresholdMarkers":true,"sizing":"auto"
//					},
//					"pluginVersion":"10.3.1",
//					"targets":[
//						{
//							"datasource":{
//								"type":"prometheus","uid":"f5976bf5-7c7a-4606-b2f5-311e2c9a02d9"
//							},
//							"disableTextWrap":false,"editorMode":"code",
//							"expr":"max(container_memory_usage_bytes{pod=~\"$App.*\"})/max(container_spec_memory_limit_bytes{pod=~\"$App.*\"})",
//							"fullMetaSearch":false,"includeNullMetadata":true,"instant":false,
//							"legendFormat":"Used","range":true,"refId":"A","useBackend":false
//						}
//					],
//					"title":"Memory usage","type":"gauge"
//				},
//				...
//			],
//			"time":{"from":"now-24h","to":"now"},"timepicker":{},"timezone":"",
//			"title":"App Debug","uid":"c0be4e42-43fc-4f37-8e5f-f7d70f58284e",
//			"version":31,"weekStart":""
//		}
//	}
type DashDataSource struct {
	Type string `json:"type"`
	UID  string `json:"uid"`
}

type Target struct {
	DataSource          DashDataSource `json:"datasource"`
	DisableTextWrap     bool           `json:"disableTextWrap"`
	EditorMode          string         `json:"editorMode"`
	Expr                string         `json:"expr"`
	FullMetaSearch      bool           `json:"fullMetaSearch"`
	IncludeNullMetadata bool           `json:"includeNullMetadata"`
	Instant             bool           `json:"instant"`
	LegendFormat        string         `json:"legendFormat"`
	Range               bool           `json:"range"`
	RefId               string         `json:"refId"`
	UseBackend          bool           `json:"useBackend"`
}

type Panel struct {
	DataSource    interface{} `json:"datasource"` // string or DashDataSource
	FieldConfig   interface{} `json:"fieldConfig"`
	GridPos       interface{} `json:"gridPos"`
	Id            int         `json:"id"`
	Options       interface{} `json:"options"`
	PluginVersion string      `json:"pluginVersion"`
	Targets       []Target    `json:"targets"`
	Title         string      `json:"title"`
	Type          string      `json:"type"`
	Panels        []Panel     `json:"panels"`
}

type DashboardJSON struct {
	Meta      interface{} `json:"meta"`
	Dashboard struct {
		Annotations           interface{} `json:"annotations"`
		Editable              bool        `json:"editable"`
		FiscalYearStartsMonth int         `json:"fiscalYearStartMonth"`
		GraphTooltip          int         `json:"graphTooltip"`
		Id                    int         `json:"id"`
		Links                 interface{} `json:"links"`
		LiveNow               bool        `json:"liveNow"`
		Panels                []Panel     `json:"panels"`
		Time                  interface{} `json:"time"`
		TimePicker            interface{} `json:"timepicker"`
		TimeZone              string      `json:"timezone"`
		Title                 string      `json:"title"`
		UID                   string      `json:"uid"`
		Version               int         `json:"version"`
		WeekStart             string      `json:"weekStart"`
	} `json:"dashboard"`
}

func GetDashboards(grafana config.Grafana) ([]Dashboard, error) {
	body, err := getBody(grafana, "/api/search?query=&type=dash-db")
	if err != nil {
		return nil, err
	}
	sources := []Dashboard{}
	err = json.Unmarshal(body, &sources)
	if err != nil {
		return sources, err
	}
	for key := range sources {
		sources[key].grafana = grafana
	}
	return sources, nil
}

func (board *Dashboard) GetJSON() (DashboardJSON, error) {
	path := fmt.Sprintf("/api/dashboards/uid/%s", board.UID)
	source := DashboardJSON{}
	body, err := getBody(board.grafana, path)
	if err != nil {
		return source, err
	}
	err = json.Unmarshal(body, &source)
	if err != nil {
		slog.Error("GetDashboard", "body", body, "err", err)
	}
	return source, err
}

func (panel *Panel) Flatten() []Panel {
	flat := []Panel{*panel}
	if panel.Panels != nil {
		for _, item := range panel.Panels {
			flat = append(flat, item.Flatten()...)
		}
	}
	return flat
}
func (dashboard *DashboardJSON) Flatten() []Panel {
	flat := []Panel{}
	for _, item := range dashboard.Dashboard.Panels {
		flat = append(flat, item.Flatten()...)
	}
	return flat
}
