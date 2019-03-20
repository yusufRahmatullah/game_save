package service

import (
	"errors"
)

const (
	gitRepoMock = "local.git"
)

var (
	gitOptionNormal = map[string]bool{
		"branch_exist": true,
		"repo_url":     true,
	}
	gitOptionRepoInvalid = map[string]bool{
		"branch_exist": true,
		"repo_url":     false,
	}
	gitOptionsBranchInvalid = map[string]bool{
		"branch_exist": false,
		"repo_url":     true,
	}
)

type GitRepositoryMock struct {
	currentBranch string
	options       map[string]bool
}
type OsRepositoryMock struct {
	gameName string
	savePath string
}

func NewGitRepositoryMock(options map[string]bool) *GitRepositoryMock {
	return &GitRepositoryMock{
		currentBranch: "",
		options:       options,
	}
}

func (g *GitRepositoryMock) Checkout(branch string) error {
	g.currentBranch = branch
	return nil
}

func (g *GitRepositoryMock) Commit(message string) error {
	if val, _ := g.options["repo_url"]; !val {
		return errors.New("")
	}
	return nil
}

func (g *GitRepositoryMock) Clone(repoURL string) error {
	if val, _ := g.options["repo_url"]; !val {
		return errors.New("")
	}
	return nil
}

func (g *GitRepositoryMock) FetchBranch(branch string) error {
	return nil
}

func (g *GitRepositoryMock) GetCurrentBranch() (string, error) {
	return g.currentBranch, nil
}

func (g *GitRepositoryMock) GetRepoURL() (string, error) {
	return gitRepoMock, nil
}

func (g *GitRepositoryMock) Pull(gameName string) error {
	if val, _ := g.options["repo_url"]; !val {
		if val2, _ := g.options["branch_exist"]; !val2 {
			return errors.New("")
		}
	}
	return nil
}

func (g *GitRepositoryMock) Push(gameName string) error {
	if val, _ := g.options["repo_url"]; !val {
		if val2, _ := g.options["branch_exist"]; !val2 {
			return errors.New("")
		}
	}
	return nil
}

func (g *GitRepositoryMock) SetRepoURL(repoURL string) error {
	return nil
}

func NewOsRepositoryMock() *OsRepositoryMock {
	return &OsRepositoryMock{}
}

func (o *OsRepositoryMock) Copy(src, dst string) error {
	return nil
}

func (o *OsRepositoryMock) GetConfig(key string) string {
	value := ""
	switch key {
	case "game_name":
		value = o.gameName
	case "save_path":
		value = o.savePath
	}
	return value
}

func (o *OsRepositoryMock) SetConfig(key, value string) error {
	switch key {
	case "game_name":
		o.gameName = value
	case "save_path":
		o.savePath = value
	}
	return nil
}
