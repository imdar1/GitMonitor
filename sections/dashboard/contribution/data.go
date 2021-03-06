package contribution

import (
	"gitmonitor/services/git"
	"gitmonitor/services/utils"

	"github.com/go-git/go-git/v5/plumbing/object"
)

type authorTable struct {
	git.Author
	git.AuthorInfo
}

type ContributorData struct {
	authorMap map[git.Author]git.AuthorInfo
}

func InitContributorData(commits []*object.Commit, r git.GitConfig) ContributorData {
	contributorMap, err := r.GetAuthorInfoByAuthor(commits)
	utils.CheckErr("InitContributorData", err)

	return ContributorData{
		authorMap: contributorMap,
	}
}
