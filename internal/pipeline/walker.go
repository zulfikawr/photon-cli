package pipeline

import (
	"os"
	"path/filepath"
	"strings"
)

// FileInfo represents a file to be processed
type FileInfo struct {
	Path string
	Type string // "image" or "svg"
}

// Walker scans a directory recursively and sends files to the jobs channel
type Walker struct {
	rootDir        string
	jobsCh         chan<- FileInfo
	ignorePatterns []string
	maxDepth       int
	minSize        int64
}

// NewWalker creates a new Walker
func NewWalker(rootDir string, jobsCh chan<- FileInfo, ignorePatterns []string, maxDepth int, minSize int64) *Walker {
	return &Walker{
		rootDir:        rootDir,
		jobsCh:         jobsCh,
		ignorePatterns: ignorePatterns,
		maxDepth:       maxDepth,
		minSize:        minSize,
	}
}

// Walk recursively scans the directory and sends valid files to the jobs channel
func (w *Walker) Walk() error {
	return w.walkDir(w.rootDir, 0)
}

// walkDir recursively walks directories respecting depth limit and ignore patterns
func (w *Walker) walkDir(dir string, depth int) error {
	// Check depth limit (root is depth 0)
	// For example: maxDepth=0 means only root, maxDepth=1 means root + level1, etc.
	if w.maxDepth > 0 && depth > w.maxDepth-1 {
		return nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())

		// Check if path matches ignore patterns
		if w.shouldIgnore(path) {
			continue
		}

		if entry.IsDir() {
			// Recursively walk subdirectories
			if err := w.walkDir(path, depth+1); err != nil {
				return err
			}
		} else {
			// Check file size minimum
			info, err := entry.Info()
			if err != nil {
				continue
			}
			if w.minSize > 0 && info.Size() < w.minSize {
				continue
			}

			// Check if file matches supported extensions
			ext := strings.ToLower(filepath.Ext(path))
			fileType := getFileType(ext)

			if fileType != "" {
				w.jobsCh <- FileInfo{
					Path: path,
					Type: fileType,
				}
			}
		}
	}

	return nil
}

// shouldIgnore checks if a path matches any ignore patterns
func (w *Walker) shouldIgnore(path string) bool {
	for _, pattern := range w.ignorePatterns {
		pattern = strings.TrimSpace(pattern)
		// Match both full path and basename
		if strings.Contains(path, pattern) || strings.HasPrefix(filepath.Base(path), pattern) {
			return true
		}
	}
	return false
}

// getFileType returns "image", "svg", or "" based on file extension
func getFileType(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image"
	case ".png":
		return "image"
	case ".svg":
		return "svg"
	default:
		return ""
	}
}
