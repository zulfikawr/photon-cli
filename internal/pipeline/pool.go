package pipeline

import (
	"os"
	"strings"
	"sync"

	"github.com/zulfikawr/bitrim/internal/config"
	"github.com/zulfikawr/bitrim/internal/optimizer"
)

// ProcessResult wraps the optimizer result with additional pipeline metadata
type ProcessResult struct {
	Result optimizer.Result
	// Additional fields can be added here as needed
}

// WorkerPool manages a pool of workers processing files concurrently
type WorkerPool struct {
	numWorkers int
	jobsCh     <-chan FileInfo
	resultsCh  chan<- ProcessResult
	opts       config.Options
	wg         sync.WaitGroup
	outputDir  string
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(numWorkers int, jobsCh <-chan FileInfo, resultsCh chan<- ProcessResult, opts config.Options, outputDir string) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		jobsCh:     jobsCh,
		resultsCh:  resultsCh,
		opts:       opts,
		outputDir:  outputDir,
	}
}

// Start spawns worker goroutines and begins processing jobs
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

// Wait blocks until all workers have finished processing
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

// worker processes jobs from the jobs channel
func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for job := range wp.jobsCh {
		var result optimizer.Result

		// Skip directory creation in dry-run mode
		if !wp.opts.DryRun {
			os.MkdirAll(wp.outputDir, 0755)
		}

		// Process based on file type
		if job.Type == "image" {
			result = optimizer.ProcessImage(job.Path, wp.outputDir, wp.opts, wp.opts.DryRun)
		} else if job.Type == "svg" {
			result = optimizer.ProcessSVG(job.Path, wp.outputDir, wp.opts.DryRun)
		}

		// Send result to results channel
		wp.resultsCh <- ProcessResult{
			Result: result,
		}
	}
}

// Coordinator manages the entire pipeline: walker, worker pool, and results collection
type Coordinator struct {
	inputDir  string
	outputDir string
	opts      config.Options
}

// NewCoordinator creates a new Coordinator
func NewCoordinator(inputDir string, outputDir string, opts config.Options) *Coordinator {
	return &Coordinator{
		inputDir:  inputDir,
		outputDir: outputDir,
		opts:      opts,
	}
}

// Run executes the full pipeline and returns aggregated results
func (c *Coordinator) Run() (PipelineStats, error) {
	stats := PipelineStats{
		ProcessedFiles: make([]optimizer.Result, 0),
	}

	// Create channels
	jobsCh := make(chan FileInfo, 100)          // Buffered channel for jobs
	resultsCh := make(chan ProcessResult, 100)  // Buffered channel for results

	// Parse ignore patterns
	var ignorePatterns []string
	if c.opts.IgnorePatterns != "" {
		ignorePatterns = strings.Split(c.opts.IgnorePatterns, ",")
	}

	// Create and start walker (producer)
	walker := NewWalker(c.inputDir, jobsCh, ignorePatterns, c.opts.MaxDepth, c.opts.MinSize)

	// Create and start worker pool (consumers)
	wp := NewWorkerPool(c.opts.Concurrency, jobsCh, resultsCh, c.opts, c.outputDir)
	wp.Start()

	// Start a goroutine to walk the directory
	go func() {
		walker.Walk()
		close(jobsCh) // Signal workers that no more jobs are coming
	}()

	// Collect results
	go func() {
		wp.Wait()
		close(resultsCh) // Signal that all results have been sent
	}()

	// Aggregate results
	for result := range resultsCh {
		if result.Result.Success {
			stats.SuccessfulFiles++
			stats.TotalBytesSaved += result.Result.BytesSaved
		} else {
			stats.FailedFiles++
		}
		stats.ProcessedFiles = append(stats.ProcessedFiles, result.Result)
	}

	return stats, nil
}
