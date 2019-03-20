package repository

import (
	"fmt"
	"os/exec"
	"os/user"
	"path"
	"strings"
)

var (
	// GameSaveRoot is path to gamesave local Git repository
	GameSaveRoot string
)

func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	GameSaveRoot = path.Join(usr.HomeDir, ".gamesave")
}

// IGitRepository is interface for interaction with files
// that stored in Git using Git's commands
type IGitRepository interface {
	Checkout(branch string) error
	Commit(message string) error
	Clone(repoURL string) error
	FetchBranch(branch string) error
	GetCurrentBranch() (string, error)
	GetRepoURL() (string, error)
	Pull(branch string) error
	Push(branch string) error
	SetRepoURL(repoURL string) error
}

// GitRepository is the implementation of IGitRepository
type GitRepository struct{}

// Checkout change branch of Git repository
func (g *GitRepository) Checkout(branch string) error {
	cmd := exec.Command("git", "checkout", "-B", branch)
	cmd.Dir = GameSaveRoot
	err := cmd.Run()
	return err
}

// Commit adds all file and commit into remote
func (g *GitRepository) Commit(message string) error {
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = GameSaveRoot
	err := cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("git", "commit", "-m", message)
	cmd.Dir = GameSaveRoot
	return cmd.Run()
}

// Clone download repository from remote on repoURL
func (g *GitRepository) Clone(repoURL string) error {
	cmd := exec.Command("git", "clone", repoURL, GameSaveRoot)
	return cmd.Run()
}

// FetchBranch fetch specific branch from remote
func (g *GitRepository) FetchBranch(branch string) error {
	cmd := exec.Command(
		"git", "fetch", "origin",
		fmt.Sprintf("%s:%s", branch, branch),
	)
	cmd.Dir = GameSaveRoot
	return cmd.Run()
}

// GetCurrentBranch get current active branch
func (g *GitRepository) GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = GameSaveRoot
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), err
}

// GetRepoURL get URL of Git repository
func (g *GitRepository) GetRepoURL() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = GameSaveRoot
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), err
}

// Pull download repository from remote on specific branch
func (g *GitRepository) Pull(branch string) error {
	err := g.Checkout(branch)
	if err != nil {
		return err
	}
	cmd := exec.Command("git", "pull", "origin", branch)
	cmd.Dir = GameSaveRoot
	return cmd.Run()
}

// Push upload repository to remote on specific branch
func (g *GitRepository) Push(branch string) error {
	err := g.Checkout(branch)
	if err != nil {
		return err
	}
	cmd := exec.Command("git", "push", "origin", branch)
	cmd.Dir = GameSaveRoot
	return cmd.Run()
}

// SetRepoURL set URL of Git repository
func (g *GitRepository) SetRepoURL(repoURL string) error {
	cmd := exec.Command("git", "remote", "set-url", "origin", repoURL)
	cmd.Dir = GameSaveRoot
	return cmd.Run()
}
