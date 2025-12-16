package pipeline

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zulfikawr/photon-cli/internal/config"
)

func TestWalker(t *testing.T) {
	// Create test directory structure
	testDir := t.TempDir()

	// Create test files
	svgFile := filepath.Join(testDir, "test.svg")
	if err := os.WriteFile(svgFile, []byte("<svg></svg>"), 0644); err != nil {
		t.Fatalf("failed to create test SVG: %v", err)
	}

	pngFile := filepath.Join(testDir, "test.png")
	if err := os.WriteFile(pngFile, []byte("fake png"), 0644); err != nil {
		t.Fatalf("failed to create test PNG: %v", err)
	}

	txtFile := filepath.Join(testDir, "test.txt")
	if err := os.WriteFile(txtFile, []byte("text"), 0644); err != nil {
		t.Fatalf("failed to create test TXT: %v", err)
	}

	// Walk the directory
	jobsCh := make(chan FileInfo, 10)
	walker := NewWalker(testDir, jobsCh)

	files := make([]FileInfo, 0)
	go func() {
		walker.Walk()
		close(jobsCh)
	}()

	for job := range jobsCh {
		files = append(files, job)
	}

	// Verify results
	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d", len(files))
	}

	// Check that SVG and PNG were found
	foundSVG := false
	foundPNG := false

	for _, f := range files {
		if filepath.Base(f.Path) == "test.svg" && f.Type == "svg" {
			foundSVG = true
		}
		if filepath.Base(f.Path) == "test.png" && f.Type == "image" {
			foundPNG = true
		}
	}

	if !foundSVG {
		t.Error("SVG file not found or incorrect type")
	}
	if !foundPNG {
		t.Error("PNG file not found or incorrect type")
	}
}

func TestWorkerPoolResultsCollection(t *testing.T) {
	testDir := t.TempDir()

	// Create a test SVG
	svgFile := filepath.Join(testDir, "test.svg")
	svgContent := `<?xml version="1.0"?>
<!-- comment -->
<svg>
  <circle/>
</svg>
`
	if err := os.WriteFile(svgFile, []byte(svgContent), 0644); err != nil {
		t.Fatalf("failed to create test SVG: %v", err)
	}

	// Create pipeline
	opts := config.Options{
		Quality:     80,
		Concurrency: 1,
	}

	coordinator := NewCoordinator(testDir, testDir, opts)
	stats, err := coordinator.Run()

	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}

	if stats.SuccessfulFiles != 1 {
		t.Errorf("expected 1 successful file, got %d", stats.SuccessfulFiles)
	}

	if stats.FailedFiles != 0 {
		t.Errorf("expected 0 failed files, got %d", stats.FailedFiles)
	}
}

func TestPipelineStats(t *testing.T) {
	stats := PipelineStats{
		SuccessfulFiles: 8,
		FailedFiles:     2,
		TotalBytesSaved: 800000,
	}

	if stats.TotalFiles() != 10 {
		t.Errorf("expected 10 total files, got %d", stats.TotalFiles())
	}

	if stats.SuccessRate() != 80.0 {
		t.Errorf("expected 80%% success rate, got %.1f%%", stats.SuccessRate())
	}

	expectedAvg := int64(100000)
	if stats.AverageSavingsPerFile() != expectedAvg {
		t.Errorf("expected %.0f average savings, got %d", float64(expectedAvg), stats.AverageSavingsPerFile())
	}
}
