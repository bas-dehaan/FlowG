package FlowG

import (
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"
)

func TestGetLogLvLID(t *testing.T) {
	cases := []struct {
		name      string
		logName   string
		wantValue uint8
		wantExist bool
	}{
		{"Loglvl exist", "INFO", 1, true},
		{"Loglvl not exist", "not_exist", 0, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotValue, gotExist := GetLogLvLID(c.logName)
			if gotValue != c.wantValue || gotExist != c.wantExist {
				t.Errorf("%s: expected %d, %t, got %d, %t", c.name, c.wantValue, c.wantExist, gotValue, gotExist)
			}
		})
	}
}

func TestLogging(t *testing.T) {
	// updating config for tests
	config.logDir = os.TempDir()
	config.logLvl = WARNING
	config.logPrefix = "Test"

	defer func() {
		// reset to default values after test, so not to break other tests
		config.logDir = ""
		config.logLvl = 1
		config.logPrefix = ""
	}()

	cases := []struct {
		name        string
		msg         string
		lvl         uint8
		expLog      bool
		invalidConf bool
	}{
		{"LvL below cutoff", "This is an info log entry", INFO, false, false},
		{"LvL at cutoff", "This is a warning log entry ", WARNING, true, false},
		{"LvL above cutoff", "This is a critical log entry ", CRITICAL, true, false},
		{"Invalid LvL", "This is an invalid log entry ", 99, true, false},
		{"Invalid config", "This is an invalid logLvL config ", CRITICAL, true, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.invalidConf {
				config.logLvl = 99
			}

			Logging(c.msg, c.lvl)

			logFileName := fmt.Sprintf("%s/%s_%s.txt", config.logDir, config.logPrefix, time.Now().Format("2006-01-02"))

			if _, err := os.Stat(logFileName); err == nil {
				if !c.expLog {
					t.Fatalf("Expected no log file creation, got 1 log file")
				}

				data, _ := os.ReadFile(logFileName)
				logEntry := string(data)

				if c.lvl <= 4 {
					re := regexp.MustCompile(fmt.Sprintf("20(?:\\d{2}-){2}\\d{2} (?:\\d{2}:){2}\\d{2}\\.\\d{3} - \\[%s\\] %s", levelNames[c.lvl], c.msg))

					if len(re.FindStringSubmatch(logEntry)) != 1 {
						t.Errorf("log file doesn't contain expected entry, ")
					}
				} else {
					reLog := regexp.MustCompile(fmt.Sprintf("20(?:\\d{2}-){2}\\d{2} (?:\\d{2}:){2}\\d{2}\\.\\d{3} - \\[%s\\] %s", levelNames[WARNING], c.msg))
					reWarn := regexp.MustCompile(fmt.Sprintf("20(?:\\d{2}-){2}\\d{2} (?:\\d{2}:){2}\\d{2}\\.\\d{3} - \\[WARNING\\] Unknown log level %d used by application, defaulting to WARNING", c.lvl))

					if len(reLog.FindStringSubmatch(logEntry)) != 1 {
						t.Errorf("log file doesn't contain expected log entry")
					}
					if len(reWarn.FindStringSubmatch(logEntry)) != 1 {
						t.Errorf("log file doesn't contain expected invalid loglvl warning")
					}
				}

				if c.invalidConf {
					reWarn := regexp.MustCompile(fmt.Sprintf("20(?:\\d{2}-){2}\\d{2} (?:\\d{2}:){2}\\d{2}\\.\\d{3} - \\[WARNING\\] Loglevel was set to invalid level 99, defaulting to INFO \\(%d\\)", INFO))

					if len(reWarn.FindStringSubmatch(logEntry)) != 1 {
						t.Errorf("log file doesn't contain expected invalid configuration warning, got: %s", logEntry)
					}
				}

				// removing log file
				err = os.Remove(logFileName)
				if err != nil {
					t.Fatalf("Failed to remove log file: %v", err)
				}
			} else if os.IsNotExist(err) && c.expLog {
				t.Errorf("Expected log file creation, got no log files")
			}
		})
	}
}
