package FlowG

import (
	"errors"
	"os"
	"reflect"
	"testing"
)

func TestSetConfig(t *testing.T) {
	cases := []struct {
		name    string
		key     string
		value   interface{}
		mkdir   bool
		wantErr error
	}{
		{"Setting glimsDir", "glimsDir", "./glimsDir", true, nil},
		{"Setting processedDir", "processedDir", "./processedDir", true, nil},
		{"Setting errorDir", "errorDir", "./errorDir", true, nil},
		{"Setting logDir", "logDir", "./logDir", true, nil},
		{"Setting logLvl", "logLvl", INFO, false, nil},
		{"Setting wrong key", "wrongKey", "value", false, errors.New(`unknown config key (wrongKey): use 'glimsDir', 'processedDir', 'errorDir', 'logDir', or 'logLvl'`)},
		{"Setting glimsDir to non-existing directory", "glimsDir", "/does/not/exist", false, errors.New("cannot find or access directory: /does/not/exist")},
		{"Setting glimsDir to non-string value", "glimsDir", 123, false, errors.New("glimsDir requires a string value")},
		{"Setting processedDir to non-string value", "processedDir", 123, false, errors.New("processedDir requires a string value")},
		{"Setting errorDir to non-string value", "errorDir", 123, false, errors.New("errorDir requires a string value")},
		{"Setting logDir to non-string value", "logDir", 123, false, errors.New("logDir requires a string value")},
		{"Setting logLvl to string value", "logLvl", "non-integer", false, errors.New("logLvl requires an integer value (int or uint8)")},
		{"Setting logLvl to int value", "logLvl", 3, false, nil},
		{"Setting logLvl to out-of-range value", "logLvl", uint8(5), false, errors.New("logLvl requires a valid log level, use 0 (DEBUG), 1 (INFO), 2 (WARN), 3 (ERROR), or 4 (CRITICAL)")},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.mkdir {
				err := os.Mkdir(c.value.(string), os.ModePerm)
				if err != nil {
					t.Fatalf("Could not create directory '%s'", c.value.(string))
				}
				defer func(path string) {
					err = os.RemoveAll(path)
					if err != nil {
						t.Fatalf("Could not remove directory '%s'", c.value.(string))
					}
				}(c.value.(string))
			}

			err := SetConfig(c.key, c.value)
			if (err != nil && c.wantErr == nil) || (err == nil && c.wantErr != nil) || (err != nil && c.wantErr != nil && err.Error() != c.wantErr.Error()) {
				t.Errorf("SetConfig(%q, %v) returned error %q, wanted error %q", c.key, c.value, err, c.wantErr)
			}
		})
	}
}

func TestGetConfig(t *testing.T) {
	cases := []struct {
		name    string
		key     string
		wantVal interface{}
		wantErr error
	}{
		{"Getting Glimsdir", "glimsDir", "./glimsDir", nil},
		{"Getting processedDir", "processedDir", "./processedDir", nil},
		{"Getting errorDir", "errorDir", "./errorDir", nil},
		{"Getting logDir", "logDir", "./logDir", nil},
		{"Getting logLvl", "logLvl", INFO, nil},
		{"Getting wrong key", "wrongKey", nil, errors.New(`unknown config key (wrongKey): use 'glimsDir', 'processedDir', 'errorDir', 'logDir', or 'logLvl'`)},
	}

	config = &configStruct{
		glimsDir:     "./glimsDir",
		processedDir: "./processedDir",
		errorDir:     "./errorDir",
		logDir:       "./logDir",
		logLvl:       INFO,
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			val, err := GetConfig(c.key)
			if (err != nil && c.wantErr == nil) || (err == nil && c.wantErr != nil) || (err != nil && c.wantErr != nil && err.Error() != c.wantErr.Error()) {
				t.Errorf("GetConfig(%q) returned error %q, wanted error %q", c.key, err, c.wantErr)
			}

			if !reflect.DeepEqual(val, c.wantVal) {
				t.Errorf("GetConfig(%q) returned value %v, wanted value %v", c.key, val, c.wantVal)
			}
		})
	}
}
