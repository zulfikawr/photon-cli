package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/zulfikawr/photon-cli/internal/config"
	"github.com/zulfikawr/photon-cli/internal/metadata"
	"github.com/zulfikawr/photon-cli/internal/pipeline"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "photon [flags] <input-directory>",
	Short: "Photon CLI - High-concurrency asset optimizer",
	Long: `Photon CLI is a powerful cross-platform asset optimizer that:
  - Recursively scans directories
  - Compresses images (JPG, PNG)
  - Converts images to WebP format
  - Minifies SVG files
  
Uses a worker pool pattern for maximum CPU concurrency.`,
	Args: cobra.ExactArgs(1),
	RunE: runOptimizer,
}

// Options for the optimizer
var opts config.Options

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Define flags
	rootCmd.Flags().StringVarP(
		&opts.Output,
		"out", "o",
		"",
		"Output directory (default: photon-output)",
	)

	rootCmd.Flags().IntVarP(
		&opts.Quality,
		"quality", "q",
		80,
		"JPEG/PNG/WebP quality (1-100, overridden by format-specific flags)",
	)

	rootCmd.Flags().IntVar(
		&opts.JPEGQuality,
		"jpeg-quality",
		0,
		"JPEG-specific quality (1-100, overrides --quality)",
	)

	rootCmd.Flags().IntVar(
		&opts.PNGQuality,
		"png-quality",
		0,
		"PNG-specific quality (1-100, overrides --quality)",
	)

	rootCmd.Flags().IntVarP(
		&opts.Width,
		"width", "w",
		0,
		"Resize width in pixels (0 = no resize)",
	)

	rootCmd.Flags().BoolVar(
		&opts.WebP,
		"webp",
		false,
		"Generate WebP copies alongside originals",
	)

	rootCmd.Flags().IntVar(
		&opts.Concurrency,
		"concurrency",
		runtime.NumCPU(),
		"Number of concurrent workers",
	)

	rootCmd.Flags().BoolVarP(
		&opts.Replace,
		"replace", "r",
		false,
		"Replace original files with compressed versions (requires confirmation)",
	)

	rootCmd.Flags().BoolVar(
		&opts.DryRun,
		"dry-run",
		false,
		"Calculate savings without writing files",
	)

	rootCmd.Flags().Int64Var(
		&opts.MinSize,
		"min-size",
		0,
		"Minimum file size to process (e.g., 10kb, 1mb). Smaller files are skipped",
	)

	rootCmd.Flags().IntVar(
		&opts.MaxDepth,
		"depth",
		0,
		"Maximum recursion depth (0 = unlimited)",
	)

	rootCmd.Flags().StringVar(
		&opts.IgnorePatterns,
		"ignore",
		"",
		"Comma-separated patterns to ignore (e.g., 'node_modules,dist,.git')",
	)

	rootCmd.Flags().BoolVar(
		&opts.KeepExif,
		"keep-exif",
		false,
		"Preserve EXIF metadata in JPEG files",
	)
}

