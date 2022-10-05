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

	cfg, err := config.LoadConfig(args.ConfigPath.Get(), utils.FileReader)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	svc := internal.NewService(internal.NewFileWatcher(cfg.FileProviders))
	svc.Init()

	iconCache := internal.NewIconCache(args.IconCachePath.Get())

	r := api.NewRouter(cfg, svc, iconCache)
	r.Run()
}
