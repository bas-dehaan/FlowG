package FlowG

import (
	"os"
	"strings"
	"testing"
)

func TestGlimsOutput(t *testing.T) {
	cases := []struct {
		name       string
		FileName   string
		expectOk   bool
		SampleList []SampleStruct
	}{
		{
			name:       "empty fileName and empty sampleList",
			FileName:   "",
			expectOk:   false,
			SampleList: []SampleStruct{},
		},
		{
			name:       "valid fileName and empty sampleList",
			FileName:   "testFile",
			expectOk:   false,
			SampleList: []SampleStruct{},
		},
		{
			name:     "empty fileName and valid full sampleList",
			FileName: "",
			expectOk: false,
			SampleList: []SampleStruct{
				{
					Barcode:           "Sample1",
					TestName:          "Compound1",
					IsolationSequence: "1",
					Result:            ptrFloat64(0.123),
					ResultINT:         ptrFloat64(0.456),
					ResultCT:          ptrFloat64(0.789),
					InstrumentID:      "Instrument1",
				},
			},
		},
		{
			name:     "valid fileName and valid full sampleList",
			FileName: "output",
			expectOk: true,
			SampleList: []SampleStruct{
				{
					Barcode:           "Sample1",
					TestName:          "Compound1",
					IsolationSequence: "1",
					Result:            ptrFloat64(0.123),
					ResultINT:         ptrFloat64(0.456),
					ResultCT:          ptrFloat64(0.789),
					InstrumentID:      "Instrument1",
				},
			},
		},
		{
			name:     "valid fileName and valid partial sampleList",
			FileName: "output",
			expectOk: true,
			SampleList: []SampleStruct{
				{
					Barcode:  "Sample1",
					TestName: "Compound1",
					// Missing IsolationSequence
					Result: ptrFloat64(0.123),
					// Missing ResultINT
					// Missing ResultCT
					InstrumentID: "Instrument1",
				},
			},
		},
		{
			name:     "valid fileName and invalid partial sampleList",
			FileName: "output",
			expectOk: false,
			SampleList: []SampleStruct{
				{
					// Missing Barcode, making it invalid
					TestName:          "Compound1",
					IsolationSequence: "1",
					Result:            ptrFloat64(0.123),
					ResultINT:         ptrFloat64(0.456),
					ResultCT:          ptrFloat64(0.789),
					InstrumentID:      "Instrument1",
				},
			},
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

			ok := GlimsOutput(c.FileName, c.SampleList)

			if ok != c.expectOk {
				t.Errorf("Unexpected status, expected %v, got %v", c.expectOk, ok)
			}

			outputFiles, _ := os.ReadDir(config.glimsDir)
			logFiles, _ := os.ReadDir(config.logDir)
			if c.expectOk {
				if len(outputFiles) != 1 {
					t.Errorf("Expected 1 output file, got %v", len(outputFiles))
				}
				if len(logFiles) > 0 {
					t.Errorf("Expected no log files, got %v", len(outputFiles))
				}
			} else {
				if len(outputFiles) > 0 {
					t.Errorf("Expected no output files, got %v", len(outputFiles))
				}
				if len(logFiles) != 1 {
					t.Errorf("Expected 1 log file, got %v", len(outputFiles))
				}
			}

			if len(outputFiles) > 0 {
				fileBytes, _ := os.ReadFile(config.glimsDir + "/" + outputFiles[0].Name())
				fileContent := string(fileBytes)

				expectedContent := c.SampleList[0].Barcode + ";" + c.SampleList[0].TestName + ";" + c.SampleList[0].IsolationSequence + ";" + convertToString(c.SampleList[0].Result) + ";" + convertToString(c.SampleList[0].ResultINT) + ";" + convertToString(c.SampleList[0].ResultCT) + ";" + c.SampleList[0].InstrumentID + "\n"
				if !strings.Contains(fileContent, expectedContent) {
					t.Errorf("Expected '%v' in output file, got '%v'", expectedContent, fileContent)
				}
			}
		})
	}
}

func TestConvertToString(t *testing.T) {
	// `int` tests
	intCases := []struct {
		name     string
		input    *int
		expected string
	}{
		{
			name:     "positive integer",
			input:    ptrInt(123),
			expected: "123",
		},
		{
			name:     "negative integer",
			input:    ptrInt(-456),
			expected: "-456",
		},
		{
			name:     "integer zero",
			input:    ptrInt(0),
			expected: "0",
		},
		{
			name:     "nil integer",
			input:    nil,
			expected: "",
		},
	}

	for _, c := range intCases {
		t.Run(c.name, func(t *testing.T) {
			actual := convertToString(c.input)
			if actual != c.expected {
				t.Errorf("convertToString(%d): expected %s, got %s", *c.input, c.expected, actual)
			}
		})
	}

	// `float64` tests
	floatCases := []struct {
		name     string
		input    *float64
		expected string
	}{
		{
			name:     "positive float",
			input:    ptrFloat64(78.90),
			expected: "78.90",
		},
		{
			name:     "negative float",
			input:    ptrFloat64(-32.10),
			expected: "-32.10",
		},
		{
			name:     "float zero",
			input:    ptrFloat64(0.00),
			expected: "0.00",
		},
		{
			name:     "nil float",
			input:    nil,
			expected: "",
		},
	}

	for _, c := range floatCases {
		t.Run(c.name, func(t *testing.T) {
			actual := convertToString(c.input)
			if actual != c.expected {
				t.Errorf("convertToString(%f): expected %s, got %s", *c.input, c.expected, actual)
			}
		})
	}
}

// Helper functions for creating pointers
func ptrInt(v int) *int {
	return &v
}

func ptrFloat64(v float64) *float64 {
	return &v
}
