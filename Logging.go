package FlowG

import (
	"fmt"
	"os"
	"time"
)

// Constants for log levels
const (
	DEBUG    = uint8(0)
	INFO     = uint8(1)
	WARNING  = uint8(2)
	ERROR    = uint8(3)
	CRITICAL = uint8(4)
)

// Map to convert level integers to strings
var levelNames = map[uint8]string{
	DEBUG:    "DEBUG",
	INFO:     "INFO",
	WARNING:  "WARNING",
	ERROR:    "ERROR",
	CRITICAL: "CRITICAL",
}

// Reverse map
var levelValues = make(map[string]uint8)

func init() {
	// Populate the reverse map using the forward map
	for k, v := range levelNames {
		levelValues[v] = k
	}
}

// GetLogLvLID returns the ID of the log level associated with the given name, along with a boolean indicating existence.
func GetLogLvLID(name string) (uint8, bool) {
	value, exists := levelValues[name]
	return value, exists
}

// Logging logs a message with a specified severity level. Messages are written to a log file specific to the current date.
func Logging(msg string, lvl uint8) {
	if config.logDir == "" {
		panic("Logging path undefined")
	}
	if config.logLvl > 4 {
		invalidLevel := config.logLvl
		config.logLvl = 1
		Logging(fmt.Sprintf("Loglevel was set to invalid level %d, defaulting to %s (%d)", invalidLevel, levelNames[config.logLvl], config.logLvl), WARNING)
	}
	if config.logLvl > lvl {
		return
	}

	// Set the log file name with today's date
	logFileName := fmt.Sprintf("%s/%s_%s.txt", config.logDir, config.logPrefix, time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// Handle the error silently (or print to stderr if needed)
		fmt.Println("Logging error: failed to open log file")
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Println("Logging error: failed to close log file")
		}
	}(file)

	// Determine the log level, default to WARNING if unknown
	logLevelName, ok := levelNames[lvl]
	if !ok {
		logLevelName = levelNames[WARNING] // Default to WARNING if unknown
		// Log an extra warning about the unknown level
		Logging(fmt.Sprintf("Unknown log level %d used by application, defaulting to %s", lvl, logLevelName), WARNING)
	}

	// Get the current timestamp and format the log entry
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	logEntry := fmt.Sprintf("%s - [%s] %s\n", timestamp, logLevelName, msg)

	// Write the log entry to the file
	_, _ = file.WriteString(logEntry)
}
