package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kosuke9809/yadokari/sandbox"
	"github.com/kosuke9809/yadokari/ui"
)

func main() {
	var model ui.Model
	if os.Getenv("YADOKARI_MOCK") != "" {
		model = ui.NewWithClient(&sandbox.MockClient{
			Sandboxes: sandbox.SampleSandboxes(),
		})
	} else {
		model = ui.New()
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
