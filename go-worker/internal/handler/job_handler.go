package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/suhas-km/distributed-mini-data-harmonizer/go-worker/internal/config"
	"github.com/suhas-km/distributed-mini-data-harmonizer/go-worker/internal/model"
	"github.com/suhas-km/distributed-mini-data-harmonizer/go-worker/internal/worker"
)

// JobHandler handles job-related HTTP requests
type JobHandler struct {
	config *config.Config
	pool   *worker.Pool
}

// NewJobHandler creates a new job handler
func NewJobHandler(config *config.Config, pool *worker.Pool) *JobHandler {
	return &JobHandler{
		config: config,
		pool:   pool,
	}
}

// RegisterRoutes registers HTTP routes
func (h *JobHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/jobs", h.handleJobs)
	mux.HandleFunc("/api/v1/jobs/", h.handleJobByID)
	mux.HandleFunc("/health", h.handleHealth)
}

// handleJobs handles POST /api/v1/jobs
func (h *JobHandler) handleJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req model.JobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	// Validate request
	if req.JobID == "" {
		http.Error(w, "Job ID is required", http.StatusBadRequest)
		return
	}
	if req.InputFile == "" {
		http.Error(w, "Input file is required", http.StatusBadRequest)
		return
	}
	if req.HarmonizationType == "" {
		http.Error(w, "Harmonization type is required", http.StatusBadRequest)
		return
	}

	// Create job
	job := model.HarmonizationJob{
		JobID:            req.JobID,
		InputFile:        req.InputFile,
		HarmonizationType: req.HarmonizationType,
		Result:           make(chan model.JobResult, 1),
	}

	// Submit job to worker pool
	log.Printf("Submitting job %s to worker pool", req.JobID)
	h.pool.Submit(job)

	// Return accepted response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "accepted",
		"job_id": req.JobID,
	})

	// Process result asynchronously
	go h.processJobResult(job)
}

// handleJobByID handles GET /api/v1/jobs/{id}
func (h *JobHandler) handleJobByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract job ID from URL
	jobID := r.URL.Path[len("/api/v1/jobs/"):]
	if jobID == "" {
		http.Error(w, "Job ID is required", http.StatusBadRequest)
		return
	}

	// This endpoint is just a placeholder for now
	// In a real implementation, we would check the job status in a database
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"job_id": jobID,
		"status": "unknown",
	})
}

// handleHealth handles GET /health
func (h *JobHandler) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "0.1.0",
	})
}

// processJobResult processes the result of a job and sends it back to the Python API
func (h *JobHandler) processJobResult(job model.HarmonizationJob) {
	// Wait for job result
	result := <-job.Result
	log.Printf("Job %s completed with status %s", job.JobID, result.Status)

	// Create status update
	update := model.JobStatusUpdate{
		Status:      result.Status,
		OutputFile:  result.OutputFile,
		Error:       result.Error,
		CompletedAt: result.ProcessedAt,
	}

	// Send status update to Python API
	url := fmt.Sprintf("%s/jobs/%s/status", h.config.PythonAPIURL, job.JobID)
	log.Printf("Sending status update for job %s to %s", job.JobID, url)

	// Convert update to JSON
	updateJSON, err := json.Marshal(update)
	if err != nil {
		log.Printf("Failed to marshal status update: %v", err)
		return
	}

	// Send HTTP request
	// In a real implementation, we would add retries and error handling
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(updateJSON))
	if err != nil {
		log.Printf("Failed to send status update: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Status update failed with status code %d", resp.StatusCode)
		return
	}

	log.Printf("Status update for job %s sent successfully", job.JobID)
}
