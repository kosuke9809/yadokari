package sandbox_test

import (
	"context"
	"testing"

	"github.com/kosuke9809/yadokari/sandbox"
)

// MockClient が Client インターフェースを実装していることをコンパイル時に確認
var _ sandbox.Client = &sandbox.MockClient{}

func TestMockClientList(t *testing.T) {
	mock := &sandbox.MockClient{
		Sandboxes: sandbox.SampleSandboxes(),
	}
	got, err := mock.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Errorf("want 2 sandboxes, got %d", len(got))
	}
}

func TestMockClientOperations(t *testing.T) {
	mock := &sandbox.MockClient{
		Sandboxes: sandbox.SampleSandboxes(),
	}

	ctx := context.Background()

	if err := mock.Start(ctx, "abc"); err != nil {
		t.Fatal(err)
	}
	if err := mock.Stop(ctx, "def"); err != nil {
		t.Fatal(err)
	}
	if err := mock.Remove(ctx, "ghi"); err != nil {
		t.Fatal(err)
	}

	if len(mock.StartedIDs) != 1 || mock.StartedIDs[0] != "abc" {
		t.Errorf("StartedIDs: want [abc], got %v", mock.StartedIDs)
	}
	if len(mock.StoppedIDs) != 1 || mock.StoppedIDs[0] != "def" {
		t.Errorf("StoppedIDs: want [def], got %v", mock.StoppedIDs)
	}
	if len(mock.RemovedIDs) != 1 || mock.RemovedIDs[0] != "ghi" {
		t.Errorf("RemovedIDs: want [ghi], got %v", mock.RemovedIDs)
	}
}
