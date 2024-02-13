package api

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/jylitalo/grafana-dashboard-sync/config"
)

func getBody(target config.Grafana, path string) ([]byte, error) {
	bearer := "Bearer " + target.Bearer
	url := target.URL + path
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
	slog.Debug("api.getBody", "body", body, "err", err)
	return body, err
}
