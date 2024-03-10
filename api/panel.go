package api

import (
	"encoding/json"
)

// Panel is part of DashboardJSON.Dashboard.Panels
type Panel struct {
	DataSource      interface{} `json:"datasource,omitempty"` // string or DashDataSource
	Description     interface{} `json:"description,omitempty"`
	FieldConfig     interface{} `json:"fieldConfig,omitempty"`
	Collapsed       interface{} `json:"collapsed,omitempty"`
	GridPos         interface{} `json:"gridPos"`
	Id              int         `json:"id"`
	Links           interface{} `json:"links,omitempty"`
	MaxDataPoints   interface{} `json:"maxDataPoints,omitempty"`
	Options         interface{} `json:"options,omitempty"`
	PluginVersion   string      `json:"pluginVersion,omitempty"`
	Targets         []Target    `json:"targets,omitempty"`
	Panels          []Panel     `json:"panels"`
	Title           string      `json:"title"`
	Transformations interface{} `json:"transformations,omitempty"`
	Type            string      `json:"type"`
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

// ToMarshalJSON deals with the fact that `row` should have
// `"panels": []` if panels is empty or nil.
// Otherwise panels should be omitted from output.
func (panel *Panel) MarshalJSON() ([]byte, error) {
	if panel.Type == "row" {
		type Alias Panel
		if panel.Panels == nil {
			panel.Panels = []Panel{}
		}
		return json.Marshal(&struct {
			*Alias
		}{
			Alias: (*Alias)(panel),
		})
	}
	// cpanel should be copy of Panel without Panels field
	cpanel := struct {
		DataSource      interface{} `json:"datasource,omitempty"` // string or DashDataSource
		Description     interface{} `json:"description,omitempty"`
		FieldConfig     interface{} `json:"fieldConfig,omitempty"`
		Collapsed       interface{} `json:"collapsed,omitempty"`
		GridPos         interface{} `json:"gridPos"`
		Id              int         `json:"id"`
		Links           interface{} `json:"links,omitempty"`
		MaxDataPoints   interface{} `json:"maxDataPoints,omitempty"`
		Options         interface{} `json:"options,omitempty"`
		PluginVersion   string      `json:"pluginVersion,omitempty"`
		Targets         []Target    `json:"targets,omitempty"`
		Title           string      `json:"title,omitempty"`
		Transformations interface{} `json:"transformations,omitempty"`
		Type            string      `json:"type"`
	}{
		DataSource:      panel.DataSource,
		Description:     panel.Description,
		FieldConfig:     panel.FieldConfig,
		Collapsed:       panel.Collapsed,
		GridPos:         panel.GridPos,
		Id:              panel.Id,
		Links:           panel.Links,
		MaxDataPoints:   panel.MaxDataPoints,
		Options:         panel.Options,
		PluginVersion:   panel.PluginVersion,
		Targets:         panel.Targets,
		Title:           panel.Title,
		Transformations: panel.Transformations,
		Type:            panel.Type,
	}
	return json.Marshal(cpanel)
}
