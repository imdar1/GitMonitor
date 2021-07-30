package contribution

import (
	"fmt"
	"gitmonitor/sections/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type authorTable struct {
	author     Author
	authorInfo AuthorInfo
}

func InitContributionTab() fyne.CanvasObject {
	return widget.NewLabel("Contribution Information")
}

func getFeatureBranchesListCanvas(data ContributorData, appData *data.AppData) fyne.CanvasObject {
	var branches []string
	for _, task := range data.tasks {
		branches = append(branches, appData.Database.GetBranchById(task.BranchId).Name)
	}

	taskList := widget.NewList(
		func() int {
			return len(data.tasks)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Feature branch name"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data.tasks[id].Name)
		},
	)
	taskDetailCard := widget.NewCard(
		"",
		"",
		nil,
	)

	taskList.OnSelected = func(id widget.ListItemID) {
		taskDetailCard.SetTitle(branches[id])
		commits, err := appData.Repo.GetLogTwoBranches(
			data.defaultBranchName,
			branches[id],
			data.defaultRemoteName,
		)

		var commitsString string
		if err != nil || len(commits) == 0 {
			commitsString = "No commit found"
		}
		for _, commit := range commits {
			commitsString = fmt.Sprintf("%s%s\n", commitsString, commit.String())
		}
		grid := widget.NewTextGridFromString(commitsString)
		taskDetailCard.SetContent(grid)
	}

	taskList.OnUnselected = func(id widget.ListItemID) {
		taskDetailCard.SetTitle("")
		taskDetailCard.SetContent(nil)
	}

	featureContent := container.NewHSplit(
		container.NewVScroll(taskList),
		container.NewVScroll(taskDetailCard),
	)

	return featureContent
}

func RenderContributorTab(wrapper fyne.CanvasObject, data ContributorData, appData *data.AppData) {
	authorList := []authorTable{}
	for key, value := range data.authorMap {
		authorList = append(authorList, authorTable{key, value})
	}

	minLabelWidth := []float32{0, 0, 0, 0, 0}

	var table *widget.Table
	table = widget.NewTable(
		func() (int, int) {
			return len(data.authorMap), 5
		},
		func() fyne.CanvasObject {
			return widget.NewLabel(" ")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			if id.Row == 0 {
				// set header
				switch id.Col {
				case 0:
					label.SetText("Author Name")
				case 1:
					label.SetText("Author Mail")
				case 2:
					label.SetText("Commits")
				case 3:
					label.SetText("First commit")
				case 4:
					label.SetText("Last commit")
				default:
					label.SetText(" ")
				}
			} else {
				// set content
				switch id.Col {
				case 0:
					label.SetText(authorList[id.Row-1].author.Name)
				case 1:
					label.SetText(authorList[id.Row-1].author.Email)
				case 2:
					label.SetText(fmt.Sprintf("%d", authorList[id.Row-1].authorInfo.TotalCommit))
				case 3:
					label.SetText(string(authorList[id.Row-1].authorInfo.FirstCommit.Format("2 Jan 2006 15:04:05")))
				case 4:
					label.SetText(string(authorList[id.Row-1].authorInfo.LastCommit.Format("2 Jan 2006 15:04:05")))
				default:
					label.SetText(" ")
				}
			}
			if minLabelWidth[id.Col] < label.MinSize().Width {
				minLabelWidth[id.Col] = label.MinSize().Width
				go func(id int) {
					table.SetColumnWidth(id, minLabelWidth[id])
					table.Refresh()
				}(id.Col)
			}
		},
	)

	featureContent := getFeatureBranchesListCanvas(data, appData)
	contributorContent := container.NewVSplit(featureContent, table)
	contributorWrapper := wrapper.(*widget.Card)
	contributorWrapper.SetContent(contributorContent)
}
