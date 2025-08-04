package processor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCSVProcessor(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "csv-processor-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test CSV file
	testCSV := filepath.Join(tmpDir, "test.csv")
	testData := []byte("patient_id,name,date_of_birth\n123,John Doe,1980-01-01\n456,Jane Smith,1990-02-15\n")
	if err := os.WriteFile(testCSV, testData, 0644); err != nil {
		t.Fatalf("Failed to write test CSV: %v", err)
	}

	// Create output file path
	outputCSV := filepath.Join(tmpDir, "output.csv")

	// Create a CSV processor
	processor := &CSVProcessor{}

	// Test ReadCSV
	records, err := processor.ReadCSV(testCSV)
	if err != nil {
		t.Fatalf("ReadCSV failed: %v", err)
	}

	// Check record count
	if len(records) != 3 { // Header + 2 data rows
		t.Errorf("Expected 3 records, got %d", len(records))
	}

	// Check header
	if records[0][0] != "patient_id" || records[0][1] != "name" || records[0][2] != "date_of_birth" {
		t.Errorf("Unexpected header: %v", records[0])
	}

	// Test WriteCSV
	if err := processor.WriteCSV(outputCSV, records); err != nil {
		t.Fatalf("WriteCSV failed: %v", err)
	}

	// Check that output file exists
	if _, err := os.Stat(outputCSV); os.IsNotExist(err) {
		t.Errorf("Output file was not created")
	}

	// Read back the output file
	outputRecords, err := processor.ReadCSV(outputCSV)
	if err != nil {
		t.Fatalf("Failed to read output CSV: %v", err)
	}

	// Check that output matches input
	if len(outputRecords) != len(records) {
		t.Errorf("Expected %d records in output, got %d", len(records), len(outputRecords))
	}

	// Test StandardizeHeader
	header := []string{"Patient ID", "Full Name", "Date of Birth"}
	standardized := processor.StandardizeHeader(header)
	
	expected := []string{"patient_id", "full_name", "date_of_birth"}
	for i, h := range standardized {
		if h != expected[i] {
			t.Errorf("Expected standardized header[%d] to be %s, got %s", i, expected[i], h)
		}
	}

	// Test StandardizeDate
	dates := map[string]string{
		"2020-01-15":   "2020-01-15", // Already standard format
		"01/15/2020":   "2020-01-15", // MM/DD/YYYY
		"15/01/2020":   "2020-01-15", // DD/MM/YYYY
		"Jan 15, 2020": "2020-01-15", // Month DD, YYYY
	}

	for input, expected := range dates {
		output := processor.StandardizeDate(input)
		if output != expected {
			t.Errorf("Expected StandardizeDate(%s) to be %s, got %s", input, expected, output)
		}
	}
}
