package config

import (
	"log/slog"
	"testing"

	"github.com/jylitalo/grafana-dashboard-sync/pkg/logging"
)

func Test2GrafanasWithDebug(t *testing.T) {
	optFn := func(opts *Options) {
		opts.Path = "test-data"
		opts.Name = "test-1"
	}
	slog.SetDefault(logging.SetupSlog(false, false))
	ctx, err := Read(optFn)
	if err != nil {
		t.Errorf("Read failed due to %v", err)
	}
	data, err := Get(ctx)
	if err != nil {
		t.Errorf("Get failed due to %v", err)
	}
	if len(data) != 2 {
		t.Errorf("Wrong number of Grafanas found (%d vs. 2)", len(data))
	}
	if !slog.Default().Enabled(ctx, slog.LevelDebug) {
		t.Errorf("Debug logging has not been enabled")
	}
}
