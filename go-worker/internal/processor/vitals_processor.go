package processor

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// VitalsProcessor processes vital signs data
type VitalsProcessor struct {
	CSVProcessor
}

// Process harmonizes vital signs data
func (p *VitalsProcessor) Process(inputFile, outputFile string) error {
	// Read CSV
	records, err := p.ReadCSV(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read vitals data: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("vitals data is empty or missing header")
	}

	// Get header and standardize it
	header := records[0]
	standardizedHeader := p.StandardizeHeader(header)

	// Find required columns
	patientIDIdx, err := p.getColumnIndex(standardizedHeader, "patient_id")
	if err != nil {
		return err
	}

	dateTimeIdx, err := p.getColumnIndex(standardizedHeader, "datetime")
	if err != nil {
		return err
	}

	// Find vital sign columns
	bpSystolicIdx, _ := p.getColumnIndex(standardizedHeader, "systolic")
	bpDiastolicIdx, _ := p.getColumnIndex(standardizedHeader, "diastolic")
	heartRateIdx, _ := p.getColumnIndex(standardizedHeader, "heart_rate")
	temperatureIdx, _ := p.getColumnIndex(standardizedHeader, "temperature")
	respiratoryRateIdx, _ := p.getColumnIndex(standardizedHeader, "respiratory_rate")
	oxygenSaturationIdx, _ := p.getColumnIndex(standardizedHeader, "oxygen_saturation")

	// Process data
	for i := 1; i < len(records); i++ {
		row := records[i]
		
		// Standardize patient ID
		if patientIDIdx < len(row) {
			row[patientIDIdx] = p.standardizePatientID(row[patientIDIdx])
		}
		
		// Standardize datetime
		if dateTimeIdx < len(row) {
			row[dateTimeIdx] = p.standardizeDateTime(row[dateTimeIdx])
		}
		
		// Standardize blood pressure (systolic)
		if bpSystolicIdx >= 0 && bpSystolicIdx < len(row) {
			row[bpSystolicIdx] = p.standardizeNumeric(row[bpSystolicIdx])
		}
		
		// Standardize blood pressure (diastolic)
		if bpDiastolicIdx >= 0 && bpDiastolicIdx < len(row) {
			row[bpDiastolicIdx] = p.standardizeNumeric(row[bpDiastolicIdx])
		}
		
		// Standardize heart rate
		if heartRateIdx >= 0 && heartRateIdx < len(row) {
			row[heartRateIdx] = p.standardizeNumeric(row[heartRateIdx])
		}
		
		// Standardize temperature
		if temperatureIdx >= 0 && temperatureIdx < len(row) {
			row[temperatureIdx] = p.standardizeTemperature(row[temperatureIdx])
		}
		
		// Standardize respiratory rate
		if respiratoryRateIdx >= 0 && respiratoryRateIdx < len(row) {
			row[respiratoryRateIdx] = p.standardizeNumeric(row[respiratoryRateIdx])
		}
		
		// Standardize oxygen saturation
		if oxygenSaturationIdx >= 0 && oxygenSaturationIdx < len(row) {
			row[oxygenSaturationIdx] = p.standardizePercentage(row[oxygenSaturationIdx])
		}
		
		records[i] = row
	}

	// Write harmonized data
	if err := p.WriteCSV(outputFile, records); err != nil {
		return fmt.Errorf("failed to write harmonized vitals data: %w", err)
	}

	return nil
}

