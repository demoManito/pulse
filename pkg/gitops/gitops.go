package gitops

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	// 展开 ~ 为用户主目录
	cfg.LocalPath = expandHome(cfg.LocalPath)

	// 从 URL 提取仓库名，拼接到 LocalPath 下
	// 例如 https://git.bilibili.co/tangjingyu/prd.git + ~/DocRepos → ~/DocRepos/prd
	repoName := repoNameFromURL(cfg.URL)
	if repoName != "" {
		cfg.LocalPath = filepath.Join(cfg.LocalPath, repoName)
	}

	auth, err := newAuth(cfg.Auth)
	if err != nil {
		return nil, err
	}

	return &GitOps{cfg: cfg, auth: auth}, nil
}

// expandHome 将路径中的 ~ 展开为用户主目录。
func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") || path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[1:])
	}
	return path
}

// repoNameFromURL 从 Git URL 中提取仓库名（去掉 .git 后缀）。
func repoNameFromURL(url string) string {
	// 处理末尾斜杠
	url = strings.TrimRight(url, "/")
	// 取最后一段
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return ""
	}
	name := parts[len(parts)-1]
	// 也处理 ssh 格式 git@host:org/repo.git
	if idx := strings.LastIndex(name, ":"); idx != -1 {
		name = name[idx+1:]
	}
	return strings.TrimSuffix(name, ".git")
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
