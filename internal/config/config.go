package config

// Options holds all CLI flags and configuration
type Options struct {
	// Input directory path
	Input string

	// Output directory path
	Output string

	// JPEG/WebP quality (1-100)
	Quality int

	// JPEG-specific quality (overrides Quality if set)
	JPEGQuality int

	// PNG-specific quality (overrides Quality if set)
	PNGQuality int

	// Resize width in pixels (0 = no resize)
	Width int

	// Generate WebP copies
	WebP bool

	// Number of concurrent workers
	Concurrency int

	// Replace original files (replace input with output)
	Replace bool

	// Dry run mode (no files written)
	DryRun bool

	// Minimum file size to process (in bytes, 0 = no minimum)
	MinSize int64

	// Maximum recursion depth (0 = unlimited)
	MaxDepth int

	// Patterns to ignore (comma-separated)
	IgnorePatterns string

	// Preserve EXIF metadata in JPG files
	KeepExif bool
}
