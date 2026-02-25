package sandbox

import (
	"context"
	"io"
	"os/exec"
)

type State string

const (
	StateRunning State = "running"
	StateStopped State = "stopped"
	StateCrashed State = "crashed"
)

type Risk int

const (
	RiskStrict Risk = iota // 🔒 network none + read-only
	RiskDev                // 🟡 限定ネットワーク or workspace rw
	RiskHigh               // 🔴 full network / privileged
)

func (r Risk) String() string {
	switch r {
	case RiskStrict:
		return "🔒"
	case RiskDev:
		return "🟡"
	case RiskHigh:
		return "🔴"
	default:
		return "🟡"
	}
}

type Mount struct {
	Path     string
	ReadOnly bool
}

type Sandbox struct {
	ID        string
	Name      string
	Agent     string
	State     State
	CPU       string
	Mem       string
	Risk      Risk
	Uptime    string
	Workspace string
	Template  string
	Network   string
	Mounts    []Mount
	Labels    map[string]string
	LastError string
	ExitCode  int
	Raw       []byte
}

type Client interface {
	List(ctx context.Context) ([]Sandbox, error)
	Inspect(ctx context.Context, id string) (Sandbox, error)
	Start(ctx context.Context, id string) error
	Stop(ctx context.Context, id string) error
	Restart(ctx context.Context, id string) error
	Remove(ctx context.Context, id string) error
	Logs(ctx context.Context, id string) (io.ReadCloser, error)
	Exec(ctx context.Context, id string) (*exec.Cmd, error)
}
