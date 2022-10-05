package internal

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"simplydash/internal/config"
	"simplydash/internal/models"
	"simplydash/internal/providers"
)

type FileWatcher struct {
	fileProviders []*providers.FileProvider
	watcher       *fsnotify.Watcher
}

func NewFileWatcher(configs []config.FileProviderConfig) *FileWatcher {
	var fpArr []*providers.FileProvider
	for _, fpc := range configs {
		p, err := providers.NewFileProvider(fpc)
		if err != nil {
			log.Printf("Failed to create file provider for file = %s", fpc.Path)
		}
		fpArr = append(fpArr, p)
	}
	return &FileWatcher{fileProviders: fpArr, watcher: nil}
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
	var categories []models.Category

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

func (fw *FileWatcher) Watch(updateChannels *[]chan struct{}) {
	var toWatch []*providers.FileProvider
	for _, fp := range fw.fileProviders {
		if fp.Watch {
			toWatch = append(toWatch, fp)
		}
	}

	if len(toWatch) == 0 {
		return
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Print("Failed to create watcher.")
		return
	}

	for _, fp := range toWatch {
		err := w.Add(fp.Path)
		if err != nil {
			log.Printf("Failed to add file = %s to watcher err = %+v", fp.Path, err)
		}
	}

	go fw.watch(w, toWatch, updateChannels)
}

func (fw *FileWatcher) watch(
	w *fsnotify.Watcher,
	watch []*providers.FileProvider,
	updateChannels *[]chan struct{},
) {
	for {
		select {
		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			log.Printf("Received error = %+v", err)
		case e, ok := <-w.Events:
			if !ok {
				return
			}

			var matched *providers.FileProvider
			for _, fp := range watch {
				if fp.Path == e.Name {
					matched = fp
					break
				}
			}

			if matched == nil || !(e.Op == fsnotify.Write) {
				continue
			}

			err := matched.Load()
			if err != nil {
				log.Printf("Failed to reload file %s err = %+v", matched.Path, err)
			}

			for _, c := range *updateChannels {
				c <- struct{}{}
			}
		}
	}
}
