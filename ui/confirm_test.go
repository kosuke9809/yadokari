package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestConfirmShow(t *testing.T) {
	m := confirmModel{}
	if m.visible {
		t.Error("should not be visible initially")
	}

	m = m.show("Remove sandbox?", func() tea.Cmd { return nil })
	if !m.visible {
		t.Error("should be visible after show")
	}
	if m.message != "Remove sandbox?" {
		t.Errorf("want message 'Remove sandbox?', got %q", m.message)
	}
}

func TestConfirmYes(t *testing.T) {
	called := false
	m := confirmModel{}
	m = m.show("confirm?", func() tea.Cmd {
		called = true
		return nil
	})

	// y キーで確認
	updated, cmd := m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})
	if cmd != nil {
		cmd() // onYes を実行
	}
	_ = updated
	if !called {
		t.Error("onYes should have been called")
	}
	if updated.visible {
		t.Error("should be hidden after yes")
	}
}

func TestConfirmCancel(t *testing.T) {
	m := confirmModel{}
	m = m.show("confirm?", func() tea.Cmd { return nil })

	// n キーでキャンセル
	updated, _ := m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	if updated.visible {
		t.Error("should be hidden after cancel")
	}
}
