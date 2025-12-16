package metadata

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/zulfikawr/photon-cli/internal/config"
	"github.com/zulfikawr/photon-cli/internal/pipeline"
)

// ProcessingRecord represents a single file's processing record
type ProcessingRecord struct {
	InputFile        string `json:"input_file"`
	OutputFile       string `json:"output_file"`
	FileType         string `json:"file_type"`
	OriginalSize     int64  `json:"original_size_bytes"`
	ProcessedSize    int64  `json:"processed_size_bytes"`
	BytesSaved       int64  `json:"bytes_saved"`
	CompressionRatio string `json:"compression_ratio"`
	Success          bool   `json:"success"`
	Error            string `json:"error,omitempty"`
}

// MetadataFile represents the complete metadata document
type MetadataFile struct {
	CreatedAt        time.Time          `json:"created_at"`
	ProcessingConfig ProcessingConfig   `json:"processing_config"`
	Summary          SummaryStats       `json:"summary"`
	ProcessedFiles   []ProcessingRecord `json:"processed_files"`
}

// ProcessingConfig stores the options used for processing
type ProcessingConfig struct {
	Quality     int    `json:"quality"`
	Width       int    `json:"width"`
	WebP        bool   `json:"webp"`
	Concurrency int    `json:"concurrency"`
	InputDir    string `json:"input_directory"`
	OutputDir   string `json:"output_directory"`
}

// SummaryStats stores aggregated statistics
type SummaryStats struct {
	TotalFiles       int     `json:"total_files"`
	SuccessfulFiles  int     `json:"successful_files"`
	FailedFiles      int     `json:"failed_files"`
	TotalBytesSaved  int64   `json:"total_bytes_saved"`
	TotalOriginalSize int64  `json:"total_original_size_bytes"`
	TotalProcessedSize int64 `json:"total_processed_size_bytes"`
	SuccessRate      float64 `json:"success_rate_percent"`
}

// Create generates a metadata file from processing results
func Create(inputDir string, outputDir string, opts config.Options, stats pipeline.PipelineStats) MetadataFile {
	records := make([]ProcessingRecord, 0, len(stats.ProcessedFiles))

	totalOriginal := int64(0)
	totalProcessed := int64(0)

	for _, result := range stats.ProcessedFiles {
		ratio := "N/A"
		if result.OriginalSize > 0 {
			ratio = formatPercentage(float64(result.BytesSaved) / float64(result.OriginalSize) * 100)
		}

		records = append(records, ProcessingRecord{
			InputFile:        result.FilePath,
			OutputFile:       result.OutputPath,
			FileType:         result.FileType,
			OriginalSize:     result.OriginalSize,
			ProcessedSize:    result.ProcessedSize,
			BytesSaved:       result.BytesSaved,
			CompressionRatio: ratio,
			Success:          result.Success,
			Error:            result.Error,
		})

		totalOriginal += result.OriginalSize
		totalProcessed += result.ProcessedSize
	}

	return MetadataFile{
		CreatedAt: time.Now(),
		ProcessingConfig: ProcessingConfig{
			Quality:     opts.Quality,
			Width:       opts.Width,
			WebP:        opts.WebP,
			Concurrency: opts.Concurrency,
			InputDir:    inputDir,
			OutputDir:   outputDir,
		},
		Summary: SummaryStats{
			TotalFiles:        stats.TotalFiles(),
			SuccessfulFiles:   stats.SuccessfulFiles,
			FailedFiles:       stats.FailedFiles,
			TotalBytesSaved:   stats.TotalBytesSaved,
			TotalOriginalSize: totalOriginal,
			TotalProcessedSize: totalProcessed,
			SuccessRate:       stats.SuccessRate(),
		},
		ProcessedFiles: records,
	}
}

// WriteToFile saves the metadata to a JSON file
func (m *MetadataFile) WriteToFile(filePath string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// formatPercentage formats a percentage value as a string
func formatPercentage(val float64) string {
	return fmt.Sprintf("%.1f%%", val)
}
