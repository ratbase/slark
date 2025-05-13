package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"slark/internal/core"
	"slark/internal/version"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Define command line flags
	versionFlag := flag.Bool("version", false, "Print the version")
	listTemplatesFlag := flag.Bool("list-templates", false, "List available templates")
	debugFlag := flag.Bool("debug", false, "Enable debug mode")

	// Parse the flags
	flag.Parse()
	var handler slog.Handler = slog.Default().Handler()

	if *debugFlag {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()

		opts := &slog.HandlerOptions{Level: slog.LevelDebug}
		handler = slog.NewJSONHandler(f, opts)
		slog.SetDefault(slog.New(handler))
		slog.Info("Debug logging enabled", "file", "debug.log")
	}

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
	slog.Info("Starting Slark")
	p := tea.NewProgram(core.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		slog.Error("error running program", "error", err)
		os.Exit(1)
	}
}
