package processor

import (
	"fmt"
	"strings"
)

// GenericProcessor is a fallback processor for data types without a specific processor
type GenericProcessor struct {
	CSVProcessor
}

// Process performs basic harmonization on generic CSV data
func (p *GenericProcessor) Process(inputFile, outputFile string) error {
	// Read CSV
	records, err := p.ReadCSV(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("data is empty or missing header")
	}

	// Get header and standardize it
	header := records[0]
	standardizedHeader := p.StandardizeHeader(header)
	records[0] = standardizedHeader

	// Look for patient ID column
	patientIDIdx := -1
	for i, h := range standardizedHeader {
		if strings.Contains(h, "patient") && (strings.Contains(h, "id") || strings.Contains(h, "number") || strings.Contains(h, "mrn")) {
			patientIDIdx = i
			break
		}
	}

	// Look for date columns
	dateColumns := make([]int, 0)
	for i, h := range standardizedHeader {
		if strings.Contains(h, "date") || strings.Contains(h, "time") {
			dateColumns = append(dateColumns, i)
		}
	}

	// Process data
	for i := 1; i < len(records); i++ {
		row := records[i]

		// Standardize patient ID if found
		if patientIDIdx >= 0 && patientIDIdx < len(row) {
			row[patientIDIdx] = p.standardizePatientID(row[patientIDIdx])
		}

		// Standardize dates if found
		for _, dateIdx := range dateColumns {
			if dateIdx < len(row) {
				row[dateIdx] = p.StandardizeDate(row[dateIdx])
			}
		}

		records[i] = row
	}

	// Write harmonized data
	if err := p.WriteCSV(outputFile, records); err != nil {
		return fmt.Errorf("failed to write harmonized data: %w", err)
	}

	return nil
}

// standardizePatientID ensures patient IDs have a consistent format
func (p *GenericProcessor) standardizePatientID(patientID string) string {
	// If it's already in PT-XXXX-XX format, keep it
	if strings.HasPrefix(patientID, "PT-") {
		return patientID
	}

	// Otherwise, ensure it's clean and consistent
	patientID = strings.TrimSpace(patientID)

	// If it's a short ID (likely a numeric ID), convert to PT format
	if len(patientID) < 8 {
		return fmt.Sprintf("PT-%s", patientID)
	}

	return patientID
}
