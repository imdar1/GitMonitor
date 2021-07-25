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
	Author
	AuthorInfo
}

func InitContributionTab() fyne.CanvasObject {
	return widget.NewLabel("Contribution Information")
}

func getFeatureBranchesListCanvas(data ContributorData, appData *data.AppData) fyne.CanvasObject {
	list := widget.NewList(
		func() int {
			return len(data.tasks)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Feature branch name"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			branch := appData.Database.GetBranchById(data.tasks[id].BranchId)
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(branch.Name)
		},
	)

	return list
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
					label.SetText(authorList[id.Row-1].Name)
				case 1:
					label.SetText(authorList[id.Row-1].Email)
				case 2:
					label.SetText(fmt.Sprintf("%d", authorList[id.Row-1].TotalCommit))
				case 3:
					label.SetText(string(authorList[id.Row-1].FirstCommit.Format("2 Jan 2006 15:04:05")))
				case 4:
					label.SetText(string(authorList[id.Row-1].LastCommit.Format("2 Jan 2006 15:04:05")))
				default:
					label.SetText(" ")
				}
			}
			if minLabelWidth[id.Col] < label.MinSize().Width {
				minLabelWidth[id.Col] = label.MinSize().Width
				go func() {
					table.SetColumnWidth(id.Col, minLabelWidth[id.Col])
					table.Refresh()
				}()
			}
		},
	)

	for i := 0; i <= 4; i++ {
		table.SetColumnWidth(i, 200)
	}

	featureContent := container.NewHSplit(
		container.NewVScroll(getFeatureBranchesListCanvas(data, appData)),
		container.NewVScroll(widget.NewCard(
			"Nama task-nya",
			"",
			widget.NewLabel("Informasi-informasi dari commit di feature branch tersebut"),
		)),
	)
	contributorContent := container.NewVSplit(featureContent, table)
	contributorWrapper := wrapper.(*widget.Card)
	contributorWrapper.SetContent(contributorContent)
}
