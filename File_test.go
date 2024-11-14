package FlowG

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func createTestFolders() error {
	// Glims folder
	err := os.Mkdir(config.glimsDir, os.ModePerm)
	if err != nil {
		return err
	}

	// Import folder
	err = os.Mkdir(config.importDir, os.ModePerm)
	if err != nil {
		return err
	}

	// Processed folder
	err = os.Mkdir(config.processedDir, os.ModePerm)
	if err != nil {
		return err
	}

	// Error folder
	err = os.Mkdir(config.errorDir, os.ModePerm)
	if err != nil {
		return err
	}

	// Log folder
	err = os.Mkdir(config.logDir, os.ModePerm)
	return err
}

func destroyTestFolders() error {
	// Glims folder
	err := os.RemoveAll(config.glimsDir)
	if err != nil {
		return err
	}

	// Import folder
	err = os.RemoveAll(config.importDir)
	if err != nil {
		return err
	}

	// Processed folder
	err = os.RemoveAll(config.processedDir)
	if err != nil {
		return err
	}

	// Error folder
	err = os.RemoveAll(config.errorDir)
	if err != nil {
		return err
	}

	// Log folder
	err = os.RemoveAll(config.logDir)
	return err
}

func TestFileMove(t *testing.T) {
	cases := []struct {
		name         string
		validPath    bool
		okProcessing bool
	}{
		{"Valid path - Ok processing", true, true},
		{"Valid path - Not ok processing", true, false},
		{"Invalid path", false, false},
	}

	config = &configStruct{
		glimsDir:     "./glims",
		importDir:    "./import",
		processedDir: "./processed",
		errorDir:     "./error",
		logDir:       "./log",
		logLvl:       WARNING,
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// create testfolders
			err := createTestFolders()
			defer func() {
				err = destroyTestFolders()
				if err != nil {
					t.Fatalf("Error cleaning up test folders: %v", err)
				}
			}()
			if err != nil {
				t.Fatalf("Error creating test folders: %v", err)
			}

			// Create dummy file
			originalPath := filepath.Join(config.importDir, "testFile.txt")
			if c.validPath {
				file, err := os.Create(originalPath)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				err = file.Close()
				if err != nil {
					t.Fatalf("Failed to close test file: %v", err)
				}
			}

			FileMove(originalPath, c.okProcessing)

			processedFileList, err := filepath.Glob(filepath.Join(config.processedDir, "*_"+filepath.Base(originalPath)))
			if err != nil {
				t.Fatalf("Error reading 'processed' folder: %v", err)
			}
			errorFileList, err := filepath.Glob(filepath.Join(config.errorDir, "*_"+filepath.Base(originalPath)))
			if err != nil {
				t.Fatalf("Error reading 'error' folder: %v", err)
			}
			logFileList, err := filepath.Glob(filepath.Join(config.logDir, "*.txt"))
			if err != nil {
				t.Fatalf("Error reading 'log' folder: %v", err)
			}

			if c.validPath {
				if c.okProcessing {
					if len(processedFileList) != 1 && len(errorFileList) != 0 && len(logFileList) != 0 {
						t.Errorf("Expected: 1 processed, 0 errors, 0 logs; Got: %d processed, %d errors, %d logs",
							len(processedFileList), len(errorFileList), len(logFileList))
					}
				} else {
					if len(processedFileList) != 0 && len(errorFileList) != 1 && len(logFileList) != 0 {
						t.Errorf("Expected: 0 processed, 1 error, 0 logs; Got: %d processed, %d errors, %d logs",
							len(processedFileList), len(errorFileList), len(logFileList))
					}
				}
			} else {
				if len(processedFileList) != 0 && len(errorFileList) != 0 && len(logFileList) != 1 {
					t.Errorf("Expected: 0 processed, 0 errors, 1 log; Got: %d processed, %d errors, %d logs",
						len(processedFileList), len(errorFileList), len(logFileList))
				}
			}

			// Check if file is gone on original path
			_, err = os.Stat(originalPath)
			if err == nil {
				t.Errorf("File is unexpectedly still present on original location")
				err = os.Remove(originalPath)
				if err != nil {
					t.Fatalf("Failed to remove original file: %v", err)
				}
			} else if !os.IsNotExist(err) {
				t.Fatalf("Unexpected error while checking file removal: %v", err)
			}

			// Check if unexpected logfiles were created
			files, err := filepath.Glob(filepath.Join(config.logDir, "*.txt"))
			if err != nil {
				t.Fatalf("Unexpected error while checking logfiles: %v", err)
			}
			if len(files) != 0 && c.validPath {
				t.Errorf("Unexpected logfiles were created: %v", files)
				err = os.RemoveAll(config.processedDir)
				if err != nil {
					t.Fatalf("Error cleaning up test folders: %v", err)
				}
				err = os.Mkdir(config.processedDir, os.ModePerm)
				if err != nil {
					t.Fatalf("Error cleaning up test folders: %v", err)
				}
			} else if len(files) == 0 && !c.validPath {
				t.Errorf("The expected logfile was not created")
			}
		})
	}
}

