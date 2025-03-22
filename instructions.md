# Slark

## Overview

call cli -> get input -> create vercel/cloudflare pages -> return id -> create cicd github files -> commit code -> push

## Technical Stack

- Language: Go
- Frameworks: Vite, Next.js, React, etc.
- CICD Provider: GitHub Actions
- Deployment Platforms: Vercel, Cloudflare
- Notification Channel: Telegram

## Features

- Generate GitHub Actions workflow files from templates
- Configure branch and folder tracking for deployment
- Set up required GitHub secrets for deployment platforms and notifications
- Automatically create projects on deployment platforms via their APIs
- Send build and deployment status notifications to Telegram
- Commit changes directly to the repository

## Command Line Interface

### Installation

```
go install github.com/ratbase/slark@latest
```

### Basic Usage

```
slark --project-path=/path/to/project --platform=vercel --telegram-chat-id=123456789
```

### Commands

- `init`: Initialize a new CICD pipeline
- `list-templates`: Show available workflow templates
- `update`: Update an existing CICD configuration

### Flags

- `--project-path`: Path to the frontend project (required)
- `--platform`: Deployment platform to use (vercel, cloudflare) (required)
- `--project-name`: Name of the project (defaults to directory name)
- `--auto-create-project`: Automatically create project on deployment platform (default: true)
- `--deploy-branch`: Branch to track for deployment (default: main)
- `--build-folder`: Folder containing deployable assets (default: dist)
- `--template`: Workflow template to use (default: basic)
- `--telegram-chat-id`: Telegram chat ID for notifications (required for notifications)
- `--telegram-thread-id`: Telegram thread ID for notifications in groups (optional)
- `--dry-run`: Preview changes without committing them
- `--no-commit`: Generate files without committing to repository

## Architecture Design

### Core Components

1. **Project Analyzer**
   - Detect project type (Vite, Next.js, React)
   - Identify build folders and commands

2. **Template Processor**
   - Load template files
   - Replace variables with user input
   - Generate final workflow files

3. **Secret Manager**
   - Generate or request API keys
   - Configure GitHub repository secrets for:
     - Deployment platform credentials (Vercel, Cloudflare)
     - Telegram Bot Token
     - Other environment variables

4. **Platform Integrator**
   - Create new projects via platform APIs
   - Retrieve project IDs and deployment URLs
   - Configure platform-specific settings

5. **Notification Setup**
   - Configure Telegram notification workflows
   - Set up notification templates for different statuses

6. **Git Interface**
   - Create branch for CICD changes
   - Commit and push changes

### Data Flow

1. User provides project details, platform choice, and notification preferences
2. Tool analyzes project to determine framework type
3. If auto-create is enabled, tool creates project on deployment platform
4. Appropriate templates are selected and populated
5. GitHub secrets are configured for the deployment platform and notifications
6. Changes are committed to the repository

## Secret Management

### Required Secrets

The tool will set up the following GitHub secrets:

1. **Platform Secrets**
   - `VERCEL_API_TOKEN`: API token for Vercel project creation and deployment
   - `CLOUDFLARE_API_TOKEN`: API token for Cloudflare Pages deployments

2. **Notification Secrets**
   - `TELEGRAM_BOT_TOKEN`: Token for the Telegram bot that sends notifications
   - `TELEGRAM_CHAT_ID`: Chat ID where notifications will be sent

### Secret Handling

Secrets are never stored in the code or committed to the repository. The tool:

1. Prompts for secrets via CLI (with optional environment variable fallbacks)
2. Uses GitHub's REST API to securely set repository secrets
3. Encrypts secrets before transmission using GitHub's public key
4. Verifies secret creation without retrieving actual values

## Workflow Structure

The tool generates two primary workflow files:

1. **Main CICD Workflow**
   - Triggered on push to specified branch
   - Runs tests, builds, and deploys the application
   - Includes conditional steps based on project type
   - Contains deployment configuration for chosen platform

2. **Notification Workflow**
   - Called by the main workflow at key points
   - Sends status updates to Telegram
   - Customizable message templates for different statuses
   - Includes build info and deployment URLs

## Template Structure

Templates should be stored in `templates/` directory with the following structure:

```
templates/
├── vercel/
│   ├── basic.yml
│   └── advanced.yml
├── cloudflare/
│   ├── basic.yml
│   └── advanced.yml
├── notifications/
    ├── telegram.yml
    └── telegram-message-templates.md
```

Each template contains placeholders like `{{project_name}}`, `{{project_id}}`, `{{deploy_branch}}`, and `{{build_folder}}`.

## Implementation Plan

### Phase 1: Core Functionality
- Project structure creation
- Command line parsing
- Template loading and processing
- Basic Git operations

### Phase 2: Platform Integration
- Vercel API integration for project creation
- Cloudflare API integration
- Secret management for deployment platforms

### Phase 3: Notification System
- Telegram notification workflow templates
- Secret management for Telegram
- Status message customization

### Phase 4: Enhanced Features
- Project auto-detection
- Multiple template support
- Configuration validation

## Project Structure

The tool follows this project structure:

```
cicd-initializer/
├── cmd/
│   └── slark/
│       └── main.go                  # Entry point for the CLI application
├── internal/
│   ├── analyzer/                    # Project type detection
│   ├── config/                      # Configuration handling
│   ├── git/                         # Git operations
│   ├── platform/
│   │   ├── platform.go              # Platform interface
│   │   ├── vercel.go                # Vercel API integration
│   │   └── cloudflare.go            # Cloudflare API integration
│   ├── secrets/
│   │   ├── manager.go               # Secrets management interface
│   │   ├── github.go                # GitHub secrets API client
│   │   └── encryption.go            # Encryption utilities
│   ├── notification/
│   │   ├── telegram.go              # Telegram notification setup
│   │   └── templates.go             # Notification message templates
│   └── template/                    # Template processing
├── pkg/
│   ├── logger/                      # Logging utilities
│   └── utils/                       # General utilities
├── templates/                       # Workflow templates
