package models

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/huh"
)

type Model struct {
	Form    *huh.Form
	Spinner spinner.Model
	Stage   int // 0: form, 1: processing, 2: results
	Err     error
	Success bool
	Result  string
	Width   int
	Height  int
}

type PlatformData struct {
	ApiKey string
	TeamId string
}

type TelegramData struct {
	BotToken string
	ChatId   string
}

type ProcessFinishedMsg struct {
	Success bool
	Result  string
	Err     error
}

// ProjectConfig represents the configuration for a project setup
type ProjectConfig struct {
	Name         string
	DeployBranch string
	BuildFolder  string
	Platform     string
	CreatedAt    time.Time
}
