package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type confirmModel struct {
	visible bool
	message string
	onYes   func() tea.Cmd
}

func (m confirmModel) show(msg string, onYes func() tea.Cmd) confirmModel {
	m.visible = true
	m.message = msg
	m.onYes = onYes
	return m
}

func (m confirmModel) hide() confirmModel {
	m.visible = false
	m.message = ""
	m.onYes = nil
	return m
}

func (m confirmModel) update(msg tea.KeyMsg) (confirmModel, tea.Cmd) {
	if !m.visible {
		return m, nil
	}
	switch msg.String() {
	case "y", "Y":
		cmd := m.onYes()
		return m.hide(), cmd
	default:
		return m.hide(), nil
	}
}

func (m confirmModel) view() string {
	if !m.visible {
		return ""
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 3).
		Render(m.message + "\n\n[y] Yes  [any] Cancel")
}
