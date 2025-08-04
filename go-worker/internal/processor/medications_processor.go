package processor

import (
	"fmt"
	"strings"
)

// MedicationsProcessor processes medication data
type MedicationsProcessor struct {
	CSVProcessor
}

// Process harmonizes medication data
func (p *MedicationsProcessor) Process(inputFile, outputFile string) error {
	// Read CSV
	records, err := p.ReadCSV(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read medications data: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("medications data is empty or missing header")
	}

	// Get header and standardize it
	header := records[0]
	standardizedHeader := p.StandardizeHeader(header)

	// Find required columns
	patientIDIdx, err := p.getColumnIndex(standardizedHeader, "patient_id")
	if err != nil {
		return err
	}

	medicationNameIdx, err := p.getColumnIndex(standardizedHeader, "medication_name")
	if err != nil {
		return err
	}

	// Find other important columns
	prescribedDateIdx, _ := p.getColumnIndex(standardizedHeader, "prescribed_date")
	dosageIdx, _ := p.getColumnIndex(standardizedHeader, "dosage")
	frequencyIdx, _ := p.getColumnIndex(standardizedHeader, "frequency")
	routeIdx, _ := p.getColumnIndex(standardizedHeader, "route")

	// Process data
	for i := 1; i < len(records); i++ {
		row := records[i]
		
		// Standardize patient ID
		if patientIDIdx < len(row) {
			row[patientIDIdx] = p.standardizePatientID(row[patientIDIdx])
		}
		
		// Standardize medication name
		if medicationNameIdx < len(row) {
			row[medicationNameIdx] = p.standardizeMedicationName(row[medicationNameIdx])
		}
		
		// Standardize prescribed date
		if prescribedDateIdx >= 0 && prescribedDateIdx < len(row) {
			row[prescribedDateIdx] = p.StandardizeDate(row[prescribedDateIdx])
		}
		
		// Standardize dosage
		if dosageIdx >= 0 && dosageIdx < len(row) {
			row[dosageIdx] = p.standardizeDosage(row[dosageIdx])
		}
		
		// Standardize frequency
		if frequencyIdx >= 0 && frequencyIdx < len(row) {
			row[frequencyIdx] = p.standardizeFrequency(row[frequencyIdx])
		}
		
		// Standardize route
		if routeIdx >= 0 && routeIdx < len(row) {
			row[routeIdx] = p.standardizeRoute(row[routeIdx])
		}
		
		records[i] = row
	}

	// Write harmonized data
	if err := p.WriteCSV(outputFile, records); err != nil {
		return fmt.Errorf("failed to write harmonized medications data: %w", err)
	}

	return nil
}

// Helper function to find column index with fallback names
func (p *MedicationsProcessor) getColumnIndex(header []string, columnName string) (int, error) {
	// Map of column name aliases
	columnAliases := map[string][]string{
		"patient_id":      {"patient_id", "patientid", "id", "patient_number", "mrn"},
		"medication_name": {"medication_name", "drug_name", "medication", "drug", "med_name"},
		"prescribed_date": {"prescribed_date", "date_prescribed", "order_date", "start_date"},
		"dosage":          {"dosage", "dose", "strength", "amount"},
		"frequency":       {"frequency", "freq", "schedule", "sig", "instructions"},
		"route":           {"route", "administration_route", "route_of_administration"},
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
func (p *MedicationsProcessor) standardizePatientID(patientID string) string {
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

// standardizeMedicationName normalizes medication names
func (p *MedicationsProcessor) standardizeMedicationName(name string) string {
	name = strings.TrimSpace(name)
	
	// Capitalize first letter of each word
	words := strings.Fields(strings.ToLower(name))
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	
	return strings.Join(words, " ")
}

// standardizeDosage normalizes medication dosages
func (p *MedicationsProcessor) standardizeDosage(dosage string) string {
	dosage = strings.TrimSpace(dosage)
	
	// Common abbreviations to standardize
	replacements := map[string]string{
		"mcg":  "μg",
		"mg":   "mg",
		"g":    "g",
		"ml":   "mL",
		"meq":  "mEq",
		"unit": "unit",
	}
	
	// Replace common abbreviations
	for abbr, replacement := range replacements {
		pattern := fmt.Sprintf(`(?i)(\d+)\s*%s\b`, abbr)
		dosage = strings.ReplaceAll(dosage, pattern, fmt.Sprintf("$1 %s", replacement))
	}
	
	return dosage
}

// standardizeFrequency normalizes medication frequencies
func (p *MedicationsProcessor) standardizeFrequency(frequency string) string {
	frequency = strings.TrimSpace(strings.ToLower(frequency))
	
	// Common abbreviations to standardize
	replacements := map[string]string{
		"qd":     "once daily",
		"od":     "once daily",
		"daily":  "once daily",
		"bid":    "twice daily",
		"tid":    "three times daily",
		"qid":    "four times daily",
		"q4h":    "every 4 hours",
		"q6h":    "every 6 hours",
		"q8h":    "every 8 hours",
		"q12h":   "every 12 hours",
		"weekly": "once weekly",
		"prn":    "as needed",
	}
	
	// Try direct replacement
	if replacement, ok := replacements[frequency]; ok {
		return replacement
	}
	
	// Try partial matches
	for abbr, replacement := range replacements {
		if strings.Contains(frequency, abbr) {
			return replacement
		}
	}
	
	return frequency
}

// standardizeRoute normalizes medication routes
func (p *MedicationsProcessor) standardizeRoute(route string) string {
	route = strings.TrimSpace(strings.ToLower(route))
	
	// Common abbreviations to standardize
	replacements := map[string]string{
		"po":       "oral",
		"oral":     "oral",
		"by mouth": "oral",
		"iv":       "intravenous",
		"ivp":      "intravenous",
		"im":       "intramuscular",
		"sc":       "subcutaneous",
		"sq":       "subcutaneous",
		"sl":       "sublingual",
		"top":      "topical",
		"inh":      "inhalation",
	}
	
	// Try direct replacement
	if replacement, ok := replacements[route]; ok {
		return replacement
	}
	
	// Try partial matches
	for abbr, replacement := range replacements {
		if strings.Contains(route, abbr) {
			return replacement
		}
	}
	
	return route
}
