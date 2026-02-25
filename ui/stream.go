package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// StreamMode はストリームの表示モード
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

func (m streamModel) setMode(mode StreamMode) streamModel {
	m.mode = mode
	return m
}

func (m streamModel) addLine(line string) streamModel {
	m.lines = append(m.lines, line)
	if !m.paused {
		m.viewport.SetContent(strings.Join(m.lines, "\n"))
		m.viewport.GotoBottom()
	}
	return m
}

func (m streamModel) clear() streamModel {
	m.lines = nil
	m.viewport.SetContent("")
	return m
}

func (m streamModel) togglePause() streamModel {
	m.paused = !m.paused
	if !m.paused {
		// resume: 最新へ
		m.viewport.SetContent(strings.Join(m.lines, "\n"))
		m.viewport.GotoBottom()
	}
	return m
}

func (m streamModel) update(msg tea.KeyMsg) (streamModel, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Start): // space キーで pause/resume（keys.Start は "s"）
		// pause/resume は space キーではなく専用に対応。ここでは viewport に委譲
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
	return header + "\n" + m.viewport.View()
}
