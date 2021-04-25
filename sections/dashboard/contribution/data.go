package contribution

import (
	"gitmonitor/services/git"
	"gitmonitor/services/utils"
)

type ContributorData struct {
	authorMap map[git.Author]git.AuthorInfo
}

func InitContributorData(r git.GitConfig) ContributorData {
	contributorMap, err := r.GetAuthorInfoByAuthor()
	utils.CheckErr(err)

	return ContributorData{
		authorMap: contributorMap,
	}
}
