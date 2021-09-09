package contribution

import (
	"gitmonitor/constants"
	"gitmonitor/models"
	"gitmonitor/sections/data"
	"gitmonitor/services/utils"
	"time"

	"fyne.io/fyne/v2"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type author struct {
	Name  string
	Email string
}

type authorInfo struct {
	TotalCommit int
	// TotalAddLines int
	// TotalDelLines int
	LastCommit  time.Time
	FirstCommit time.Time
}

type ContributorData struct {
	Wrapper fyne.CanvasObject

	tasks             []models.Task
	authorMap         map[author]authorInfo
	defaultBranchName string
	defaultRemoteName string
}

func (data ContributorData) Render(appData *data.AppData) {
	renderContributorTab(data, appData)
}

func (data *ContributorData) SetTasks(tasks []models.Task) {
	data.tasks = getInProgressTask(tasks)
}

func getAuthorInfoByAuthor(commits []*object.Commit) (map[author]authorInfo, error) {
	authorInfoMap := make(map[author]authorInfo)
	for _, c := range commits {
		currAuthor := author{
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

		if authorInfoData, ok := authorInfoMap[currAuthor]; ok {
			authorInfoData.FirstCommit = c.Author.When
			authorInfoData.TotalCommit++
			// authorInfo.TotalAddLines += currAddition
			// authorInfo.TotalDelLines += currDeletion
			authorInfoMap[currAuthor] = authorInfoData
		} else {
			authorInfo := authorInfo{
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

func getInProgressTask(tasks []models.Task) []models.Task {
	var inProgressTasks []models.Task
	for _, task := range tasks {
		if task.TaskStatus == int(constants.InProgress) {
			inProgressTasks = append(inProgressTasks, task)
		}
	}
	return inProgressTasks
}

func InitContributorData(
	wrapper fyne.CanvasObject,
	commits []*object.Commit,
	tasks []models.Task,
	defaultBranchName string,
	defaultRemoteName string,
) ContributorData {
	contributorMap, err := getAuthorInfoByAuthor(commits)
	utils.CheckErr("InitContributorData", err)

	return ContributorData{
		authorMap:         contributorMap,
		tasks:             getInProgressTask(tasks),
		defaultBranchName: defaultBranchName,
		defaultRemoteName: defaultRemoteName,
		Wrapper:           wrapper,
	}
}
