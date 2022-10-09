package main

import (
	"simplydash/internal"
)

func main() {
	args := internal.NewArgs()

	internal.Log = internal.NewLoggerStr(args.LogLevel.Get())
	internal.Log.Debug.Printf("Starting with args = %s", args)

	cfg := internal.NewConfig(args.ConfigPath.Get(), internal.FileReader, internal.FileWriter)

	wsNot := internal.NewWebsocketNotifier()
	fileSvc := internal.NewFileService(internal.NewFileWatcher(), args, cfg, wsNot, internal.FileReader)

	svc := internal.NewService(fileSvc)
	err := svc.Init()
	if err != nil {
		internal.Log.Error.Fatalf("Failed to initialize service err = %+v", err)
	}

	iconCache := internal.NewIconCache(args.IconCachePath.Get())

	r := internal.NewRouter(cfg, svc, iconCache, wsNot)

	internal.Log.Info.Printf("Starting server on port %s using config file %s and icon cache path %s",
		args.Port.Get(), args.ConfigPath.Get(), args.IconCachePath.Get())
	r.Run(args.Port.Get())
}
