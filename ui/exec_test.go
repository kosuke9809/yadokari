package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kosuke9809/yadokari/sandbox"
)

func TestExecKeyHandled(t *testing.T) {
	mock := &sandbox.MockClient{
		Sandboxes: sandbox.SampleSandboxes(),
	}
	m := NewWithClient(mock)

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(Model)
	updated, _ = m.Update(sandboxesUpdatedMsg{sandboxes: sandbox.SampleSandboxes()})
	m = updated.(Model)

	// e キーで exec（MockClient の Exec は error を返すが tea.ExecProcess は Cmd を返す）
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("e")})
	// cmd が返ってくる（nil ではない）
	if cmd == nil {
		t.Error("expected a cmd for exec")
	}
}

func TestBackToList(t *testing.T) {
	mock := &sandbox.MockClient{
		Sandboxes: sandbox.SampleSandboxes(),
	}
	m := NewWithClient(mock)

	// FocusStream に切り替える
	m.focus = FocusStream

	// ESC でリストに戻る
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(Model)

	if m.focus != FocusList {
		t.Errorf("want FocusList after ESC, got %d", m.focus)
	}
}
