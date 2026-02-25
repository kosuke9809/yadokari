package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestStreamSearch(t *testing.T) {
	m := newStreamModel()
	m = m.setSize(80, 20)
	m = m.addLine("hello world")
	m = m.addLine("goodbye world")
	m = m.addLine("hello again")

	// 検索モードに入る
	updated, _ := m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")})
	if !updated.searching {
		t.Error("should be in searching mode after /")
	}

	// "hello" と入力
	updated, _ = updated.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	updated, _ = updated.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("e")})
	updated, _ = updated.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	if updated.searchInput.Value() != "hel" {
		t.Errorf("want search value 'hel', got %q", updated.searchInput.Value())
	}

	// ESC で検索終了
	updated, _ = updated.update(tea.KeyMsg{Type: tea.KeyEsc})
	if updated.searching {
		t.Error("should exit searching mode after ESC")
	}
}

func TestRenderLinesHighlight(t *testing.T) {
	m := newStreamModel()
	m = m.setSize(80, 20)
	m = m.addLine("match this line")
	m = m.addLine("no match here")
	m.searchInput.SetValue("match")

	rendered := m.renderLines()
	if rendered == "" {
		t.Error("renderLines should not be empty")
	}
}
