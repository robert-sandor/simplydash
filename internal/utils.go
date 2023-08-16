package internal

import (
	"sort"

	"github.com/sirupsen/logrus"
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

func SetupLogging(args Args) {
	level, err := logrus.ParseLevel(args.Log.Level)
	if err != nil {
		level = logrus.WarnLevel
	}
	logrus.SetLevel(level)

	if args.Log.Type == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:    true,
			QuoteEmptyFields: true,
		})
	}
}
