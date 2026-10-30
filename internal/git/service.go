package git

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Status represents high-level Git working tree state.
type Status struct {
	Branch       string   `json:"branch"`
	IsClean      bool     `json:"isClean"`
	Ahead        int      `json:"ahead"`
	Behind       int      `json:"behind"`
	Staged       []string `json:"staged"`
	Unstaged     []string `json:"unstaged"`
	Untracked    []string `json:"untracked"`
	ConflictList []string `json:"conflictList"`
}

// Service provides Git repository orchestration across Forge IDE and Forge Git.
type Service struct {
	repoDir string
}

// NewService creates a new Git service scoped to repoDir.
func NewService(repoDir string) *Service {
	return &Service{repoDir: repoDir}
}

func (s *Service) runGit(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = s.repoDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s failed: %w (output: %s)", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

// GetStatus parses git status with machine-readable porcelain format.
func (s *Service) GetStatus(ctx context.Context) (*Status, error) {
	branch, _ := s.runGit(ctx, "rev-parse", "--abbrev-ref", "HEAD")
	if branch == "" {
		branch = "HEAD"
	}

	out, err := s.runGit(ctx, "status", "--porcelain=v1")
	if err != nil {
		return nil, err
	}

	st := &Status{
		Branch:    branch,
		IsClean:   len(out) == 0,
		Staged:    []string{},
		Unstaged:  []string{},
		Untracked: []string{},
	}

	if len(out) == 0 {
		return st, nil
	}

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if len(line) < 3 {
			continue
		}
		indexStatus := line[0]
		workTreeStatus := line[1]
		file := strings.TrimSpace(line[3:])

		if indexStatus == '?' && workTreeStatus == '?' {
			st.Untracked = append(st.Untracked, file)
			continue
		}
		if indexStatus == 'U' || workTreeStatus == 'U' {
			st.ConflictList = append(st.ConflictList, file)
			continue
		}
		if indexStatus != ' ' && indexStatus != '?' {
			st.Staged = append(st.Staged, file)
		}
		if workTreeStatus != ' ' && workTreeStatus != '?' {
			st.Unstaged = append(st.Unstaged, file)
		}
	}

	return st, nil
}

// Stage stages specified files or all files if empty.
func (s *Service) Stage(ctx context.Context, files ...string) error {
	if len(files) == 0 {
		_, err := s.runGit(ctx, "add", "-A")
		return err
	}
	args := append([]string{"add", "--"}, files...)
	_, err := s.runGit(ctx, args...)
	return err
}

// Commit creates a new Git commit.
func (s *Service) Commit(ctx context.Context, message string) error {
	if strings.TrimSpace(message) == "" {
		return fmt.Errorf("commit message cannot be empty")
	}
	_, err := s.runGit(ctx, "commit", "-m", message)
	return err
}

// Init initializes a new git repository in repoDir.
func (s *Service) Init(ctx context.Context) error {
	_, err := s.runGit(ctx, "init")
	return err
}

// LogCommit represents a node in the visual Git commit graph.
type LogCommit struct {
	Hash      string    `json:"hash"`
	ShortHash string    `json:"shortHash"`
	Author    string    `json:"author"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Parents   []string  `json:"parents"`
}

// GetLog returns the recent topological commit history with parent hashes.
func (s *Service) GetLog(ctx context.Context, maxCount int) ([]LogCommit, error) {
	if maxCount <= 0 {
		maxCount = 50
	}
	format := "%H%x00%h%x00%an%x00%ae%x00%at%x00%s%x00%P"
	out, err := s.runGit(ctx, "log", fmt.Sprintf("-n%d", maxCount), fmt.Sprintf("--pretty=format:%s", format))
	if err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return []LogCommit{}, nil
	}

	var commits []LogCommit
	entries := strings.Split(out, "\n")
	for _, entry := range entries {
		parts := strings.Split(entry, "\x00")
		if len(parts) < 7 {
			continue
		}
		var parents []string
		if len(parts[6]) > 0 {
			parents = strings.Fields(parts[6])
		}

		commits = append(commits, LogCommit{
			Hash:      parts[0],
			ShortHash: parts[1],
			Author:    parts[2],
			Email:     parts[3],
			Message:   parts[5],
			Parents:   parents,
		})
	}

	return commits, nil
}
