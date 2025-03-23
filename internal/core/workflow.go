package core

import (
	"fmt"
	"strings"

	"slark/internal/models"
)

// GenerateWorkflows creates workflow files based on the project configuration
// and platform-specific settings.
func GenerateWorkflows(config models.ProjectConfig, platformData models.PlatformData, telegramData models.TelegramData) ([]string, error) {
	// List to store the paths of generated workflow files
	var generatedFiles []string

	// Generate platform-specific workflows
	switch config.Platform {
	case "vercel":
		files, err := generateVercelWorkflow(config, platformData)
		if err != nil {
			return nil, err
		}
		generatedFiles = append(generatedFiles, files...)

	case "cloudflare":
		files, err := generateCloudflareWorkflow(config, platformData)
		if err != nil {
			return nil, err
		}
		generatedFiles = append(generatedFiles, files...)

	default:
		return nil, fmt.Errorf("unsupported platform: %s", config.Platform)
	}

	// Add notification workflows if enabled
	if telegramData.BotToken != "" && telegramData.ChatId != "" {
		files, err := generateNotificationWorkflow(config, telegramData)
		if err != nil {
			return nil, err
		}
		generatedFiles = append(generatedFiles, files...)
	}

	return generatedFiles, nil
}

// generateVercelWorkflow creates GitHub Actions workflow files for Vercel deployments
func generateVercelWorkflow(config models.ProjectConfig, platformData models.PlatformData) ([]string, error) {
	// Validate Vercel-specific requirements
	if platformData.ApiKey == "" {
		return nil, fmt.Errorf("Vercel API key is required")
	}

	// Use team ID if provided
	teamConfig := ""
	if platformData.TeamId != "" {
		teamConfig = fmt.Sprintf("--scope %s", platformData.TeamId)
	}

	// Create workflow content
	// Marked as _ to avoid unused variable warning while keeping the code for reference
	_ = fmt.Sprintf(`
name: Deploy to Vercel

on:
  push:
    branches:
      - %s

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          
      - name: Install Vercel CLI
        run: npm install --global vercel@latest
        
      - name: Deploy to Vercel
        run: |
          vercel --token ${{ secrets.VERCEL_TOKEN }} %s --prod
        env:
          VERCEL_TOKEN: ${{ secrets.VERCEL_TOKEN }}
          VERCEL_PROJECT_ID: ${{ secrets.VERCEL_PROJECT_ID }}
          VERCEL_ORG_ID: ${{ secrets.VERCEL_ORG_ID }}
`, config.DeployBranch, teamConfig)

	// Define the workflow file path
	workflowPath := ".github/workflows/vercel-deploy.yml"

	// TODO: In a real implementation, write this content to the file
	// For now, just return the path that would be created

	return []string{workflowPath}, nil
}

// generateCloudflareWorkflow creates GitHub Actions workflow files for Cloudflare deployments
func generateCloudflareWorkflow(config models.ProjectConfig, platformData models.PlatformData) ([]string, error) {
	// Validate Cloudflare-specific requirements
	if platformData.ApiKey == "" {
		return nil, fmt.Errorf("Cloudflare API key is required")
	}

	// Create workflow content
	// Marked as _ to avoid unused variable warning while keeping the code for reference
	_ = fmt.Sprintf(`
name: Deploy to Cloudflare Pages

on:
  push:
    branches:
      - %s

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          
      - name: Deploy to Cloudflare Pages
        uses: cloudflare/pages-action@v1
        with:
          apiToken: ${{ secrets.CLOUDFLARE_API_TOKEN }}
          accountId: ${{ secrets.CLOUDFLARE_ACCOUNT_ID }}
          projectName: %s
          directory: %s
          gitHubToken: ${{ secrets.GITHUB_TOKEN }}
`, config.DeployBranch, config.Name, config.BuildFolder)

	// Define the workflow file path
	workflowPath := ".github/workflows/cloudflare-deploy.yml"

	// TODO: In a real implementation, write this content to the file
	// For now, just return the path that would be created

	return []string{workflowPath}, nil
}

// generateNotificationWorkflow creates workflow files for notifications
func generateNotificationWorkflow(config models.ProjectConfig, telegramData models.TelegramData) ([]string, error) {
	// Create workflow content for Telegram notifications
	// Marked as _ to avoid unused variable warning while keeping the code for reference
	_ = fmt.Sprintf(`
name: Deployment Notifications

on:
  workflow_run:
    workflows: ["Deploy to %s"]
    types:
      - completed

jobs:
  notify:
    runs-on: ubuntu-latest
    steps:
      - name: Send Telegram notification
        uses: appleboy/telegram-action@master
        with:
          to: %s
          token: %s
          message: |
            Project: %s
            Deployment ${{ github.event.workflow_run.conclusion }} 
            See details: ${{ github.event.workflow_run.html_url }}
`, strings.Title(config.Platform), telegramData.ChatId, "${{ secrets.TELEGRAM_TOKEN }}", config.Name)

	// Define the workflow file path
	workflowPath := ".github/workflows/notification.yml"

	// TODO: In a real implementation, write this content to the file
	// For now, just return the path that would be created

	return []string{workflowPath}, nil
}
