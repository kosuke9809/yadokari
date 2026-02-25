package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kosuke9809/yadokari/sandbox"
)

func TestModelInitWithMock(t *testing.T) {
	mock := &sandbox.MockClient{
		Sandboxes: sandbox.SampleSandboxes(),
	}
	m := NewWithClient(mock)
	if m.client == nil {
		t.Error("client should not be nil")
	}
}

func TestModelViewLoading(t *testing.T) {
	mock := &sandbox.MockClient{}
	m := NewWithClient(mock)
	got := m.View()
	if !strings.Contains(got, "loading") {
		t.Errorf("want 'loading' in view, got: %s", got)
	}
}

func TestModelUpdateSandboxes(t *testing.T) {
	mock := &sandbox.MockClient{
		Sandboxes: sandbox.SampleSandboxes(),
	}
	m := NewWithClient(mock)

	// WindowSizeMsg でサイズ設定
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(Model)

	// サンドボックス更新メッセージ
	updated, _ = m.Update(sandboxesUpdatedMsg{sandboxes: sandbox.SampleSandboxes()})
	m = updated.(Model)

	got := m.View()
	if !strings.Contains(got, "yadokari") {
		t.Errorf("want 'yadokari' title in view, got: %s", got)
	}
	// 一覧にサンドボックスが表示されているか
	if !strings.Contains(got, "claude") {
		t.Errorf("want 'claude' agent in view, got: %s", got)
	}
}
