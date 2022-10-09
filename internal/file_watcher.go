package internal

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"path"
)

type FileWatcher struct {
	files   []*file
	watcher *fsnotify.Watcher
}

func NewFileWatcher() *FileWatcher {
	return &FileWatcher{files: make([]*file, 0), watcher: nil}
}

type file struct {
	path    string
	onWrite func()
}

func (w *FileWatcher) Init() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		Log.Error.Printf("Failed to create fsnotify watcher err = %+v", err)
		return err
	}
	w.watcher = watcher

	Log.Debug.Printf("Successfully initialized fsnotify watcher. Starting watch loop ...")
	go w.watchLoop()
	return nil
}

func (w *FileWatcher) Add(path string, onWrite func()) error {
	if w.watcher == nil {
		return errors.New("file watcher not initialized")
	}

	for _, f := range w.files {
		if f.path == path {
			return nil
		}
	}

	err := w.watcher.Add(path)
	if err != nil {
		return err
	}

	w.files = append(w.files, &file{path: path, onWrite: onWrite})
	return nil
}

func (w *FileWatcher) Remove(path string) error {
	if w.watcher == nil {
		return errors.New("file watcher not initialized")
	}

	for i, f := range w.files {
		if f.path == path {
			err := w.watcher.Remove(path)
			if err != nil {
				return err
			}
			w.files = append(w.files[:i], w.files[i+1:]...)
			return nil
		}
	}
	return nil
}

func (w *FileWatcher) watchLoop() {
	for {
		select {
		case err, ok := <-w.watcher.Errors:
			if !ok {
				Log.Info.Println("Watcher closed, stopping watch loop.")
				return
			}
			Log.Debug.Printf("fsnotify: got error = %+v", err)
		case event, ok := <-w.watcher.Events:
			if !ok {
				Log.Info.Println("Watcher closed, stopping watch loop.")
				return
			}
			Log.Debug.Printf("fsnotify: got event = %+v", event)

			if event.Op == fsnotify.Write {
				for _, f := range w.files {
					if path.Clean(f.path) == path.Clean(event.Name) {
						f.onWrite()
					}
				}
			}
		}
	}
}
