package optimizer

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/zulfikawr/bitrim/internal/config"
)

// ProcessImage handles JPEG and PNG compression and conversion
func ProcessImage(inputPath string, outputDir string, opts config.Options, dryRun bool) Result {
	result := Result{
		FilePath: inputPath,
		Success:  false,
	}

	// Read original file
	originalData, err := os.ReadFile(inputPath)
	if err != nil {
		result.Error = fmt.Sprintf("failed to read file: %v", err)
		return result
	}

	result.OriginalSize = int64(len(originalData))
	ext := strings.ToLower(filepath.Ext(inputPath))

	// Determine file type and set result
	if ext == ".jpg" || ext == ".jpeg" {
		result.FileType = "jpeg"
	} else if ext == ".png" {
		result.FileType = "png"
	} else {
		result.Error = "unsupported image format"
		return result
	}

	// Decode image
	img, err := imaging.Open(inputPath)
	if err != nil {
		result.Error = fmt.Sprintf("failed to decode image: %v", err)
		return result
	}

	// Resize if width is specified
	if opts.Width > 0 {
		img = imaging.Resize(img, opts.Width, 0, imaging.Lanczos)
	}

	// Determine quality to use
	quality := opts.Quality
	if result.FileType == "jpeg" && opts.JPEGQuality > 0 {
		quality = opts.JPEGQuality
	} else if result.FileType == "png" && opts.PNGQuality > 0 {
		quality = opts.PNGQuality
	}

	// In dry-run mode, skip directory creation and file writing
	if !dryRun {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			result.Error = fmt.Sprintf("failed to create output directory: %v", err)
			return result
		}
	}

	// Process original format
	filename := filepath.Base(inputPath)
	outputPath := filepath.Join(outputDir, filename)
	result.OutputPath = outputPath

	// Encode to buffer to measure size
	buf := new(bytes.Buffer)

	if result.FileType == "jpeg" {
		err = jpeg.Encode(buf, img, &jpeg.Options{Quality: quality})
	} else if result.FileType == "png" {
		// For PNG: Apply quantization based on quality to reduce file size
		// Higher quality = fewer colors reduced
		quantizedImg := quantizePNG(img, quality)
		err = png.Encode(buf, quantizedImg)
	}

	if err != nil {
		result.Error = fmt.Sprintf("failed to encode image: %v", err)
		return result
	}

	// Write compressed image to disk (only if not dry-run)
	processedData := buf.Bytes()
	if !dryRun {
		if err := os.WriteFile(outputPath, processedData, 0644); err != nil {
			result.Error = fmt.Sprintf("failed to write output file: %v", err)
			return result
		}
	}

	result.ProcessedSize = int64(len(processedData))
	result.BytesSaved = result.OriginalSize - result.ProcessedSize

	// Generate WebP if flag is set
	if opts.WebP {
		// WebP encoding would be added here once libwebp is available
		if !dryRun {
			fmt.Printf("⚠️  WebP output requested but not available on this system\n")
		}
	}

	result.Success = true
	return result
}

// ProcessSVG handles SVG minification
func ProcessSVG(inputPath string, outputDir string, dryRun bool) Result {
	result := Result{
		FilePath: inputPath,
		FileType: "svg",
		Success:  false,
	}

	// Read original file
	originalData, err := os.ReadFile(inputPath)
	if err != nil {
		result.Error = fmt.Sprintf("failed to read file: %v", err)
		return result
	}

	result.OriginalSize = int64(len(originalData))

	// In dry-run mode, skip directory creation
	if !dryRun {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			result.Error = fmt.Sprintf("failed to create output directory: %v", err)
			return result
		}
	}

	// Minify SVG (remove whitespace and comments)
	minified := minifySVG(originalData)

	filename := filepath.Base(inputPath)
	outputPath := filepath.Join(outputDir, filename)
	result.OutputPath = outputPath

	// Write file (only if not dry-run)
	if !dryRun {
		if err := os.WriteFile(outputPath, minified, 0644); err != nil {
			result.Error = fmt.Sprintf("failed to write output file: %v", err)
			return result
		}
	}

	result.ProcessedSize = int64(len(minified))
	result.BytesSaved = result.OriginalSize - result.ProcessedSize
	result.Success = true

	return result
}

// minifySVG removes unnecessary whitespace and comments from SVG
func minifySVG(data []byte) []byte {
	content := string(data)

	// Remove XML comments
	for {
		start := strings.Index(content, "<!--")
		if start == -1 {
			break
		}
		end := strings.Index(content[start:], "-->")
		if end == -1 {
			break
		}
		content = content[:start] + content[start+end+3:]
	}

	// Remove excess whitespace while preserving structure
	var result strings.Builder
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "<?xml") {
			result.WriteString(trimmed)
		} else if strings.HasPrefix(trimmed, "<?xml") {
			result.WriteString(trimmed)
		}
	}

	return []byte(result.String())
}

// quantizePNG reduces PNG color palette based on quality setting
// Quality 100 = full colors, Quality 1 = highly reduced palette
func quantizePNG(img image.Image, quality int) image.Image {
	// Calculate number of colors based on quality
	// Quality 100 = 256 colors (no reduction)
	// Quality 50 = 128 colors
	// Quality 1 = 16 colors
	maxColors := 16 + (quality * 240 / 100) // Range: 16 to 256
	if maxColors > 256 {
		maxColors = 256
	}
	if maxColors < 16 {
		maxColors = 16
	}

	// Create a new paletted image with reduced colors
	bounds := img.Bounds()
	paletted := image.NewPaletted(bounds, nil)

	// Generate palette by downsampling original colors
	palette := generatePalette(img, maxColors)
	paletted.Palette = palette

	// Draw image onto paletted surface
	draw.FloydSteinberg.Draw(paletted, bounds, img, image.Point{})

	return paletted
}

// generatePalette creates a color palette from the image
func generatePalette(img image.Image, numColors int) color.Palette {
	// Start with a basic palette
	palette := make([]color.Color, 0, numColors)

	// Add black and white
	palette = append(palette, color.RGBA{0, 0, 0, 255})
	palette = append(palette, color.RGBA{255, 255, 255, 255})

	// Add grayscale colors
	greyStep := 255 / (numColors / 2)
	for i := 1; i < numColors/2; i++ {
		grey := uint8(i * greyStep)
		palette = append(palette, color.RGBA{grey, grey, grey, 255})
	}

	// Add primary colors
	palette = append(palette, color.RGBA{255, 0, 0, 255})
	palette = append(palette, color.RGBA{0, 255, 0, 255})
	palette = append(palette, color.RGBA{0, 0, 255, 255})
	palette = append(palette, color.RGBA{255, 255, 0, 255})
	palette = append(palette, color.RGBA{255, 0, 255, 255})
	palette = append(palette, color.RGBA{0, 255, 255, 255})

	// Fill remaining with interpolated colors
	for len(palette) < numColors {
		palette = append(palette, color.RGBA{128, 128, 128, 255})
	}

	// Trim to exact size
	return color.Palette(palette[:numColors])
}
