package sandbox

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// MockClient はテスト用の Client 実装
type MockClient struct {
	Sandboxes  []Sandbox
	Err        error
	LogsOutput string

	// 呼び出し記録
	StartedIDs   []string
	StoppedIDs   []string
	RestartedIDs []string
	RemovedIDs   []string
}

func (m *MockClient) List(_ context.Context) ([]Sandbox, error) {
	return m.Sandboxes, m.Err
}

func (m *MockClient) Inspect(_ context.Context, id string) (Sandbox, error) {
	if m.Err != nil {
		return Sandbox{}, m.Err
	}
	for _, s := range m.Sandboxes {
		if s.ID == id || s.Name == id {
			return s, nil
		}
	}
	return Sandbox{}, fmt.Errorf("sandbox not found: %s", id)
}

func (m *MockClient) Start(_ context.Context, id string) error {
	m.StartedIDs = append(m.StartedIDs, id)
	return m.Err
}

func (m *MockClient) Stop(_ context.Context, id string) error {
	m.StoppedIDs = append(m.StoppedIDs, id)
	return m.Err
}

func (m *MockClient) Restart(_ context.Context, id string) error {
	m.RestartedIDs = append(m.RestartedIDs, id)
	return m.Err
}

func (m *MockClient) Remove(_ context.Context, id string) error {
	m.RemovedIDs = append(m.RemovedIDs, id)
	return m.Err
}

func (m *MockClient) Logs(_ context.Context, _ string) (io.ReadCloser, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return io.NopCloser(strings.NewReader(m.LogsOutput)), nil
}

func (m *MockClient) Exec(_ context.Context, _ string) (*exec.Cmd, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return nil, errors.New("exec not available in mock")
}

// SampleSandboxes は開発・テスト用のサンプルデータを返す
func SampleSandboxes() []Sandbox {
	return []Sandbox{
		{
			ID:        "claude-yadokari",
			Name:      "claude-yadokari",
			Agent:     "claude",
			State:     StateRunning,
			CPU:       "23%",
			Mem:       "1.2G",
			Risk:      RiskDev,
			Uptime:    "2h",
			Workspace: "/home/user/project",
		},
		{
			ID:        "codex-yadokari",
			Name:      "codex-yadokari",
			Agent:     "codex",
			State:     StateStopped,
			CPU:       "-",
			Mem:       "-",
			Risk:      RiskStrict,
			Uptime:    "-",
			Workspace: "/home/user/project",
		},
	}
}
