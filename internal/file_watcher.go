package internal

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"simplydash/internal/config"
	"simplydash/internal/models"
	"simplydash/internal/providers"
	"simplydash/internal/utils"
)

type FileWatcher struct {
	cfgPath       string
	cfg           *config.Config
	fileProviders []*providers.FileProvider
	watcher       *fsnotify.Watcher
}

func NewFileWatcher(cfgPath string, cfg *config.Config, providerConfigs []config.FileProviderConfig) *FileWatcher {
	var fpArr []*providers.FileProvider
	for _, fpc := range providerConfigs {
		p, err := providers.NewFileProvider(fpc)
		if err != nil {
			log.Printf("Failed to create file provider for file = %s", fpc.Path)
			continue
		}
		fpArr = append(fpArr, p)
	}

	return &FileWatcher{cfgPath: cfgPath, cfg: cfg, fileProviders: fpArr, watcher: newWatcher()}
}

func (fw *FileWatcher) Load() {
	for _, fp := range fw.fileProviders {
		err := fp.Load()
		if err != nil {
			log.Printf("Failed to load file provider for file = %s", fp.Path)
		}
	}
}

func (fw *FileWatcher) Get() []models.Category {
	createdCategories := make(map[string]int, 0)
	categories := make([]models.Category, 0)

	for _, fp := range fw.fileProviders {
		for _, pCat := range fp.Get() {
			if index, ok := createdCategories[pCat.Name]; ok {
				categories[index].Items = append(categories[index].Items, pCat.Items...)
				continue
			}

			createdCategories[pCat.Name] = len(categories)
			categories = append(categories, pCat)
		}
	}

	return categories
}

func (fw *FileWatcher) Watch(updateChannels *[]chan string) {
	if fw.watcher == nil {
		fw.watcher = newWatcher()
	}

	err := fw.watcher.Add(fw.cfgPath)
	if err != nil {
		log.Printf("Failed to add config file path %s to watcher")
	}

	var fpToWatch []*providers.FileProvider
	for _, fp := range fw.fileProviders {
		if fp.Watch {
			fpToWatch = append(fpToWatch, fp)
		}
	}

	for _, fp := range fpToWatch {
		err := fw.watcher.Add(fp.Path)
		if err != nil {
			log.Printf("Failed to add file = %s to watcher err = %+v", fp.Path, err)
		}
	}

	go fw.watch(fpToWatch, updateChannels)
}

func (fw *FileWatcher) watch(
	fileProviders []*providers.FileProvider,
	updateChannels *[]chan string,
) {
	for {
		select {
		case err, ok := <-fw.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Received error = %+v", err)
		case e, ok := <-fw.watcher.Events:
			if !ok {
				return
			}

			if e.Op == fsnotify.Write {
				if fw.cfgPath == e.Name {
					fw.updateCfg(updateChannels)
					break
				}

				for _, fp := range fileProviders {
					if fp.Path == e.Name {
						fw.updateFileProvider(fp, updateChannels)
						break
					}
				}
			}
		}
	}
}

func (fw *FileWatcher) updateFileProvider(fp *providers.FileProvider, updateChannels *[]chan string) {
	err := fp.Load()
	if err != nil {
		log.Printf("Failed to reload file %s err = %+v", fp.Path, err)
		return
	}

	for _, c := range *updateChannels {
		c <- "update-categories"
	}
}

func (fw *FileWatcher) updateCfg(updateChannels *[]chan string) {
	err := fw.cfg.Load(fw.cfgPath, utils.FileReader)
	if err != nil {
		log.Printf("Failed to reload config from path %s err = %+v", fw.cfgPath, err)
		return
	}

	for _, c := range *updateChannels {
		c <- "update-config"
	}
}

func newWatcher() *fsnotify.Watcher {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Failed to create fsnotify watcher err = %+v", err)
		return nil
	}
	return w
}
