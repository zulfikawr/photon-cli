# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-12-16

### Added
- Initial release of Bitrim
- High-concurrency image and SVG optimization
- Support for JPEG, PNG, and SVG formats
- Intelligent PNG color quantization for effective compression
- JPEG quality-based compression
- SVG minification
- Optional image resizing with aspect ratio preservation
- Worker pool pattern for concurrent processing
- Safe-by-default operation with `bitrim-output` folder
- `--replace` flag for explicit original file replacement (with confirmation)
- `--dry-run` mode to preview savings without writing files
- Detailed statistics and compression metrics
- JSON metadata tracking for audit trail
- Ignore patterns for selective folder exclusion
- Depth limiting for directory traversion control
- Minimum file size filtering
- Format-specific quality overrides (`--jpeg-quality`, `--png-quality`)
- EXIF metadata preservation option (`--keep-exif`)
- Configurable concurrency levels
- Cross-platform support (Windows, macOS, Linux)

### Features
- Batch process entire directory structures recursively
- Real-time compression statistics
- Complete audit trail in `metadata.json`
- Command-line configuration with sensible defaults
- Detailed help text and usage examples

### Technical
- Built with Go 1.21+
- Uses `spf13/cobra` for CLI
- Uses `disintegration/imaging` for image processing
- Uses `tdewolff/minify` for SVG optimization
- Concurrent processing with goroutines and channels
