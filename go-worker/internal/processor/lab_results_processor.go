package processor

import (
	"fmt"
	"strconv"
	"strings"
)

// LabResultsProcessor processes laboratory results data
type LabResultsProcessor struct {
	CSVProcessor
}

// Process harmonizes laboratory results data
func (p *LabResultsProcessor) Process(inputFile, outputFile string) error {
	// Read CSV
	records, err := p.ReadCSV(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read lab results data: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("lab results data is empty or missing header")
	}

	// Get header and standardize it
	header := records[0]
	standardizedHeader := p.StandardizeHeader(header)

	// Find required columns
	patientIDIdx, err := p.getColumnIndex(standardizedHeader, "patient_id")
	if err != nil {
		return err
	}

	testNameIdx, err := p.getColumnIndex(standardizedHeader, "test_name")
	if err != nil {
		return err
	}

	// Find other important columns
	testDateIdx, _ := p.getColumnIndex(standardizedHeader, "test_date")
	resultValueIdx, _ := p.getColumnIndex(standardizedHeader, "result_value")
	unitsIdx, _ := p.getColumnIndex(standardizedHeader, "units")
	referenceRangeIdx, _ := p.getColumnIndex(standardizedHeader, "reference_range")
	abnormalFlagIdx, _ := p.getColumnIndex(standardizedHeader, "abnormal_flag")

	// Process data
	for i := 1; i < len(records); i++ {
		row := records[i]
		
		// Standardize patient ID
		if patientIDIdx < len(row) {
			row[patientIDIdx] = p.standardizePatientID(row[patientIDIdx])
		}
		
		// Standardize test name
		if testNameIdx < len(row) {
			row[testNameIdx] = p.standardizeTestName(row[testNameIdx])
		}
		
		// Standardize test date
		if testDateIdx >= 0 && testDateIdx < len(row) {
			row[testDateIdx] = p.StandardizeDate(row[testDateIdx])
		}
		
		// Standardize result value
		if resultValueIdx >= 0 && resultValueIdx < len(row) && unitsIdx >= 0 && unitsIdx < len(row) {
			row[resultValueIdx] = p.standardizeResultValue(row[resultValueIdx], row[unitsIdx])
		}
		
		// Standardize units
		if unitsIdx >= 0 && unitsIdx < len(row) {
			row[unitsIdx] = p.standardizeUnits(row[unitsIdx])
		}
		
		// Standardize reference range
		if referenceRangeIdx >= 0 && referenceRangeIdx < len(row) {
			row[referenceRangeIdx] = p.standardizeReferenceRange(row[referenceRangeIdx])
		}
		
		// Standardize abnormal flag
		if abnormalFlagIdx >= 0 && abnormalFlagIdx < len(row) {
			row[abnormalFlagIdx] = p.standardizeAbnormalFlag(row[abnormalFlagIdx])
		}
		
		records[i] = row
	}

	// Write harmonized data
	if err := p.WriteCSV(outputFile, records); err != nil {
		return fmt.Errorf("failed to write harmonized lab results data: %w", err)
	}

	return nil
}

