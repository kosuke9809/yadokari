package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StreamMode int

const (
	StreamLogs StreamMode = iota
	StreamExec
)

type streamModel struct {
	viewport    viewport.Model
	mode        StreamMode
	paused      bool
	lines       []string
	searching   bool
	searchInput textinput.Model
}

func newStreamModel() streamModel {
	ti := textinput.New()
	ti.Placeholder = "search..."
	ti.CharLimit = 100
	return streamModel{
		viewport:    viewport.New(0, 0),
		searchInput: ti,
	}
}

func (m streamModel) setSize(width, height int) streamModel {
	m.viewport.Width = width
	m.viewport.Height = height
	return m
}

func (m streamModel) setMode(mode StreamMode) streamModel {
	m.mode = mode
	return m
}

func (m streamModel) addLine(line string) streamModel {
	m.lines = append(m.lines, line)
	if !m.paused {
		m.viewport.SetContent(m.renderLines())
		m.viewport.GotoBottom()
	}
	return m
}

func (m streamModel) renderLines() string {
	query := m.searchInput.Value()
	if query == "" {
		return strings.Join(m.lines, "\n")
	}
	highlight := lipgloss.NewStyle().Background(lipgloss.Color("3")).Foreground(lipgloss.Color("0"))
	var result []string
	for _, line := range m.lines {
		if strings.Contains(line, query) {
			result = append(result, highlight.Render(line))
		} else {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

func (m streamModel) clear() streamModel {
	m.lines = nil
	m.searching = false
	m.searchInput.SetValue("")
	m.viewport.SetContent("")
	return m
}

func (m streamModel) togglePause() streamModel {
	m.paused = !m.paused
	if !m.paused {
		m.viewport.SetContent(m.renderLines())
		m.viewport.GotoBottom()
	}
	return m
}

func (m streamModel) update(msg tea.KeyMsg) (streamModel, tea.Cmd) {
	if m.searching {
		switch msg.Type {
		case tea.KeyEsc, tea.KeyEnter:
			m.searching = false
			m.searchInput.Blur()
			m.viewport.SetContent(m.renderLines())
			return m, nil
		}
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		m.viewport.SetContent(m.renderLines())
		return m, cmd
	}

	switch {
	case key.Matches(msg, keys.Search):
		m.searching = true
		m.searchInput.Focus()
		return m, textinput.Blink
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m streamModel) view(_, _ int) string {
	modeStr := "Logs"
	if m.mode == StreamExec {
		modeStr = "Exec"
	}
	paused := ""
	if m.paused {
		paused = lipgloss.NewStyle().
			Foreground(lipgloss.Color("3")).
			Render(" [PAUSED]")
	}
	header := lipgloss.NewStyle().Bold(true).Render("STREAM (mode: "+modeStr+")") + paused

	searchBar := ""
	if m.searching {
		searchBar = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("/") + m.searchInput.View()
	} else if m.searchInput.Value() != "" {
		searchBar = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("filter: "+m.searchInput.Value())
	}

	return header + searchBar + "\n" + m.viewport.View()
}
