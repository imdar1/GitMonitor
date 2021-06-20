package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

type GitConfig struct {
	repo *git.Repository
	auth transport.AuthMethod
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
