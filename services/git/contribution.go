package git

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Author struct {
	Name  string
	Email string
}

type AuthorInfo struct {
	TotalCommit   int
	TotalAddLines int
	TotalDelLines int
	LastCommit    time.Time
	FirstCommit   time.Time
}

func (r *GitConfig) GetAuthorInfoByAuthor() (map[Author]AuthorInfo, error) {
	cIter, err := r.repo.Log(&git.LogOptions{Order: git.LogOrderCommitterTime})
	if err != nil {
		return nil, err
	}

	authorInfoMap := make(map[Author]AuthorInfo)
	err = cIter.ForEach(func(c *object.Commit) error {
		currAuthor := Author{
			Name:  c.Author.Name,
			Email: c.Author.Email,
		}
		stats, err := c.Stats()
		if err != nil {
			return err
		}

		currAddition := 0
		currDeletion := 0
		for _, val := range stats {
			currAddition += val.Addition
			currDeletion += val.Deletion
		}

		if authorInfo, ok := authorInfoMap[currAuthor]; ok {
			authorInfo.LastCommit = c.Author.When
			authorInfo.TotalCommit++
			authorInfo.TotalAddLines += currAddition
			authorInfo.TotalDelLines += currDeletion
			authorInfoMap[currAuthor] = authorInfo
		} else {
			authorInfo := AuthorInfo{
				TotalCommit:   1,
				TotalAddLines: currAddition,
				TotalDelLines: currDeletion,
				LastCommit:    c.Author.When,
				FirstCommit:   c.Author.When,
			}
			authorInfoMap[currAuthor] = authorInfo
		}

		return nil
	})
	return authorInfoMap, err
}
