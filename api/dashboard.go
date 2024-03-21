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

// AnnotationsPermissions is part of DashboardJSON.Meta
type AnnotationsPermissions struct {
	CanAdd    bool `json:"canAdd"`
	CanEdit   bool `json:"canEdit"`
	CanDelete bool `json:"canDelete"`
}

// DashDataSource is part of Target and Variable
type DashDataSource struct {
	Type string `json:"type"`
	UID  string `json:"uid"`
}

// Variable is dashboard variable, part of DashboardJSON
type Variable struct {
	Current        interface{}   `json:"current"`
	DataSource     interface{}   `json:"datasource"` // can be string or DashDataSoure
	Definition     string        `json:"definition"`
	Hide           int           `json:"hide"`
	IncludeAll     bool          `json:"includeAll"`
	Label          string        `json:"label,omitempty"`
	Multi          bool          `json:"multi"`
	Name           string        `json:"name"`
	Options        []interface{} `json:"options"`
	Query          interface{}   `json:"query"` // can be string or struct with qryType (int), query (string) and refId (string)
	Refresh        int           `json:"refresh"`
	Regex          string        `json:"regex"`
	SkipUrlSync    bool          `json:"skipUrlSync"`
	Sort           int           `json:"sort"`
	TagValuesQuery interface{}   `json:"tagValuesQuery,omitempty"`
	TagsQuery      interface{}   `json:"tagsQuery,omitempty"`
	Type           string        `json:"type"`
	UseTags        interface{}   `json:"useTags,omitempty"`
}

// DashboardJSON is JSON presentation of actual dashboard.
// You can find examples from
// - test-data/instances_closed.json
// - test-case/instances_open.json
type DashboardJSON struct {
	Meta struct {
		Type                   string `json:"type"`
		CanSave                bool   `json:"canSave"`
		CanEdit                bool   `json:"canEdit"`
		CanAdmin               bool   `json:"canAdmin"`
		CanStar                bool   `json:"canStar"`
		CanDelete              bool   `json:"canDelete"`
		Slug                   string `json:"slug"`
		URL                    string `json:"url"`
		Expires                string `json:"expires"`
		Created                string `json:"created"`
		Updated                string `json:"updated"`
		UpdatedBy              string `json:"updatedBy"`
		CreatedBy              string `json:"createdBy"`
		Version                int    `json:"version"`
		HasACL                 bool   `json:"hasAcl"`
		IsFolder               bool   `json:"isFolder"`
		FolderId               int    `json:"folderId"`
		FolderUID              string `json:"folderUid"`
		FolderTitle            string `json:"folderTitle"`
		FolderURL              string `json:"folderUrl"`
		Provisioned            bool   `json:"provisioned"`
		ProvisionedExternalId  string `json:"provisionedExternalId"`
		AnnotationsPermissions struct {
			Dashboard    AnnotationsPermissions `json:"dashboard"`
			Organization AnnotationsPermissions `json:"organization"`
		} `json:"annotationsPermissions"`
	} `json:"meta"`
	Dashboard struct {
		Annotations           interface{} `json:"annotations"`
		Description           string      `json:"description,omitempty"`
		Editable              bool        `json:"editable"`
		FiscalYearStartsMonth int         `json:"fiscalYearStartMonth"`
		GnetId                int         `json:"gnetId,omitempty"`
		GraphTooltip          int         `json:"graphTooltip"`
		Id                    int         `json:"id"`
		Links                 interface{} `json:"links"`
		LiveNow               bool        `json:"liveNow"`
		Panels                []Panel     `json:"panels"`
		Refresh               interface{} `json:"refresh"`
		SchemaVersion         int         `json:"schemaVersion"`
		Tags                  interface{} `json:"tags"`
		Templating            struct {
			List []Variable `json:"list"`
		} `json:"templating"`
		Time       interface{} `json:"time"`
		TimePicker interface{} `json:"timepicker"`
		TimeZone   string      `json:"timezone"`
		Title      string      `json:"title"`
		UID        string      `json:"uid"`
		Version    int         `json:"version"`
		WeekStart  string      `json:"weekStart"`
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
	body, err := getBody(board.grafana, path)
	if err != nil {
		return DashboardJSON{}, err
	}
	// fmt.Println("---")
	// fmt.Println(string(body))
	// fmt.Println("---")
	return parseDashboardJSON(body)
}

func parseDashboardJSON(body []byte) (DashboardJSON, error) {
	source := DashboardJSON{}
	err := json.Unmarshal(body, &source)
	if err != nil {
		slog.Error("parseDashboardJSON", "body", string(body), "err", err)
	}
	return source, err
}

func (dashboard *DashboardJSON) Flatten() []Panel {
	flat := []Panel{}
	for _, item := range dashboard.Dashboard.Panels {
		flat = append(flat, item.Flatten()...)
	}
	return flat
}
