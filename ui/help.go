package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type helpModel struct {
	visible bool
}

func (m helpModel) toggle() helpModel {
	m.visible = !m.visible
	return m
}

func (m helpModel) view() string {
	if !m.visible {
		return ""
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 3).
		Width(50)

	rows := []string{
		lipgloss.NewStyle().Bold(true).Render("Keyboard Shortcuts"),
		"",
		"j / ↓    Move down",
		"k / ↑    Move up",
		"s        Start / Stop toggle",
		"r        Restart",
		"d        Remove (with confirm)",
		"i        Inspect (refresh detail)",
		"I        Raw inspect JSON",
		"l        Logs mode",
		"e        Exec shell",
		"E        Exec command (TODO)",
		"f        Cycle filter",
		"R        Manual refresh",
		"/        Search (in logs)",
		"ESC      Back to list",
		"?        Toggle this help",
		"q        Quit",
	}

	return style.Render(strings.Join(rows, "\n"))
}
