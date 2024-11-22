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
	Barcode           string
	TestName          string
	IsolationSequence string
	Result            float64
	ResultINT         float64
	ResultCT          float64
	InstrumentID      string
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
		if len(sample.Barcode) == 0 || len(sample.TestName) == 0 || len(sample.InstrumentID) == 0 {
			Logging("Incomplete sample send to GlimsOutput, skipping", WARNING)
			continue
		}

		record := []string{
			sample.Barcode,                    // Column 01, SPECIMEN_ID
			sample.TestName,                   // Column 02, TEST_ID
			sample.IsolationSequence,          // Column 03, ISOLATION_SEQUENCE
			convertToString(sample.Result),    // Column 04, RESULT
			convertToString(sample.ResultINT), // Column 05, RSLTTYPE_INT
			convertToString(sample.ResultCT),  // Column 06, RSLTTYPE_CT
			sample.InstrumentID,               // Column 07, INSTRUMENT_ID
		}
		if err = writer.Write(record); err != nil {
			Logging(fmt.Sprintf("Cannot write to Glims-output file '%s': %v", FileName, err), ERROR)
			return false
		}

		successCounter++
		Logging(fmt.Sprintf("GlimsOutput - Sample '%s' was processed correcly", sample.Barcode), DEBUG)
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
