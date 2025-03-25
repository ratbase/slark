package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ListTemplates prints all available templates categorized by platform.
func ListTemplates() {
	// Path to templates directory
	templatesPath := "templates"

	// Check if templates directory exists
	if _, err := os.Stat(templatesPath); os.IsNotExist(err) {
		fmt.Println("No templates found.")
		return
	}

	// Map to store templates by category
	templatesByCategory := make(map[string][]string)

	// Walk through template directories
	err := filepath.WalkDir(templatesPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root templates directory
		if path == templatesPath {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(templatesPath, path)
		if err != nil {
			return err
		}

		// Split the path to get category
		parts := strings.Split(relPath, string(os.PathSeparator))
		category := strings.Title(parts[0])

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Skip hidden files
		if strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		// Get template name without extension
		templateName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		templateName = strings.ReplaceAll(templateName, "-", " ")
		templateName = strings.ReplaceAll(templateName, "_", " ")
		templateName = strings.Title(templateName)

		// Add to map
		templatesByCategory[category] = append(templatesByCategory[category], templateName)

		return nil
	})

	if err != nil {
		fmt.Printf("Error reading templates: %v\n", err)
		return
	}

	// Print templates by category
	fmt.Println("Available templates:")
	for category, templates := range templatesByCategory {
		fmt.Printf("\n%s:\n", category)
		for _, tmpl := range templates {
			fmt.Printf("  - %s\n", tmpl)
		}
	}
}
