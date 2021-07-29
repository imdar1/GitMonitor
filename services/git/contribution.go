package git

import (
	"errors"
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
	startCommitHash, err := r.repo.ResolveRevision(
		plumbing.Revision(fmt.Sprintf("refs/remotes/%s/%s", remoteName, featureBranch)),
	)

	if err != nil {
		utils.CheckErr(serviceName, err)
		return nil, err
	}
	endCommitHash, err := r.repo.ResolveRevision(
		plumbing.Revision(fmt.Sprintf("refs/remotes/%s/%s", remoteName, defaultBranch)),
	)
	if err != nil {
		utils.CheckErr(serviceName, err)
		return nil, err
	}

	startCommit, err := r.repo.CommitObject(*startCommitHash)
	if err != nil {
		utils.CheckErr(serviceName, err)
		return nil, err
	}

	endCommit, err := r.repo.CommitObject(*endCommitHash)
	if err != nil {
		utils.CheckErr(serviceName, err)
		return nil, err
	}

	if err != nil {
		utils.CheckErr(serviceName, err)
		return nil, err
	}

	commonAncestorHash, err := endCommit.MergeBase(startCommit)
	if len(commonAncestorHash) == 0 {
		err := errors.New("common ancestor is not found")
		return nil, err
	}
	if err != nil {
		utils.CheckErr(serviceName, err)
		return nil, err
	}

	cIter, err := r.repo.Log(&git.LogOptions{
		From: *startCommitHash,
	})
	if err != nil {
		utils.CheckErr(serviceName, err)
		return nil, err
	}

	var commits []*object.Commit
	err = cIter.ForEach(func(c *object.Commit) error {
		if c.Hash == commonAncestorHash[0].Hash {
			return storer.ErrStop
		}
		commits = append(commits, c)
		return nil
	})
	return commits, err
}
