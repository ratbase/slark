package core

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"slark/internal/models"

	tea "github.com/charmbracelet/bubbletea"
)

// SetupProject handles the core project setup logic
// It validates inputs, creates necessary project configurations,
// and prepares everything needed for generating workflows
func SetupProject(projectName, deployBranch, buildFolder, platform string) (models.ProjectConfig, error) {
	// Validate project inputs
	if err := validateProjectInputs(projectName, platform); err != nil {
		return models.ProjectConfig{}, err
	}

	// Set default values for optional fields
	if deployBranch == "" {
		deployBranch = "main"
	}

	if buildFolder == "" {
		buildFolder = "./"
	}

	// Clean up build folder path
	buildFolder = filepath.Clean(buildFolder)

	// Create project configuration
	config := models.ProjectConfig{
		Name:         projectName,
		DeployBranch: deployBranch,
		BuildFolder:  buildFolder,
		Platform:     platform,
		CreatedAt:    time.Now(),
	}

	return config, nil
}

// validateProjectInputs performs validation on required project inputs
func validateProjectInputs(projectName, platform string) error {
	if projectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	validPlatforms := map[string]bool{
		"vercel":     true,
		"cloudflare": true,
	}

	if !validPlatforms[platform] {
		return fmt.Errorf("unsupported platform: %s", platform)
	}

	return nil
}

// ProcessProject is the main function that processes project setup and returns a tea.Cmd
// It's used by the TUI to handle the asynchronous project setup process
func ProcessProject(projectName, deployBranch, buildFolder, platform string, platformData models.PlatformData, telegramData models.TelegramData) tea.Cmd {
	return func() tea.Msg {
		// Initialize result builder
		var resultBuilder strings.Builder

		// Setup project
		config, err := SetupProject(projectName, deployBranch, buildFolder, platform)
		if err != nil {
			return models.ProcessFinishedMsg{
				Success: false,
				Result:  "",
				Err:     err,
			}
		}

		// Generate workflows based on platform
		workflowFiles, err := GenerateWorkflows(config, platformData, telegramData)
		if err != nil {
			return models.ProcessFinishedMsg{
				Success: false,
				Result:  "",
				Err:     err,
			}
		}

		// Build success message
		resultBuilder.WriteString(fmt.Sprintf("Project: %s\n", config.Name))
		resultBuilder.WriteString(fmt.Sprintf("Deploy Branch: %s\n", config.DeployBranch))
		resultBuilder.WriteString(fmt.Sprintf("Build Folder: %s\n", config.BuildFolder))
		resultBuilder.WriteString(fmt.Sprintf("Platform: %s\n", config.Platform))
		resultBuilder.WriteString("\nGenerated workflow files:\n")

		for _, file := range workflowFiles {
			resultBuilder.WriteString(fmt.Sprintf("- %s\n", file))
		}

		resultBuilder.WriteString("\nCI/CD pipeline configured successfully!")

		return models.ProcessFinishedMsg{
			Success: true,
			Result:  resultBuilder.String(),
			Err:     nil,
		}
	}
}
