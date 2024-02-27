package api

import (
	"encoding/json"

	"github.com/jylitalo/grafana-dashboard-sync/config"
)

// [
//
//	{
//		"id":4,"uid":"bcd5eed3-7234-49dc-96b7-4acb031759e1","orgId":1,"name":"Jaeger",
//		"type":"jaeger","typeName":"Jaeger",
//		"typeLogoUrl":"/public/app/plugins/datasource/jaeger/img/jaeger_logo.svg",
//		"access":"proxy","url":"http://foo.com/","user":"","database":"","basicAuth":false,
//		"isDefault":false,
//		"jsonData":{
//			"tracesToLogsV2":{
//				"customQuery":false,
//				"datasourceUid":"ae732f43-47ae-4862-a3c2-5cbbd4347919"
//			}
//		},
//		"readOnly":false
//	},
//	...
//
// ]
type DataSource struct {
	Name string `json:"name"`
	Type string `json:"type"`
	UID  string `json:"uid"`
}

func GetDataSources(target config.Grafana) ([]DataSource, error) {
	body, err := getBody(target, "/api/datasources")
	if err != nil {
		return nil, err
	}
	sources := []DataSource{}
	err = json.Unmarshal(body, &sources)
	return sources, err
}
