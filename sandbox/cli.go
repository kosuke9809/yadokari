package sandbox

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"time"
)

const defaultTimeout = 5 * time.Second

// CLIClient は docker sandbox CLI をラップする
type CLIClient struct {
	timeout time.Duration
}

// NewCLIClient は新しい CLIClient を返す
func NewCLIClient() *CLIClient {
	return &CLIClient{timeout: defaultTimeout}
}

// run は docker sandbox コマンドを実行し stdout を返す
func (c *CLIClient) run(ctx context.Context, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "docker", append([]string{"sandbox"}, args...)...)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%w: %s", err, errOut.String())
	}
	return out.Bytes(), nil
}

// List は docker sandbox ls --json を実行してパースする
func (c *CLIClient) List(ctx context.Context) ([]Sandbox, error) {
	data, err := c.run(ctx, "ls", "--json")
	if err != nil {
		return nil, err
	}
	return ParseSandboxList(data)
}

// Inspect は ls --json から対象サンドボックスを探して返す
// inspect コマンドが存在しないため List の結果を使う
func (c *CLIClient) Inspect(ctx context.Context, id string) (Sandbox, error) {
	sandboxes, err := c.List(ctx)
	if err != nil {
		return Sandbox{}, err
	}
	for _, s := range sandboxes {
		if s.ID == id || s.Name == id {
			return s, nil
		}
	}
	return Sandbox{}, fmt.Errorf("sandbox not found: %s", id)
}

// Logs は未サポート（docker sandbox logs コマンドが存在しない）
func (c *CLIClient) Logs(_ context.Context, _ string) (io.ReadCloser, error) {
	return nil, errors.New("logs not supported by docker sandbox CLI")
}

// Exec は docker sandbox exec -it <name> /bin/sh を返す
// タイムアウトなし
func (c *CLIClient) Exec(_ context.Context, id string) (*exec.Cmd, error) {
	cmd := exec.Command("docker", "sandbox", "exec", "-i", "-t", id, "/bin/sh")
	return cmd, nil
}

// Start は docker sandbox run <name> で停止中のサンドボックスを起動する
func (c *CLIClient) Start(ctx context.Context, id string) error {
	_, err := c.run(ctx, "run", id)
	return err
}

// Stop は docker sandbox stop <name> でサンドボックスを停止する
func (c *CLIClient) Stop(ctx context.Context, id string) error {
	_, err := c.run(ctx, "stop", id)
	return err
}

// Restart は Stop → Start の順で再起動する
func (c *CLIClient) Restart(ctx context.Context, id string) error {
	if err := c.Stop(ctx, id); err != nil {
		return err
	}
	return c.Start(ctx, id)
}

// Remove は docker sandbox rm <name> でサンドボックスを削除する
func (c *CLIClient) Remove(ctx context.Context, id string) error {
	_, err := c.run(ctx, "rm", id)
	return err
}
