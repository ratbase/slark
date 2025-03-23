package main

import (
	"fmt"
	"os"

	"slark/internal/core"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(core.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