func TestFileWatch(t *testing.T) {
	cases := []struct {
		name        string
		validFolder bool
		createFile  bool
	}{
		{
			name:        "Valid folder, file creation",
			validFolder: true,
			createFile:  true,
		},
		{
			name:        "Valid folder, no file creation",
			validFolder: true,
			createFile:  false,
		},
		{
			name:        "Invalid folder",
			validFolder: false,
			createFile:  false,
		},
	}

	config = &configStruct{
		glimsDir:     "./glims",
		importDir:    "./import",
		processedDir: "./processed",
		errorDir:     "./error",
		logDir:       "./log",
		logLvl:       WARNING,
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var err error
			var file *os.File

			if c.validFolder {
				err = createTestFolders()
				if err != nil {
					t.Fatalf("Error creating test folders: %v", err)
				}
			} else {
				// Always create a log folder to test the logging capabilities
				err = os.Mkdir(config.logDir, os.ModePerm)
				if err != nil {
					t.Fatalf("Error creating test log folder: %v", err)
				}
			}

			// Channel to capture callback invocation
			callbackInvoked := make(chan bool, 1)

			callback := func(path string) bool {
				callbackInvoked <- true
				return true
			}

			go func() {
				FileWatch(callback)
			}()

			// Allow time for FileWatch to initialize
			time.Sleep(100 * time.Millisecond)

			if c.createFile {
				file, err = os.Create(filepath.Join(config.importDir, "testFile.txt"))
				if err != nil {
					t.Fatalf("Error creating test file: %v", err)
				}
			}

			select {
			case invoked := <-callbackInvoked:
				if invoked != c.createFile {
					t.Errorf("FileWatch() callback invoked = %v, expected %v", invoked, c.createFile)
				}
			case <-time.After(2 * time.Second):
				if !c.validFolder {
					// If no folder is created, the watch should fail and a log should be created
					logFiles, _ := os.ReadDir(config.logDir)
					if len(logFiles) != 1 {
						t.Errorf("FileWatch() did not create log entry when expected to")
					}
				} else if c.createFile {
					// If the folder and file is created, we would have expected a callback trigger by now: fail test
					t.Errorf("FileWatch() got timeout, expected file detection")
				}
				// If the folder is valid and no file is created, a timeout is expected: do not fail test
			}

			if c.createFile {
				err = file.Close()
				if err != nil {
					t.Fatalf("Error closing testFile.txt: %v", err)
				}
			}
			if c.validFolder {
				err = destroyTestFolders()
				if err != nil {
					t.Fatalf("Error cleaning up test folders: %v", err)
				}
			} else {
				err = os.RemoveAll(config.logDir)
				if err != nil {
					t.Fatalf("Error cleaning up test log folder: %v", err)
				}
			}
		})
	}
}
