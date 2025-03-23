package core

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"slark/internal/models"
	"strings"
)

type Model struct {
	models.Model
}

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}

		if m.Stage == 2 && msg.Type == tea.KeyEnter {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		if m.Form != nil {
			m.Form = m.Form.WithWidth(msg.Width)
			return m, m.Form.Init()
		}

	case models.ProcessFinishedMsg:
		// Processing finished
		m.Stage = 2
		m.Success = msg.Success
		m.Result = msg.Result
		m.Err = msg.Err
		return m, nil

	case spinner.TickMsg:
		if m.Stage == 1 {
			// Update spinner while processing
			var cmd tea.Cmd
			m.Spinner, cmd = m.Spinner.Update(msg)
			return m, cmd
		}
	}

	// Only handle form updates when in form stage
	if m.Stage == 0 {
		form, cmd := m.Form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.Form = f

			// Check if the form has been completed
			if m.Form.State == huh.StateCompleted {
				projectName := m.Form.GetString("projectName")
				deployBranch := m.Form.GetString("deployBranch")
				buildFolder := m.Form.GetString("buildFolder")
				platform := m.Form.GetString("platform")
				platformData := models.PlatformData{
					ApiKey: m.Form.GetString("vercelToken"),
					TeamId: m.Form.GetString("vercelTeamName"),
				}
				telegramData := models.TelegramData{
					BotToken: m.Form.GetString("telegramToken"),
					ChatId:   m.Form.GetString("telegramChatId"),
				}

				if deployBranch == "" {
					deployBranch = "main"
				}

				if buildFolder == "" {
					buildFolder = "./"
				}

				m.Stage = 1
				return m, tea.Batch(
					m.Spinner.Tick,
					ProcessProject(projectName, deployBranch, buildFolder, platform, platformData, telegramData),
				)
			}
		}
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if m.Stage == 0 {
		var b strings.Builder
		b.WriteString(lipgloss.NewStyle().Bold(true).Render("Slark"))
		b.WriteString("\n\n")
		b.WriteString(m.Form.View())
		return b.String()
	} else if m.Stage == 1 {
		return m.ProcessingView()
	} else {
		return m.ResultsView()
	}
}
func (m Model) ProcessingView() string {
	return fmt.Sprintf("\n  %s Setting up your project...\n\n  This won't take long.", m.Spinner.View())
}

func (m Model) ResultsView() string {
	if m.Success {
		return fmt.Sprintf("\n%s\n\n%s\n\n%s",
			successStyle.Render("✓ Success!"),
			m.Result,
			helpStyle.Render("Press Enter to exit"))
	}

	errorMessage := "An unknown error occurred"
	if m.Err != nil {
		errorMessage = m.Err.Error()
	}
	return fmt.Sprintf("\n%s\n\n%s\n\n%s",
		errorStyle.Render("✗ Error!"),
		errorStyle.Render(errorMessage),
		helpStyle.Render("Press Enter to exit"))
}

func InitialModel() Model {
	// Setup spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))

	// Create a form with custom styling for each field
	projectNameInput := huh.NewInput().
		Key("projectName").
		Title("Project Name").
		Placeholder("slark").
		Validate(func(s string) error {
			if s == "" {
				return fmt.Errorf("project name cannot be empty")
			}
			return nil
		})

	deployBranchInput := huh.NewInput().
		Key("deployBranch").
		Title("Deploy Branch").
		Placeholder("dev")

	buildFolderInput := huh.NewInput().
		Key("buildFolder").
		Title("Build Folder").
		Placeholder("./")

	platformSelect := huh.NewSelect[string]().
		Key("platform").
		Title("Deployment Platform").
		Options(
			huh.NewOption("Vercel", "vercel"),
			huh.NewOption("Cloudflare Pages", "cloudflare"),
			// huh.NewOption("GitHub Pages", "github-pages"),
		)

	vercelProjectInput := huh.NewGroup(
		huh.NewInput().
			Key("vercelTeamName").
			Title("Your Vercel Team ID").
			Placeholder("team_xxxx"),
		huh.NewInput().
			Key("vercelToken").
			Title("Your Vercel API Token").
			EchoMode(huh.EchoModePassword),
	)

	telegramInput := huh.NewGroup(
		huh.NewInput().
			Key("BotToken").
			Title("Your Telegram Bot ID").
			EchoMode(huh.EchoModePassword),
		huh.NewInput().
			Key("ChatId").
			Title("Your Telegram Chat ID").
			Placeholder("-100"),
	)
	form := huh.NewForm(
		huh.NewGroup(
			projectNameInput,
			deployBranchInput,
			buildFolderInput,
			platformSelect,
		),
		vercelProjectInput,
    telegramInput,
	).WithShowHelp(true)

	return Model{
		models.Model{
			Form:    form,
			Spinner: s,
			Stage:   0,
		},
	}
}

func (m Model) Init() tea.Cmd {
	return m.Form.Init()
}
