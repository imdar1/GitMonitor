package git

import (
	"fmt"
	"gitmonitor/services/utils"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

func (r *GitConfig) GetLogTwoBranches(
	defaultBranch string,
	featureBranch string,
	remoteName string,
) ([]*object.Commit, error) {
	const serviceName = "GetLogTwoBranches"
	startCommit, err := r.repo.ResolveRevision(
		plumbing.Revision(fmt.Sprintf("refs/remotes/%s/%s", remoteName, defaultBranch)),
	)
	if err != nil {
		utils.CheckErr(serviceName, err)
		return nil, err
	}
	endCommit, err := r.repo.ResolveRevision(
		plumbing.Revision(fmt.Sprintf("refs/remotes/%s/%s", remoteName, defaultBranch)),
	)
	if err != nil {
		utils.CheckErr(serviceName, err)
		return nil, err
	}

	cIter, err := r.repo.Log(&git.LogOptions{
		From: *endCommit,
	})
	if err != nil {
		utils.CheckErr(serviceName, err)
		return nil, err
	}

	var commits []*object.Commit
	err = cIter.ForEach(func(c *object.Commit) error {
		if c.Hash == *startCommit {
			return storer.ErrStop
		}
		commits = append(commits, c)
		return nil
	})
	return commits, err
}
