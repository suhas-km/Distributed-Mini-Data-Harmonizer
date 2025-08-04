package worker

import (
	"context"
	"log"
	"sync"

	"github.com/suhas-km/distributed-mini-data-harmonizer/go-worker/internal/model"
)

// Pool represents a worker pool for processing harmonization jobs
type Pool struct {
	workers    int
	jobQueue   chan model.HarmonizationJob
	dispatcher *Dispatcher
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewPool creates a new worker pool
func NewPool(workers, queueSize int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Pool{
		workers:    workers,
		jobQueue:   make(chan model.HarmonizationJob, queueSize),
		dispatcher: NewDispatcher(),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start starts the worker pool
func (p *Pool) Start() {
	log.Printf("Starting worker pool with %d workers", p.workers)
	
	// Start workers
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.startWorker(i)
	}
}

// Stop stops the worker pool
func (p *Pool) Stop() {
	log.Println("Stopping worker pool")
	p.cancel()
	close(p.jobQueue)
	p.wg.Wait()
}

// Submit submits a job to the worker pool
func (p *Pool) Submit(job model.HarmonizationJob) {
	p.jobQueue <- job
}

// startWorker starts a worker
func (p *Pool) startWorker(id int) {
	defer p.wg.Done()
	
	log.Printf("Worker %d started", id)
	
	for {
		select {
		case job, ok := <-p.jobQueue:
			if !ok {
				log.Printf("Worker %d shutting down", id)
				return
			}
			
			log.Printf("Worker %d processing job %s", id, job.JobID)
			
			// Process job
			result, err := p.dispatcher.Process(job)
			if err != nil {
				log.Printf("Worker %d failed to process job %s: %v", id, job.JobID, err)
				job.Result <- model.JobResult{
					JobID:       job.JobID,
					Status:      model.JobStatusFailed,
					Error:       err.Error(),
					ProcessedAt: result.ProcessedAt,
				}
			} else {
				log.Printf("Worker %d completed job %s", id, job.JobID)
				job.Result <- result
			}
			
		case <-p.ctx.Done():
			log.Printf("Worker %d context cancelled", id)
			return
		}
	}
}
