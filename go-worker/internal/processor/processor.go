package processor

// Processor defines the interface for data harmonization processors
type Processor interface {
	// Process harmonizes data from inputFile and writes to outputFile
	Process(inputFile, outputFile string) error
}
