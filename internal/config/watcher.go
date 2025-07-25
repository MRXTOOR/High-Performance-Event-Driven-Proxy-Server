package config

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

type OnChangeFunc func(string) bool

func StartConfigWatcher(configPath string, onChange OnChangeFunc) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("[HOTRELOAD] Ошибка создания watcher: %v", err)
		return
	}
	defer watcher.Close()

	dir, _ := splitDirFile(configPath)
	if err := watcher.Add(dir); err != nil {
		log.Printf("[HOTRELOAD] Ошибка добавления watcher: %v", err)
		return
	}

	lastReload := time.Now()
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 && event.Name == configPath {
				if time.Since(lastReload) > 500*time.Millisecond {
					log.Printf("[HOTRELOAD] Обнаружено изменение конфига: %s", event.Name)
					if onChange(configPath) {
						lastReload = time.Now()
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("[HOTRELOAD] Ошибка watcher: %v", err)
		}
	}
}

func splitDirFile(path string) (string, string) {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[:i], path[i+1:]
		}
	}
	return ".", path
}
