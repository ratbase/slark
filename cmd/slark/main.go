package main

import (
	// "flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
)

// Model represents the application state
type Model struct {
	form    *huh.Form
	spinner spinner.Model
	stage   int // 0: form, 1: processing, 2: results
	err     error
	success bool
	result  string
	width   int
	height  int
}

// Initialize the model
func initialModel() Model {
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

	form := huh.NewForm(
		huh.NewGroup(
			projectNameInput,
			deployBranchInput,
			buildFolderInput,
			platformSelect,
		),
	).WithShowHelp(true)

	return Model{
		form:    form,
		spinner: s,
		stage:   0,
	}
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

type processFinishedMsg struct {
	success bool
	result  string
	err     error
}

func processProject(projectName, deployBranch, buildFolder, platform string) tea.Cmd {
	return func() tea.Msg {
		// Simulate processing
		success := true

		// Create a detailed result message
		var resultBuilder strings.Builder
		resultBuilder.WriteString(fmt.Sprintf("Project: %s\n", projectName))
		resultBuilder.WriteString(fmt.Sprintf("Deploy Branch: %s\n", deployBranch))
		resultBuilder.WriteString(fmt.Sprintf("Build Folder: %s\n", buildFolder))
		resultBuilder.WriteString(fmt.Sprintf("Platform: %s\n", platform))
		resultBuilder.WriteString("\nCI/CD pipeline configured successfully!")

		var err error

		// In a real implementation, this would call core.Setup or similar function

		return processFinishedMsg{
			success: success,
			result:  resultBuilder.String(),
			err:     err,
		}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}

		if m.stage == 2 && msg.Type == tea.KeyEnter {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.form != nil {
			m.form = m.form.WithWidth(msg.Width)
			return m, m.form.Init()
		}

	case processFinishedMsg:
		// Processing finished
		m.stage = 2
		m.success = msg.success
		m.result = msg.result
		m.err = msg.err
		return m, nil

	case spinner.TickMsg:
		if m.stage == 1 {
			// Update spinner while processing
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	// Only handle form updates when in form stage
	if m.stage == 0 {
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f

			// Check if the form has been completed
			if m.form.State == huh.StateCompleted {
				projectName := m.form.GetString("projectName")
				deployBranch := m.form.GetString("deployBranch")
				buildFolder := m.form.GetString("buildFolder")
				platform := m.form.GetString("platform")

				if deployBranch == "" {
					deployBranch = "main"
				}

				if buildFolder == "" {
					buildFolder = "./"
				}

				m.stage = 1
				return m, tea.Batch(
					m.spinner.Tick,
					processProject(projectName, deployBranch, buildFolder, platform),
				)
			}
		}
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if m.stage == 0 {
		var b strings.Builder
		b.WriteString(lipgloss.NewStyle().Bold(true).Render("Slark"))
		b.WriteString("\n\n")
		b.WriteString(m.form.View())
		return b.String()
	} else if m.stage == 1 {
		return m.processingView()
	} else {
		return m.resultsView()
	}
}

func (m Model) processingView() string {
	return fmt.Sprintf("\n  %s Setting up your project...\n\n  This won't take long.", m.spinner.View())
}

func (m Model) resultsView() string {
	if m.success {
		return fmt.Sprintf("\n%s\n\n%s\n\n%s",
			successStyle.Render("✓ Success!"),
			m.result,
			helpStyle.Render("Press Enter to exit"))
	}

	errorMessage := "An unknown error occurred"
	if m.err != nil {
		errorMessage = m.err.Error()
	}

	return fmt.Sprintf("\n%s\n\n%s\n\n%s",
		errorStyle.Render("✗ Error!"),
		errorStyle.Render(errorMessage),
		helpStyle.Render("Press Enter to exit"))
}

func main() {
	// // Check for CLI flags
	// cliMode := flag.Bool("cli", false, "Run in CLI mode")
	// projectName := flag.String("name", "", "Project name")
	// platform := flag.String("platform", "vercel", "Deployment platform")
	// branch := flag.String("branch", "", "Deploy branch")
	// buildFolder := flag.String("build-folder", "", "Build folder")
	//
	// flag.Parse()
	//
	// // Handle CLI mode
	// if *cliMode {
	// 	if *projectName == "" {
	// 		fmt.Println("Error: Project name is required in CLI mode")
	// 		fmt.Println("Usage: slark --cli --name=<project_name>")
	// 		os.Exit(1)
	// 	}
	//
	// 	// Set default values if empty
	// 	deployBranch := *branch
	// 	if deployBranch == "" {
	// 		deployBranch = "main"
	// 	}
	//
	// 	buildDir := *buildFolder
	// 	if buildDir == "" {
	// 		buildDir = "dist"
	// 	}
	//
	// 	fmt.Printf("Setting up project '%s' in CLI mode\n", *projectName)
	// 	fmt.Printf("Platform: %s, Branch: %s, Build folder: %s\n", *platform, deployBranch, buildDir)
	// 	fmt.Println("Project created successfully!")
	// 	return
	// }

	// Run TUI mode
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
