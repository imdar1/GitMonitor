package config

import "github.com/go-git/go-git/v5"

type GitConfig struct {
	Repo *git.Repository
}

func InitGit(path string) (GitConfig, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return GitConfig{}, err
	}

	g := GitConfig{
		Repo: r,
	}
	return g, err
}
