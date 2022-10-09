package internal

type FileService struct {
	watcher  *FileWatcher
	args     *Args
	config   *Config
	notifier *WebsocketNotifier
	reader   func(string) ([]byte, error)

	providers []*FileProvider
}

func NewFileService(watcher *FileWatcher, args *Args, config *Config, notifier *WebsocketNotifier, reader func(string) ([]byte, error)) *FileService {
	return &FileService{
		watcher:   watcher,
		args:      args,
		config:    config,
		notifier:  notifier,
		reader:    reader,
		providers: make([]*FileProvider, 0),
	}
}

func (s *FileService) Init() error {
	Log.Debug.Printf("Initializing file service ...")
	err := s.watcher.Init()
	if err != nil {
		return err
	}

	Log.Debug.Printf("Adding config file to watcher ...")
	err = s.watcher.Add(s.args.ConfigPath.Get(), s.updateConfig)
	if err != nil {
		return err
	}

	Log.Debug.Printf("Loading file providers ...")
	s.updateFileProviders()
	return nil
}

func (s *FileService) Get() []Category {
	createdCategories := make(map[string]int, 0)
	categories := make([]Category, 0)

	Log.Debug.Printf("Retrieving categories from file providers ...")
	for _, fp := range s.providers {
		for _, pCat := range fp.Get() {
			Log.Debug.Printf("Retrieving categories from file provider for path = %s , categories = %+v", fp.Path, pCat)
			if index, ok := createdCategories[pCat.Name]; ok {
				categories[index].Items = append(categories[index].Items, pCat.Items...)
				continue
			}

			createdCategories[pCat.Name] = len(categories)
			categories = append(categories, pCat)
		}
	}

	Log.Debug.Printf("Returning categories = %+v", categories)
	return categories
}

func (s *FileService) updateConfig() {
	Log.Debug.Printf("Starting config load ...")
	err := s.config.Load(s.args.ConfigPath.Get(), FileReader)
	if err != nil {
		Log.Error.Printf("Failed to reload config err = %+v", err)
		return
	}

	Log.Debug.Printf("Sending message on config updated")
	s.notifier.Update(UpdateConfig)

	Log.Debug.Printf("Reloading file providers")
	s.updateFileProviders()
}

func (s *FileService) updateFileProviders() {
	s.removeProviders()
	s.addNewProviders()
}

func (s *FileService) addNewProviders() {
	for _, fpc := range s.config.FileProviders {
		found := false
		for _, fp := range s.providers {
			if fp.Path == fpc.Path {
				Log.Debug.Printf("File provider already exists fpc = %+v", fpc)
				found = true
				break
			}
		}

		if !found {
			Log.Debug.Printf("Adding new file provider with config = %+v", fpc)
			provider, err := NewFileProvider(fpc, s.notifier, s.reader)
			if err != nil {
				Log.Info.Printf("Failed to create file provider for config = %+v err = %+v", fpc, err)
				break
			}

			Log.Debug.Printf("Initializing file provider ...")
			err = provider.Init()
			if err != nil {
				Log.Info.Printf("Failed to init file provider for config = %+v err = %+v", fpc, err)
				break
			}

			Log.Debug.Printf("File provider initialized.")
			s.providers = append(s.providers, provider)

			if fpc.Watch {
				err = s.watcher.Add(provider.Path, provider.Update)
				if err != nil {
					Log.Info.Printf("Failed to watch file provider for config = %+v err = %+v", fpc, err)
					break
				}
				Log.Debug.Printf("File provider added to watcher. fpc = %+v", fpc)
			}
		}
	}
}

func (s *FileService) removeProviders() {
	for i, fp := range s.providers {
		found := false
		for _, fpc := range s.config.FileProviders {
			if fp.Path == fpc.Path {
				Log.Debug.Printf("File provider for path = %s found in config ... skip removal", fp.Path)
				found = true
				break
			}
		}

		if !found {
			Log.Debug.Printf("Removing file provider for path = %s", fp.Path)
			s.providers = append(s.providers[:i], s.providers[i+1:]...)
			err := s.watcher.Remove(fp.Path)
			if err != nil {
				Log.Error.Printf("Failed to remove file provider from watcher path = %s err = %+v", fp.Path, err)
			}
		}
	}
}
