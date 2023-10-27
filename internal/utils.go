package internal

import (
	"context"
	"log/slog"
	"os"
	"sort"
)

func insertOrdered(apps []App, app App) []App {
	i := sort.Search(len(apps), func(i int) bool { return apps[i].Name > app.Name })

	if i == len(apps) {
		return append(apps, app)
	}

	apps = append(apps[:i+1], apps[i:]...)
	apps[i] = app
	return apps
}

func SetupSlog(args Args) {
	level := slog.LevelInfo
	err := level.UnmarshalText([]byte(args.Log.Level))
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "invalid log level", slog.String("logLevel", args.Log.Level))
	}

	var handler slog.Handler
	if args.Log.Type == "json" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: level <= slog.LevelDebug,
			Level:     level,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: level <= slog.LevelDebug,
			Level:     level,
		})
	}

	slog.SetDefault(slog.New(handler))
}
