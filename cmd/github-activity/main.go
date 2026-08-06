package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github-activity/internal/github"
	"github-activity/internal/ui"
)

func main() {
	client := github.NewClient()
	
	m := ui.NewModel(client)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting application: %v\n", err)
		os.Exit(1)
	}
}
