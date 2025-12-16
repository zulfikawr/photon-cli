package optimizer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProcessSVG(t *testing.T) {
	// Create a test SVG file
	testDir := t.TempDir()
	svgPath := filepath.Join(testDir, "test.svg")
	outputDir := filepath.Join(testDir, "output")

	// Create a simple SVG with extra whitespace
	svgContent := `<?xml version="1.0" encoding="UTF-8"?>
<!-- This is a comment -->
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100">
  <!-- Another comment -->
  <circle cx="50" cy="50" r="40" fill="blue"/>
</svg>
`

	if err := os.WriteFile(svgPath, []byte(svgContent), 0644); err != nil {
		t.Fatalf("failed to create test SVG: %v", err)
	}

	// Process the SVG
	result := ProcessSVG(svgPath, outputDir, false)

	if !result.Success {
		t.Fatalf("ProcessSVG failed: %s", result.Error)
	}

	if result.FileType != "svg" {
		t.Errorf("expected FileType 'svg', got '%s'", result.FileType)
	}

	if result.BytesSaved <= 0 {
		t.Errorf("expected positive BytesSaved, got %d", result.BytesSaved)
	}

	// Verify output file exists
	if _, err := os.Stat(result.OutputPath); err != nil {
		t.Fatalf("output file not found: %v", err)
	}
}

func TestMinifySVG(t *testing.T) {
	input := `<?xml version="1.0"?>
<!-- comment -->
<svg>
  <circle/>
</svg>
`
	output := minifySVG([]byte(input))
	outputStr := string(output)

	if len(output) >= len(input) {
		t.Error("minified SVG should be smaller than original")
	}

	if string(output) != outputStr {
		t.Error("minified SVG content mismatch")
	}
}
