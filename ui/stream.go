package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/viewport"
)

type StreamMode int

const (
	StreamLogs StreamMode = iota
	StreamExec
)

type streamModel struct {
	viewport viewport.Model
	mode     StreamMode
	paused   bool
	lines    []string
}

func newStreamModel() streamModel {
	return streamModel{viewport: viewport.New(0, 0)}
}

func (m streamModel) setSize(width, height int) streamModel {
	m.viewport.Width = width
	m.viewport.Height = height
	return m
}

func (m streamModel) clear() streamModel {
	m.lines = nil
	m.viewport.SetContent("")
	return m
}

func (m streamModel) update(msg tea.KeyMsg) (streamModel, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m streamModel) view(width, height int) string {
	return "STREAM\n" + m.viewport.View()
}
