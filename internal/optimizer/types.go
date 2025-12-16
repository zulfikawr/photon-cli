package optimizer

// Result represents the outcome of processing a file
type Result struct {
	// Original file path
	FilePath string

	// File type (jpg, png, svg)
	FileType string

	// Original file size in bytes
	OriginalSize int64

	// Processed file size in bytes
	ProcessedSize int64

	// Bytes saved
	BytesSaved int64

	// Output file path
	OutputPath string

	// Whether processing was successful
	Success bool

	// Error message if processing failed
	Error string
}

// ImageFormat represents supported image formats
type ImageFormat string

const (
	FormatJPEG ImageFormat = "jpeg"
	FormatPNG  ImageFormat = "png"
	FormatWebP ImageFormat = "webp"
	FormatSVG  ImageFormat = "svg"
)
