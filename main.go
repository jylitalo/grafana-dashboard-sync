package main

import (
	"log"

	"github.com/jylitalo/grafana-dashboard-sync/cmd"
	"github.com/jylitalo/grafana-dashboard-sync/config"
)

func main() {
	ctx, err := config.Read()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	if err := cmd.Execute(ctx); err != nil {
		log.Fatalf("Error: %s", err)
	}
}
