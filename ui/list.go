package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kosuke9809/yadokari/sandbox"
)

type listModel struct {
	sandboxes []sandbox.Sandbox
	cursor    int
}

func newListModel() listModel { return listModel{} }

func (m listModel) setSandboxes(sandboxes []sandbox.Sandbox) listModel {
	m.sandboxes = sandboxes
	if m.cursor >= len(sandboxes) && len(sandboxes) > 0 {
		m.cursor = len(sandboxes) - 1
	}
	return m
}

func (m listModel) selected() *sandbox.Sandbox {
	if len(m.sandboxes) == 0 || m.cursor >= len(m.sandboxes) {
		return nil
	}
	s := m.sandboxes[m.cursor]
	return &s
}

func (m listModel) update(msg tea.KeyMsg) (listModel, tea.Cmd) {
	return m, nil
}

func (m listModel) view(width, height int) string {
	return "SANDBOX LIST\n"
}
