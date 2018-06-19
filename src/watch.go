package yolo

import (
	"github.com/azer/logger"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
	"strings"
)

func NewWatch(_include, _exclude *Patterns) (*Watch, error) {
	include, err := _include.Expand()
	if err != nil {
		return nil, err
	}

	exclude, err := _exclude.Expand()
	if err != nil {
		return nil, err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	watch := Watch{
		Include: include,
		Exclude: exclude,
		Watcher: *watcher,
	}

	return &watch, nil
}

type Watch struct {
	Include *Patterns
	Exclude *Patterns
	Watcher fsnotify.Watcher
}

func (watch *Watch) Add(path string) error {
	log.Info("Watching %s", path)

	if err := watch.Watcher.Add(path); err != nil {
		log.Error("Can not watch", logger.Attrs{
			"error": err,
			"path":  path,
		})
	}

	return nil
}

func (watch *Watch) Remove(path string) error {
	if err := watch.Watcher.Remove(path); err != nil {
		log.Error("Can not unwatch", logger.Attrs{
			"error": err,
			"path":  path,
		})
	}

	return nil
}

func (watch *Watch) AddRecursively(root string) error {
	fileInfo, err := os.Stat(root)
	if err != nil {
		return err
	}

	if !fileInfo.IsDir() {
		return watch.Add(root)
	}

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() && !strings.HasPrefix(path, ".git") && !strings.Contains(path, "/.git") {
			watch.Add(path)
		}

		return nil
	})

	return nil
}

func (watch *Watch) RemoveRecursively(root string) error {
	fileInfo, err := os.Stat(root)
	if err != nil {
		return err
	}

	if !fileInfo.IsDir() {
		return watch.Remove(root)
	}

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() {
			watch.Remove(path)
		}

		return nil
	})

	return nil
}

func (watch Watch) Start(callback func(*WatchEvent)) {
	defer watch.Watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watch.Watcher.Events:
				if watch.Exclude.Has(event.Name) {
					continue
				}

				we := &WatchEvent{
					Filename: event.Name,
					Create:   event.Op == fsnotify.Create,
					Write:    event.Op == fsnotify.Write,
					Rename:   event.Op == fsnotify.Rename,
					Remove:   event.Op == fsnotify.Remove,
				}

				if we.Create {
					watch.AddRecursively(we.Filename)
				} else if we.Remove {
					watch.Remove(we.Filename)
				}

				if we.ShouldBeNotified() {
					log.Info("New event", logger.Attrs{
						"name":    we.Filename,
						"create?": we.Create,
						"write?":  we.Write,
						"rename?": we.Rename,
						"remove?": we.Remove,
					})

					callback(we)
				}

			case err := <-watch.Watcher.Errors:
				log.Error("File watcher failed.", logger.Attrs{
					"error": err,
				})
			}
		}
	}()

	for _, path := range *watch.Include {
		watch.AddRecursively(path)
	}

	<-done
}
