package git

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGitService_LifecycleAndGraph(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	svc := NewService(tmpDir)
	if err := svc.Init(ctx); err != nil {
		t.Fatalf("git init failed: %v", err)
	}

	// configure dummy author in test repo
	_ = exec.Command("git", "-C", tmpDir, "config", "user.name", "ForgeNexa Tester").Run()
	_ = exec.Command("git", "-C", tmpDir, "config", "user.email", "test@forgenexa.org").Run()

	// Initial status
	st, err := svc.GetStatus(ctx)
	if err != nil {
		t.Fatalf("get status failed: %v", err)
	}
	if !st.IsClean {
		t.Fatalf("expected clean repo after init")
	}

	// Create test file
	filePath := filepath.Join(tmpDir, "hello.txt")
	if err := os.WriteFile(filePath, []byte("Hello ForgeNexa!"), 0644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	st, err = svc.GetStatus(ctx)
	if err != nil {
		t.Fatalf("get status failed: %v", err)
	}
	if len(st.Untracked) != 1 || st.Untracked[0] != "hello.txt" {
		t.Fatalf("expected untracked hello.txt, got %v", st.Untracked)
	}

	// Stage
	if err := svc.Stage(ctx, "hello.txt"); err != nil {
		t.Fatalf("stage failed: %v", err)
	}

	st, err = svc.GetStatus(ctx)
	if err != nil {
		t.Fatalf("get status failed: %v", err)
	}
	if len(st.Staged) != 1 || st.Staged[0] != "hello.txt" {
		t.Fatalf("expected staged hello.txt, got %v", st.Staged)
	}

	// Commit
	if err := svc.Commit(ctx, "feat: initial commit"); err != nil {
		t.Fatalf("commit failed: %v", err)
	}

	// Verify Commit Log
	logs, err := svc.GetLog(ctx, 10)
	if err != nil {
		t.Fatalf("get log failed: %v", err)
	}
	if len(logs) != 1 {
		t.Fatalf("expected 1 commit in log, got %d", len(logs))
	}
	if logs[0].Message != "feat: initial commit" {
		t.Fatalf("commit message mismatch: expected 'feat: initial commit', got '%s'", logs[0].Message)
	}
}
