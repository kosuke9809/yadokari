package ui

import (
	"strings"
	"testing"
)

func TestStreamAddLine(t *testing.T) {
	m := newStreamModel()
	m = m.setSize(80, 20)
	m = m.addLine("line 1")
	m = m.addLine("line 2")
	m = m.addLine("line 3")

	if len(m.lines) != 3 {
		t.Errorf("want 3 lines, got %d", len(m.lines))
	}
}

func TestStreamClear(t *testing.T) {
	m := newStreamModel()
	m = m.setSize(80, 20)
	m = m.addLine("line 1")
	m = m.clear()

	if len(m.lines) != 0 {
		t.Errorf("want 0 lines after clear, got %d", len(m.lines))
	}
}

func TestStreamTogglePause(t *testing.T) {
	m := newStreamModel()
	if m.paused {
		t.Error("should not be paused initially")
	}
	m = m.togglePause()
	if !m.paused {
		t.Error("should be paused after toggle")
	}
	m = m.togglePause()
	if m.paused {
		t.Error("should not be paused after second toggle")
	}
}

func TestStreamView(t *testing.T) {
	m := newStreamModel()
	m = m.setSize(80, 10)
	m = m.addLine("hello from sandbox")
	got := m.view(80, 10)
	if !strings.Contains(got, "STREAM") {
		t.Errorf("view should contain STREAM, got: %s", got)
	}
}
