package git

import (
	"fmt"
	"gitmonitor/constants"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetAuthors(commits []*object.Commit) []string {
	set := make(map[string]bool)
	var signatures []string
	for _, v := range commits {
		authorFormat := fmt.Sprintf("%s%s%s", v.Author.Email, constants.Separator, v.Author.Name)
		set[authorFormat] = true
	}

	for k := range set {
		signatures = append(signatures, k)
	}

	return signatures
}

func (r *GitConfig) GetOriginUrl() string {
	origin, err := r.repo.Remote("origin")
	if err != nil {
		return ""
	}
	return origin.Config().URLs[0]
}

func (r *GitConfig) GetPaths() ([]string, error) {
	refs, err := r.repo.Head()
	if err != nil {
		return []string{}, err
	}

	commit, err := r.repo.CommitObject(refs.Hash())
	if err != nil {
		return []string{}, err
	}
	tree, err := commit.Tree()
	if err != nil {
		return []string{}, err
	}

	var paths []string
	for _, entry := range tree.Entries {
		paths = append(paths, entry.Name)
	}
	fmt.Println(paths)

	return paths, nil
}

func (r *GitConfig) GetCommitObjects() ([]*object.Commit, error) {
	cIter, err := r.repo.Log(&git.LogOptions{Order: git.LogOrderCommitterTime})
	if err != nil {
		return nil, nil
	}
	var commits []*object.Commit
	cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})
	return commits, err
}
