package config

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/jylitalo/grafana-dashboard-sync/pkg/logging"
)

type Grafana struct {
	Name   string
	URL    string
	Bearer string
}

type Options struct {
	Path string
	Name string
}

type Config map[string]Grafana
type ctxType string

const ctxKey ctxType = "grafana-dashboard-sync"

func Get(ctx context.Context) (Config, error) {
	if ctx == nil {
		return Config{}, errors.New("context is nil")
	}
	value := ctx.Value(ctxKey)
	if value == nil {
		return Config{}, errors.New("context doesn't have Config")
	}
	return value.(Config), nil
}

func Read(optFns ...func(*Options)) (context.Context, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error in UserHomeDir: %w", err)
	}
	opts := &Options{
		Path: home,
		Name: ".grafana-dashboard-sync",
	}
	for _, optFn := range optFns {
		optFn(opts)
	}
	vip := viper.GetViper()
	vip.AddConfigPath(opts.Path)
	vip.SetConfigName(opts.Name)
	if err = vip.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error in ReadInConfig: %w", err)
	}

	booleans := map[string]bool{"color": true, "debug": false}
	value := Config{}
	for _, keyName := range vip.AllKeys() {
		if _, found := booleans[keyName]; found {
			booleans[keyName] = vip.GetBool(keyName)
			continue
		}
		fields := strings.SplitN(keyName, ".", 2)
		server := fields[0]
		subKey := strings.ToLower(fields[1])
		val, ok := value[server]
		if !ok {
			val = Grafana{Name: server}
		}
		s := vip.GetString(keyName)
		switch {
		case subKey == "bearer":
			val.Bearer = s
		case subKey == "url":
			val.URL = s
		default:
			return nil, fmt.Errorf("unknown key (%s) in config file", keyName)
		}
		value[fields[0]] = val
	}
	slog.SetDefault(logging.SetupSlog(booleans["debug"], booleans["color"]))
	slog.Debug("config", "value", value)
	ctx := context.WithValue(context.Background(), ctxKey, value)
	return ctx, nil
}
