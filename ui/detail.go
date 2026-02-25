package ui

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kosuke9809/yadokari/sandbox"
)

type detailModel struct {
	sandbox sandbox.Sandbox
	showRaw bool
}

func newDetailModel() detailModel { return detailModel{} }

func (m detailModel) setSandbox(s sandbox.Sandbox) detailModel {
	m.sandbox = s
	return m
}

func (m detailModel) toggleRaw() detailModel {
	m.showRaw = !m.showRaw
	return m
}

func (m detailModel) view(width, _ int) string {
	s := m.sandbox
	if s.ID == "" {
		return "No sandbox selected\n"
	}

	if m.showRaw {
		return m.rawView(width)
	}

	labelStyle := lipgloss.NewStyle().Bold(true).Width(12)
	row := func(label, value string) string {
		return labelStyle.Render(label+":") + " " + value + "\n"
	}

	var sb strings.Builder
	sb.WriteString(row("Name", s.Name))
	sb.WriteString(row("Agent", s.Agent))
	sb.WriteString(row("Status", string(s.State)))
	sb.WriteString(row("Risk", s.Risk.String()))
	if s.Workspace != "" {
		sb.WriteString(row("Workspace", s.Workspace))
	}
	if s.Template != "" {
		sb.WriteString(row("Template", s.Template))
	}
	if s.Network != "" {
		sb.WriteString(row("Network", s.Network))
	}
	for _, mount := range s.Mounts {
		mode := "ro"
		if !mount.ReadOnly {
			mode = "rw"
		}
		sb.WriteString(row("Mount", fmt.Sprintf("%s (%s)", mount.Path, mode)))
	}
	if s.LastError != "" {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
		sb.WriteString(labelStyle.Render("LastError:") + " " + errStyle.Render(s.LastError) + "\n")
	}
	return sb.String()
}

func (m detailModel) rawView(_ int) string {
	if len(m.sandbox.Raw) == 0 {
		// Raw フィールドがなければ JSON に変換して表示
		data, err := json.MarshalIndent(m.sandbox, "", "  ")
		if err != nil {
			return "failed to marshal sandbox\n"
		}
		return string(data) + "\n"
	}
	var pretty []byte
	var err error
	var v any
	if err = json.Unmarshal(m.sandbox.Raw, &v); err == nil {
		pretty, err = json.MarshalIndent(v, "", "  ")
	}
	if err != nil {
		return string(m.sandbox.Raw) + "\n"
	}
	return string(pretty) + "\n"
}