// Helper function to find column index with fallback names
func (p *LabResultsProcessor) getColumnIndex(header []string, columnName string) (int, error) {
	// Map of column name aliases
	columnAliases := map[string][]string{
		"patient_id":      {"patient_id", "patientid", "id", "patient_number", "mrn"},
		"test_name":       {"test_name", "lab_test", "test", "procedure_name", "component"},
		"test_date":       {"test_date", "date", "collection_date", "result_date", "date_performed"},
		"result_value":    {"result_value", "value", "result", "numeric_result", "observation_value"},
		"units":           {"units", "unit", "unit_of_measure", "uom"},
		"reference_range": {"reference_range", "ref_range", "normal_range", "normal_values", "reference_interval"},
		"abnormal_flag":   {"abnormal_flag", "flag", "abnormal", "result_flag", "status"},
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
func (p *LabResultsProcessor) standardizePatientID(patientID string) string {
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

// standardizeTestName normalizes lab test names
func (p *LabResultsProcessor) standardizeTestName(name string) string {
	name = strings.TrimSpace(name)
	
	// Map of common test name variations
	testNameMap := map[string]string{
		"wbc": "White Blood Cell Count",
		"rbc": "Red Blood Cell Count",
		"hgb": "Hemoglobin",
		"hb":  "Hemoglobin",
		"hct": "Hematocrit",
		"plt": "Platelet Count",
		"gluc": "Glucose",
		"bun": "Blood Urea Nitrogen",
		"crea": "Creatinine",
		"na": "Sodium",
		"k": "Potassium",
		"cl": "Chloride",
		"co2": "Carbon Dioxide",
		"ca": "Calcium",
		"phos": "Phosphorus",
		"mg": "Magnesium",
		"ast": "Aspartate Aminotransferase",
		"alt": "Alanine Aminotransferase",
		"alp": "Alkaline Phosphatase",
		"tbil": "Total Bilirubin",
		"dbil": "Direct Bilirubin",
		"tprot": "Total Protein",
		"alb": "Albumin",
		"a1c": "Hemoglobin A1C",
		"tsh": "Thyroid Stimulating Hormone",
		"ft4": "Free Thyroxine",
		"chol": "Total Cholesterol",
		"trig": "Triglycerides",
		"hdl": "HDL Cholesterol",
		"ldl": "LDL Cholesterol",
	}
	
	// Check for exact match (case-insensitive)
	nameLower := strings.ToLower(name)
	if standardName, ok := testNameMap[nameLower]; ok {
		return standardName
	}
	
	// Check for partial match
	for abbr, standardName := range testNameMap {
		if strings.Contains(nameLower, abbr) {
			return standardName
		}
	}
	
	// If no match found, capitalize first letter of each word
	words := strings.Fields(strings.ToLower(name))
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	
	return strings.Join(words, " ")
}

// standardizeResultValue normalizes lab result values
func (p *LabResultsProcessor) standardizeResultValue(value, unit string) string {
	value = strings.TrimSpace(value)
	
	// If empty, return empty
	if value == "" {
		return ""
	}
	
	// Handle non-numeric results
	if strings.EqualFold(value, "positive") || strings.EqualFold(value, "pos") {
		return "Positive"
	}
	if strings.EqualFold(value, "negative") || strings.EqualFold(value, "neg") {
		return "Negative"
	}
	if strings.EqualFold(value, "normal") || strings.EqualFold(value, "nrm") {
		return "Normal"
	}
	if strings.EqualFold(value, "abnormal") || strings.EqualFold(value, "abn") {
		return "Abnormal"
	}
	
	// Try to parse as float
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		// Format with appropriate precision based on unit
		switch strings.ToLower(unit) {
		case "g/dl", "g/l", "mg/dl", "mg/l":
			return strconv.FormatFloat(f, 'f', 1, 64)
		case "%":
			return strconv.FormatFloat(f, 'f', 1, 64)
		default:
			return strconv.FormatFloat(f, 'f', 2, 64)
		}
	}
	
	return value
}

// standardizeUnits normalizes lab result units
func (p *LabResultsProcessor) standardizeUnits(unit string) string {
	unit = strings.TrimSpace(unit)
	
	// Map of common unit variations
	unitMap := map[string]string{
		"g/dl": "g/dL",
		"g/l": "g/L",
		"mg/dl": "mg/dL",
		"mg/l": "mg/L",
		"mmol/l": "mmol/L",
		"umol/l": "μmol/L",
		"u/l": "U/L",
		"iu/l": "IU/L",
		"meq/l": "mEq/L",
		"ng/ml": "ng/mL",
		"pg/ml": "pg/mL",
		"k/ul": "K/μL",
		"m/ul": "M/μL",
		"thou/ul": "K/μL",
		"mill/ul": "M/μL",
	}
	
	// Check for exact match (case-insensitive)
	unitLower := strings.ToLower(unit)
	if standardUnit, ok := unitMap[unitLower]; ok {
		return standardUnit
	}
	
	return unit
}

// standardizeReferenceRange normalizes reference ranges
func (p *LabResultsProcessor) standardizeReferenceRange(range_ string) string {
	range_ = strings.TrimSpace(range_)
	
	// If empty, return empty
	if range_ == "" {
		return ""
	}
	
	// Replace common separators with a dash
	range_ = strings.ReplaceAll(range_, " to ", "-")
	range_ = strings.ReplaceAll(range_, "to", "-")
	range_ = strings.ReplaceAll(range_, "~", "-")
	
	return range_
}

// standardizeAbnormalFlag normalizes abnormal flags
func (p *LabResultsProcessor) standardizeAbnormalFlag(flag string) string {
	flag = strings.TrimSpace(strings.ToUpper(flag))
	
	switch flag {
	case "H", "HIGH", "ELEVATED", "ABOVE NORMAL":
		return "H"
	case "L", "LOW", "DECREASED", "BELOW NORMAL":
		return "L"
	case "N", "NORMAL", "WNL", "WITHIN NORMAL LIMITS":
		return "N"
	case "A", "ABNORMAL", "ABN":
		return "A"
	case "C", "CRITICAL", "CRIT", "PANIC":
		return "C"
	default:
		return flag
	}
}
