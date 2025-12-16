package pipeline

import "github.com/zulfikawr/bitrim/internal/optimizer"

// PipelineStats aggregates statistics from a pipeline run
type PipelineStats struct {
	// Total number of files successfully processed
	SuccessfulFiles int

	// Total number of files that failed
	FailedFiles int

	// Total bytes saved across all files
	TotalBytesSaved int64

	// Individual results for each file
	ProcessedFiles []optimizer.Result
}

// TotalFiles returns the total number of files processed
func (ps *PipelineStats) TotalFiles() int {
	return ps.SuccessfulFiles + ps.FailedFiles
}

// AverageSavingsPerFile returns the average bytes saved per file
func (ps *PipelineStats) AverageSavingsPerFile() int64 {
	if ps.SuccessfulFiles == 0 {
		return 0
	}
	return ps.TotalBytesSaved / int64(ps.SuccessfulFiles)
}

// SuccessRate returns the percentage of successful files
func (ps *PipelineStats) SuccessRate() float64 {
	total := ps.TotalFiles()
	if total == 0 {
		return 0
	}
	return float64(ps.SuccessfulFiles) / float64(total) * 100
}
