package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TemplateManager handles loading and managing PR templates
type TemplateManager struct {
	templateDir string
}

// NewTemplateManager creates a new template manager
func NewTemplateManager() *TemplateManager {
	return &TemplateManager{
		templateDir: "templates",
	}
}

// LoadTemplate loads a template by name or file path
func LoadTemplate(templateName string) (string, error) {
	tm := NewTemplateManager()

	// If it's a file path, load directly
	if strings.Contains(templateName, "/") || strings.Contains(templateName, "\\") {
		return tm.loadTemplateFile(templateName)
	}

	// Otherwise, try to load from templates directory
	return tm.loadTemplateByName(templateName)
}

// loadTemplateFile loads a template from a specific file path
func (tm *TemplateManager) loadTemplateFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %w", filePath, err)
	}
	return string(content), nil
}

// loadTemplateByName loads a template from the templates directory
func (tm *TemplateManager) loadTemplateByName(templateName string) (string, error) {
	// Try different file extensions
	extensions := []string{".md", ".txt", ""}

	for _, ext := range extensions {
		filename := templateName + ext
		paths := tm.getTemplatePaths(filename)

		for _, path := range paths {
			if content, err := os.ReadFile(path); err == nil {
				return string(content), nil
			}
		}
	}

	return "", fmt.Errorf("template '%s' not found in templates directory", templateName)
}

// getTemplatePaths returns possible paths for a template file
func (tm *TemplateManager) getTemplatePaths(filename string) []string {
	return []string{
		filename, // Current directory
		filepath.Join(tm.templateDir, filename),
		filepath.Join(".github", tm.templateDir, filename),
		filepath.Join(".gitlab", tm.templateDir, filename),
	}
}

// GetAvailableTemplates returns list of available templates
func GetAvailableTemplates() ([]string, error) {
	tm := NewTemplateManager()
	var templates []string

	// Check templates directory
	dirs := tm.getTemplateDirs()
	for _, dir := range dirs {
		if files, err := os.ReadDir(dir); err == nil {
			for _, file := range files {
				if !file.IsDir() && isTemplateFile(file.Name()) {
					templates = append(templates, strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())))
				}
			}
		}
	}

	return templates, nil
}

// getTemplateDirs returns possible template directories
func (tm *TemplateManager) getTemplateDirs() []string {
	return []string{
		tm.templateDir,
		filepath.Join(".github", tm.templateDir),
		filepath.Join(".gitlab", tm.templateDir),
		".",
	}
}

// isTemplateFile checks if a file is a template file
func isTemplateFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".md" || ext == ".txt" || ext == ".tmpl"
}

// CreateDefaultTemplate creates a default template if it doesn't exist
func CreateDefaultTemplate() error {
	tm := NewTemplateManager()

	// Create templates directory if it doesn't exist
	if err := os.MkdirAll(tm.templateDir, 0755); err != nil {
		return err
	}

	defaultTemplatePath := filepath.Join(tm.templateDir, "default.md")
	if _, err := os.Stat(defaultTemplatePath); os.IsNotExist(err) {
		defaultContent := `# Description

<!-- AI will generate a description of the changes here -->

# Changelog

<!-- AI will list the key changes and updates here -->

# Test Evidence

<!-- AI will describe testing performed here -->

# Additional Notes

<!-- AI will add any additional context or notes here -->
`
		if err := os.WriteFile(defaultTemplatePath, []byte(defaultContent), 0644); err != nil {
			return err
		}
	}

	return nil
}
