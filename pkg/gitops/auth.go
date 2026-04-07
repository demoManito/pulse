package gitops

import (
	"fmt"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

// newAuth creates a transport.AuthMethod from the given config.
func newAuth(cfg AuthConfig) (transport.AuthMethod, error) {
	switch cfg.Type {
	case "http":
		return &http.BasicAuth{
			Username: "token",
			Password: cfg.Token,
		}, nil
	case "ssh":
		keys, err := ssh.NewPublicKeysFromFile("git", cfg.SSHKey, cfg.Password)
		if err != nil {
			return nil, fmt.Errorf("gitops: failed to load SSH key %s: %w", cfg.SSHKey, err)
		}
		return keys, nil
	case "":
		return nil, nil
	default:
		return nil, fmt.Errorf("gitops: unsupported auth type: %s", cfg.Type)
	}
}
