package gitops

// Config holds git repository configuration.
type Config struct {
	// URL is the remote repository URL.
	URL string `yaml:"url"`
	// LocalPath is the local directory to clone into.
	LocalPath string `yaml:"local_path"`
	// Branch is the branch name. Default: "main".
	Branch string `yaml:"branch"`
	// Auth holds authentication configuration.
	Auth AuthConfig `yaml:"auth"`
}

// AuthConfig holds git authentication configuration.
type AuthConfig struct {
	// Type is the auth type: "ssh" or "http".
	Type string `yaml:"type"`
	// Token is the HTTP personal access token (used when Type="http").
	Token string `yaml:"token"`
	// SSHKey is the path to the SSH private key file (used when Type="ssh").
	SSHKey string `yaml:"ssh_key"`
	// Password is the SSH private key passphrase (optional).
	Password string `yaml:"password"`
}
