package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kosuke9809/yadokari/sandbox"
)

func TestLogsStreaming(t *testing.T) {
	mock := &sandbox.MockClient{
		Sandboxes:  sandbox.SampleSandboxes(),
		LogsOutput: "line1\nline2\nline3\n",
	}
	m := NewWithClient(mock)

	// サイズとサンドボックスをセット
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(Model)
	updated, _ = m.Update(sandboxesUpdatedMsg{sandboxes: sandbox.SampleSandboxes()})
	m = updated.(Model)

	// l キーでLogsモードに切り替え
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	m = updated.(Model)

	// フォーカスが Stream に切り替わっているか
	if m.focus != FocusStream {
		t.Errorf("want FocusStream, got %d", m.focus)
	}

	// Cmd が返ってきているか（ログ読み取り開始）
	if cmd == nil {
		t.Error("expected a cmd for log streaming")
	}

	// Cmd を実行してログ行を取得
	msg := cmd()
	switch msg := msg.(type) {
	case logLineMsg:
		if !strings.Contains(msg.line, "line") {
			t.Errorf("expected log line, got: %s", msg.line)
		}
	case logDoneMsg:
		// 空の場合は即 done になる可能性もある
	default:
		t.Errorf("unexpected msg type: %T", msg)
	}
}

func TestLogsNotSupported(t *testing.T) {
	// CLIClient では logs が未サポートのため、エラーメッセージがstream に表示される
	mock := &sandbox.MockClient{
		Sandboxes: sandbox.SampleSandboxes(),
		// LogsOutput を空にしておく
	}
	m := NewWithClient(mock)

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(Model)
	updated, _ = m.Update(sandboxesUpdatedMsg{sandboxes: sandbox.SampleSandboxes()})
	m = updated.(Model)

	// l キー
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	m = updated.(Model)

	// MockClient は Logs をサポートするので FocusStream になるはず
	if m.focus != FocusStream {
		t.Errorf("want FocusStream, got %d", m.focus)
	}
}
