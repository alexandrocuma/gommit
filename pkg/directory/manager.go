package directory

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// LoadFile loads a single file from the specified path
func LoadFile(filePath string) ([]byte, error) {
	resolvedPath, err := ResolvePath(filePath)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(resolvedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	return content, nil
}

// LoadFileString loads a single file and returns content as string
func LoadFileString(filePath string) (string, error) {
	content, err := LoadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// LoadFileFromDir loads a file by name from a specific directory
func LoadFileFromDir(dirPath, fileName string) ([]byte, error) {
	return LoadFile(filepath.Join(dirPath, fileName))
}

// LoadFileFromDirString loads a file by name and returns content as string
func LoadFileFromDirString(dirPath, fileName string) (string, error) {
	return LoadFileString(filepath.Join(dirPath, fileName))
}

// =============================================================================
// Directory Validation and Checking
// =============================================================================

// DirExists checks if a directory exists and is accessible
func DirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	return err == nil && info.IsDir()
}

// FileExists checks if a file exists and is accessible
func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	return err == nil && !info.IsDir()
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(dirPath string, perm os.FileMode) error {
	if DirExists(dirPath) {
		return nil
	}
	return os.MkdirAll(dirPath, perm)
}

// ------------------------------------------------------------------
// 1.  List every **sub-directory** (non-recursive) inside a directory
// ------------------------------------------------------------------
func ListDirs(dirPath string) ([]string, error) {
	resolved, err := ResolvePath(dirPath)
	if err != nil {
		return nil, err
	}
	return listEntries(resolved, func(_ string, isDir bool) bool { return isDir })
}

// ------------------------------------------------------------------
// 2.  List only directories whose name matches a given prefix
// ------------------------------------------------------------------
func ListDirsWithPrefix(dirPath, prefix string) ([]string, error) {
	resolved, err := ResolvePath(dirPath)
	if err != nil {
		return nil, err
	}
	return listEntries(resolved, func(name string, isDir bool) bool {
		return isDir && strings.HasPrefix(name, prefix)
	})
}

// ------------------------------------------------------------------
// 3.  Generic directory filter (used internally)
// ------------------------------------------------------------------
func listEntries(dirPath string, want func(name string, isDir bool) bool) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read directory %s: %w", dirPath, err)
	}
	var out []string
	for _, e := range entries {
		if want(e.Name(), e.IsDir()) {
			out = append(out, e.Name())
		}
	}
	return out, nil
}


// =============================================================================
// File Listing and Filtering
// =============================================================================

// ListFiles lists all files (non-directories) in a directory
func ListFiles(dirPath string) ([]string, error) {
	resolvedPath, err := ResolvePath(dirPath)
	if err != nil {
		return nil, err
	}

	return ListFilesWithFilter(resolvedPath, func(name string, isDir bool) bool {
		return !isDir
	})
}

// ListFilesByExtension lists files with specific extensions
func ListFilesByExtension(dirPath string, extensions ...string) ([]string, error) {
	resolvedPath, err := ResolvePath(dirPath)
	if err != nil {
		return nil, err
	}

	extMap := make(map[string]bool)
	for _, ext := range extensions {
		extMap[strings.ToLower(ext)] = true
	}

	return ListFilesWithFilter(resolvedPath, func(name string, isDir bool) bool {
		if isDir {
			return false
		}
		ext := strings.ToLower(filepath.Ext(name))
		return extMap[ext]
	})
}

