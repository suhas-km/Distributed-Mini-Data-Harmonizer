package model

import (
	"time"
)

// JobStatus represents the status of a job
type JobStatus string

// Job status constants
const (
	JobStatusQueued     JobStatus = "queued"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

// JobRequest represents a job request from the Python API
type JobRequest struct {
	JobID            string `json:"job_id"`
	InputFile        string `json:"input_file"`
	HarmonizationType string `json:"harmonization_type"`
}

// JobResult represents the result of a job
type JobResult struct {
	JobID       string    `json:"job_id"`
	Status      JobStatus `json:"status"`
	OutputFile  string    `json:"output_file,omitempty"`
	Error       string    `json:"error,omitempty"`
	ProcessedAt time.Time `json:"processed_at"`
}

// JobStatusUpdate represents a job status update to be sent to the Python API
type JobStatusUpdate struct {
	Status      JobStatus `json:"status"`
	OutputFile  string    `json:"output_file,omitempty"`
	Error       string    `json:"error,omitempty"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

// HarmonizationJob represents a job to be processed by a worker
type HarmonizationJob struct {
	JobID            string
	InputFile        string
	OutputFile       string
	HarmonizationType string
	Result           chan JobResult
}
