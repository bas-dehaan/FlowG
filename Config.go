package FlowG

import (
	"errors"
	"fmt"
	"os"
)

type configStruct struct {
	glimsDir     string
	processedDir string
	errorDir     string
	logDir       string
	logLvl       uint8
}

var config = &configStruct{}

// SetConfig sets the specified configuration key to the provided value, performing type-checking and validation.
// Available keys are: 'glimsDir', 'processedDir', 'errorDir', 'logDir', 'logLvl'.
func SetConfig(key string, value interface{}) error {
	isDir := false

	switch key {
	case "glimsDir":
		v, ok := value.(string)
		if !ok {
			return errors.New("glimsDir requires a string value")
		}
		config.glimsDir = v
		isDir = true

	case "processedDir":
		v, ok := value.(string)
		if !ok {
			return errors.New("processedDir requires a string value")
		}
		config.processedDir = v
		isDir = true

	case "errorDir":
		v, ok := value.(string)
		if !ok {
			return errors.New("errorDir requires a string value")
		}
		config.errorDir = v
		isDir = true

	case "logDir":
		v, ok := value.(string)
		if !ok {
			return errors.New("logDir requires a string value")
		}
		config.logDir = v
		isDir = true

	case "logLvl":
		v, ok := value.(uint8)

		if !ok {
			// Attempt repair for non-uint8 int's
			val, ok := value.(int)
			v = uint8(val)

			if !ok {
				return errors.New("logLvl requires an integer value (int or uint8)")
			}
		}

		if _, exists := levelNames[v]; !exists {
			return errors.New("logLvl requires a valid log level, use 0 (DEBUG), 1 (INFO), 2 (WARN), 3 (ERROR), or 4 (CRITICAL)")
		}
		config.logLvl = v

	default:
		return fmt.Errorf("unknown config key (%s): use 'glimsDir', 'processedDir', 'errorDir', 'logDir', or 'logLvl'", key)
	}

	// Check directory existence only for path keys
	if isDir {
		if _, err := os.Stat(value.(string)); os.IsNotExist(err) {
			return fmt.Errorf("cannot find or access directory: %s", value)
		}
	}

	return nil
}

// GetConfig retrieves the configuration value associated with the given key.
// Available keys are: 'glimsDir', 'processedDir', 'errorDir', 'logDir', 'logLvl'.
// Returns the configuration value and a nil error if the key is found,
// otherwise returns nil and an error indicating that the key is unknown.
func GetConfig(key string) (interface{}, error) {
	switch key {
	case "glimsDir":
		return config.glimsDir, nil
	case "processedDir":
		return config.processedDir, nil
	case "errorDir":
		return config.errorDir, nil
	case "logDir":
		return config.logDir, nil
	case "logLvl":
		return config.logLvl, nil
	default:
		return nil, fmt.Errorf("unknown config key (%s): use 'glimsDir', 'processedDir', 'errorDir', 'logDir', or 'logLvl'", key)
	}
}