package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kosuke9809/yadokari/sandbox"
)

func TestHelpToggle(t *testing.T) {
	m := helpModel{}
	if m.visible {
		t.Error("help should not be visible initially")
	}
	m = m.toggle()
	if !m.visible {
		t.Error("help should be visible after toggle")
	}
	got := m.view()
	if !strings.Contains(got, "Keyboard Shortcuts") {
		t.Errorf("help view should contain 'Keyboard Shortcuts', got: %s", got)
	}
	m = m.toggle()
	if m.visible {
		t.Error("help should be hidden after second toggle")
	}
	if m.view() != "" {
		t.Error("hidden help view should be empty")
	}
}

func TestHelpKeyInModel(t *testing.T) {
	mock := &sandbox.MockClient{}
	m := NewWithClient(mock)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
	m = updated.(Model)
	if !m.help.visible {
		t.Error("help should be visible after ? key")
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
	m = updated.(Model)
	if m.help.visible {
		t.Error("help should be hidden after second ? key")
	}
}