func runOptimizer(cmd *cobra.Command, args []string) error {
	opts.Input = args[0]

	// Validate input directory exists
	inputInfo, err := os.Stat(opts.Input)
	if err != nil {
		return fmt.Errorf("input directory error: %w", err)
	}
	if !inputInfo.IsDir() {
		return fmt.Errorf("input must be a directory")
	}

	// Handle replace flag
	if opts.Replace {
		// Show confirmation prompt
		fmt.Printf("\n⚠️  WARNING: Replace Mode Enabled\n")
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("You will replace ALL images in this directory:\n")
		fmt.Printf("  📁 %s\n", opts.Input)
		fmt.Printf("\n⚠️  IMPORTANT:\n")
		fmt.Printf("  • Original files WILL BE OVERWRITTEN\n")
		fmt.Printf("  • This action CANNOT be undone\n")
		fmt.Printf("  • Make sure you have a backup!\n")
		fmt.Printf("\n")

		// Get user confirmation
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Type 'yes' to continue or press Enter to cancel: ")
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "yes" {
			fmt.Printf("\n❌ Operation cancelled. Your files are safe.\n")
			return nil
		}

		// Set output to input for replacement
		opts.Output = opts.Input
		fmt.Printf("✓ Confirmed. Proceeding with replacement...\n\n")
	} else {
		// If output not specified, create photon-output in the current directory
		if opts.Output == "" {
			opts.Output = filepath.Join(".", "photon-output")
		}
	}

	fmt.Printf("🚀 Photon CLI initialized\n")
	fmt.Printf("   Input:       %s\n", opts.Input)
	
	if opts.DryRun {
		fmt.Printf("   Mode:        🏃 DRY RUN (no files will be written)\n")
	}
	
	fmt.Printf("   Output:      %s\n", opts.Output)
	fmt.Printf("   Quality:     %d%%\n", opts.Quality)
	
	if opts.JPEGQuality > 0 {
		fmt.Printf("   JPEG Quality: %d%%\n", opts.JPEGQuality)
	}
	if opts.PNGQuality > 0 {
		fmt.Printf("   PNG Quality: %d%%\n", opts.PNGQuality)
	}
	
	fmt.Printf("   WebP:        %v\n", opts.WebP)
	fmt.Printf("   Concurrency: %d workers\n", opts.Concurrency)
	
	if opts.Width > 0 {
		fmt.Printf("   Resize:      %dpx width\n", opts.Width)
	}
	if opts.MinSize > 0 {
		fmt.Printf("   Min Size:    %s\n", formatBytes(opts.MinSize))
	}
	if opts.MaxDepth > 0 {
		fmt.Printf("   Max Depth:   %d levels\n", opts.MaxDepth)
	}
	if opts.IgnorePatterns != "" {
		fmt.Printf("   Ignore:      %s\n", opts.IgnorePatterns)
	}
	if opts.KeepExif {
		fmt.Printf("   Keep EXIF:   true\n")
	}
	if opts.Replace {
		fmt.Printf("   Mode:        🔴 REPLACE (files will be overwritten)\n")
	}
	fmt.Printf("\n")

	// Create and run the pipeline coordinator
	coordinator := pipeline.NewCoordinator(opts.Input, opts.Output, opts)
	stats, err := coordinator.Run()
	if err != nil {
		return fmt.Errorf("pipeline error: %w", err)
	}

	// Display summary
	fmt.Printf("✨ Processing complete!\n")
	fmt.Printf("   Total files:      %d\n", stats.TotalFiles())
	fmt.Printf("   Successful:       %d\n", stats.SuccessfulFiles)
	fmt.Printf("   Failed:           %d\n", stats.FailedFiles)
	fmt.Printf("   Total saved:      %s\n", formatBytes(stats.TotalBytesSaved))
	if stats.SuccessfulFiles > 0 {
		fmt.Printf("   Average per file: %s\n", formatBytes(stats.AverageSavingsPerFile()))
	}
	fmt.Printf("   Success rate:     %.1f%%\n", stats.SuccessRate())
	fmt.Printf("\n")

	// Display output folder in a terminal-friendly format
	absOutputPath, err := filepath.Abs(opts.Output)
	if err != nil {
		absOutputPath = opts.Output
	}
	fmt.Printf("📁 Output folder: %s\n", absOutputPath)
	fmt.Printf("\n")

	// Create and write metadata file (only if not dry-run)
	if !opts.DryRun {
		metaData := metadata.Create(opts.Input, opts.Output, opts, stats)
		metadataPath := filepath.Join(opts.Output, "metadata.json")
		if err := metaData.WriteToFile(metadataPath); err != nil {
			fmt.Printf("⚠️  Warning: Could not write metadata: %v\n", err)
		} else {
			fmt.Printf("📄 Metadata:      %s\n", metadataPath)
		}
	} else {
		fmt.Printf("📄 Metadata:      (skipped in dry-run mode)\n")
	}

	return nil
}

// formatBytes formats bytes into human-readable format
func formatBytes(bytes int64) string {
	units := []string{"B", "KB", "MB", "GB"}
	size := float64(bytes)
	unitIndex := 0

	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}

	return fmt.Sprintf("%.1f%s", size, units[unitIndex])
}
