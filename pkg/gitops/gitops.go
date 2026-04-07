package gitops

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

// GitOps provides git operations on a repository.
type GitOps struct {
	cfg  Config
	repo *git.Repository
	auth transport.AuthMethod
}

// New creates a new GitOps instance.
func New(cfg Config) (*GitOps, error) {
	if cfg.Branch == "" {
		cfg.Branch = "main"
	}

	auth, err := newAuth(cfg.Auth)
	if err != nil {
		return nil, err
	}

	return &GitOps{cfg: cfg, auth: auth}, nil
}

// LocalPath 返回本地仓库路径。
func (g *GitOps) LocalPath() string {
	return g.cfg.LocalPath
}

// Clone clones the remote repository to LocalPath.
// If LocalPath already contains a valid repository, it opens and pulls instead.
func (g *GitOps) Clone(ctx context.Context) error {
	if _, err := os.Stat(g.cfg.LocalPath); err == nil {
		if err := g.Open(); err == nil {
			return g.Pull(ctx)
		}
	}

	repo, err := git.PlainCloneContext(ctx, g.cfg.LocalPath, false, &git.CloneOptions{
		URL:           g.cfg.URL,
		Auth:          g.auth,
		ReferenceName: plumbing.NewBranchReferenceName(g.cfg.Branch),
		SingleBranch:  true,
	})
	if err != nil {
		return fmt.Errorf("gitops: clone failed: %w", err)
	}
	g.repo = repo
	return nil
}

// Open opens an existing local repository.
func (g *GitOps) Open() error {
	repo, err := git.PlainOpen(g.cfg.LocalPath)
	if err != nil {
		return fmt.Errorf("gitops: open failed: %w", err)
	}
	g.repo = repo
	return nil
}

// Pull fetches and merges the latest changes from the remote.
func (g *GitOps) Pull(ctx context.Context) error {
	w, err := g.repo.Worktree()
	if err != nil {
		return fmt.Errorf("gitops: worktree failed: %w", err)
	}

	err = w.PullContext(ctx, &git.PullOptions{
		Auth:          g.auth,
		ReferenceName: plumbing.NewBranchReferenceName(g.cfg.Branch),
		SingleBranch:  true,
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return fmt.Errorf("gitops: pull failed: %w", err)
	}
	return nil
}

// Add stages files matching the given patterns to the index.
func (g *GitOps) Add(patterns ...string) error {
	w, err := g.repo.Worktree()
	if err != nil {
		return fmt.Errorf("gitops: worktree failed: %w", err)
	}

	for _, pattern := range patterns {
		if _, err := w.Add(pattern); err != nil {
			return fmt.Errorf("gitops: add %q failed: %w", pattern, err)
		}
	}
	return nil
}

// Commit creates a new commit with the staged changes.
func (g *GitOps) Commit(message string) error {
	w, err := g.repo.Worktree()
	if err != nil {
		return fmt.Errorf("gitops: worktree failed: %w", err)
	}

	_, err = w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "pulse",
			Email: "pulse@localhost",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("gitops: commit failed: %w", err)
	}
	return nil
}

// Push pushes local commits to the remote.
func (g *GitOps) Push(ctx context.Context) error {
	err := g.repo.PushContext(ctx, &git.PushOptions{
		Auth: g.auth,
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return fmt.Errorf("gitops: push failed: %w", err)
	}
	return nil
}
