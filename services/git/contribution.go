package git

import (
	"fmt"
	"gitmonitor/services/utils"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type Author struct {
	Name  string
	Email string
}

type AuthorInfo struct {
	TotalCommit int
	// TotalAddLines int
	// TotalDelLines int
	LastCommit  time.Time
	FirstCommit time.Time
}

func (r *GitConfig) GetAuthorInfoByAuthor(commits []*object.Commit) (map[Author]AuthorInfo, error) {
	authorInfoMap := make(map[Author]AuthorInfo)
	for _, c := range commits {
		currAuthor := Author{
			Name:  c.Author.Name,
			Email: c.Author.Email,
		}

		// Disable stats for now cuz it's slow
		// stats, err := c.Stats()
		// if err != nil {
		// 	return authorInfoMap, err
		// }

		// currAddition := 0
		// currDeletion := 0
		// for _, val := range stats {
		// 	currAddition += val.Addition
		// 	currDeletion += val.Deletion
		// }

		if authorInfo, ok := authorInfoMap[currAuthor]; ok {
			authorInfo.FirstCommit = c.Author.When
			authorInfo.TotalCommit++
			// authorInfo.TotalAddLines += currAddition
			// authorInfo.TotalDelLines += currDeletion
			authorInfoMap[currAuthor] = authorInfo
		} else {
			authorInfo := AuthorInfo{
				TotalCommit: 1,
				// TotalAddLines: currAddition,
				// TotalDelLines: currDeletion,
				LastCommit:  c.Author.When,
				FirstCommit: c.Author.When,
			}
			authorInfoMap[currAuthor] = authorInfo
		}
	}
	return authorInfoMap, nil
}

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
