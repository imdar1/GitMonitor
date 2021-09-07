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
) ([]*object.Commit, *object.Commit, error) {
	const serviceName = "GetLogTwoBranches"
	startCommitHash, err := r.repo.ResolveRevision(
		plumbing.Revision(fmt.Sprintf("refs/remotes/%s/%s", remoteName, featureBranch)),
	)

	if err != nil {
		utils.CheckErr(serviceName, err)
		return nil, nil, err
	}
	endCommitHash, err := r.repo.ResolveRevision(
		plumbing.Revision(fmt.Sprintf("refs/remotes/%s/%s", remoteName, defaultBranch)),
	)
	if err != nil {
		utils.CheckErr(serviceName, err)
		return nil, nil, err
	}

	startCommit, err := r.repo.CommitObject(*startCommitHash)
	if err != nil {
		utils.CheckErr(serviceName, err)
		return nil, nil, err
	}

	endCommit, err := r.repo.CommitObject(*endCommitHash)
	if err != nil {
		utils.CheckErr(serviceName, err)
		return nil, nil, err
	}

	commonAncestors, err := endCommit.MergeBase(startCommit)
	if len(commonAncestors) == 0 {
		err := errors.New("common ancestor is not found")
		utils.CheckErr(serviceName, err)
		return nil, nil, err
	}
	if err != nil {
		utils.CheckErr(serviceName, err)
		return nil, nil, err
	}

	cIter, err := r.repo.Log(&git.LogOptions{
		From: *startCommitHash,
	})
	if err != nil {
		utils.CheckErr(serviceName, err)
		return nil, commonAncestors[0], err
	}

	var commits []*object.Commit
	err = cIter.ForEach(func(c *object.Commit) error {
		if c.Hash == commonAncestors[0].Hash {
			return storer.ErrStop
		}
		commits = append(commits, c)
		return nil
	})
	return commits, commonAncestors[0], err
}

func (r *GitConfig) GetDiff(
	mainCommit *object.Commit,
	baseCommit *object.Commit,
) (object.FileStats, error) {
	const serviceName = "GetDiff"

	patch, err := baseCommit.Patch(mainCommit)
	if err != nil {
		utils.CheckErr(serviceName, err)
		return nil, err
	}

	return patch.Stats(), nil
}
