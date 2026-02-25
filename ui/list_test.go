package ui

import (
	"testing"

	"github.com/kosuke9809/yadokari/sandbox"
)

func TestListFilter(t *testing.T) {
	sandboxes := []sandbox.Sandbox{
		{ID: "1", Name: "a", State: sandbox.StateRunning, Risk: sandbox.RiskDev},
		{ID: "2", Name: "b", State: sandbox.StateStopped, Risk: sandbox.RiskStrict},
		{ID: "3", Name: "c", State: sandbox.StateRunning, Risk: sandbox.RiskHigh},
	}

	m := newListModel()
	m = m.setSandboxes(sandboxes)

	// FilterAll
	m.filter = FilterAll
	if got := len(m.filtered()); got != 3 {
		t.Errorf("FilterAll: want 3, got %d", got)
	}

	// FilterRunning
	m.filter = FilterRunning
	if got := len(m.filtered()); got != 2 {
		t.Errorf("FilterRunning: want 2, got %d", got)
	}

	// FilterStopped
	m.filter = FilterStopped
	if got := len(m.filtered()); got != 1 {
		t.Errorf("FilterStopped: want 1, got %d", got)
	}

	// FilterRisky
	m.filter = FilterRisky
	if got := len(m.filtered()); got != 1 {
		t.Errorf("FilterRisky: want 1, got %d", got)
	}
}

func TestListSelected(t *testing.T) {
	sandboxes := []sandbox.Sandbox{
		{ID: "1", Name: "first"},
		{ID: "2", Name: "second"},
	}
	m := newListModel()
	m = m.setSandboxes(sandboxes)

	s := m.selected()
	if s == nil {
		t.Fatal("expected selected sandbox, got nil")
	}
	if s.Name != "first" {
		t.Errorf("want first, got %s", s.Name)
	}
}

func TestListCursorBounds(t *testing.T) {
	sandboxes := []sandbox.Sandbox{
		{ID: "1", Name: "only"},
	}
	m := newListModel()
	m = m.setSandboxes(sandboxes)

	// カーソルが範囲外にならないことを確認
	m.cursor = 10
	m = m.setSandboxes(sandboxes) // setSandboxes でカーソルを調整
	if m.cursor != 0 {
		t.Errorf("cursor should be adjusted to 0, got %d", m.cursor)
	}
}
