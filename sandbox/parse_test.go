package sandbox_test

import (
	"testing"

	"github.com/kosuke9809/yadokari/sandbox"
)

func TestParseSandboxList(t *testing.T) {
	input := []byte(`{
		"vms": [
			{
				"name": "refactor1",
				"agent": "claude",
				"status": "running",
				"socket_path": "/tmp/docker-abc",
				"workspaces": ["/home/user/project"]
			},
			{
				"name": "test-exp",
				"agent": "codex",
				"status": "stopped",
				"socket_path": "/tmp/docker-def",
				"workspaces": []
			}
		]
	}`)

	got, err := sandbox.ParseSandboxList(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 sandboxes, got %d", len(got))
	}

	// 1つ目
	if got[0].Name != "refactor1" {
		t.Errorf("want name=refactor1, got %s", got[0].Name)
	}
	if got[0].Agent != "claude" {
		t.Errorf("want agent=claude, got %s", got[0].Agent)
	}
	if got[0].State != sandbox.StateRunning {
		t.Errorf("want state=running, got %s", got[0].State)
	}
	if got[0].Workspace != "/home/user/project" {
		t.Errorf("want workspace=/home/user/project, got %s", got[0].Workspace)
	}

	// 2つ目
	if got[1].State != sandbox.StateStopped {
		t.Errorf("want state=stopped, got %s", got[1].State)
	}
	if got[1].Workspace != "" {
		t.Errorf("want workspace='', got %s", got[1].Workspace)
	}
}

func TestNormalizeState(t *testing.T) {
	// NormalizeState を外部テストするために ParseSandboxList 経由でテスト
	tests := []struct {
		status string
		want   sandbox.State
	}{
		{"running", sandbox.StateRunning},
		{"stopped", sandbox.StateStopped},
		{"exited", sandbox.StateStopped},
		{"crashed", sandbox.StateCrashed},
	}
	for _, tt := range tests {
		input := []byte(`{"vms":[{"name":"t","agent":"a","status":"` + tt.status + `","workspaces":[]}]}`)
		got, err := sandbox.ParseSandboxList(input)
		if err != nil {
			t.Fatal(err)
		}
		if got[0].State != tt.want {
			t.Errorf("status=%s: want %s, got %s", tt.status, tt.want, got[0].State)
		}
	}
}
