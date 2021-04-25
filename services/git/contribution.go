package git

import (
	"time"

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

func (r *GitConfig) GetAuthorInfoByAuthor(commits []*object.Commit) (map[Author]AuthorInfo, error) {
	authorInfoMap := make(map[Author]AuthorInfo)
	for _, c := range commits {
		currAuthor := Author{
			Name:  c.Author.Name,
			Email: c.Author.Email,
		}
		stats, err := c.Stats()
		if err != nil {
			return authorInfoMap, err
		}

		currAddition := 0
		currDeletion := 0
		for _, val := range stats {
			currAddition += val.Addition
			currDeletion += val.Deletion
		}

		if authorInfo, ok := authorInfoMap[currAuthor]; ok {
			authorInfo.FirstCommit = c.Author.When
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
	}
	return authorInfoMap, nil
}
