package worker

import (
	"fmt"
	"testing"
	"time"

	"github.com/suhas-km/distributed-mini-data-harmonizer/go-worker/internal/model"
)

func TestPoolBasicFunctionality(t *testing.T) {
	// Create a new pool with 2 workers and a queue size of 5
	pool := NewPool(2, 5)

	// Start the pool
	pool.Start()
	defer pool.Stop()

	// Create a test job
	job := model.HarmonizationJob{
		JobID:             "test-job-1",
		InputFile:         "test-input.csv",
		HarmonizationType: "generic",
		Result:            make(chan model.JobResult, 1),
	}

	// Submit the job
	pool.Submit(job)

	// Wait for the result with a timeout
	select {
	case result := <-job.Result:
		// We expect the job to fail since the input file doesn't exist
		if result.Status != model.JobStatusFailed {
			t.Errorf("Expected job status to be failed, got %s", result.Status)
		}
		if result.Error == "" {
			t.Error("Expected error message, got empty string")
		}
	case <-time.After(5 * time.Second):
		t.Error("Timed out waiting for job result")
	}
}

func TestPoolMultipleJobs(t *testing.T) {
	// Create a new pool with 2 workers and a queue size of 5
	pool := NewPool(2, 5)

	// Start the pool
	pool.Start()
	defer pool.Stop()

	// Create multiple test jobs
	jobCount := 5
	jobs := make([]model.HarmonizationJob, jobCount)

	for i := 0; i < jobCount; i++ {
		jobs[i] = model.HarmonizationJob{
			JobID:            fmt.Sprintf("test-job-%d", i+1),
			InputFile:         "test-input.csv",
			HarmonizationType: "generic",
			Result:            make(chan model.JobResult, 1),
		}

		// Submit the job
		pool.Submit(jobs[i])
	}

	// Wait for all results with a timeout
	for i := 0; i < jobCount; i++ {
		select {
		case result := <-jobs[i].Result:
			// We expect the job to fail since the input file doesn't exist
			if result.Status != model.JobStatusFailed {
				t.Errorf("Job %d: Expected job status to be failed, got %s", i, result.Status)
			}
		case <-time.After(5 * time.Second):
			t.Errorf("Job %d: Timed out waiting for job result", i)
		}
	}
}
