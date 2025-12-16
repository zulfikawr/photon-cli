# 📷 Bitrim

Bitrim (short for bit trim) is a blazing-fast, cross-platform CLI tool for batch optimizing images and SVG files with high-concurrency processing. Compress JPEGs, PNGs, and minify SVGs with a single command.

[![Go 1.21+](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![GitHub Release](https://img.shields.io/github/v/release/zulfikawr/bitrim?include_prereleases)](https://github.com/zulfikawr/bitrim/releases)

## ✨ Features

- **🔄 Batch Processing**: Recursively scan and optimize entire directory structures
- **🎯 Format-Specific Optimization**:
  - **PNG**: Intelligent color quantization for meaningful compression
  - **JPEG**: Adjustable quality settings with optional EXIF preservation
  - **SVG**: Minification to reduce file sizes
- **⚡ High-Concurrency**: Worker pool pattern for maximum CPU utilization
- **🛡️ Safe by Default**: Creates `bitrim-output` folder, never overwrites originals without explicit confirmation
- **📊 Detailed Statistics**: Real-time compression metrics and success rates
- **🏷️ Metadata Tracking**: JSON audit trail of all processed files
- **🎨 Smart Resizing**: Optional width-based resize with aspect ratio preservation
- **🔍 Selective Processing**: Ignore patterns, minimum file size filters, depth limiting
- **🔒 Data Protection**: Original files preserved, replaceable only with explicit `--replace` flag
- **💾 Dry-Run Mode**: Preview savings without writing files

## 🚀 Quick Start

### Installation

#### From Source
```bash
git clone https://github.com/zulfikawr/bitrim.git
cd bitrim
go build -o bitrim
```

#### With Go
```bash
go install github.com/zulfikawr/bitrim@latest
```

### Basic Usage

```bash
# Optimize all images in a directory (outputs to ./bitrim-output/)
bitrim ./path/to/images

# Custom output directory
bitrim -o ./compressed ./path/to/images

# Adjust compression quality (default: 80%)
bitrim -q 60 ./path/to/images

# Resize images to 1200px width
bitrim -w 1200 ./path/to/images

# Combine multiple options
bitrim -q 50 -w 1600 -o ./optimized ./path/to/images

# See all options
bitrim --help
```

## 📋 Command Line Options

### Core Options

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--out` | `-o` | `bitrim-output` | Output directory for optimized files |
| `--quality` | `-q` | `80` | JPEG/PNG quality (1-100) |
| `--width` | `-w` | `0` | Resize images to width (px), 0=no resize |
| `--webp` | - | `false` | Convert images to WebP format |
| `--concurrency` | - | `2` | Number of worker threads |

### Format-Specific Quality

| Flag | Default | Description |
|------|---------|-------------|
| `--jpeg-quality` | `0` | Override JPEG quality (overrides `--quality` for JPEGs) |
| `--png-quality` | `0` | Override PNG quality (overrides `--quality` for PNGs) |

### Advanced Options

| Flag | Default | Description |
|------|---------|-------------|
| `--replace` / `-r` | `false` | Replace original files (requires confirmation) |
| `--dry-run` | `false` | Calculate savings without writing files |
| `--min-size` | `0` | Minimum file size to process (e.g., `1mb`, `100kb`) |
| `--depth` | `0` | Maximum recursion depth (0=unlimited) |
| `--ignore` | `` | Comma-separated patterns to ignore |
| `--keep-exif` | `false` | Preserve EXIF metadata in JPEG files |

## 💡 Usage Examples

### Basic Optimization
```bash
bitrim ./images
# Outputs optimized files to ./bitrim-output/
# Default: 80% quality, no resize
```

### Aggressive Compression
```bash
bitrim -q 50 -w 1200 ./images
# Saves: ~70% file size reduction
# Quality: Noticeable but acceptable
```

### Quality Focus
```bash
bitrim -q 90 ./images
# Saves: ~30% file size reduction
# Quality: Minimal difference from originals
```

### Preview Changes (Dry Run)
```bash
bitrim --dry-run -q 60 -w 1600 ./images
# Shows estimated savings without creating files
# Perfect for testing settings before committing
```

### Replace Originals
```bash
bitrim --replace -q 75 ./images
# ⚠️ Replaces files in ./images (requires "yes" confirmation)
# Original files will be overwritten
```

### Selective Processing
```bash
bitrim --ignore "node_modules,dist,.git" --depth 3 ./project
# Ignore specific folders
# Limit recursion to 3 levels deep
# Skip files smaller than 1MB
```

### Format-Specific Quality
```bash
# Use different quality for PNG vs JPEG
bitrim --jpeg-quality 85 --png-quality 60 ./images

# Only apply quality to PNGs (ignore general --quality)
bitrim --png-quality 55 ./images
```

### Preserve EXIF Data
```bash
bitrim --keep-exif -q 85 ./photos
# Maintains EXIF metadata while compressing JPEGs
```

## 📊 Output

Bitrim provides detailed feedback:

```
🚀 Bitrim initialized
   Input:       ./images
   Output:      bitrim-output
   Quality:     80%
   WebP:        false
   Concurrency: 2 workers
   Resize:      1200px width

✨ Processing complete!
   Total files:      150
   Successful:       150
   Failed:           0
   Total saved:      142.5MB
   Average per file: 0.95MB
   Success rate:     100.0%

📁 Output folder: /absolute/path/bitrim-output

📄 Metadata:      bitrim-output/metadata.json
```

### Metadata File

A `metadata.json` file is generated in the output folder with complete audit trail:

```json
{
  "timestamp": "2024-12-16T14:30:45Z",
  "config": {
    "quality": 80,
    "width": 1200,
    "replace": false,
    "dry_run": false
  },
  "summary": {
    "total_files": 150,
    "successful": 150,
    "failed": 0,
    "total_bytes_original": 500000000,
    "total_bytes_processed": 142500000,
    "total_bytes_saved": 357500000,
    "compression_ratio": 71.5
  },
  "files": [
    {
      "input_path": "images/photo1.jpg",
      "output_path": "bitrim-output/photo1.jpg",
      "file_type": "jpeg",
      "size_original": 5242880,
      "size_processed": 1500000,
      "bytes_saved": 3742880,
      "compression_ratio": 71.4,
      "success": true,
      "error": ""
    }
  ]
}
```

## 🔧 Technical Details

### Architecture

- **Language**: Go 1.21+
- **Pattern**: Worker pool with goroutines and channels
- **Concurrency**: Configurable worker threads (default: 2)
- **Dependencies**:
  - `github.com/spf13/cobra` - CLI framework
  - `github.com/disintegration/imaging` - Image processing
  - `github.com/tdewolff/minify/v2` - SVG minification

### How It Works

1. **Walker**: Recursively scans input directory respecting depth limits and ignore patterns
2. **Coordinator**: Orchestrates worker pool and manages pipeline flow
3. **Workers**: Process files concurrently using available CPU cores
4. **Optimizer**: 
   - Decodes image formats
   - Applies color quantization (PNG) or quality reduction (JPEG)
   - Optionally resizes to specified width
   - Encodes with optimized settings
5. **Stats**: Aggregates metrics and generates metadata

### Compression Strategy

**PNG Files**:
- Uses color quantization to reduce palette
- Quality % maps to color count (80% quality ≈ 204 colors)
- Lower quality = more aggressive color reduction
- Effective for graphics and illustrations

**JPEG Files**:
- Adjusts compression quality directly
- Preserves all pixels but reduces detail
- Quality % directly controls compression level
- Best for photographs

**SVG Files**:
- Minifies XML structure
- Removes unnecessary attributes and whitespace
- Preserves visual appearance
- Best for logos and vector graphics

## 🛡️ Safety Features

### Original Files Protected
- By default, creates `bitrim-output/` folder
- Original files in input directory never modified
- No data loss risk

### Explicit Replacement
```bash
bitrim --replace ./images
# Shows warning:
# ⚠️ WARNING: Replace Mode Enabled
# Type 'yes' to continue or press Enter to cancel:
```

### Dry-Run Preview
```bash
bitrim --dry-run -q 60 -w 1600 ./images
# Calculates savings without creating any files
# Perfect for testing before committing
```

### Metadata Audit Trail
- Complete record of all operations
- Input/output paths tracked
- Compression ratios recorded
- Success/failure status documented

## 📈 Performance

On modern hardware with 2 worker threads:
- **Small images** (<1MB): 50-100 images/second
- **Medium images** (1-10MB): 5-15 images/second
- **Large images** (>10MB): 1-5 images/second
- **Concurrency**: Scales with `--concurrency` flag (2-8 recommended)

### Example Results

Compressing a batch of mixed-resolution photos:
```
Input:  2,450 images (15.3GB)
Settings: -q 75 -w 1600
Output: 4.2GB (72.5% reduction)
Time: ~4 minutes with 4 workers
```

## 🐛 Troubleshooting

### Issue: "Input directory error"
```bash
# Ensure path exists and is a directory
ls -la ./images
```

### Issue: "Failed to write metadata"
```bash
# Output folder not writable - check permissions
ls -la ./bitrim-output
```

### Issue: "No files processed"
```bash
# Check ignore patterns and file extensions
bitrim --ignore "" ./images  # Disable ignore patterns
```

### Issue: Large quality drop
```bash
# Use higher quality setting
bitrim -q 85 ./images  # Instead of -q 50
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📦 Related Projects

- [imagemin](https://github.com/imagemin/imagemin) - Node.js image optimization
- [ImageOptim](https://imageoptim.com/command-line.html) - macOS image optimizer
- [OptiPNG](http://optipng.sourceforge.net/) - PNG optimizer

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/zulfikawr/bitrim/issues)
- **Discussions**: [GitHub Discussions](https://github.com/zulfikawr/bitrim/discussions)
- **Documentation**: Full feature docs in `README.md`

## 🎯 Roadmap

- [ ] WebP format support
- [ ] AVIF format support  
- [ ] Parallel batch processing across directories
- [ ] Configuration files (.bitrimrc)
- [ ] Progress bar with ETA
- [ ] EXIF auto-rotation for photos
- [ ] Smart quality based on image content
- [ ] Output format conversion (JPEG → WebP)

---

**Made with ❤️ by zulfikar** | Optimizing images, one file at a time.