// Helper function to find column index with fallback names
func (p *VitalsProcessor) getColumnIndex(header []string, columnName string) (int, error) {
	// Map of column name aliases
	columnAliases := map[string][]string{
		"patient_id": {"patient_id", "patientid", "id", "patient_number", "mrn"},
		"datetime":   {"datetime", "date_time", "timestamp", "recorded_at", "measurement_time", "date"},
		"systolic":   {"systolic", "bp_systolic", "sbp", "systolic_bp"},
		"diastolic":  {"diastolic", "bp_diastolic", "dbp", "diastolic_bp"},
		"heart_rate": {"heart_rate", "hr", "pulse"},
		"temperature": {"temperature", "temp", "body_temperature"},
		"respiratory_rate": {"respiratory_rate", "resp_rate", "rr"},
		"oxygen_saturation": {"oxygen_saturation", "o2_sat", "spo2", "oxygen_sat"},
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
func (p *VitalsProcessor) standardizePatientID(patientID string) string {
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

// standardizeDateTime normalizes datetime values to ISO format
func (p *VitalsProcessor) standardizeDateTime(dateTime string) string {
	// Try different datetime formats
	formats := []string{
		"2006-01-02 15:04:05",     // YYYY-MM-DD HH:MM:SS
		"2006-01-02T15:04:05",     // YYYY-MM-DDThh:mm:ss
		"01/02/2006 15:04:05",     // MM/DD/YYYY HH:MM:SS
		"02/01/2006 15:04:05",     // DD/MM/YYYY HH:MM:SS
		"2006-01-02",              // YYYY-MM-DD (date only)
		"01/02/2006",              // MM/DD/YYYY (date only)
		"02/01/2006",              // DD/MM/YYYY (date only)
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateTime); err == nil {
			return t.Format("2006-01-02T15:04:05")
		}
	}

	// If only date is provided, append midnight time
	if t, err := time.Parse("2006-01-02", p.StandardizeDate(dateTime)); err == nil {
		return t.Format("2006-01-02T00:00:00")
	}

	// Return original if no format matches
	return dateTime
}

// standardizeNumeric normalizes numeric values
func (p *VitalsProcessor) standardizeNumeric(value string) string {
	value = strings.TrimSpace(value)
	
	// If empty, return empty
	if value == "" {
		return ""
	}
	
	// Try to parse as float
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		// Format with 1 decimal place
		return strconv.FormatFloat(f, 'f', 1, 64)
	}
	
	// Remove any non-numeric characters except decimal point
	var result strings.Builder
	hasDecimal := false
	
	for _, c := range value {
		if c >= '0' && c <= '9' {
			result.WriteRune(c)
		} else if c == '.' && !hasDecimal {
			result.WriteRune(c)
			hasDecimal = true
		}
	}
	
	// If we managed to extract a number, parse and format it
	if result.Len() > 0 {
		if f, err := strconv.ParseFloat(result.String(), 64); err == nil {
			return strconv.FormatFloat(f, 'f', 1, 64)
		}
	}
	
	return value
}

// standardizeTemperature normalizes temperature values to Celsius
func (p *VitalsProcessor) standardizeTemperature(temp string) string {
	temp = p.standardizeNumeric(temp)
	
	// If empty or failed to parse, return as is
	if temp == "" {
		return temp
	}
	
	// Parse as float
	f, err := strconv.ParseFloat(temp, 64)
	if err != nil {
		return temp
	}
	
	// Convert Fahrenheit to Celsius if it appears to be Fahrenheit
	if f > 45 {
		celsius := (f - 32) * 5 / 9
		return strconv.FormatFloat(celsius, 'f', 1, 64)
	}
	
	return temp
}

// standardizePercentage normalizes percentage values
func (p *VitalsProcessor) standardizePercentage(value string) string {
	value = strings.TrimSpace(value)
	
	// If empty, return empty
	if value == "" {
		return ""
	}
	
	// Remove % sign if present
	value = strings.TrimSuffix(value, "%")
	
	// Try to parse as float
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		// Ensure value is between 0 and 100
		if f > 1 && f <= 100 {
			return strconv.FormatFloat(f, 'f', 1, 64)
		} else if f >= 0 && f <= 1 {
			// Convert from decimal to percentage
			return strconv.FormatFloat(f*100, 'f', 1, 64)
		}
	}
	
	return value
}
