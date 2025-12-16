# Bitrim - Bun Package

This is the Bun package for **Bitrim**, a blazing-fast, cross-platform image and SVG optimizer.

## Installation

```bash
bun install -g bitrim
```

This will download and install the appropriate binary for your platform (Windows, macOS, or Linux).

## Quick Start

```bash
# Optimize images in current directory
bitrim .

# Optimize with specific quality
bitrim . --quality 75

# Dry run to see what would be optimized
bitrim . --dry-run

# Get help
bitrim --help
```

## Features

- 🚀 **High-Performance**: Multi-threaded concurrent processing with worker pools
- 🎨 **Multiple Formats**: JPEG, PNG, SVG, and WebP support
- 📦 **Smart Compression**: 
  - PNG color quantization (80% quality = 70% size reduction)
  - JPEG quality-based compression
  - SVG minification
- 🛡️ **Safe-by-Default**: Creates `bitrim-output` folder, never overwrites originals
- ⚙️ **Flexible**: Dry-run mode, batch processing, depth control, ignore patterns
- 🔧 **Cross-Platform**: Works on Windows, macOS, and Linux

## Full Documentation

For complete documentation, examples, and troubleshooting, see:
https://github.com/zulfikawr/bitrim#readme

## License

MIT - See LICENSE file for details

## Contributing

Contributions welcome! See CONTRIBUTING.md on the main repository.

## Source Repository

https://github.com/zulfikawr/bitrim
