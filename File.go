package FlowG

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileWatch continuously watches the directory specified in config.glimsDir for new file creation events.
// When a new file is detected, it waits for 1 second before executing the provided callback function with
// the file path as an argument. If the callback returns true, the file is moved to the processed directory;
// otherwise, it is moved to the error directory. Logging is performed for critical errors during the process.
func FileWatch(callback func(string) bool) {
	if _, err := os.Stat(config.glimsDir); os.IsNotExist(err) {
		Logging(fmt.Sprintf("Cannot find importDir: %v", err), CRITICAL)
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		Logging(fmt.Sprintf("Error while starting watching importDir: %v", err), CRITICAL)
		return
	}
	defer func(watcher *fsnotify.Watcher) {
		err = watcher.Close()
		if err != nil {
			Logging(fmt.Sprintf("Error while closing watch on importDir: %v", err), CRITICAL)
		}
	}(watcher)

	err = watcher.Add(config.glimsDir)
	if err != nil {
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				Logging("An unknown fatal error occurred while watching importDir", CRITICAL)
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				filePath := event.Name
				go func(filePath string) {
					time.Sleep(1 * time.Second) // Wait 1 second before triggering to ensure completion of file write
					fileOk := callback(filePath)
					FileMove(filePath, fileOk)
				}(filePath)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				Logging(fmt.Sprintf("Fatal error while watching importDir: %v", err), CRITICAL)
				return
			}
			Logging(fmt.Sprintf("Non-fatal error while watching importDir: %v", err), ERROR)
		}
	}
}

// FileMove moves a file from the given path to a processed or error directory based on the status flag (ok).
func FileMove(path string, ok bool) {
	timestamp := strings.ReplaceAll(time.Now().Format("20060102150405.000"), ".", "")
	FileName := fmt.Sprintf("%s_%s", timestamp, filepath.Base(path))
	var destPath string

	if ok {
		destPath = filepath.Join(config.processedDir, FileName)
	} else {
		destPath = filepath.Join(config.errorDir, FileName)
	}

	err := os.Rename(path, destPath)
	if err != nil {
		Logging(fmt.Sprintf("Error while moving file: %v", err), ERROR)
	}
}
