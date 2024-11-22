package FlowG

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type SampleStruct struct {
	SampleName             string
	Compound               string
	ResultCalculatedAmount float64
	ResultInterceptAmount  float64
	InstrumentUsed         string
	PeakIDForOutput        int
	DilutionFactor         float64
	OutputToGlims          string
}

// GlimsOutput processes a list of samples and outputs them to a CSV file with the provided filename according to the FlowG standard.
func GlimsOutput(FileName string, SampleList []SampleStruct) bool {
	if len(FileName) == 0 {
		Logging("Invalid or no FileName was given to GlimsOutput, doing nothing", ERROR)
		return false
	}
	if len(SampleList) == 0 {
		Logging("Empty SampleList was given to GlimsOutput, doing nothing", WARNING)
		return false
	}

	timestamp := strings.ReplaceAll(time.Now().Format("20060102150405.000"), ".", "")
	FileName = fmt.Sprintf("%s_%s.txt", timestamp, FileName)

	file, err := os.Create(filepath.Join(config.glimsDir, FileName))
	if err != nil {
		Logging(fmt.Sprintf("Cannot create Glims-output file '%s': %v", FileName, err), ERROR)
		return false
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			Logging(fmt.Sprintf("Cannot close Glims-output file '%s': %v", FileName, err), ERROR)
			return
		}
		Logging(fmt.Sprintf("GlimsOutput successfully closed file '%s'", FileName), DEBUG)
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Comma = ';'

	successCounter := 0
	for _, sample := range SampleList {
		Logging(fmt.Sprintf("GlimsOutput - Processing sample: %v", sample), DEBUG)
		if len(sample.SampleName) == 0 || len(sample.Compound) == 0 || len(sample.InstrumentUsed) == 0 {
			Logging("Incomplete sample send to GlimsOutput, skipping", WARNING)
			continue
		}

		record := []string{
			sample.SampleName, // Column 01, Sample name
			sample.Compound,   // Column 02, Compound
			"",                // Column 03, Empty column
			"",                // Column 04, Empty column
			convertToString(sample.ResultCalculatedAmount), // Column 05, Final concentration
			convertToString(sample.ResultInterceptAmount),  // Column 06, Intercept concentration
			sample.InstrumentUsed,                          // Column 07, Instrument Used
			convertToString(sample.PeakIDForOutput),        // Column 08, Peak ID for Output
			convertToString(sample.DilutionFactor),         // Column 09, Dilution factor
			sample.OutputToGlims,                           // Column 10, Output to Glims
		}
		if err = writer.Write(record); err != nil {
			Logging(fmt.Sprintf("Cannot write to Glims-output file '%s': %v", FileName, err), ERROR)
			return false
		}

		successCounter++
		Logging(fmt.Sprintf("GlimsOutput - Sample '%s' was processed correcly", sample.SampleName), DEBUG)
	}

	// Delete the outputfile if there were no samples successfully added to it
	if successCounter == 0 {
		// Force a file closure
		err = file.Close()
		if err != nil {
			Logging(fmt.Sprintf("Cannot close Glims-output file '%s': %v", FileName, err), ERROR)
			return false
		}
		Logging(fmt.Sprintf("GlimsOutput successfully closed file '%s'", FileName), DEBUG)

		// Delete file
		err = os.Remove(filepath.Join(config.glimsDir, FileName))
		Logging(fmt.Sprintf("The file '%s' didn't contain any valid sampled. The empty Glims-output was deleted", FileName), INFO)
		if err != nil {
			Logging(fmt.Sprintf("Cannot remove empty/obsolete Glims-output file '%s': %v", FileName, err), WARNING)
		}
		return false
	}
	return true
}

// convertToString converts an integer or float64 value to a string. Returns an empty string if the value is zero.
func convertToString[T int | float64](value T) string {
	if value == 0 {
		return ""
	}
	switch any(value).(type) {
	case int:
		return strconv.Itoa(any(value).(int))
	case float64:
		return strconv.FormatFloat(any(value).(float64), 'f', 2, 64)
	default:
		Logging(fmt.Sprintf("Cannot convert value to T. Got type '%T'", any(value)), ERROR)
		return ""
	}
}
