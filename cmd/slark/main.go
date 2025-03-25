package main

import (
	"flag"
	"fmt"
	"os"

	"slark/internal/core"
	"slark/internal/version"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Define command line flags
	versionFlag := flag.Bool("version", false, "Print the version")
	listTemplatesFlag := flag.Bool("list-templates", false, "List available templates")

	// Parse the flags
	flag.Parse()

	// Check for version flag
	if *versionFlag {
		fmt.Printf("slark version %s\n", version.GetVersion())
		return
	}

	// Check for list templates flag
	if *listTemplatesFlag {
		core.ListTemplates()
		return
	}

	// Run the main program
	p := tea.NewProgram(core.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
