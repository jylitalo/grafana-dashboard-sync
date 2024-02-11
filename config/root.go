package config

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Grafana struct {
	Name   string
	URL    string
	Bearer string
}
type Config map[string]Grafana

const ctxKey string = "grafana-dashboard-sync"

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

func Read() (context.Context, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("Error in UserHomeDir: %w", err)
	}
	viper.AddConfigPath(home)
	viper.SetConfigName(".grafana-dashboard-sync")
	if err = viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Error in ReadInConfig: %w", err)
	}

	value := Config{}
	for _, keyName := range viper.AllKeys() {
		fields := strings.SplitN(keyName, ".", 2)
		server := fields[0]
		subKey := strings.ToLower(fields[1])
		v, ok := value[server]
		if !ok {
			v = Grafana{Name: server}
		}
		s := viper.GetString(keyName)
		switch {
		case subKey == "bearer":
			v.Bearer = s
		case subKey == "url":
			v.URL = s
		}
		value[fields[0]] = v
	}
	slog.Debug("config", "value", value)
	ctx := context.WithValue(context.Background(), ctxKey, value)
	return ctx, nil
}
