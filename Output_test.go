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
					SampleName:             "Sample1",
					Compound:               "Compound1",
					ResultCalculatedAmount: 0.123,
					ResultInterceptAmount:  0.456,
					InstrumentUsed:         "Instrument1",
					PeakIDForOutput:        1,
					DilutionFactor:         1.23,
					OutputToGlims:          "Output1",
				},
			},
		},
		{
			name:     "valid fileName and valid full sampleList",
			FileName: "output",
			expectOk: true,
			SampleList: []SampleStruct{
				{
					SampleName:             "Sample1",
					Compound:               "Compound1",
					ResultCalculatedAmount: 0.123,
					ResultInterceptAmount:  0.456,
					InstrumentUsed:         "Instrument1",
					PeakIDForOutput:        1,
					DilutionFactor:         1.23,
					OutputToGlims:          "Output1",
				},
			},
		},
		{
			name:     "valid fileName and valid partial sampleList",
			FileName: "output",
			expectOk: true,
			SampleList: []SampleStruct{
				{
					SampleName:             "Sample1",
					Compound:               "Compound1",
					ResultCalculatedAmount: 0.123,
					// Missing ResultInterceptAmount
					InstrumentUsed: "Instrument1",
					// Missing PeakIDForOutput
					// Missing DilutionFactor
					// Missing OutputToGlims
				},
			},
		},
		{
			name:     "valid fileName and invalid partial sampleList",
			FileName: "output",
			expectOk: false,
			SampleList: []SampleStruct{
				{
					// Missing SampleName, making it invalid
					Compound:               "Compound1",
					ResultCalculatedAmount: 0.123,
					ResultInterceptAmount:  0.456,
					InstrumentUsed:         "Instrument1",
					PeakIDForOutput:        1,
					DilutionFactor:         1.23,
					OutputToGlims:          "Output1",
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
			if !c.expectOk {
				if len(outputFiles) > 0 {
					t.Errorf("Expected no output files, got %v", len(outputFiles))
				}
				if len(logFiles) != 1 {
					t.Errorf("Expected 1 log file, got %v", len(outputFiles))
				}
			}
			if c.expectOk {
				if len(outputFiles) != 1 {
					t.Errorf("Expected 1 output file, got %v", len(outputFiles))
				}
				if len(logFiles) > 0 {
					t.Errorf("Expected no log files, got %v", len(outputFiles))
				}
			}

			if len(outputFiles) > 0 {
				fileBytes, _ := os.ReadFile(config.glimsDir + "/" + outputFiles[0].Name())
				fileContent := string(fileBytes)

				expectedContent := c.SampleList[0].SampleName + ";" + c.SampleList[0].Compound + ";;;" + convertToString(c.SampleList[0].ResultCalculatedAmount) + ";" + convertToString(c.SampleList[0].ResultInterceptAmount) + ";" + c.SampleList[0].InstrumentUsed + ";" + convertToString(c.SampleList[0].PeakIDForOutput) + ";" + convertToString(c.SampleList[0].DilutionFactor) + ";" + c.SampleList[0].OutputToGlims + "\n"
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
		input    int
		expected string
	}{
		{
			name:     "integer zero",
			input:    0,
			expected: "",
		},
		{
			name:     "positive integer",
			input:    123,
			expected: "123",
		},
		{
			name:     "negative integer",
			input:    -456,
			expected: "-456",
		},
	}

	for _, c := range intCases {
		t.Run(c.name, func(t *testing.T) {
			actual := convertToString(c.input)
			if actual != c.expected {
				t.Errorf("convertToString(%d): expected %s, got %s", c.input, c.expected, actual)
			}
		})
	}

	// `float64` tests
	floatCases := []struct {
		name     string
		input    float64
		expected string
	}{
		{
			name:     "float zero",
			input:    0.0,
			expected: "",
		},
		{
			name:     "positive float",
			input:    78.90,
			expected: "78.90",
		},
		{
			name:     "negative float",
			input:    -32.10,
			expected: "-32.10",
		},
	}

	for _, c := range floatCases {
		t.Run(c.name, func(t *testing.T) {
			actual := convertToString(c.input)
			if actual != c.expected {
				t.Errorf("convertToString(%f): expected %s, got %s", c.input, c.expected, actual)
			}
		})
	}
}
