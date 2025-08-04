package processor

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CSVProcessor is a base processor for CSV files
type CSVProcessor struct{}

// ReadCSV reads a CSV file and returns its contents
func (p *CSVProcessor) ReadCSV(filePath string) ([][]string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	return records, nil
}

// WriteCSV writes data to a CSV file
func (p *CSVProcessor) WriteCSV(filePath string, data [][]string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write all records
	if err := writer.WriteAll(data); err != nil {
		return fmt.Errorf("failed to write CSV: %w", err)
	}

	return nil
}

// StandardizeHeader standardizes header names (lowercase, trim spaces, replace spaces with underscores)
func (p *CSVProcessor) StandardizeHeader(header []string) []string {
	standardized := make([]string, len(header))
	for i, h := range header {
		// Lowercase, trim spaces, replace spaces with underscores
		standardized[i] = strings.ReplaceAll(strings.TrimSpace(strings.ToLower(h)), " ", "_")
	}
	return standardized
}

// GetColumnIndex returns the index of a column in the header
func (p *CSVProcessor) GetColumnIndex(header []string, columnName string) (int, error) {
	for i, h := range header {
		if strings.EqualFold(h, columnName) {
			return i, nil
		}
	}
	return -1, fmt.Errorf("column not found: %s", columnName)
}

// StandardizeDate standardizes date formats to YYYY-MM-DD
func (p *CSVProcessor) StandardizeDate(date string) string {
	// Try different date formats
	formats := []string{
		"2006-01-02",       // YYYY-MM-DD
		"01/02/2006",       // MM/DD/YYYY
		"02/01/2006",       // DD/MM/YYYY
		"01-02-2006",       // MM-DD-YYYY
		"02-01-2006",       // DD-MM-YYYY
		"Jan 02, 2006",     // Mon DD, YYYY
		"January 02, 2006", // Month DD, YYYY
	}

	for _, format := range formats {
		if t, err := parseDate(date, format); err == nil {
			return t.Format("2006-01-02")
		}
	}

	// Return original if no format matches
	return date
}

// Helper function to parse dates with different formats
func parseDate(date string, format string) (time.Time, error) {
	return time.Parse(format, date)
}
