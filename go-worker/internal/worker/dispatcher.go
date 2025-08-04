package worker

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/suhas-km/distributed-mini-data-harmonizer/go-worker/internal/model"
	"github.com/suhas-km/distributed-mini-data-harmonizer/go-worker/internal/processor"
)

// Dispatcher routes jobs to the appropriate processor
type Dispatcher struct {
	processors map[string]processor.Processor
}

// NewDispatcher creates a new dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		processors: map[string]processor.Processor{
			"patients":     &processor.PatientProcessor{},
			"vitals":       &processor.VitalsProcessor{},
			"medications":  &processor.MedicationsProcessor{},
			"lab_results":  &processor.LabResultsProcessor{},
			"generic":      &processor.GenericProcessor{},
		},
	}
}

// Process processes a job using the appropriate processor
func (d *Dispatcher) Process(job model.HarmonizationJob) (model.JobResult, error) {
	// Get processor for job type
	proc, exists := d.processors[job.HarmonizationType]
	if !exists {
		return model.JobResult{}, fmt.Errorf("no processor found for harmonization type: %s", job.HarmonizationType)
	}

	// Generate output file path if not provided
	if job.OutputFile == "" {
		ext := filepath.Ext(job.InputFile)
		baseDir := filepath.Dir(job.InputFile)
		baseName := filepath.Base(job.InputFile)
		fileName := baseName[:len(baseName)-len(ext)]
		job.OutputFile = filepath.Join(baseDir, "..", "results", fmt.Sprintf("%s_harmonized%s", fileName, ext))
	}

	// Process the job
	err := proc.Process(job.InputFile, job.OutputFile)
	
	// Create result
	result := model.JobResult{
		JobID:       job.JobID,
		ProcessedAt: time.Now(),
	}
	
	if err != nil {
		result.Status = model.JobStatusFailed
		result.Error = err.Error()
		return result, err
	}
	
	result.Status = model.JobStatusCompleted
	result.OutputFile = job.OutputFile
	
	return result, nil
}
