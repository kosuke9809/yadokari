package ui

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kosuke9809/yadokari/sandbox"
)

// Focus はどのペインがアクティブかを示す
type Focus int

const (
	FocusList Focus = iota
	FocusStream
)

// Bubble Tea メッセージ型
type sandboxesUpdatedMsg struct{ sandboxes []sandbox.Sandbox }
type sandboxErrMsg struct{ err error }
type tickMsg time.Time
type toastClearMsg struct{}

// Model はトップレベルの Bubble Tea モデル
type Model struct {
	list   listModel
	detail detailModel
	stream streamModel
	focus  Focus
	toast  string
	client sandbox.Client
	width  int
	height int
}

// New は CLIClient を使う本番用コンストラクタ
func New() Model {
	return NewWithClient(sandbox.NewCLIClient())
}

// NewWithClient は任意の Client を使うコンストラクタ（テスト用）
func NewWithClient(client sandbox.Client) Model {
	return Model{
		client: client,
		list:   newListModel(),
		detail: newDetailModel(),
		stream: newStreamModel(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tick(),
		fetchSandboxes(m.client),
	)
}

func tick() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func fetchSandboxes(client sandbox.Client) tea.Cmd {
	return func() tea.Msg {
		sandboxes, err := client.List(context.Background())
		if err != nil {
			return sandboxErrMsg{err}
		}
		return sandboxesUpdatedMsg{sandboxes}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.stream = m.stream.setSize(m.streamWidth(), m.streamHeight())
		return m, nil

	case tickMsg:
		return m, tea.Batch(tick(), fetchSandboxes(m.client))

	case sandboxesUpdatedMsg:
		m.list = m.list.setSandboxes(msg.sandboxes)
		if selected := m.list.selected(); selected != nil {
			m.detail = m.detail.setSandbox(*selected)
		}
		return m, nil

	case sandboxErrMsg:
		return m.showToast(msg.err.Error())

	case toastClearMsg:
		m.toast = ""
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Refresh):
			return m, fetchSandboxes(m.client)
		}
		if m.focus == FocusList {
			return m.updateList(msg)
		}
		return m.updateStream(msg)
	}
	return m, nil
}

func (m Model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Task 8〜12 で実装
	var cmd tea.Cmd
	m.list, cmd = m.list.update(msg)
	if selected := m.list.selected(); selected != nil {
		m.detail = m.detail.setSandbox(*selected)
	}
	return m, cmd
}

func (m Model) updateStream(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Task 10〜14 で実装
	var cmd tea.Cmd
	m.stream, cmd = m.stream.update(msg)
	return m, cmd
}

func (m Model) showToast(msg string) (Model, tea.Cmd) {
	m.toast = msg
	return m, tea.Tick(3*time.Second, func(time.Time) tea.Msg {
		return toastClearMsg{}
	})
}

func (m Model) View() string {
	if m.width == 0 {
		return "yadokari loading...\n"
	}

	lw := m.listWidth()
	rw := m.rightWidth()
	dh := m.detailHeight()
	sh := m.streamHeight()

	listPane := lipgloss.NewStyle().
		Width(lw).
		Height(m.height - 3).
		BorderStyle(lipgloss.NormalBorder()).
		BorderRight(true).
		Render(m.list.view(lw, m.height-3))

	detailPane := lipgloss.NewStyle().
		Width(rw).
		Height(dh).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Render(m.detail.view(rw, dh))

	streamPane := lipgloss.NewStyle().
		Width(rw).
		Height(sh).
		Render(m.stream.view(rw, sh))

	right := lipgloss.JoinVertical(lipgloss.Left, detailPane, streamPane)
	body := lipgloss.JoinHorizontal(lipgloss.Top, listPane, right)

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	title := titleStyle.Render("yadokari - Docker Sandboxes TUI")

	toast := ""
	if m.toast != "" {
		toast = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("⚠ "+m.toast)
	}

	return title + "\n" + body + toast
}

// レイアウト計算ヘルパー
func (m Model) listWidth() int    { return m.width / 3 }
func (m Model) rightWidth() int   { return m.width - m.listWidth() - 1 }
func (m Model) detailHeight() int { return (m.height - 3) * 2 / 3 }
func (m Model) streamHeight() int { return m.height - 3 - m.detailHeight() }
func (m Model) streamWidth() int  { return m.rightWidth() }
