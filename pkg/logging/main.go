package logging

import (
	"log/slog"
	"os"

	"github.com/jylitalo/tint"
	"github.com/mattn/go-isatty"
)

func SetupSlog(debug bool, color bool) *slog.Logger {
	logLevel := map[bool]slog.Level{
		true:  slog.LevelDebug,
		false: slog.LevelInfo,
	}[debug]
	w := os.Stderr
	return slog.New(tint.NewHandler(w, &tint.Options{
		Level:   logLevel,
		NoColor: !isatty.IsTerminal(w.Fd()) || !color,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.Attr{}
			}
			return a
		},
	}))
}
