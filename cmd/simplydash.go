package main

import (
	"log"
	"simplydash/internal"
	"simplydash/internal/api"
	"simplydash/internal/config"
	"simplydash/internal/utils"
)

func main() {
	args := internal.NewArgs()

	cfg := config.NewConfig(args.ConfigPath.Get(), utils.FileReader, utils.FileWriter)

	svc := internal.NewService(internal.NewFileWatcher(args.ConfigPath.Get(), cfg, cfg.FileProviders))
	svc.Init()

	iconCache := internal.NewIconCache(args.IconCachePath.Get())

	r := api.NewRouter(cfg, svc, iconCache)

	log.Printf("Starting server on port %s using config file %s and icon cache path %s",
		args.Port.Get(), args.ConfigPath.Get(), args.IconCachePath.Get())
	r.Run(args.Port.Get())
}
