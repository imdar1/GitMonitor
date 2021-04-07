package services

import (
	"github.com/go-git/go-git/v5"
)

type GitConfig struct {
	repo *git.Repository
}

func InitGit(path string) (GitConfig, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return GitConfig{}, err
	}

	g := GitConfig{
		repo: r,
	}
	return g, err
}

func (g *GitConfig) OnRepositoryLoaded() {
	// TODO: retrieve general information

	// TODO: retrieve commit information (contributitors & commit)

	// TODO: retrieve task data from DB
}
