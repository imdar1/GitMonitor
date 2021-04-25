package contribution

import (
	"gitmonitor/services/git"
	"gitmonitor/services/utils"

	"github.com/go-git/go-git/v5/plumbing/object"
)

type ContributorData struct {
	authorMap map[git.Author]git.AuthorInfo
}

func InitContributorData(commits []*object.Commit, r git.GitConfig) ContributorData {
	contributorMap, err := r.GetAuthorInfoByAuthor(commits)
	utils.CheckErr(err)

	return ContributorData{
		authorMap: contributorMap,
	}
}
