package ui

import (
	"bufio"
	"context"
	"io"
	"strings"
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
type logLineMsg struct{ line string }
type logErrMsg struct{ err error }
type logDoneMsg struct{}

// Model はトップレベルの Bubble Tea モデル
type Model struct {
	list      listModel
	detail    detailModel
	stream    streamModel
	confirm   confirmModel
	help      helpModel
	focus     Focus
	toast     string
	client    sandbox.Client
	width     int
	height    int
	logReader io.ReadCloser
	logBuf    *bufio.Reader
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

func readNextLine(r io.Reader, buf *bufio.Reader) tea.Cmd {
	return func() tea.Msg {
		line, err := buf.ReadString('\n')
		if len(line) > 0 {
			return logLineMsg{line: strings.TrimRight(line, "\n")}
		}
		if err != nil {
			if err == io.EOF {
				return logDoneMsg{}
			}
			return logErrMsg{err}
		}
		return logDoneMsg{}
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

	case fetchSandboxesMsg:
		return m, fetchSandboxes(m.client)

	case sandboxErrMsg:
		return m.showToast(msg.err.Error())

	case toastClearMsg:
		m.toast = ""
		return m, nil

	case logLineMsg:
		m.stream = m.stream.addLine(msg.line)
		return m, readNextLine(m.logReader, m.logBuf)

	case logErrMsg:
		if msg.err.Error() != "logs not supported by docker sandbox CLI" {
			return m.showToast(msg.err.Error())
		}
		return m, nil

	case logDoneMsg:
		if m.logReader != nil {
			m.logReader.Close()
			m.logReader = nil
			m.logBuf = nil
		}
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Refresh):
			return m, fetchSandboxes(m.client)
		case key.Matches(msg, keys.Help):
			m.help = m.help.toggle()
			return m, nil
		}
		if m.focus == FocusList {
			return m.updateList(msg)
		}
		return m.updateStream(msg)
	}
	return m, nil
}

func (m Model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// confirm が表示中はダイアログに委譲
	if m.confirm.visible {
		var cmd tea.Cmd
		m.confirm, cmd = m.confirm.update(msg)
		return m, cmd
	}

	switch {
	case key.Matches(msg, keys.Start):
		if s := m.list.selected(); s != nil {
			return m, m.toggleStartStop(s)
		}
	case key.Matches(msg, keys.Restart):
		if s := m.list.selected(); s != nil {
			return m, sandboxOp(m.client.Restart, s.ID)
		}
	case key.Matches(msg, keys.Remove):
		if s := m.list.selected(); s != nil {
			id, name := s.ID, s.Name
			m.confirm = m.confirm.show(
				"Remove \""+name+"\"?",
				func() tea.Cmd { return sandboxOp(m.client.Remove, id) },
			)
		}
	case key.Matches(msg, keys.RawInspect):
		m.detail = m.detail.toggleRaw()
	case key.Matches(msg, keys.Logs):
		if s := m.list.selected(); s != nil {
			// 既存のログリーダーを閉じる
			if m.logReader != nil {
				m.logReader.Close()
				m.logReader = nil
				m.logBuf = nil
			}
			m.focus = FocusStream
			m.stream = m.stream.clear().setMode(StreamLogs)
			r, err := m.client.Logs(context.Background(), s.ID)
			if err != nil {
				if err.Error() != "logs not supported by docker sandbox CLI" {
					return m.showToast(err.Error())
				}
				m.stream = m.stream.addLine("[logs not supported for this sandbox]")
				return m, nil
			}
			m.logReader = r
			m.logBuf = bufio.NewReader(r)
			return m, readNextLine(m.logReader, m.logBuf)
		}
	case key.Matches(msg, keys.Exec):
		if s := m.list.selected(); s != nil {
			return m, m.execShell(s.ID)
		}
	case key.Matches(msg, keys.ExecCmd):
		// TODO: 1行コマンド実行モード（将来実装）
		return m.showToast("ExecCmd not yet implemented")
	default:
		var cmd tea.Cmd
		m.list, cmd = m.list.update(msg)
		if selected := m.list.selected(); selected != nil {
			m.detail = m.detail.setSandbox(*selected)
		}
		return m, cmd
	}
	return m, nil
}

func sandboxOp(fn func(context.Context, string) error, id string) tea.Cmd {
	return func() tea.Msg {
		if err := fn(context.Background(), id); err != nil {
			return sandboxErrMsg{err}
		}
		return fetchSandboxesMsg{}
	}
}

type fetchSandboxesMsg struct{}

func (m Model) toggleStartStop(s *sandbox.Sandbox) tea.Cmd {
	if s.State == sandbox.StateRunning {
		return sandboxOp(m.client.Stop, s.ID)
	}
	return sandboxOp(m.client.Start, s.ID)
}

func (m Model) execShell(id string) tea.Cmd {
	cmd, err := m.client.Exec(context.Background(), id)
	if err != nil {
		return func() tea.Msg { return sandboxErrMsg{err} }
	}
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return sandboxErrMsg{err}
		}
		return nil
	})
}

func (m Model) updateStream(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Back):
		m.focus = FocusList
		return m, nil
	}
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

	result := title + "\n" + body + toast
	if m.confirm.visible {
		result += "\n\n" + m.confirm.view()
	}
	if m.help.visible {
		result += "\n\n" + m.help.view()
	}
	return result
}

// レイアウト計算ヘルパー
func (m Model) listWidth() int    { return m.width / 3 }
func (m Model) rightWidth() int   { return m.width - m.listWidth() - 1 }
func (m Model) detailHeight() int { return (m.height - 3) * 2 / 3 }
func (m Model) streamHeight() int { return m.height - 3 - m.detailHeight() }
func (m Model) streamWidth() int  { return m.rightWidth() }
