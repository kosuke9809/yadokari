package ui

import (
	"strings"
	"testing"

	"github.com/kosuke9809/yadokari/sandbox"
)

func TestDetailViewEmpty(t *testing.T) {
	m := newDetailModel()
	got := m.view(80, 40)
	if !strings.Contains(got, "No sandbox selected") {
		t.Errorf("want 'No sandbox selected', got: %s", got)
	}
}

func TestDetailViewSandbox(t *testing.T) {
	m := newDetailModel()
	m = m.setSandbox(sandbox.Sandbox{
		ID:        "abc",
		Name:      "refactor1",
		Agent:     "claude",
		State:     sandbox.StateRunning,
		Workspace: "/home/user/project",
	})
	got := m.view(80, 40)
	if !strings.Contains(got, "refactor1") {
		t.Errorf("want name in view, got: %s", got)
	}
	if !strings.Contains(got, "claude") {
		t.Errorf("want agent in view, got: %s", got)
	}
}

func TestDetailToggleRaw(t *testing.T) {
	m := newDetailModel()
	m = m.setSandbox(sandbox.Sandbox{ID: "x", Name: "test"})
	if m.showRaw {
		t.Error("showRaw should be false initially")
	}
	m = m.toggleRaw()
	if !m.showRaw {
		t.Error("showRaw should be true after toggle")
	}
	got := m.view(80, 40)
	// raw view は JSON を含む
	if !strings.Contains(got, "test") {
		t.Errorf("raw view should contain sandbox name, got: %s", got)
	}
}
