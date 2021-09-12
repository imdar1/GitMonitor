package contribution

import (
	"fmt"
	"gitmonitor/sections/data"
	"gitmonitor/services/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type authorTable struct {
	author     author
	authorInfo authorInfo
}

func fillCommitString(
	textGrid *widget.TextGrid,
	data ContributorData,
	appData *data.AppData,
	featureBranchName string,
	selectedIndex int,
) {
	commits, commonAncestorCommit, err := appData.Repo.GetLogTwoBranches(
		data.defaultBranchName,
		featureBranchName,
		data.defaultRemoteName,
	)

	var commitsString string
	if err != nil || len(commits) == 0 {
		commitsString = "No commit found"
	}
	totalAddition := 0
	totalDeletion := 0
	for index, commit := range commits {
		baseCommit := commonAncestorCommit
		if index < len(commits)-1 {
			baseCommit = commits[index+1]
		}
		stats, err := appData.Repo.GetDiff(commit, baseCommit)

		commitsString = fmt.Sprintf("%s%s\n", commitsString, commit.String())
		utils.CheckErr("getFeatureBranchesListCanvas", err)
		for _, stat := range stats {
			totalAddition += stat.Addition
			totalDeletion += stat.Deletion
		}
	}
	if len(commits) == 0 {
		commitsString = "No commit found"
	}
	commitsString = fmt.Sprintf(
		"Author: %s\nTotal Addition: %d\nTotal Deletion: %d\nCommit History:\n\n%s",
		data.tasks[selectedIndex].AssigneeName,
		totalAddition,
		totalDeletion,
		commitsString,
	)
	textGrid.SetText(commitsString)
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
		selectedBranchName := branches[id]
		if selectedBranchName == "" {
			grid := widget.NewTextGridFromString("No associated branch found")
			taskDetailCard.SetContent(grid)
			return
		}

		taskDetailCard.SetTitle(selectedBranchName)
		grid := widget.NewTextGridFromString(fmt.Sprintf("Analyzing %s branch...", selectedBranchName))
		taskDetailCard.SetContent(grid)
		go fillCommitString(grid, data, appData, selectedBranchName, id)
	}

	taskList.OnUnselected = func(id widget.ListItemID) {
		taskDetailCard.SetTitle("")
		taskDetailCard.SetContent(nil)
	}

	featureContent := container.NewHSplit(
		container.NewScroll(taskList),
		container.NewScroll(taskDetailCard),
	)

	return featureContent
}

func renderContributorTab(data ContributorData, appData *data.AppData) {
	authorList := []authorTable{}
	for key, value := range data.authorMap {
		authorList = append(authorList, authorTable{key, value})
	}

	// TODO: Change this hardcoded size into the real min width after the race condition mitigated
	minLabelWidth := []float32{208, 373, 80, 171, 163}

	table := widget.NewTable(
		func() (int, int) {
			return len(data.authorMap) + 1, 5
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
		},
	)
	for i, width := range minLabelWidth {
		table.SetColumnWidth(i, width)
	}

	featureContent := getFeatureBranchesListCanvas(data, appData)
	contributorContent := container.NewVSplit(featureContent, table)
	contributorWrapper := data.Wrapper.(*widget.Card)
	contributorWrapper.SetContent(contributorContent)
}