// ListFilesWithFilter lists files using a custom filter function
func ListFilesWithFilter(dirPath string, filter func(name string, isDir bool) bool) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	var files []string
	for _, entry := range entries {
		if filter(entry.Name(), entry.IsDir()) {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

// WalkFiles walks a directory tree and returns all files matching filter
func WalkFiles(rootDir string, filter func(path string, info fs.FileInfo) bool) ([]string, error) {
	var files []string

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if !info.IsDir() && filter(path, info) {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", rootDir, err)
	}

	return files, nil
}

// =============================================================================
// Template-Specific Loading (From Previous Requirements)
// =============================================================================

// TemplateResult holds the result from loading a template from a directory
type TemplateResult struct {
	DirPath string
	Content string
	Error   error
}

// LoadTemplate loads a template by name from a directory or direct file path
func LoadTemplate(dirPath, templateName string) (string, error) {
	if strings.Contains(templateName, "/") || strings.Contains(templateName, "\\") {
		return loadTemplateFile(templateName)
	}
	return loadTemplateByName(dirPath, templateName)
}

func loadTemplateFile(filePath string) (string, error) {
	resolvedPath, err := ResolvePath(filePath)
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(resolvedPath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %w", filePath, err)
	}
	return string(content), nil
}

func loadTemplateByName(dirPath, templateName string) (string, error) {
	resolvedPath, err := ResolvePath(dirPath)
	if err != nil {
		return "", err
	}

	extensions := []string{".md", ".txt", ""}

	for _, ext := range extensions {
		filename := templateName + ext
		paths := getTemplatePaths(resolvedPath, filename)

		for _, path := range paths {
			content, err := os.ReadFile(path)
			if err == nil {
				return string(content), nil
			}
		}
	}

	return "", fmt.Errorf("template '%s' not found in directory '%s'", templateName, dirPath)
}

func getTemplatePaths(dirPath, filename string) []string {
	
	return []string{
		filepath.Join(dirPath, filename),
		filepath.Join(".github", dirPath, filename),
		filepath.Join(".gitlab", dirPath, filename),
		filename,
	}
}

// =============================================================================
// Multi-Directory Operations (Concurrent)
// =============================================================================

// LoadTemplatesParallel loads templates from multiple directories concurrently
// Results are returned in input order
func LoadTemplatesParallel(dirPaths []string, templateName string) []TemplateResult {
	results := make([]TemplateResult, len(dirPaths))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, dirPath := range dirPaths {
		wg.Add(1)
		go func(idx int, path string) {
			defer wg.Done()
			content, err := LoadTemplate(path, templateName)

			mu.Lock()
			results[idx] = TemplateResult{DirPath: path, Content: content, Error: err}
			mu.Unlock()
		}(i, dirPath)
	}

	wg.Wait()
	return results
}

// LoadFilesFromDirs loads a specific file from multiple directories in parallel
func LoadFilesFromDirs(dirPaths []string, fileName string) []TemplateResult {
	results := make([]TemplateResult, len(dirPaths))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, dirPath := range dirPaths {
		wg.Add(1)
		go func(idx int, path string) {
			defer wg.Done()
			content, err := LoadFileFromDirString(path, fileName)

			mu.Lock()
			results[idx] = TemplateResult{DirPath: path, Content: content, Error: err}
			mu.Unlock()
		}(i, dirPath)
	}

	wg.Wait()
	return results
}

// GetSuccessfulContents extracts only successful results
func GetSuccessfulContents(results []TemplateResult) []string {
	var contents []string
	for _, r := range results {
		if r.Error == nil {
			contents = append(contents, r.Content)
		}
	}
	return contents
}

// HasErrors checks if any operations failed
func HasErrors(results []TemplateResult) bool {
	for _, r := range results {
		if r.Error != nil {
			return true
		}
	}
	return false
}

// GetErrors collects all errors from results
func GetErrors(results []TemplateResult) []error {
	var errors []error
	for _, r := range results {
		if r.Error != nil {
			errors = append(errors, fmt.Errorf("%s: %w", r.DirPath, r.Error))
		}
	}
	return errors
}

// =============================================================================
// Helper Utilities
// =============================================================================

// SafeFileName sanitizes a filename to remove dangerous characters
func SafeFileName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, "..", "")
	name = strings.ReplaceAll(name, string(filepath.Separator), "_")
	name = strings.ReplaceAll(name, "/", "_")
	return filepath.Base(name)
}

// IsTemplateFile checks if a file has a template extension
func IsTemplateFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".md" || ext == ".txt" || ext == ".tmpl" || ext == ".html"
}

// JoinPath safely joins path components
func JoinPath(components ...string) string {
	cleaned := make([]string, len(components))
	for i, comp := range components {
		cleaned[i] = filepath.Clean(comp)
	}
	return filepath.Join(cleaned...)
}

// DefaultDir returns the first existing directory from the list
func DefaultDir(candidates []string, fallback string) string {
	for _, dir := range candidates {
		if DirExists(dir) {
			return dir
		}
	}
	return fallback
}

// ResolvePath expands ~ to home directory and cleans the path
func ResolvePath(path string) (string, error) {
	if path == "~" {
		return os.UserHomeDir()
	}
	
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		
		// Join home dir with the rest of the path (after "~/")
		return filepath.Join(homeDir, path[2:]), nil
	}
	
	// For relative or absolute paths, just clean them
	return filepath.Clean(path), nil
}
