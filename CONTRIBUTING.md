# Contributing to Bitrim

Thank you for your interest in contributing to Bitrim! We appreciate your help in making this project better.

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Git
- Basic knowledge of Go and image optimization

### Setting Up Development Environment

1. Clone the repository:
```bash
git clone https://github.com/zulfikawr/bitrim.git
cd bitrim
```

2. Install dependencies:
```bash
go mod download
```

3. Build the project:
```bash
go build -o bitrim
```

4. Run tests:
```bash
go test ./...
```

## Development Workflow

### Branch Naming
- Feature: `feature/description`
- Bug fix: `fix/description`
- Documentation: `docs/description`

### Commit Messages
- Use clear, descriptive commit messages
- Start with a verb (Add, Fix, Update, Remove, etc.)
- Example: "Add PNG color quantization for better compression"

### Code Style
- Follow Go conventions and best practices
- Use `gofmt` for formatting
- Run `go vet` to check for common errors
- Write tests for new functionality

### Testing
Before submitting a PR:
1. Run all tests: `go test ./...`
2. Run linter: `go vet ./...`
3. Test with various image formats and sizes
4. Verify no regressions with existing functionality

## Project Structure

```
bitrim/
├── cmd/
│   └── root.go           # CLI command definitions
├── internal/
│   ├── config/           # Configuration and options
│   ├── metadata/         # Metadata generation and tracking
│   ├── optimizer/        # Image and SVG processing
│   └── pipeline/         # Worker pool and coordination
├── main.go               # Entry point
├── go.mod
└── README.md
```

## Adding New Features

### New Image Format Support
1. Add decoder/encoder in `internal/optimizer/processor.go`
2. Update `internal/pipeline/walker.go` to recognize file type
3. Add tests in `internal/optimizer/processor_test.go`
4. Update README.md with new format

### New Command-Line Flag
1. Add field to `internal/config/Options` struct
2. Define flag in `cmd/root.go` with `rootCmd.Flags()`
3. Pass option through pipeline to optimizer
4. Update README.md with documentation

### Performance Improvements
1. Benchmark existing code: `go test -bench ./...`
2. Implement improvement
3. Verify benchmarks don't regress
4. Document changes in commit message

## Bug Reports

When reporting bugs, please include:
- Operating system and Go version
- Bitrim version
- Command used and flags
- Input file type and size
- Error message or unexpected behavior
- Steps to reproduce

## Feature Requests

When suggesting features:
- Describe the use case clearly
- Explain how it benefits users
- Consider performance implications
- Check if similar features exist

## Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make your changes
4. Add/update tests as needed
5. Update README.md if adding user-facing features
6. Commit with clear messages
7. Push to your fork
8. Create a Pull Request with:
   - Clear title and description
   - Reference to related issues
   - Explanation of changes
   - Any breaking changes noted

## Code Review Guidelines

- Be respectful and constructive
- Focus on the code, not the person
- Ask questions rather than make demands
- Test locally before approving
- Verify tests pass and coverage is maintained

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Questions?

Feel free to open an issue or discussion on GitHub for questions about the project or contribution process.

---

Thank you for helping improve Bitrim! 🎉
