package contribution

import (
	"gitmonitor/constants"
	"gitmonitor/models"
	"gitmonitor/services/utils"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
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

type ContributorData struct {
	authorMap         map[Author]AuthorInfo
	tasks             []models.Task
	defaultBranchName string
	defaultRemoteName string
}

func getAuthorInfoByAuthor(commits []*object.Commit) (map[Author]AuthorInfo, error) {
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

func getInProgressAndDoneTask(tasks []models.Task) []models.Task {
	var inProgressAndDoneTasks []models.Task
	for _, task := range tasks {
		if task.TaskStatus == int(constants.InProgress) || task.TaskStatus == int(constants.Done) ||
			task.TaskStatus == int(constants.DoneLate) {
			inProgressAndDoneTasks = append(inProgressAndDoneTasks, task)
		}
	}
	return inProgressAndDoneTasks
}

func InitContributorData(
	commits []*object.Commit,
	tasks []models.Task,
	defaultBranchName string,
	defaultRemoteName string,
) ContributorData {
	contributorMap, err := getAuthorInfoByAuthor(commits)
	utils.CheckErr("InitContributorData", err)

	return ContributorData{
		authorMap:         contributorMap,
		tasks:             getInProgressAndDoneTask(tasks),
		defaultBranchName: defaultBranchName,
		defaultRemoteName: defaultRemoteName,
	}
}
