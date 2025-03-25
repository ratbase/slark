package core

import (
	"fmt"
	"os"

	"slark/internal/models"
)

// GenerateWorkflows creates workflow files based on the project configuration
// and platform-specific settings.
func GenerateWorkflows(config models.ProjectConfig, platformData models.PlatformData) ([]string, error) {
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
	if platformData.BotToken != "" && platformData.ChatId != "" {
		files, err := generateNotificationWorkflow(config, platformData)
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

	// Create workflow content
	// Marked as _ to avoid unused variable warning while keeping the code for reference
	template := fmt.Sprintf(`
name: %s - GitHub Actions Vercel Deployment - branch %s
env:
  VERCEL_ORG_ID: ${{ secrets.VERCEL_ORG_ID }}
  VERCEL_PROJECT_ID: ${{ secrets.VERCEL_IMURA_LANDING }}
on:
  push:
    branches:
      - %s
    paths:
      - %s/**
      - .github/workflows/%s.%s.yml
jobs:
  Deploy-Production:
    runs-on: self-hosted
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v4
        with:
          node-version: 22
      - name: Install Vercel CLI
        run: | 
          npm install --global vercel@canary
          # npm install --global yarn
          npm install -g pnpm
      - name: Pull Vercel Environment Information
        run: vercel pull --yes --environment=production --token=${{ secrets.VERCEL_TOKEN }}
      - name: Build Project Artifacts
        id: build
        run: vercel build --prod --token=${{ secrets.VERCEL_TOKEN }}

      - name: Deploy Project Artifacts to Vercel
        id: deploy
        run: vercel deploy --prebuilt --prod --token=${{ secrets.VERCEL_TOKEN }}

      - name: "set result"
        id: deploy-task-result
        if: always()
        run: |
          if ${{ steps.build.outcome == 'success' && (steps.deploy.outcome == 'success' || steps.deploy.outcome == null) }}; then # Check both build and deploy
            echo "deploy_result=success" >> "$GITHUB_OUTPUT"
          else
            echo "deploy_result=failure" >> "$GITHUB_OUTPUT"
          fi
    outputs:
      deploy_result: ${{ steps.deploy-task-result.outputs.deploy_result }}
    `, config.Name, config.DeployBranch, config.DeployBranch, config.BuildFolder, config.Name, config.DeployBranch)

	//if telegram token is set append this to above string
	if platformData.BotToken != "" && platformData.ChatId != "" {
		template += fmt.Sprintf(`
  noti-tele:
    name: Notify Telegram
    uses: "./.github/workflows/.telegram-noti.yml"
    needs: Deploy-Production
    if: |
      always()
    with:
      main_job_name: Deploy-Production
      results: Deploy ${{ needs.Deploy-Production.outputs.deploy_result }}
      service_name: %s
      `, config.Name)
	}

	// Define the workflow file path
	workflowPath := fmt.Sprintf(`.github/workflows/%s.%s.yml`, config.Name, config.DeployBranch)

	err := os.MkdirAll(".github/workflows", 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create .github/workflows directory: %w", err)
	}

	// Write the workflow content to the file
	err = os.WriteFile(workflowPath, []byte(template), 0644)
	if err != nil {
		return nil, err
	}

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
func generateNotificationWorkflow(config models.ProjectConfig, platformData models.PlatformData) ([]string, error) {
	// Create workflow content from notification template
	template := `on:
  workflow_call:
    inputs:
      main_job_name:
        required: true
        type: string
      results:
        required: true
        type: string
      service_name:
        required: true
        type: string
      dev_id:
        required: false
        type: string
        default: "U061KBS9HDK"
jobs:
  telegram_message:
    runs-on: self-hosted
    if: always()
    steps:
      - name: send telegram message on push
        uses: appleboy/telegram-action@master
        with:
          to: ${{ secrets.TELEGRAM_CHAT_ID }}
          token: ${{ secrets.TELEGRAM_BOT_TOKEN }}
          message: |
            ${{ github.actor }} created commit:
            Commit message: ${{ github.event.commits[0].message }}
            Repository: ${{ github.repository }}
            Project: ${{ inputs.service_name }}
            GitHub Action build result: ${{ inputs.results }}
            See changes: https://github.com/${{ github.repository }}/commit/${{github.sha}}`

	// Define the workflow file path
	workflowPath := ".github/workflows/.telegram-noti.yml"

	// Create the .github/workflows directory if it doesn't exist
	err := os.MkdirAll(".github/workflows", 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create .github/workflows directory: %w", err)
	}

	// Write the workflow content to the file
	err = os.WriteFile(workflowPath, []byte(template), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write notification workflow file: %w", err)
	}

	return []string{workflowPath}, nil
}
