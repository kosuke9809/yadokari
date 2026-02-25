package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kosuke9809/yadokari/sandbox"
)

// Filter はリストのフィルター種別
type Filter int

const (
	FilterAll Filter = iota
	FilterRunning
	FilterStopped
	FilterRisky
)

func (f Filter) String() string {
	switch f {
	case FilterRunning:
		return "running"
	case FilterStopped:
		return "stopped"
	case FilterRisky:
		return "risky 🔴"
	default:
		return "all"
	}
}

type listModel struct {
	sandboxes []sandbox.Sandbox
	cursor    int
	filter    Filter
}

func newListModel() listModel {
	return listModel{}
}

func (m listModel) setSandboxes(sandboxes []sandbox.Sandbox) listModel {
	m.sandboxes = sandboxes
	filtered := m.filtered()
	if m.cursor >= len(filtered) {
		if len(filtered) > 0 {
			m.cursor = len(filtered) - 1
		} else {
			m.cursor = 0
		}
	}
	return m
}

func (m listModel) filtered() []sandbox.Sandbox {
	var result []sandbox.Sandbox
	for _, s := range m.sandboxes {
		switch m.filter {
		case FilterRunning:
			if s.State == sandbox.StateRunning {
				result = append(result, s)
			}
		case FilterStopped:
			if s.State == sandbox.StateStopped {
				result = append(result, s)
			}
		case FilterRisky:
			if s.Risk == sandbox.RiskHigh {
				result = append(result, s)
			}
		default:
			result = append(result, s)
		}
	}
	return result
}

func (m listModel) selected() *sandbox.Sandbox {
	items := m.filtered()
	if len(items) == 0 || m.cursor >= len(items) {
		return nil
	}
	s := items[m.cursor]
	return &s
}

func (m listModel) update(msg tea.KeyMsg) (listModel, tea.Cmd) {
	items := m.filtered()
	switch {
	case key.Matches(msg, keys.Up):
		if m.cursor > 0 {
			m.cursor--
		}
	case key.Matches(msg, keys.Down):
		if m.cursor < len(items)-1 {
			m.cursor++
		}
	case key.Matches(msg, keys.Filter):
		m.filter = (m.filter + 1) % 4
		m.cursor = 0
	}
	return m, nil
}

func (m listModel) view(width, _ int) string {
	var sb strings.Builder

	// ヘッダー
	header := fmt.Sprintf("%-20s %-8s %-8s %4s",
		"NAME", "AGENT", "STATUS", "RISK")
	sb.WriteString(lipgloss.NewStyle().Bold(true).Render(header) + "\n")
	sb.WriteString(strings.Repeat("─", width) + "\n")

	// フィルター表示
	filterLine := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Render(fmt.Sprintf("[filter: %s]", m.filter.String()))
	sb.WriteString(filterLine + "\n")

	items := m.filtered()
	if len(items) == 0 {
		sb.WriteString("  (no sandboxes)\n")
		return sb.String()
	}

	cursor := lipgloss.NewStyle().Reverse(true)
	normal := lipgloss.NewStyle()

	for i, s := range items {
		line := fmt.Sprintf("%-20s %-8s %-8s %4s",
			truncate(s.Name, 20),
			truncate(s.Agent, 8),
			string(s.State),
			s.Risk.String(),
		)
		if i == m.cursor {
			sb.WriteString(cursor.Render("> "+line) + "\n")
		} else {
			sb.WriteString(normal.Render("  "+line) + "\n")
		}
	}
	return sb.String()
}

func truncate(s string, n int) string {
	if len([]rune(s)) <= n {
		return s
	}
	return string([]rune(s)[:n-1]) + "…"
}
