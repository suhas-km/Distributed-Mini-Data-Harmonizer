package processor

import (
	"fmt"
	"strings"
)

// PatientProcessor processes patient data
type PatientProcessor struct {
	CSVProcessor
}

// Process harmonizes patient data
func (p *PatientProcessor) Process(inputFile, outputFile string) error {
	// Read CSV
	records, err := p.ReadCSV(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read patient data: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("patient data is empty or missing header")
	}

	// Get header and standardize it
	header := records[0]
	standardizedHeader := p.StandardizeHeader(header)

	// Find required columns
	patientIDIdx, err := p.getColumnIndex(standardizedHeader, "patient_id")
	if err != nil {
		return err
	}

	dobIdx, err := p.getColumnIndex(standardizedHeader, "date_of_birth")
	if err != nil {
		return err
	}

	genderIdx, err := p.getColumnIndex(standardizedHeader, "gender")
	if err != nil {
		return err
	}

	// Process data
	for i := 1; i < len(records); i++ {
		row := records[i]
		
		// Standardize patient ID (ensure consistent format)
		if patientIDIdx < len(row) {
			row[patientIDIdx] = p.standardizePatientID(row[patientIDIdx])
		}
		
		// Standardize date of birth
		if dobIdx < len(row) {
			row[dobIdx] = p.StandardizeDate(row[dobIdx])
		}
		
		// Standardize gender
		if genderIdx < len(row) {
			row[genderIdx] = p.standardizeGender(row[genderIdx])
		}
		
		records[i] = row
	}

	// Write harmonized data
	if err := p.WriteCSV(outputFile, records); err != nil {
		return fmt.Errorf("failed to write harmonized patient data: %w", err)
	}

	return nil
}

// Helper function to find column index with fallback names
func (p *PatientProcessor) getColumnIndex(header []string, columnName string) (int, error) {
	// Map of column name aliases
	columnAliases := map[string][]string{
		"patient_id":    {"patient_id", "patientid", "id", "patient_number", "mrn"},
		"date_of_birth": {"date_of_birth", "dob", "birth_date", "birthdate"},
		"gender":        {"gender", "sex"},
	}

	// Try direct match first
	for i, h := range header {
		if strings.EqualFold(h, columnName) {
			return i, nil
		}
	}

	// Try aliases
	if aliases, ok := columnAliases[columnName]; ok {
		for _, alias := range aliases {
			for i, h := range header {
				if strings.EqualFold(h, alias) {
					return i, nil
				}
			}
		}
	}

	return -1, fmt.Errorf("column not found: %s", columnName)
}

// standardizePatientID ensures patient IDs have a consistent format
func (p *PatientProcessor) standardizePatientID(patientID string) string {
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

// standardizeGender normalizes gender values
func (p *PatientProcessor) standardizeGender(gender string) string {
	gender = strings.TrimSpace(strings.ToLower(gender))
	
	switch gender {
	case "m", "male", "man":
		return "Male"
	case "f", "female", "woman":
		return "Female"
	case "u", "unknown", "unspecified", "":
		return "Unknown"
	case "o", "other":
		return "Other"
	default:
		return "Unknown"
	}
}
